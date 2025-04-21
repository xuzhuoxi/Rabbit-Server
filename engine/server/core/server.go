package core

import (
	"fmt"
	"github.com/xuzhuoxi/Rabbit-Home/core"
	homeClient "github.com/xuzhuoxi/Rabbit-Home/core/client"
	"github.com/xuzhuoxi/Rabbit-Server/engine/config"
	"github.com/xuzhuoxi/Rabbit-Server/engine/server"
	"github.com/xuzhuoxi/Rabbit-Server/engine/server/extension"
	"github.com/xuzhuoxi/Rabbit-Server/engine/server/status"
	"github.com/xuzhuoxi/infra-go/cryptox"
	"github.com/xuzhuoxi/infra-go/cryptox/asymmetric"
	"github.com/xuzhuoxi/infra-go/cryptox/symmetric"
	"github.com/xuzhuoxi/infra-go/eventx"
	"github.com/xuzhuoxi/infra-go/filex"
	"github.com/xuzhuoxi/infra-go/logx"
	"github.com/xuzhuoxi/infra-go/netx"
	"github.com/xuzhuoxi/infra-go/slicex"
	"net/http"
	"strings"
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

	homeUrl              string
	homeRsaPrivateCipher asymmetric.IRSAPrivateCipher
	homeInternalCipher   cryptox.ICipher
	updating             bool
	rate                 time.Duration
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

func (o *RabbitServer) GetPlatformId() string {
	return o.Config.PlatformId
}

func (o *RabbitServer) GetTypeName() string {
	return o.Config.TypeName
}

func (o *RabbitServer) Init(cfg config.CfgRabbitServerItem) {
	o.Config = cfg
	fmt.Println("Init", cfg)
	o.StatusDetail = status.NewServerStatusDetail(cfg.Id, DefaultStatsInterval)
	o.ExtManager = NewCustomRabbitManager(o.StatusDetail)
	o.homeUrl = fmt.Sprintf("%s://%s", cfg.Home.Network, cfg.Home.NetAddr)

	// 计算更新频率
	o.rate = cfg.Home.RateDuration()
	if o.rate <= 0 {
		o.rate = time.Minute * 2
		o.GetLogger().Warnln("[RabbitServer.Init] Lack Home.rate value", cfg.Home.Rate)
	}

	// 设置SockServer信息
	s, err := netx.ParseSockNetwork(o.Config.Client.Network).NewServer()
	if nil != err {
		panic(err.Error() + o.Config.Client.Network)
	}
	o.SockServer = s.(netx.ISockEventServer)
	o.SockServer.SetName(o.Config.Client.NetName)
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
	list := o.getExtensionList()
	o.registerExtensionList(list)
}

func (o *RabbitServer) getExtensionList() []string {
	cfgExtension := o.Config.Client.Extension
	list := server.GetAllExtensions()
	if !cfgExtension.Custom {
		return list
	}
	if len(cfgExtension.Blocks) > 0 {
		for _, block := range cfgExtension.Blocks {
			if strings.ToLower(block) == "all" {
				return nil
			}
			newList, suc := slicex.RemoveValueString(list, block)
			if suc {
				list = newList
			}
		}
		return list
	}
	if len(cfgExtension.Allows) > 0 {
		var rs []string
		for _, allow := range cfgExtension.Allows {
			if strings.ToLower(allow) == "all" {
				return list
			}
			if slicex.ContainsString(list, allow) {
				rs = append(rs, allow)
			}
		}
		return rs
	}
	return nil
}

func (o *RabbitServer) registerExtensionList(list []string) {
	if len(list) == 0 {
		return
	}
	for _, extName := range list {
		es, err := server.NewRabbitExtension(extName)
		if err != nil {
			o.GetLogger().Errorln(err)
			continue
		}
		o.ExtContainer.AppendExtension(es)
	}
}

func (o *RabbitServer) Start() {
	o.loadPrivateKey()
	o.StatusDetail.Start()
	o.SockServer.AddEventListener(netx.ServerEventStart, o.onSockServerStart)
	o.SockServer.AddEventListener(netx.ServerEventStop, o.onSockServerStop)
	o.SockServer.AddEventListener(netx.ServerEventConnOpened, o.onConnOpened)
	o.SockServer.AddEventListener(netx.ServerEventConnClosed, o.onConnClosed)
	o.ExtManager.StartManager()
	_ = o.SockServer.StartServer(netx.SockParams{
		Network:      netx.ParseSockNetwork(o.Config.Client.Network),
		LocalAddress: o.Config.Client.NetAddr}) //这里会阻塞
}

func (o *RabbitServer) Stop() {
	_ = o.SockServer.StopServer()
	o.ExtManager.StopManager()
	o.SockServer.RemoveEventListener(netx.ServerEventConnOpened, o.onConnOpened)
	o.SockServer.RemoveEventListener(netx.ServerEventConnClosed, o.onConnClosed)
	o.SockServer.RemoveEventListener(netx.ServerEventStop, o.onSockServerStop)
	o.SockServer.RemoveEventListener(netx.ServerEventStart, o.onSockServerStart)
	o.StatusDetail.ReStats()
	o.homeRsaPrivateCipher = nil
}

func (o *RabbitServer) Restart() {
	o.Stop()
	o.Save()
	o.Start()
}

func (o *RabbitServer) Save() {
	o.GetLogger().Infoln("[RabbitServer.Save]", "()")
}

func (o *RabbitServer) loadPrivateKey() {
	home := o.Config.Home
	if !home.Encrypt {
		return
	}
	if len(home.KeyPath) == 0 {
		o.GetLogger().Warnln("[RabbitServer.loadPrivateKey] lack key-path value.")
		return
	}
	rsa, err := asymmetric.LoadPrivateCipherPEM(filex.FixFilePath(home.KeyPath))
	if nil != err {
		o.GetLogger().Errorln("[RabbitServer.loadPrivateKey]", err)
		return
	}
	o.homeRsaPrivateCipher = rsa
}

func (o *RabbitServer) onSockServerStart(evd *eventx.EventData) {
	evd.StopImmediatePropagation()
	o.GetLogger().Infoln("[RabbitServer.onSockServerStart]", "SockServer Start...")
	if !o.Config.Home.Enable {
		return
	}
	o.doLink()
	o.DispatchEvent(evd.EventType, o, evd.Data)
}

func (o *RabbitServer) onSockServerStop(evd *eventx.EventData) {
	evd.StopImmediatePropagation()
	o.GetLogger().Infoln("[RabbitServer.onSockServerStop]", "SockServer Stop...")
	if !o.Config.Home.Enable {
		return
	}
	o.doUnlink()
	o.DispatchEvent(evd.EventType, o, evd.Data)
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
	linkInfo := o.getLinkInfo()
	var err error
	if o.Config.Home.Post {
		o.GetLogger().Infoln("[RabbitServer.doLink] post:", o.homeUrl)
		err = homeClient.LinkWithPost(o.homeUrl, linkInfo, o.StatusDetail.StatsWeight(), o.onLinkResp)
	} else {
		o.GetLogger().Infoln("[RabbitServer.doLink] get:", o.homeUrl)
		err = homeClient.LinkWithGet(o.homeUrl, linkInfo, o.StatusDetail.StatsWeight(), o.onLinkResp)
	}
	if nil != err {
		o.GetLogger().Warnln("[RabbitServer.doLink error]", err)
	}
}

func (o *RabbitServer) onLinkResp(res *http.Response, body *[]byte) {
	suc, fail, err := homeClient.ParseLinkBackInfo(res, body, o.homeRsaPrivateCipher)
	if nil != err || nil != fail {
		if nil != err {
			o.GetLogger().Warnln("[RabbitServer.onLinkResp error]", err)
		} else {
			o.GetLogger().Warnln("[RabbitServer.onLinkResp fail]", fail)
		}
		// 重连
		time.Sleep(o.rate)
		o.doLink()
		return
	}
	if o.Config.Home.Encrypt {
		if len(suc.InternalSK) == 32 {
			o.homeInternalCipher = symmetric.NewAESCipher(suc.InternalSK)
			o.GetLogger().Infoln("[RabbitServer.onLinkResp home.internalSK]:", suc.InternalSK)
		} else {
			o.GetLogger().Warnln("[RabbitServer.onLinkResp home.internalSK] size error:", suc.InternalSK)
		}
	}
	if o.Config.Client.Encrypt {
		if len(suc.OpenSK) == 32 {
			openCipher := symmetric.NewAESCipher(suc.OpenSK)
			o.ExtManager.SetPacketCipher(openCipher)
			o.GetLogger().Infoln("[RabbitServer.onLinkResp client.openSK]:", suc.OpenSK)
		} else {
			o.GetLogger().Warnln("[RabbitServer.onLinkResp client.openSK] size error:", suc.OpenSK)
		}
	}
	go o.rateUpdate()
}

func (o *RabbitServer) rateUpdate() {
	if o.updating {
		return
	}
	o.updating = true
	for o.updating && o.SockServer.IsRunning() {
		time.Sleep(o.rate)
		o.doUpdate()
	}
}

func (o *RabbitServer) doUpdate() {
	updateInfo := o.getUpdateInfo()
	var err error
	if o.Config.Home.Post {
		o.GetLogger().Infoln("[RabbitServer.doUpdate] post:", o.homeUrl)
		err = homeClient.UpdateWithPost(o.homeUrl, updateInfo, o.homeInternalCipher, o.onUpdateResp)
	} else {
		o.GetLogger().Infoln("[RabbitServer.doUpdate] get:", o.homeUrl)
		err = homeClient.UpdateWithGet(o.homeUrl, updateInfo, o.homeInternalCipher, o.onUpdateResp)
	}
	if nil != err {
		o.GetLogger().Warnln("[RabbitServer.doUpdate] error]", err)
	}
}

func (o *RabbitServer) onUpdateResp(resp *http.Response, body *[]byte) {
	// 未注册, 重连
	if resp.StatusCode == http.StatusNotFound {
		o.updating = false
		time.Sleep(o.rate)
		o.doLink()
		return
	}
	fail, err := homeClient.ParseUpdateBackInfo(resp, body)
	if nil != err || nil != fail {
		o.GetLogger().Warnln("[RabbitServer.onUpdateResp]", err, fail)
		return
	}
}

func (o *RabbitServer) doUnlink() {
	unlinkInfo := o.getUnlinkInfo()
	var err error
	if o.Config.Home.Post {
		o.GetLogger().Infoln("[RabbitServer.doUnlink] post:", o.homeUrl)
		err = homeClient.UnlinkWithPost(o.homeUrl, unlinkInfo, o.onUnlinkResp)
	} else {
		o.GetLogger().Infoln("[RabbitServer.doUnlink] get:", o.homeUrl)
		err = homeClient.UnlinkWithGet(o.homeUrl, unlinkInfo, o.onUnlinkResp)
	}
	if nil != err {
		o.GetLogger().Warnln("[RabbitServer.doUnlink] error:", err)
	}
}

func (o *RabbitServer) onUnlinkResp(resp *http.Response, body *[]byte) {
	suc, fail, err := homeClient.ParseUnlinkBackInfo(resp, body)
	if nil != err || nil != fail {
		o.GetLogger().Warnln("[RabbitServer.onUnlinkResp]", err, fail)
		return
	}
	o.updating = false
	o.GetLogger().Infoln("[RabbitServer.onUnlinkResp]", suc)
}

func (o *RabbitServer) getLinkInfo() core.LinkInfo {
	info := core.LinkInfo{
		Id:          o.Config.Id,
		PlatformId:  o.Config.PlatformId,
		TypeName:    o.Config.TypeName,
		OpenNetwork: o.Config.Client.Network,
		OpenAddr:    o.Config.Client.NetAddr,
		OpenKeyOn:   o.Config.Client.Encrypt,
	}
	if o.Config.Home.Encrypt {
		original := info.OriginalSignData()
		signature, err := o.homeRsaPrivateCipher.SignBase64(original, server.Base64Encoding)
		if nil != err {
			o.GetLogger().Errorln("[RabbitServer.getLinkInfo] SignError:", err)
		} else {
			info.Signature = signature
		}
	}
	return info
}

func (o *RabbitServer) getUnlinkInfo() core.UnlinkInfo {
	info := core.UnlinkInfo{
		Id: o.Config.Id,
	}
	if o.Config.Home.Encrypt {
		original := info.OriginalSignData()
		signature, err := o.homeRsaPrivateCipher.SignBase64(original, server.Base64Encoding)
		if nil != err {
			o.GetLogger().Errorln("[RabbitServer.getUnlinkInfo] SignError:", err)
		} else {
			info.Signature = signature
		}
	}
	return info
}
func (o *RabbitServer) getUpdateInfo() core.UpdateInfo {
	return core.UpdateInfo{
		Id:     o.Config.Id,
		Weight: o.StatusDetail.StatsWeight(),
	}
}
func (o *RabbitServer) getUpdateDetailInfo() core.UpdateDetailInfo {
	return core.UpdateDetailInfo{
		Id:    o.Config.Id,
		Links: 100,
	}
}
