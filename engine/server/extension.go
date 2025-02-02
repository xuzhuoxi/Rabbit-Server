// Package extendx
// Created by xuzhuoxi
// on 2019-02-17.
// @author xuzhuoxi
//
package server

import (
	"github.com/xuzhuoxi/infra-go/logx"
)

type IExtension interface {
	// ExtensionName 主键标识
	ExtensionName() string
}

type IInitExtension interface {
	// InitExtension 初始化
	InitExtension() error
	// DestroyExtension 反初始化
	DestroyExtension() error
}

type ISaveExtension interface {
	// SaveExtension 保存数据
	SaveExtension() error
}

type IEnableExtension interface {
	// Enable 是否启用
	Enable() bool
	// EnableExtension 启用
	EnableExtension() error
	// DisableExtension 禁用
	DisableExtension() error
}

type IGoroutineExtension interface {
	// MaxGo
	// 最大并发处理个数
	MaxGo() int
}

type IProtoExtension interface {
	logx.ILoggerSupport
	IExtension
	IInitExtension
	// CheckProtocolId
	// 检查ProtoId是否为被当前扩展支持
	CheckProtocolId(protoId string) bool
	// GetParamInfo
	// 检查ProtoId对应的设置
	GetParamInfo(protoId string) (paramType ExtensionParamType, handler IProtoParamsHandler)
}

// ExtensionHandlerNoneParam
// Extension响应函数－无参数
type ExtensionHandlerNoneParam func(resp IExtensionResponse, req IExtensionRequest)

// ExtensionHandlerBinaryParam
// Extension响应函数－二进制参数
type ExtensionHandlerBinaryParam func(resp IExtensionResponse, req IBinaryRequest)

// ExtensionHandlerStringParam
// Extension响应函数－字符串参数
type ExtensionHandlerStringParam func(resp IExtensionResponse, req IStringRequest)

// ExtensionHandlerObjectParam
// Extension响应函数－具体对象参数
type ExtensionHandlerObjectParam func(resp IExtensionResponse, req IObjectRequest)

//-------------------------------------------------------------

type IRequestExtensionSetter interface {
	// SetRequestHandler
	// 设置请求响应处理(无参数)
	SetRequestHandler(protoId string, handler ExtensionHandlerNoneParam)
	// SetRequestHandlerBinary
	// 设置请求响应处理(字节数组参数)
	SetRequestHandlerBinary(protoId string, handler ExtensionHandlerBinaryParam)
	// SetRequestHandlerString
	// 设置请求响应处理(字符串参数)
	SetRequestHandlerString(protoId string, handler ExtensionHandlerStringParam)
	// SetRequestHandlerObject
	//设置请求响应处理(对象参数)
	SetRequestHandlerObject(protoId string, handler ExtensionHandlerObjectParam,
		paramOrigin interface{}, paramHandler IProtoParamsHandler)
	// ClearRequestHandler
	// 清除设置
	ClearRequestHandler(protoId string)
}

type IOnRequestExtension interface {
	// OnRequest
	// 请求响应
	OnRequest(resp IExtensionResponse, req IExtensionRequest)
}

type IBeforeRequestExtension interface {
	// BeforeRequest
	// 执行响应前的一些处理
	BeforeRequest(req IExtensionRequest)
}

type IAfterRequestExtension interface {
	// AfterRequest
	// 响应结束前的一些处理
	AfterRequest(resp IExtensionResponse, req IExtensionRequest)
}
