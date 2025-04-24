// Package runtime
// Created by xuzhuoxi
// on 2019-05-18.
// @author xuzhuoxi
//
package runtime

import (
	"fmt"
	"github.com/xuzhuoxi/Rabbit-Server/engine/mgr"
	"github.com/xuzhuoxi/Rabbit-Server/engine/server"
	"github.com/xuzhuoxi/Rabbit-Server/engine/server/runtime/params"
	"github.com/xuzhuoxi/Rabbit-Server/engine/server/runtime/params/packets"
	"github.com/xuzhuoxi/Rabbit-Server/engine/server/runtime/status"
	"github.com/xuzhuoxi/Rabbit-Server/engine/server/runtime/verify"
	"github.com/xuzhuoxi/infra-go/cryptox"
	"github.com/xuzhuoxi/infra-go/logx"
	"github.com/xuzhuoxi/infra-go/netx"
	"github.com/xuzhuoxi/infra-go/timex"
	"sync"
	"time"
)

type ExtensionManager struct {
	HandlerContainer   netx.IPackHandlerContainer
	ExtensionContainer server.IRabbitExtensionContainer

	SockSender     netx.ISockSender
	UserConnMapper netx.IUserConnMapper

	PacketVerifier server.IPacketVerifier
	RequestPacket  packets.RequestPacket
	ParamsSupport  params.ParamsSupport
	CryptoSupport  CryptoSupport

	Logger   logx.ILogger
	MgrMutex sync.RWMutex
}

func (o *ExtensionManager) InitManager(handlerContainer netx.IPackHandlerContainer, extensionContainer server.IRabbitExtensionContainer,
	sockSender netx.ISockSender) {
	o.MgrMutex.Lock()
	defer o.MgrMutex.Unlock()

	o.HandlerContainer, o.ExtensionContainer, o.SockSender = handlerContainer, extensionContainer, sockSender
	o.ParamsSupport.SockSender = sockSender

	o.PacketVerifier = verify.NewIPacketVerifier()
	o.PacketVerifier.AppendVerifyHandler(verify.NewExtensionVerifyItem(extensionContainer).Verify)
	_, _, _, _, _, cfg := mgr.DefaultManager.GetInitManager().GetConfigs()
	o.PacketVerifier.AppendVerifyHandler(verify.NewFreqVerifyItem(cfg).Verify)
}

func (o *ExtensionManager) SetUserConnMapper(mapper netx.IUserConnMapper) {
	o.MgrMutex.Lock()
	defer o.MgrMutex.Unlock()
	o.UserConnMapper = mapper
	o.ParamsSupport.UserConnMapper = mapper
}
func (o *ExtensionManager) SetPacketCipher(cipher cryptox.ICipher) {
	o.CryptoSupport.SetPacketCipher(cipher)
}

func (o *ExtensionManager) AppendVerifyHandler(handler server.FuncVerifyPacket) {
	o.MgrMutex.Lock()
	defer o.MgrMutex.Unlock()
	o.PacketVerifier.AppendVerifyHandler(handler)
}

func (o *ExtensionManager) SetLogger(logger logx.ILogger) {
	o.MgrMutex.Lock()
	defer o.MgrMutex.Unlock()
	o.Logger = logger
}

func (o *ExtensionManager) StartManager() {
	o.MgrMutex.Lock()
	defer o.MgrMutex.Unlock()
	o.ExtensionContainer.InitExtensions()
	// 关联响应处理函数
	_ = o.HandlerContainer.AppendPackHandler(o.OnMessageUnpack)
}

func (o *ExtensionManager) StopManager() {
	o.MgrMutex.Lock()
	defer o.MgrMutex.Unlock()
	_ = o.HandlerContainer.ClearHandler(o.OnMessageUnpack)
	_ = o.ExtensionContainer.DestroyExtensions()
}

func (o *ExtensionManager) SaveExtension(extName string) error {
	o.MgrMutex.Lock()
	defer o.MgrMutex.Unlock()
	return o.ExtensionContainer.SaveExtension(extName)
}

func (o *ExtensionManager) SaveExtensions() []error {
	o.MgrMutex.Lock()
	defer o.MgrMutex.Unlock()
	return o.ExtensionContainer.SaveExtensions()
}

func (o *ExtensionManager) EnableExtension(extName string, enable bool) error {
	o.MgrMutex.Lock()
	defer o.MgrMutex.Unlock()
	return o.ExtensionContainer.EnableExtension(extName, enable)
}

func (o *ExtensionManager) EnableExtensions(enable bool) []error {
	o.MgrMutex.Lock()
	defer o.MgrMutex.Unlock()
	return o.ExtensionContainer.EnableExtensions(enable)
}

// OnMessageUnpack
// 消息处理入口，这里是并发方法
// 并发注意:本方法是否并发，取决于SockServer的并发性
// 在netx中，TCP,Quic,WebSocket为并发响应，UDP为非并发响应
// msgData非共享的，但在ParsePacket后这部分数据会发生变化
func (o *ExtensionManager) OnMessageUnpack(msgData []byte, connInfo netx.IConnInfo, _ interface{}) bool {
	packetData, err := o.CryptoSupport.DecryptPacket(msgData)
	if nil != err {
		return false
	}
	name, pid, uid, data := o.RequestPacket.ParsePacket(packetData)
	rsCode := o.PacketVerifier.Verify(name, pid, uid)
	if server.CodeSuc != rsCode {
		// 这里可以直接响应失败
		return false
	}
	extension := o.ExtensionContainer.GetExtension(name).(server.IRabbitExtension)
	//参数处理
	response, request := o.ParamsSupport.GetRecycleParams(extension, connInfo, name, pid, uid, data)
	defer func() {
		params.DefaultRequestPool.Recycle(request)
		params.DefaultResponsePool.Recycle(response)
	}()
	// 响应处理
	if be, ok := extension.(server.IBeforeRequestExtension); ok { //前置处理
		be.BeforeRequest(response, request)
	}
	if re, ok := extension.(server.IOnRequestExtension); ok {
		re.OnRequest(response, request)
	}
	if ae, ok := extension.(server.IAfterRequestExtension); ok { //后置处理
		ae.AfterRequest(response, request)
	}
	return true
}

// RabbitExtensionManager ---------- ---------- ---------- ---------- ----------

func NewRabbitExtensionManager(statusDetail *status.ServerStatusDetail) server.IRabbitExtensionManager {
	rs := &RabbitExtensionManager{
		StatusDetail: statusDetail,
	}
	return rs
}

type RabbitExtensionManager struct {
	ExtensionManager
	StatusDetail *status.ServerStatusDetail
}

func (o *RabbitExtensionManager) StartManager() {
	o.MgrMutex.Lock()
	defer o.MgrMutex.Unlock()
	o.ExtensionContainer.InitExtensions()
	_ = o.HandlerContainer.AppendPackHandler(o.OnMessageUnpack)
}

func (o *RabbitExtensionManager) StopManager() {
	o.MgrMutex.Lock()
	defer o.MgrMutex.Unlock()
	_ = o.HandlerContainer.ClearHandler(o.OnMessageUnpack)
	o.ExtensionContainer.DestroyExtensions()
}

func (o *RabbitExtensionManager) OnMessageUnpack(msgData []byte, connInfo netx.IConnInfo, _ interface{}) bool {
	//o.Logger.Debugln("[RabbitExtensionManager.OnMessageUnpack]", connInfo.GetRemoteAddress(), msgData)

	// 解密数据
	packetData, err := o.CryptoSupport.DecryptPacket(msgData)
	if nil != err {
		return false
	}

	// 解析数据，提取响应扩展名与协议Id，并进行存在性检验
	name, pid, uid, data := o.RequestPacket.ParsePacket(packetData)
	// 检验响应合法性
	rsCode := o.PacketVerifier.Verify(name, pid, uid)
	if server.CodeSuc != rsCode {
		resp := params.DefaultResponsePool.GetInstance()
		defer params.DefaultResponsePool.Recycle(resp)
		resp.SetHeader(name, pid, uid, connInfo.GetRemoteAddress())
		resp.(server.IExtensionResponseSettings).SetSockSender(o.SockSender)
		resp.(server.IExtensionResponseSettings).SetConnInfo(connInfo)
		resp.SetResultCode(rsCode)
		_ = resp.ResponseNone()
		o.Logger.Warnln("[RabbitExtensionManager.onRabbitGamePack]",
			fmt.Sprintf("Extension Settlement: Name=%s, PId=%s, FailCode=%d", name, pid, rsCode)) // 记录失败日志
		return false
	}
	// 统计请求次数
	o.StatusDetail.AddReqCount()
	ext := o.ExtensionContainer.GetExtension(name).(server.IRabbitExtension)
	// 响应参数生成
	response, request := o.ParamsSupport.GetRecycleParams(ext, connInfo, name, pid, uid, data)
	defer func() {
		params.DefaultRequestPool.Recycle(request)
		params.DefaultResponsePool.Recycle(response)
	}()

	// 响应处理 ---------- ---------- ---------- ---------- ----------

	// 前置处理逻辑
	if be, ok := ext.(server.IBeforeRequestExtension); ok {
		be.BeforeRequest(response, request)
	}
	// 请求处理逻辑
	if re, ok := ext.(server.IOnRequestExtension); ok {
		func() { //记录时间状态
			tn := time.Now().UnixNano()
			defer func() {
				un := time.Now().UnixNano() - tn
				o.Logger.Infoln("[RabbitExtensionManager.OnMessageUnpack]", fmt.Sprintf("Extension Settlement: Name=%s, PId=%s, UsedTime=%s",
					name, pid, timex.FormatMillisLocal(un/1e6, "5.999999ms"))) // 记录响应时间
				o.StatusDetail.AddRespUnixNano(un)
			}()
			re.OnRequest(response, request)
		}()
	}
	// 后置处理逻辑
	if ae, ok := ext.(server.IAfterRequestExtension); ok {
		ae.AfterRequest(response, request)
	}
	return true

}
