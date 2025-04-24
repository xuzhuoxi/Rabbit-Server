// Package server
// Created by xuzhuoxi
// on 2019-02-17.
// @author xuzhuoxi
//
package server

import (
	"github.com/xuzhuoxi/infra-go/cryptox"
	"github.com/xuzhuoxi/infra-go/encodingx"
	"github.com/xuzhuoxi/infra-go/logx"
	"github.com/xuzhuoxi/infra-go/netx"
)

// Extension ---------- ---------- ---------- ---------- ----------

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
	// SetEnable 设置启用状态
	SetEnable(enable bool) error
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
	GetParamInfo(protoId string) (paramType ExtensionParamType, handler IPacketCoding)
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

type IBeforeRequestExtension interface {
	// BeforeRequest
	// 执行响应前的一些处理
	BeforeRequest(resp IExtensionResponse, req IExtensionRequest)
}

type IOnRequestExtension interface {
	// OnRequest
	// 请求响应
	OnRequest(resp IExtensionResponse, req IExtensionRequest)
}

type IAfterRequestExtension interface {
	// AfterRequest
	// 响应结束后的一些处理
	AfterRequest(resp IExtensionResponse, req IExtensionRequest)
}

// Container ---------- ---------- ---------- ---------- ----------

type IExtensionContainer interface {
	// AppendExtension
	// 增加Extension
	AppendExtension(ext IExtension)
	// CheckExtension
	// 检查
	CheckExtension(named string) bool
	// GetExtension
	// 取Extension
	GetExtension(named string) IExtension
	// Len
	// Extension数量
	Len() int
	// Extensions
	// 列表
	Extensions() []IExtension
	// ExtensionsReversed
	// 反向列表
	ExtensionsReversed() []IExtension
	// Range
	// 按列表处理
	Range(handler func(index int, ext IExtension))
	// RangeReverse
	// 按反向列表处理
	RangeReverse(handler func(index int, ext IExtension))
	// HandleAt
	// 对指定Extension执行处理
	HandleAt(index int, handler func(index int, ext IExtension)) error
	// HandleAtName
	// 对指定Extension执行处理
	HandleAtName(name string, handler func(name string, ext IExtension)) error
}

type iInitExtensions interface {
	// InitExtensions
	// 初始化全部Extension
	InitExtensions() []error
	// DestroyExtensions
	// 销毁全部Extension
	DestroyExtensions() []error
}

type iSaveExtensions interface {
	// SaveExtension
	// 保存数据
	SaveExtension(name string) error
	// SaveExtensions
	// 保存数据
	SaveExtensions() []error
}

type iEnableExtensions interface {
	// EnableExtension
	// 设置Extension的激活状态
	EnableExtension(extName string, enable bool) error
	// EnableExtensions
	// 设置全部Extension的激活状态
	EnableExtensions(enable bool) []error
}

// IRabbitExtensionContainer
// Extension容器接口
type IRabbitExtensionContainer interface {
	IExtensionContainer
	iInitExtensions
	iSaveExtensions
	iEnableExtensions
}

// IRabbitExtensionManager
// Extension管理器接口
type IRabbitExtensionManager interface {
	logx.ILoggerSetter
	netx.IUserConnMapperSetter
	iSaveExtensions
	iEnableExtensions

	// InitManager
	// 初始化
	// handlerContainer: 解包处理
	// extensionContainer： 服务扩展
	// sockSender: 消息发送器
	InitManager(handlerContainer netx.IPackHandlerContainer, extensionContainer IRabbitExtensionContainer, sockSender netx.ISockSender)

	// StartManager
	// 开始运行
	StartManager()
	// StopManager
	// 停止运行
	StopManager()

	// SetPacketCipher
	// 设置消息包加密解密处理器
	SetPacketCipher(cipher cryptox.ICipher)
	// AppendVerifyHandler
	// 添加消息验证处理器
	AppendVerifyHandler(handler FuncVerifyPacket)
	// OnMessageUnpack
	// 消息处理入口，这里是并发方法
	OnMessageUnpack(msgData []byte, connInfo netx.IConnInfo, other interface{}) bool
}
