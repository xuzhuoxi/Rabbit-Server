package core

import (
	"fmt"
	"github.com/xuzhuoxi/Rabbit-Home/core"
	homeClient "github.com/xuzhuoxi/Rabbit-Home/core/client"
	"github.com/xuzhuoxi/Rabbit-Server/engine/config"
	"github.com/xuzhuoxi/Rabbit-Server/engine/server"
	"github.com/xuzhuoxi/Rabbit-Server/engine/server/extension"
	"github.com/xuzhuoxi/Rabbit-Server/engine/server/status"
	"github.com/xuzhuoxi/infra-go/eventx"
	"github.com/xuzhuoxi/infra-go/logx"
	"github.com/xuzhuoxi/infra-go/netx"
	"net/http"
	"time"
)

func NewIRabbitServer() server.IRabbitServer {
	return NewRabbitServer()
}

func NewRabbitServer() *RabbitServer {
	container := extension.NewIRabbitExtensionContainer()
	rs := &RabbitServer{
		ExtContainer: container,
	}
	return rs
}

type RabbitServer struct {
	eventx.EventDispatcher
	logx.LoggerSupport
	Config       config.CfgRabbitServerItem
	SockServer   netx.ISockEventServer
	ExtContainer server.IRabbitExtensionContainer
	ExtManager   server.IRabbitExtensionManager
	StatusDetail *status.ServerStatusDetail
}

func (o *RabbitServer) GetConnSet() (set netx.IServerConnSet, ok bool) {
	return o.SockServer, o.SockServer != nil
}

func (o *RabbitServer) GetExtensionManager() (mgr server.IRabbitExtensionManager, ok bool) {
	return o.ExtManager, o.ExtManager != nil
}

func (o *RabbitServer) GetId() string {
	return o.Config.Id
}

func (o *RabbitServer) GetName() string {
	return o.Config.Name
}

func (o *RabbitServer) Init(cfg config.CfgRabbitServerItem) {
	o.Config = cfg
	o.StatusDetail = status.NewServerStatusDetail(cfg.Id, DefaultStatsInterval)
	o.ExtManager = NewCustomRabbitManager(o.StatusDetail)

	// 设置SockServer信息
	server, err := netx.ParseSockNetwork(o.Config.FromUser.Network).NewServer()
	if nil != err {
		panic(err)
	}
	o.SockServer = server.(netx.ISockEventServer)
	o.SockServer.SetName(o.Config.FromUser.Name)
	o.SockServer.SetMaxConn(100)
	o.SockServer.SetLogger(o.GetLogger())
	// 注入Extension
	o.registerExtensions()
	// 初始化ExtensionManager
	// 这里把Manager、SockServer、Container进行关联
	o.ExtManager.InitManager(o.SockServer.GetPackHandlerContainer(), o.ExtContainer, o.SockServer)
	o.ExtManager.SetLogger(o.GetLogger())
	o.ExtManager.SetUserConnMapper(RabbitUserConnMapper)
}

func (o *RabbitServer) registerExtensions() {
	var list []string
	if o.Config.Extension.All {
		list = server.GetAllExtensions()
	} else {
		list = o.Config.Extension.List
	}
	if len(list) == 0 {
		return
	}
	for _, extName := range list {
		extension, err := server.NewRabbitExtension(extName)
		if err != nil {
			o.GetLogger().Errorln(err)
			continue
		}
		o.ExtContainer.AppendExtension(extension)
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
	o.GetLogger().Infoln("[RabbitServer.Save]", "()")
}

func (o *RabbitServer) onSockServerStart(evd *eventx.EventData) {
	evd.StopImmediatePropagation()
	o.GetLogger().Infoln("[RabbitServer.onSockServerStart]", "SockServer Start...")
	if o.Config.ToHome.Disable {
		return
	}
	o.doLink()
	go o.rateUpdate()
	o.DispatchEvent(evd.EventType, o, evd.Data)
}

func (o *RabbitServer) onSockServerStop(evd *eventx.EventData) {
	evd.StopImmediatePropagation()
	o.GetLogger().Infoln("[RabbitServer.onSockServerStop]", "SockServer Stop...")
	if o.Config.ToHome.Disable {
		return
	}
	o.doUnlink()
	o.DispatchEvent(evd.EventType, o, evd.Data)
}

func (o *RabbitServer) rateUpdate() {
	rate := o.Config.GetToHomeRate()
	o.GetLogger().Debugln("RabbitServer.rateUpdate：", rate)
	for o.SockServer.IsRunning() {
		time.Sleep(rate)
		o.doUpdate()
	}
}

func (o *RabbitServer) onUpdateResp(resp *http.Response, body *[]byte) {
	if resp.StatusCode == http.StatusNotFound {
		// 未注册, 重连
		o.doLink()
	}
}

func (o *RabbitServer) onConnOpened(evd *eventx.EventData) {
	evd.StopImmediatePropagation()
	o.StatusDetail.AddLinkCount()
	connInfo := evd.Data.(netx.IConnInfo)
	o.GetLogger().Infoln("[RabbitServer.onConnOpened]", "Client Connection Open:", connInfo)
	o.DispatchEvent(evd.EventType, o, evd.Data)
}

func (o *RabbitServer) onConnClosed(evd *eventx.EventData) {
	evd.StopImmediatePropagation()
	connInfo := evd.Data.(netx.IConnInfo)
	o.GetLogger().Infoln("[RabbitServer.onConnClosed]", "Client Connection Close:", connInfo)
	o.StatusDetail.RemoveLinkCount()
	o.DispatchEvent(evd.EventType, o, evd.Data)
}

// -----------------------------------

func (o *RabbitServer) doLink() {
	url := "http://" + o.Config.ToHome.Addr
	fmt.Println("RabbitServer.doLink:", url)
	err := homeClient.LinkWithGet(url, o.getLinkInfo(), o.StatusDetail.StatsWeight(), nil)
	if nil != err {
		o.GetLogger().Warnln("[RabbitServer.rateUpdate]", err)
	}
}

func (o *RabbitServer) doUnlink() {
	url := "http://" + o.Config.ToHome.Addr
	//fmt.Println("RabbitServer.doUnlink:", url)
	err := homeClient.UnlinkWithGet(url, o.GetId(), nil)
	if nil != err {
		o.GetLogger().Warnln("[RabbitServer.rateUpdate]", err)
	}
}
func (o *RabbitServer) doUpdate() {
	url := "http://" + o.Config.ToHome.Addr
	//fmt.Println("RabbitServer.doUpdate:", url)
	err := homeClient.UpdateWithGet(url, o.getUpdateStatus(), o.onUpdateResp)
	if nil != err {
		o.GetLogger().Warnln("[RabbitServer.rateUpdate]", err)
	}
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
