// Package main
// Created by xuzhuoxi
// on 2019-03-24.
// @author xuzhuoxi
//
package main

import (
	"encoding/binary"
	"fmt"
	"github.com/xuzhuoxi/Rabbit-Home/core"
	"github.com/xuzhuoxi/Rabbit-Home/core/client"
	"github.com/xuzhuoxi/Rabbit-Server/demo/client/net"
	"github.com/xuzhuoxi/Rabbit-Server/demo/client/proto/login"
	"github.com/xuzhuoxi/infra-go/bytex"
	"github.com/xuzhuoxi/infra-go/cryptox"
	"github.com/xuzhuoxi/infra-go/cryptox/asymmetric"
	"github.com/xuzhuoxi/infra-go/cryptox/key"
	"github.com/xuzhuoxi/infra-go/cryptox/symmetric"
	"github.com/xuzhuoxi/infra-go/encodingx"
	"github.com/xuzhuoxi/infra-go/encodingx/jsonx"
	"github.com/xuzhuoxi/infra-go/filex"
	"github.com/xuzhuoxi/infra-go/logx"
	"github.com/xuzhuoxi/infra-go/netx"
	"github.com/xuzhuoxi/infra-go/netx/tcpx"
	"net/http"
	"time"
)

var (
	sleep                             = make(chan struct{})
	order            binary.ByteOrder = binary.LittleEndian
	dataBlockHandler                  = bytex.NewDataBlockHandler(order, bytex.DefaultDataToBlockHandler, bytex.DefaultBlockToDataHandler)
	logger           logx.ILogger
)

var (
	homeUrl         = "http://127.0.0.1:9000"
	encrypt         = true
	publicPemPath   = "keys/x509_public.pem"
	publicRsaCipher asymmetric.IRSAPublicCipher
	queryInfo       = &core.QueryRouteInfo{
		PlatformId: "main01",
		TypeName:   "Rabbit-Server",
		TempAesKey: key.DeriveKeyPbkdf2StrDefault("000000"),
	}
	queryBackInfo = &core.QueryRouteBackInfo{
		OpenNetwork: "tcp",
		OpenAddr:    "127.0.0.1:41000",
		OpenKeyOn:   false,
	}
	openCipher cryptox.ICipher
	userClient *net.UserClient
	freq       = time.Millisecond * 500
)

func init() {
	bytex.DefaultOrder, bytex.DefaultDataBlockHandler = order, dataBlockHandler         // 包 bytex 下的大小端设置，封包处理
	encodingx.DefaultOrder, encodingx.DefaultDataBlockHandler = order, dataBlockHandler // 包 encodingx下的大小端设置，封包处理
	tcpx.TcpDataBlockHandler = dataBlockHandler                                         // Tcp封包处理
	jsonx.DefaultDataBlockHandler = dataBlockHandler
	logger = logx.NewLogger()
	logger.SetConfig(logx.LogConfig{Type: logx.TypeConsole, Level: logx.LevelAll})
}

func main() {
	if encrypt {
		loadHomePublicKey()
		doHomeQuery()
	} else {
		doConnServer()
	}
	<-sleep
}

func loadHomePublicKey() {
	public, err := asymmetric.LoadPublicCipherPEM(filex.FixFilePath(publicPemPath))
	if nil != err {
		panic(err)
	}
	publicRsaCipher = public
}

func doHomeQuery() {
	client.QueryRouteWithGet(homeUrl, *queryInfo, publicRsaCipher, func(res *http.Response, body *[]byte) {
		suc, fail, err := client.ParseQueryRouteBackInfo(res, body)
		if nil != err || nil != fail {
			logger.Warnln("RouteBackInfo Error.", err, fail)
			return
		}
		err = suc.ComputeOpenSK(queryInfo.TempAesKey)
		if nil != err {
			logger.Warnln("ComputeOpenSK Error.", err)
			return
		}
		logger.Infoln("Query Home Suc:", suc)
		queryBackInfo = suc
		openCipher = symmetric.NewAESCipher(suc.OpenSK)
		doConnServer()
	})
}

func doConnServer() {
	err := openClient()
	if nil != err {
		return
	}
	startClientListeners()
	doLogin()
}
func openClient() error {
	userId := "uid_01"
	uc := net.NewUserClient(userId)
	err := uc.OpenWitAddr(queryBackInfo.OpenAddr)
	if nil != err {
		logger.Warnln("OpenClient Error.", err)
		return err
	}
	userClient = uc
	return nil
}

func startClientListeners() {
	userClient.SockClient.GetPackHandlerContainer().SetPackHandlers([]netx.FuncPackHandler{onPack})
	go userClient.SockClient.StartReceiving()
}

func doLogin() {
	go func() {
		for {
			login.TestLoginExtension(userClient, openCipher)
			time.Sleep(freq)
		}
	}()
}

func onPack(data []byte, connInfo netx.IConnInfo, other interface{}) (catch bool) {
	fmt.Println("Rabbit-Server:Demo-Client.onPack:", len(data), connInfo)
	dataBlock := bytex.NewBuffDataBlock(bytex.NewDefaultDataBlockHandler())
	dataBlock.WriteBytes(data)
	name := dataBlock.ReadString()
	pid := dataBlock.ReadString()
	uid := dataBlock.ReadString()
	var rsCode int32
	_ = dataBlock.ReadBinary(&rsCode)
	if rsCode == 0 {
		str := dataBlock.ReadString()
		fmt.Println("Response Data suc:", name, pid, uid, rsCode, str)
	} else {
		fmt.Println("Response Data：fail:", name, pid, uid, rsCode)
	}
	return true
}
