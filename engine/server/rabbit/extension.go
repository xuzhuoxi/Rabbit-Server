// Package rabbit
// Create on 2023/6/14
// @author xuzhuoxi
package rabbit

import (
	"fmt"
	"github.com/xuzhuoxi/Rabbit-Server/engine/server"
	"github.com/xuzhuoxi/infra-go/extendx/protox"
	"github.com/xuzhuoxi/infra-go/timex"
	"time"
)

// Container

func NewRabbitExtensionContainer() server.IRabbitExtensionContainer {
	return protox.NewIProtocolExtensionContainer()
}

// Manager

func NewRabbitExtensionManager(statusDetail *ServerStatusDetail) server.IRabbitExtensionManager {
	rs := &RabbitExtensionManager{
		ExtensionManager: *protox.NewExtensionManager(),
		StatusDetail:     statusDetail,
	}
	return rs
}

type RabbitExtensionManager struct {
	protox.ExtensionManager
	StatusDetail *ServerStatusDetail
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
	m.StatusDetail.AddReqCount()
	name, pid, uid, data := m.ParseMessage(msgData)
	extension, rsCode := m.Verify(name, pid, uid)
	if protox.CodeSuc != rsCode {
		resp := protox.DefaultResponsePool.GetInstance()
		defer protox.DefaultResponsePool.Recycle(resp)
		resp.SetHeader(name, pid, uid, senderAddress)
		resp.(protox.IExtensionResponseSettings).SetSockSender(m.SockSender)
		resp.SetResultCode(rsCode)
		resp.SendNoneResponse()
		m.Logger.Warnln(fmt.Sprintf("Extension Settlement: Name=%s, PId=%s, FailCode=%d",
			name, pid, rsCode)) // 记录失败日志
		return false
	}
	// 参数处理
	response, request := m.GetRecycleParams(extension, senderAddress, name, pid, uid, data)
	defer func() {
		protox.DefaultRequestPool.Recycle(request)
		protox.DefaultResponsePool.Recycle(response)
	}()
	// 响应处理
	if be, ok := extension.(protox.IBeforeRequestExtension); ok { //前置处理
		be.BeforeRequest(request)
	}
	if re, ok := extension.(protox.IRequestExtension); ok {
		func() { //记录时间状态
			tn := time.Now().UnixNano()
			defer func() {
				un := time.Now().UnixNano() - tn
				m.Logger.Infoln(fmt.Sprintf("Extension Settlement: Name=%s, PId=%s, UsedTime=%s",
					name, pid, timex.FormatMillisLocal(un/1e6, "5.999999ms"))) // 记录响应时间
				m.StatusDetail.AddRespUnixNano(un)
			}()
			re.OnRequest(response, request)
		}()
	}
	if ae, ok := extension.(protox.IAfterRequestExtension); ok { //后置处理
		ae.AfterRequest(response, request)
	}
	return true
}
