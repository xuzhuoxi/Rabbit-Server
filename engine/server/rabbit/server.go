package rabbit

import (
	"fmt"
	"github.com/xuzhuoxi/Rabbit-Home/core"
	homeClient "github.com/xuzhuoxi/Rabbit-Home/core/client"
	"github.com/xuzhuoxi/Rabbit-Home/core/home"
	"github.com/xuzhuoxi/Rabbit-Server/engine/server"
	"github.com/xuzhuoxi/infra-go/eventx"
	"github.com/xuzhuoxi/infra-go/logx"
	"github.com/xuzhuoxi/infra-go/netx"
	"github.com/xuzhuoxi/infra-go/netx/tcpx"
	"net/http"
	"time"
)

func NewRabbitServer(cfg server.CfgRabbitServer) server.IRabbitServer {
	container := NewServerExtensionContainer()
	server := tcpx.NewTCPServer()
	statusDetail := NewServerStatusDetail(cfg.Name, DefaultStatsInterval)
	mgr := NewServerExtensionManager(statusDetail)
	logger := logx.NewLogger()
	return &RabbitServer{
		Config:       cfg,
		SockServer:   server,
		ExtContainer: container,
		ExtManager:   mgr,
		StatusDetail: statusDetail,
		Logger:       logger,
	}
}

type RabbitServer struct {
	eventx.EventDispatcher
	Config       server.CfgRabbitServer
	SockServer   tcpx.ITCPServer
	ExtContainer IServerExtensionContainer
	ExtManager   IServerExtensionManager
	StatusDetail *ServerStatusDetail
	Logger       logx.ILogger
}

func (o *RabbitServer) GetId() string {
	return o.Config.Id
}

func (o *RabbitServer) GetName() string {
	return o.Config.Name
}

func (o *RabbitServer) GetLogger() logx.ILogger {
	return o.Logger
}

func (o *RabbitServer) Init() {
	// 注入Extension
	ForeachExtensionConstructor(func(constructor FuncServerExtension) {
		o.ExtContainer.AppendExtension(constructor())
	})
	// 设置SockServer信息
	o.SockServer.SetName(o.Config.FromUser.Name)
	o.SockServer.SetMax(100)
	o.SockServer.SetLogger(o.Logger)
	// 初始化ExtensionManager
	o.ExtManager.InitManager(o.SockServer.GetPackHandlerContainer(), o.ExtContainer, o.SockServer)
	o.ExtManager.SetLogger(o.Logger)
	o.ExtManager.SetAddressProxy(AddressProxy)
	// 初始化Logger
	cfgLog := o.Config.Log
	if nil != cfgLog {
		o.Logger.SetConfig(cfgLog.ToLogConfig())
	}
}

func (o *RabbitServer) Start() {
	o.StatusDetail.Start()
	o.SockServer.AddEventListener(netx.ServerEventStart, o.onSockServerStart)
	o.SockServer.AddEventListener(netx.ServerEventStop, o.onSockServerStop)
	o.SockServer.AddEventListener(netx.ServerEventConnOpened, o.onConnOpened)
	o.SockServer.AddEventListener(netx.ServerEventConnClosed, o.onConnClosed)
	o.ExtManager.StartManager()
	o.SockServer.StartServer(netx.SockParams{
		Network: netx.ParseSockNetwork(o.Config.FromUser.Network), LocalAddress: o.Config.FromUser.Addr}) //这里会阻塞
}

func (o *RabbitServer) Stop() {
	o.SockServer.StopServer()
	o.ExtManager.StopManager()
	o.SockServer.RemoveEventListener(netx.ServerEventConnOpened, o.onConnOpened)
	o.SockServer.RemoveEventListener(netx.ServerEventConnClosed, o.onConnClosed)
	o.SockServer.RemoveEventListener(netx.ServerEventStop, o.onSockServerStop)
	o.SockServer.RemoveEventListener(netx.ServerEventStart, o.onSockServerStart)
	o.StatusDetail.ReStats()
}

func (o *RabbitServer) Restart() {
	o.Stop()
	o.Save()
	o.Start()
}

func (o *RabbitServer) Save() {
	o.Logger.Infoln("Save!")
}

func (o *RabbitServer) onSockServerStart(evd *eventx.EventData) {
	o.doLink()
	go o.rateUpdate()
}

func (o *RabbitServer) onSockServerStop(evd *eventx.EventData) {
	o.doUnlink()
}

func (o *RabbitServer) rateUpdate() {
	url := fmt.Sprintf("http://%s/%s", o.Config.ToHome.Addr, home.PatternUpdate)
	rate := o.Config.GetToHomeRate()
	for o.SockServer.IsRunning() {
		time.Sleep(rate)
		err := homeClient.UpdateWithGet(url, o.getUpdateStatus(), o.onUpdateResp)
		if nil != err {
			o.Logger.Warnln(err)
		}
	}
}

func (o *RabbitServer) onUpdateResp(resp *http.Response, body *[]byte) {
	if resp.StatusCode == http.StatusNotFound {
		// 未注册, 重连
		o.doLink()
	}
}

func (o *RabbitServer) onConnOpened(evd *eventx.EventData) {
	o.StatusDetail.AddLinkCount()
}

func (o *RabbitServer) onConnClosed(evd *eventx.EventData) {
	address := evd.Data.(string)
	AddressProxy.RemoveByAddress(address)
	o.StatusDetail.RemoveLinkCount()
}

// -----------------------------------

func (o *RabbitServer) doLink() {
	url := fmt.Sprintf("http://%s/%s", o.Config.ToHome.Addr, home.PatternLink)
	homeClient.LinkWithGet(url, o.getLinkInfo(), o.StatusDetail.StatsWeight())
}

func (o *RabbitServer) doUnlink() {
	url := fmt.Sprintf("http://%s/%s", o.Config.ToHome.Addr, home.PatternUnlink)
	homeClient.UnlinkWithGet(url, o.GetId())
}

func (o *RabbitServer) getLinkInfo() core.LinkEntity {
	return core.LinkEntity{
		Id:         o.Config.FromUser.Name,
		PlatformId: o.Config.Name,
		Name:       o.Config.Name,
		Network:    o.Config.FromUser.Network,
		Addr:       o.Config.FromUser.Addr,
	}
}
func (o *RabbitServer) getUpdateStatus() core.EntityStatus {
	return core.EntityStatus{
		Id:     o.Config.FromUser.Name,
		Weight: o.StatusDetail.StatsWeight(),
	}
}
