// Package server
// Create on 2025/2/2
// @author xuzhuoxi
package server

import (
	"github.com/xuzhuoxi/infra-go/netx"
)

type IProtoReturnMessage interface {
	IProtoHeader
	// PrepareData
	// 准备设置回复参数
	PrepareData()
	// AppendLen
	// 追加参数 - 长度值
	AppendLen(ln int) error
	// AppendBinary
	// 追加参数 - 字节数组, 自动补充长度数据
	AppendBinary(data ...[]byte) error
	// AppendCommon
	// 追加参数 - 通用数据类型
	AppendCommon(data ...interface{}) error
	// AppendString
	// 追加返回- 字符串
	AppendString(data ...string) error
	// AppendJson
	// 追加返回- Json字符串 或 可序列化的Struct
	AppendJson(data ...interface{}) error
	// AppendObject
	// 追加参数
	// data only supports pointer types
	// data 只支持指针类型
	AppendObject(data ...interface{}) error
	// GenMsgBytes
	// 生成消息
	GenMsgBytes() (msg []byte, err error)
}

type IExtensionResponseSettings interface {
	netx.IAddressProxySetter
	netx.ISockSenderSetter
}

// IExtensionResponse
// 响应对象参数集合接口
type IExtensionResponse interface {
	IProtoReturnMessage
	// SetParamInfo
	// 设置参数类型与处理器
	SetParamInfo(paramType ExtensionParamType, paramHandler IProtoParamsHandler)
	// SetResultCode
	// 设置返回状态码
	SetResultCode(rsCode int32)
	// SendResponse
	// 根据设置好的参数响应
	SendResponse() error
	// SendResponseTo
	// 根据设置好的参数响应到其它用户
	SendResponseTo(interruptOnErr bool, clientIds ...string) error
	// SendNotify
	// 根据设置好的参数响应
	SendNotify(eName string, notifyPId string) error
	// SendNotifyTo
	// 根据设置好的参数响应到其它用户
	SendNotifyTo(eName string, notifyPId string, interruptOnErr bool, clientIds ...string) error
	iExtResp
}

type iExtResp interface {
	// SendNoneResponse
	// 无额外参数响应
	SendNoneResponse() error
	// SendNoneResponseToClient
	// 无额外参数响到其它用户
	SendNoneResponseToClient(interruptOnErr bool, clientIds ...string) error
	// SendBinaryResponse
	// 响应客户端(二进制参数)
	SendBinaryResponse(data ...[]byte) error
	// SendCommonResponse
	// 响应客户端(基础数据参数)
	SendCommonResponse(data ...interface{}) error
	// SendStringResponse
	// 响应客户端(字符串参数)
	SendStringResponse(data ...string) error
	// SendJsonResponse
	// 响应客户端(Json字符串参数)
	SendJsonResponse(data ...interface{}) error
	// SendObjectResponse
	// 响应客户端(具体结构体参数)
	// data only supports pointer types
	// data 只支持指针类型
	SendObjectResponse(data ...interface{}) error
}
