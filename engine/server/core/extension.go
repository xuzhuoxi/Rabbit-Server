// Package rabbit
// Create on 2023/6/14
// @author xuzhuoxi
package core

import (
	"fmt"
	"github.com/xuzhuoxi/Rabbit-Server/engine/server"
	"github.com/xuzhuoxi/Rabbit-Server/engine/server/extension"
	"github.com/xuzhuoxi/Rabbit-Server/engine/server/status"
	"github.com/xuzhuoxi/infra-go/timex"
	"time"
)

// Container

func NewRabbitExtensionContainer() server.IRabbitExtensionContainer {
	return extension.NewIProtoExtensionContainer()
}

// Manager

func NewRabbitExtensionManager(statusDetail *status.ServerStatusDetail) server.IRabbitExtensionManager {
	rs := &RabbitExtensionManager{
		ExtensionManager: *extension.NewExtensionManager(),
		StatusDetail:     statusDetail,
	}
	return rs
}

type RabbitExtensionManager struct {
	extension.ExtensionManager
	StatusDetail *status.ServerStatusDetail
}

func (m *RabbitExtensionManager) StartManager() {
	m.Mutex.Lock()
	defer m.Mutex.Unlock()
	m.ExtensionContainer.InitExtensions()
	m.HandlerContainer.AppendPackHandler(m.onRabbitGamePack)
}

func (m *RabbitExtensionManager) StopManager() {
	m.Mutex.Lock()
	defer m.Mutex.Unlock()
	m.HandlerContainer.ClearHandler(m.onRabbitGamePack)
	m.ExtensionContainer.DestroyExtensions()
}

func (m *RabbitExtensionManager) onRabbitGamePack(msgData []byte, senderAddress string, other interface{}) bool {
	//m.Logger.Infoln("ExtManager.onPack", senderAddress, msgData)
	funcName := "[RabbitExtensionManager.onRabbitGamePack]"
	m.StatusDetail.AddReqCount()
	name, pid, uid, data := m.ParseMessage(msgData)
	ext, rsCode := m.Verify(name, pid, uid)
	if server.CodeSuc != rsCode {
		resp := extension.DefaultResponsePool.GetInstance()
		defer extension.DefaultResponsePool.Recycle(resp)
		resp.SetHeader(name, pid, uid, senderAddress)
		resp.(server.IExtensionResponseSettings).SetSockSender(m.SockSender)
		resp.SetResultCode(rsCode)
		resp.SendNoneResponse()
		m.Logger.Warnln("[RabbitExtensionManager.onRabbitGamePack]",
			fmt.Sprintf("Extension Settlement: Name=%s, PId=%s, FailCode=%d", name, pid, rsCode)) // 记录失败日志
		return false
	}
	// 参数处理
	response, request := m.GetRecycleParams(ext, senderAddress, name, pid, uid, data)
	defer func() {
		extension.DefaultRequestPool.Recycle(request)
		extension.DefaultResponsePool.Recycle(response)
	}()
	// 响应处理
	if be, ok := ext.(server.IBeforeRequestExtension); ok { //前置处理
		be.BeforeRequest(request)
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
