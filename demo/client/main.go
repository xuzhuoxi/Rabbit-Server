// Package main
// Created by xuzhuoxi
// on 2019-03-24.
// @author xuzhuoxi
//
package main

import (
	"fmt"
	"github.com/xuzhuoxi/Rabbit-Server/demo/client/net"
	"github.com/xuzhuoxi/Rabbit-Server/demo/client/proto/login"
	"github.com/xuzhuoxi/infra-go/bytex"
	"github.com/xuzhuoxi/infra-go/netx"
	"time"
)

var (
	sleep = make(chan struct{})
)

func main() {
	uc, err := openClient()
	if nil != err {
		panic(err)
	}
	doLogin(uc)
	<-sleep
}

func openClient() (client *net.UserClient, err error) {
	userId := "uid_01"
	uc := net.NewUserClient(userId)
	err = uc.OpenWitAddr(net.RemoteAddress1)
	if nil != err {
		return nil, err
	}
	return uc, nil
}

func doLogin(uc *net.UserClient) {
	uc.SockClient.GetPackHandlerContainer().SetPackHandlers([]netx.FuncPackHandler{onPack})
	//uc.SockClient.GetPackHandlerContainer().AppendPackHandler(packHandler)
	go uc.SockClient.StartReceiving()
	go func() {
		for {
			login.TestLoginExtension(uc)
			time.Sleep(time.Second * 2)
		}
	}()
}

func onPack(data []byte, senderAddress string, other interface{}) (catch bool) {
	dataBlock := bytex.NewBuffDataBlock(bytex.NewDefaultDataBlockHandler())
	dataBlock.WriteBytes(data)
	name := dataBlock.ReadString()
	pid := dataBlock.ReadString()
	uid := dataBlock.ReadString()
	ok := dataBlock.ReadString()
	str := dataBlock.ReadString()
	fmt.Println("Response Data：", name, pid, uid, ok, str)
	return true
}
