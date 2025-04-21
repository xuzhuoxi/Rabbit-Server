// Package core
// Create on 2023/6/14
// @author xuzhuoxi
package core

import (
	"fmt"
	"github.com/xuzhuoxi/Rabbit-Server/engine/server"
	"github.com/xuzhuoxi/Rabbit-Server/engine/server/extension"
	"github.com/xuzhuoxi/Rabbit-Server/engine/server/status"
	"github.com/xuzhuoxi/infra-go/netx"
	"github.com/xuzhuoxi/infra-go/timex"
	"time"
)

// Manager

func NewCustomRabbitManager(statusDetail *status.ServerStatusDetail) server.IRabbitExtensionManager {
	rs := &CustomRabbitManager{
		RabbitExtensionManager: *extension.NewRabbitExtensionManager(),
		StatusDetail:           statusDetail,
	}
	return rs
}

type CustomRabbitManager struct {
	extension.RabbitExtensionManager
	StatusDetail *status.ServerStatusDetail
}

func (m *CustomRabbitManager) StartManager() {
	m.MgrMutex.Lock()
	defer m.MgrMutex.Unlock()
	m.ExtensionContainer.InitExtensions()
	m.HandlerContainer.AppendPackHandler(m.onRabbitGamePack)
}

func (m *CustomRabbitManager) StopManager() {
	m.MgrMutex.Lock()
	defer m.MgrMutex.Unlock()
	m.HandlerContainer.ClearHandler(m.onRabbitGamePack)
	m.ExtensionContainer.DestroyExtensions()
}

func (m *CustomRabbitManager) onRabbitGamePack(msgData []byte, connInfo netx.IConnInfo, other interface{}) bool {
	//m.Logger.Debugln("[CustomRabbitManager.onRabbitGamePack]", connInfo.GetRemoteAddress(), msgData)
	funcName := "[RabbitExtensionManager.onRabbitGamePack]"
	packet, err := m.DecryptPacket(msgData)
	if nil != err {
		return false
	}
	//m.Logger.Debugln("[CustomRabbitManager.onRabbitGamePack]", connInfo.GetRemoteAddress(), packet)
	m.StatusDetail.AddReqCount()
	name, pid, uid, data := m.ParseMessage(packet)
	ext, rsCode := m.Verify(name, pid, uid)
	if server.CodeSuc != rsCode {
		resp := extension.DefaultResponsePool.GetInstance()
		defer extension.DefaultResponsePool.Recycle(resp)
		resp.SetHeader(name, pid, uid, connInfo.GetRemoteAddress())
		resp.(server.IExtensionResponseSettings).SetSockSender(m.SockSender)
		resp.(server.IExtensionResponseSettings).SetConnInfo(connInfo)
		resp.SetResultCode(rsCode)
		resp.ResponseNone()
		//m.Logger.Warnln("[RabbitExtensionManager.onRabbitGamePack]",
		//	fmt.Sprintf("Extension Settlement: Name=%s, PId=%s, FailCode=%d", name, pid, rsCode)) // 记录失败日志
		return false
	}
	// 参数处理
	response, request := m.GetRecycleParams(ext, connInfo, name, pid, uid, data)
	defer func() {
		extension.DefaultRequestPool.Recycle(request)
		extension.DefaultResponsePool.Recycle(response)
	}()
	// 响应处理
	if be, ok := ext.(server.IBeforeRequestExtension); ok { //前置处理
		be.BeforeRequest(response, request)
	}
	if re, ok := ext.(server.IOnRequestExtension); ok {
		func() { //记录时间状态
			tn := time.Now().UnixNano()
			defer func() {
				un := time.Now().UnixNano() - tn
				m.Logger.Infoln(funcName, fmt.Sprintf("Extension Settlement: Name=%s, PId=%s, UsedTime=%s",
					name, pid, timex.FormatMillisLocal(un/1e6, "5.999999ms"))) // 记录响应时间
				m.StatusDetail.AddRespUnixNano(un)
			}()
			re.OnRequest(response, request)
		}()
	}
	if ae, ok := ext.(server.IAfterRequestExtension); ok { //后置处理
		ae.AfterRequest(response, request)
	}
	return true
}
