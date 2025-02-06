// Package extension
// Created by xuzhuoxi
// on 2019-05-18.
// @author xuzhuoxi
//
package extension

import (
	"fmt"
	"github.com/xuzhuoxi/Rabbit-Server/engine/server"
	"github.com/xuzhuoxi/infra-go/bytex"
	"github.com/xuzhuoxi/infra-go/logx"
	"github.com/xuzhuoxi/infra-go/netx"
	"sync"
)

func NewIRabbitExtensionManager() server.IRabbitExtensionManager {
	return NewRabbitExtensionManager()
}

func NewRabbitExtensionManager() *RabbitExtensionManager {
	m := &RabbitExtensionManager{}
	m.SetCustomParsePacketFunc(m.ParseMessage)
	return m
}

type RabbitExtensionManager struct {
	HandlerContainer   netx.IPackHandlerContainer
	ExtensionContainer server.IRabbitExtensionContainer
	SockSender         netx.ISockSender

	Logger       logx.ILogger
	AddressProxy netx.IAddressProxy
	Mutex        sync.RWMutex

	CustomManagerSupport
}

func (m *RabbitExtensionManager) InitManager(handlerContainer netx.IPackHandlerContainer, extensionContainer server.IRabbitExtensionContainer,
	sockSender netx.ISockSender) {
	m.Mutex.Lock()
	defer m.Mutex.Unlock()
	m.HandlerContainer, m.ExtensionContainer, m.SockSender = handlerContainer, extensionContainer, sockSender
}

func (m *RabbitExtensionManager) SetAddressProxy(proxy netx.IAddressProxy) {
	m.Mutex.Lock()
	defer m.Mutex.Unlock()
	m.AddressProxy = proxy
}

func (m *RabbitExtensionManager) SetLogger(logger logx.ILogger) {
	m.Mutex.Lock()
	defer m.Mutex.Unlock()
	m.Logger = logger
}

func (m *RabbitExtensionManager) StartManager() {
	m.Mutex.Lock()
	defer m.Mutex.Unlock()
	m.ExtensionContainer.InitExtensions()
	var _ = m.HandlerContainer.AppendPackHandler(m.OnMessageUnpack)
}

func (m *RabbitExtensionManager) StopManager() {
	m.Mutex.Lock()
	defer m.Mutex.Unlock()
	var _ = m.HandlerContainer.ClearHandler(m.OnMessageUnpack)
	var _ = m.ExtensionContainer.DestroyExtensions()
}

func (m *RabbitExtensionManager) SaveExtensions() {
	m.Mutex.Lock()
	defer m.Mutex.Unlock()
	var _ = m.ExtensionContainer.SaveExtensions()
}

func (m *RabbitExtensionManager) SaveExtension(name string) {
	m.Mutex.Lock()
	defer m.Mutex.Unlock()
	var _ = m.ExtensionContainer.SaveExtension(name)
}

func (m *RabbitExtensionManager) EnableExtension(name string) {
	m.Mutex.Lock()
	defer m.Mutex.Unlock()
	var _ = m.ExtensionContainer.EnableExtension(name, true)
}

func (m *RabbitExtensionManager) DisableExtension(name string) {
	m.Mutex.Lock()
	defer m.Mutex.Unlock()
	var _ = m.ExtensionContainer.EnableExtension(name, false)
}

func (m *RabbitExtensionManager) EnableExtensions() {
	m.Mutex.Lock()
	defer m.Mutex.Unlock()
	var _ = m.ExtensionContainer.EnableExtensions(true)
}

func (m *RabbitExtensionManager) DisableExtensions() {
	m.Mutex.Lock()
	defer m.Mutex.Unlock()
	var _ = m.ExtensionContainer.EnableExtensions(false)
}

//---------------------------------

type iExtResponse interface {
	server.IExtensionResponseSettings
	server.IExtensionResponse
}

// OnMessageUnpack
// 消息处理入口，这里是并发方法
// 并发注意:本方法是否并发，取决于SockServer的并发性
// 在netx中，TCP,Quic,WebSocket为并发响应，UDP为非并发响应
// msgData非共享的，但在parsePackMessage后这部分数据会发生变化
func (m *RabbitExtensionManager) OnMessageUnpack(msgData []byte, senderAddress string, other interface{}) bool {
	//m.Logger.Infoln("ExtensionManager.onPack", senderAddress, msgData)
	m.CustomStartOnPack(senderAddress)
	name, pid, uid, data := m.CustomParsePacket(msgData) // 默认调用m.ParseMessage
	extension, rsCode := m.Verify(name, pid, uid)
	if server.CodeSuc != rsCode {
		// 这里可以直接响应失败
		return false
	}
	//参数处理
	response, request := m.GetRecycleParams(extension, senderAddress, name, pid, uid, data)
	defer func() {
		DefaultRequestPool.Recycle(request)
		DefaultResponsePool.Recycle(response)
	}()
	//响应处理
	m.DoRequest(extension, request, response)
	return true
}

// ParseMessage
// block0 : eName utf8
// block1 : pid	utf8
// block2 : uid	utf8
// block3 : data [][]byte
// [n]其它信息
func (m *RabbitExtensionManager) ParseMessage(msgBytes []byte) (extName string, pid string, uid string, data [][]byte) {
	if nil != m.FuncParsePacket {
		return m.FuncParsePacket(msgBytes)
	}
	buffToData := bytex.DefaultPoolBuffToData.GetInstance()
	defer bytex.DefaultPoolBuffToData.Recycle(buffToData)

	buffToData.WriteBytes(msgBytes)
	extName = buffToData.ReadString()
	pid = buffToData.ReadString()
	uid = buffToData.ReadString()
	if buffToData.Len() <= 0 {
		return
	}
	index := 0
	for buffToData.Len() > 0 {
		n, d := buffToData.ReadDataTo(msgBytes[index:]) //由于msgBytes前部分数据已经处理完成，可以利用这部分空间
		//h.singleCase.GetLogger().Traceln("parsePackMessage", uid, d)
		if n == 0 { // 没有读到字节，注意 n!=0时, d是有可能是nil的
			//h.singleCase.GetLogger().Warnln("data is nil")
			break
		}
		data = append(data, d)
		index += n
	}
	return extName, pid, uid, data
}

func (m *RabbitExtensionManager) Verify(extName string, pid string, uid string) (e server.IRabbitExtension, rsCode int32) {
	ext, ok := m.GetRabbitExtension(extName)
	// 有效性验证
	if !ok {
		if nil != m.Logger {
			m.Logger.Warnln("[ExtensionManager.Verify]", fmt.Sprintf("Undefined Extension(%s)! Sender(%s)", extName, uid))
		}
		return nil, server.CodeProtoFail
	}
	if !ext.CheckProtoId(pid) { //有效性检查
		if nil != m.Logger {
			m.Logger.Warnln("[ExtensionManager.Verify]", fmt.Sprintf("Undefined ProtoId(%s) Send to Extension(%s)! Sender(%s)", pid, extName, uid))
		}
		return nil, server.CodeProtoFail
	}
	return ext, m.CustomVerify(extName, pid, uid)
}

// GetRecycleParams
// 构造响应参数
func (m *RabbitExtensionManager) GetRecycleParams(extension server.IRabbitExtension, senderAddress string, name string, pid string, uid string, data [][]byte) (resp server.IExtensionResponse, req server.IExtensionRequest) {
	t, h := extension.GetParamInfo(pid)
	response := DefaultResponsePool.GetInstance().(iExtResponse)
	response.SetHeader(name, pid, uid, senderAddress)
	response.SetSockSender(m.SockSender)
	response.SetAddressProxy(m.AddressProxy)
	response.SetResultCode(server.CodeSuc)
	response.SetParamInfo(t, h)
	request := DefaultRequestPool.GetInstance()
	request.SetHeader(name, pid, uid, senderAddress)
	request.SetRequestData(t, h, data)
	return response, request
}

// GetRecycleResponse
// 构造响应参数
func (m *RabbitExtensionManager) GetRecycleResponse(extension server.IRabbitExtension, senderAddress string, name string, pid string, uid string, data [][]byte) (resp server.IExtensionResponse) {
	t, h := extension.GetParamInfo(pid)
	response := DefaultResponsePool.GetInstance().(iExtResponse)
	response.SetHeader(name, pid, uid, senderAddress)
	response.SetSockSender(m.SockSender)
	response.SetAddressProxy(m.AddressProxy)
	response.SetResultCode(server.CodeSuc)
	response.SetParamInfo(t, h)
	return response
}

// GetRecycleRequest
// 获取可回收的请求结构
func (m *RabbitExtensionManager) GetRecycleRequest(extension server.IRabbitExtension, senderAddress string, name string, pid string, uid string, data [][]byte) (req server.IExtensionRequest) {
	t, h := extension.GetParamInfo(pid)
	request := DefaultRequestPool.GetInstance()
	request.SetHeader(name, pid, uid, senderAddress)
	request.SetRequestData(t, h, data)
	return request
}

func (m *RabbitExtensionManager) DoRequest(extension server.IRabbitExtension, request server.IExtensionRequest, response server.IExtensionResponse) {
	// 响应处理
	if be, ok := extension.(server.IBeforeRequestExtension); ok { //前置处理
		be.BeforeRequest(response, request)
	}
	if re, ok := extension.(server.IOnRequestExtension); ok {
		m.CustomStartOnRequest(response, request)
		re.OnRequest(response, request)
		m.CustomFinishOnRequest(response, request)
	}
	if ae, ok := extension.(server.IAfterRequestExtension); ok { //后置处理
		ae.AfterRequest(response, request)
	}
}

func (m *RabbitExtensionManager) GetRabbitExtension(extName string) (pe server.IRabbitExtension, ok bool) {
	if pe, ok := m.ExtensionContainer.GetExtension(extName).(server.IRabbitExtension); ok {
		return pe, true
	}
	return nil, false
}
