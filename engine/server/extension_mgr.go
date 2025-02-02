// Package server
// Create on 2025/2/2
// @author xuzhuoxi
package server

import (
	"github.com/xuzhuoxi/infra-go/logx"
	"github.com/xuzhuoxi/infra-go/netx"
)

type IExtensionContainer interface {
	// AppendExtension
	// 增加Extension
	AppendExtension(ext IExtension)
	// CheckExtension
	// 检查
	CheckExtension(name string) bool
	// GetExtension
	// 取Extension
	GetExtension(name string) IExtension
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

type IProtoExtensionContainer interface {
	IExtensionContainer
	// InitExtensions
	// 初始化全部Extension
	InitExtensions() []error
	// DestroyExtensions
	// 销毁全部Extension
	DestroyExtensions() []error

	// SaveExtensions
	// 保存
	SaveExtensions() []error
	// SaveExtension
	// 保存指定
	SaveExtension(name string) error

	// EnableExtensions
	// 设置启用全部Extension
	EnableExtensions(enable bool) []error
	// EnableExtension
	// 设置启用Extension
	EnableExtension(name string, enable bool) error
}

// IExtensionManager
// Extension管理接口
type IExtensionManager interface {
	logx.ILoggerSetter
	netx.IAddressProxySetter

	// InitManager
	// 初始化
	// handlerContainer: 解包处理
	// extensionContainer： 服务扩展
	// sockSender: 消息发送器
	InitManager(handlerContainer netx.IPackHandlerContainer, extensionContainer IProtoExtensionContainer, sockSender netx.ISockSender)

	// StartManager
	// 开始运行
	StartManager()
	// StopManager
	// 停止运行
	StopManager()

	// SaveExtension
	// 保存指定Extension的临时数据
	SaveExtension(name string)
	// SaveExtensions
	// 保存全部Extension的临时数据
	SaveExtensions()

	// EnableExtension
	// 启用指定Extension的临时数据
	EnableExtension(name string)
	// EnableExtensions
	// 启用全部Extension的临时数据
	EnableExtensions()
	// DisableExtension
	// 禁用指定Extension的临时数据
	DisableExtension(name string)
	// DisableExtensions
	// 禁用全部Extension的临时数据
	DisableExtensions()

	// OnMessageUnpack
	// 消息处理入口，这里是并发方法
	OnMessageUnpack(msgData []byte, senderAddress string, other interface{}) bool
	// DoRequest
	// 消息处理入口，这里是并发方法
	DoRequest(extension IProtoExtension, req IExtensionRequest, resp IExtensionResponse)
}

// FuncStartOnPack
// 响应入口
type FuncStartOnPack func(senderAddress string)

// FuncParseMessage
// 解释二进制数据
type FuncParseMessage func(msgBytes []byte) (name string, pid string, uid string, data [][]byte)

// FuncGetExtension
// 消息处理入口，这里是并发方法
type FuncGetExtension func(name string) (extension IProtoExtension, rsCode int32)

// FuncStartOnRequest
// 响应开始
type FuncStartOnRequest func(resp IExtensionResponse, req IExtensionRequest)

// FuncFinishOnRequest
// 响应完成
type FuncFinishOnRequest func(resp IExtensionResponse, req IExtensionRequest)

type IExtensionManagerCustomizeSetting interface {
	// SetCustomStartOnPackFunc
	// 设置自定义响应开始行为
	SetCustomStartOnPackFunc(funcStartOnPack FuncStartOnPack)
	// SetCustomParseFunc
	// 设置自定义数据解释行为
	SetCustomParseFunc(funcParse FuncParseMessage)
	// SetCustomGetExtensionFunc
	// 设置自定义扩展获取
	SetCustomGetExtensionFunc(funcVerify FuncGetExtension)
	// SetCustomVerifyFunc
	// 设置自定义验证
	SetCustomVerifyFunc(funcVerify FuncVerify)
	// SetCustomVerify
	// 设置自定义验证接口
	SetCustomVerify(verify IReqVerify)
	// SetCustomStartOnRequestFunc
	// 设置自定义响应前置行为
	SetCustomStartOnRequestFunc(funcStart FuncStartOnRequest)
	// SetCustomFinishOnRequestFunc
	// 设置自定义响应完成行为
	SetCustomFinishOnRequestFunc(funcFinish FuncFinishOnRequest)
	// SetCustom
	// 设置自定义行为
	SetCustom(funcStartOnPack FuncStartOnPack, funcParse FuncParseMessage, funcVerify FuncVerify, funcStart FuncStartOnRequest, funcFinish FuncFinishOnRequest)
}

type IExtensionManagerCustomizeSupport interface {
	CustomStartOnPack(senderAddress string)
	CustomParseMessage(msgBytes []byte) (name string, pid string, uid string, data [][]byte)
	CustomGetExtension(name string) (extension IProtoExtension, rsCode int32)
	CustomVerify(name string, pid string, uid string) (rsCode int32)
	CustomStartOnRequest(resp IExtensionResponse, req IExtensionRequest)
	CustomFinishOnRequest(resp IExtensionResponse, req IExtensionRequest)
}
