// Package rabbit
// Create on 2023/6/14
// @author xuzhuoxi
package rabbit

import (
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
	m.HandlerContainer.AppendPackHandler(m.onSnailGamePack)
}

func (m *RabbitExtensionManager) StopManager() {
	m.Mutex.Lock()
	defer m.Mutex.Unlock()
	m.HandlerContainer.ClearHandler(m.onSnailGamePack)
	m.ExtensionContainer.DestroyExtensions()
}

func (m *RabbitExtensionManager) onSnailGamePack(msgData []byte, senderAddress string, other interface{}) bool {
	//m.Log.Infoln("ExtManager.onPack", senderAddress, msgData)
	m.StatusDetail.AddReqCount()
	name, pid, uid, data := m.ParseMessage(msgData)
	extension, ok := m.Verify(name, pid, uid)
	if !ok {
		return false
	}
	// 参数处理
	response, request := m.GenParams(extension, senderAddress, name, pid, uid, data)
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
				m.Logger.Infoln(name, pid, un, timex.FormatUnixMilli(un/1e6, "5.999999ms")) //记录响应时间
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
