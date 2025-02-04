// Package server
// Created by xuzhuoxi
// on 2019-02-17.
// @author xuzhuoxi
//
package server

import (
	"github.com/xuzhuoxi/infra-go/encodingx"
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

type IRabbitExtension interface {
	logx.ILoggerSupport
	IExtension
	IInitExtension
	// CheckProtoId
	// 检查ProtoId是否为被当前扩展支持
	CheckProtoId(protoId string) bool
	// GetParamInfo
	// 检查ProtoId对应的设置
	GetParamInfo(protoId string) (paramType ExtensionParamType, handler IPacketParamsHandler)
}

// FuncBeforeRequest 响应前置函数
type FuncBeforeRequest = func(resp IExtensionResponse, req IExtensionRequest)

// FuncAfterRequest 响应后置函数
type FuncAfterRequest = func(resp IExtensionResponse, req IExtensionRequest)

// FuncOnNoneParamRequest
// Extension响应函数－无参数
type FuncOnNoneParamRequest = func(resp IExtensionResponse, req IExtensionRequest)

// FuncOnBinaryRequest
// Extension响应函数－二进制参数
type FuncOnBinaryRequest = func(resp IExtensionResponse, req IBinaryRequest)

// FuncOnStringRequest
// Extension响应函数－字符串参数
type FuncOnStringRequest = func(resp IExtensionResponse, req IStringRequest)

// FuncOnObjectRequest
// Extension响应函数－具体对象参数
type FuncOnObjectRequest = func(resp IExtensionResponse, req IObjectRequest)

type IRequestExtensionSetter interface {
	// SetBeforeRequestHandler
	// 设置前置响应处理
	SetBeforeRequestHandler(protoId string, handler FuncBeforeRequest)
	// SetAfterRequestHandler
	// 设置后置响应处理
	SetAfterRequestHandler(protoId string, handler FuncAfterRequest)
	// SetOnRequestHandler
	// 设置请求响应处理(无参数)
	SetOnRequestHandler(protoId string, handler FuncOnNoneParamRequest)
	// SetOnBinaryRequestHandler
	// 设置请求响应处理(字节数组参数)
	SetOnBinaryRequestHandler(protoId string, handler FuncOnBinaryRequest)
	// SetOnStringRequestHandler
	// 设置请求响应处理(字符串参数)
	SetOnStringRequestHandler(protoId string, handler FuncOnStringRequest)
	// SetOnObjectRequestHandler
	//设置请求响应处理(对象参数)
	SetOnObjectRequestHandler(protoId string, handler FuncOnObjectRequest, ctor FuncParamObjectCtor, codingHandler encodingx.ICodingHandler)
	// ClearRequestHandler
	// 清除设置
	ClearRequestHandler(protoId string)
	// ClearRequestHandlers
	// 清除全部设置
	ClearRequestHandlers()
}

type IOnRequestExtension interface {
	// OnRequest
	// 请求响应
	OnRequest(resp IExtensionResponse, req IExtensionRequest)
}

type IBeforeRequestExtension interface {
	// BeforeRequest
	// 执行响应前的一些处理
	BeforeRequest(resp IExtensionResponse, req IExtensionRequest)
}

type IAfterRequestExtension interface {
	// AfterRequest
	// 响应结束后的一些处理
	AfterRequest(resp IExtensionResponse, req IExtensionRequest)
}
