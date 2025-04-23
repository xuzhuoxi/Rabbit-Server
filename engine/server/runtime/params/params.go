// Package params
// Create on 2025/4/23
// @author xuzhuoxi
package params

import (
	"github.com/xuzhuoxi/Rabbit-Server/engine/server"
	"github.com/xuzhuoxi/infra-go/netx"
)

type ParamsSupport struct {
	SockSender     netx.ISockSender
	UserConnMapper netx.IUserConnMapper
}

// GetRecycleParams
// 构造响应参数
func (o *ParamsSupport) GetRecycleParams(extension server.IRabbitExtension, connInfo netx.IConnInfo,
	name string, pid string, uid string, data [][]byte) (resp server.IExtensionResponse, req server.IExtensionRequest) {
	t, h := extension.GetParamInfo(pid)
	return o.getResponse(t, h, connInfo, name, pid, uid), o.getRequest(t, h, connInfo, name, pid, uid, data)
}

// GetRecycleResponse
// 构造响应参数
func (o *ParamsSupport) GetRecycleResponse(extension server.IRabbitExtension, connInfo netx.IConnInfo,
	name string, pid string, uid string) (resp server.IExtensionResponse) {
	t, h := extension.GetParamInfo(pid)
	return o.getResponse(t, h, connInfo, name, pid, uid)
}

// GetRecycleRequest
// 获取可回收的请求结构
func (o *ParamsSupport) GetRecycleRequest(extension server.IRabbitExtension, connInfo netx.IConnInfo,
	name string, pid string, uid string, data [][]byte) (req server.IExtensionRequest) {
	t, h := extension.GetParamInfo(pid)
	return o.getRequest(t, h, connInfo, name, pid, uid, data)
}

func (o *ParamsSupport) getResponse(pType server.ExtensionParamType, pHandler server.IPacketCoding, connInfo netx.IConnInfo,
	name string, pid string, uid string) (resp server.IExtensionResponse) {
	response := DefaultResponsePool.GetInstance().(iExtResponse)
	response.SetHeader(name, pid, uid, connInfo.GetRemoteAddress())
	response.SetSockSender(o.SockSender)
	response.SetUserConnMapper(o.UserConnMapper)
	response.SetConnInfo(connInfo)
	response.SetResultCode(server.CodeSuc)
	response.SetParamInfo(pType, pHandler)
	return response
}

func (o *ParamsSupport) getRequest(pType server.ExtensionParamType, pHandler server.IPacketCoding, connInfo netx.IConnInfo,
	name string, pid string, uid string, data [][]byte) (req server.IExtensionRequest) {
	request := DefaultRequestPool.GetInstance().(iExtRequest)
	request.SetHeader(name, pid, uid, connInfo.GetRemoteAddress())
	request.SetConnInfo(connInfo)
	request.SetRequestData(pType, pHandler, data)
	return request
}
