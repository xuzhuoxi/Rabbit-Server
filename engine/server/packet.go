// Package server
// Create on 2025/2/2
// @author xuzhuoxi
package server

import "github.com/xuzhuoxi/infra-go/encodingx"

// ExtensionParamType
// Extension响应参数类型
type ExtensionParamType int

const (
	// None 无参数类型
	None ExtensionParamType = iota
	// Binary 二进制参数类型
	Binary
	// String 字符串参数类型
	String
	// Object 结构体参数类型
	Object
)

// IPacketHeader
// 协议参数头接口
type IPacketHeader interface {
	// ExtName
	// 协议分组
	ExtName() string
	// ProtoId
	// 协议标识
	ProtoId() string
	// ClientId
	// 客户端标识
	ClientId() string
	// ClientAddress
	// 客户端地址
	ClientAddress() string
	// SetHeader
	// 设置参数头信息
	SetHeader(extensionName string, protoId string, clientId string, clientAddress string)
	// GetHeaderInfo 取头信息
	GetHeaderInfo() IPacketHeader
}

// FuncParamObjectCtor 通用对象参数构造器
type FuncParamObjectCtor = func() interface{}

// IPacketCoding
// 协议参数处理器接口
// 要求：并发安全
type IPacketCoding interface {
	// SetCodingHandler
	// 设置编解码器
	SetCodingHandler(handler encodingx.ICodingHandler)
	// SetCodingHandlers
	// 设置编解码器
	SetCodingHandlers(reqHandler encodingx.ICodingHandler, returnHandler encodingx.ICodingHandler)
	// DecodeRequestParam
	// 处理请求参数转换：二进制->具体数据
	DecodeRequestParam(data []byte) interface{}
	// DecodeRequestParams
	// 处理请求参数转换：二进制->具体数据(批量)
	DecodeRequestParams(data [][]byte) []interface{}
	// EncodeResponseParam
	// 处理响应参数转换：具体数据->二进制
	EncodeResponseParam(data interface{}) []byte
	// EncodeResponseParams
	// 处理响应参数转换：具体数据->二进制(批量)
	EncodeResponseParams(data []interface{}) [][]byte
}

// FuncVerifyPacket
// 消息数据包校验函数
type FuncVerifyPacket = func(extName string, pid string, uid string) (rsCode int32)

// FuncNewIPacketVerifier
// 消息数据包校验器构造函数
type FuncNewIPacketVerifier = func() IPacketVerifier

// IPacketVerifyItem 协议单项检验器
type IPacketVerifyItem interface {
	// Verify 验证请求入口
	Verify(extName string, pid string, uid string) (rsCode int32)
}

// IPacketVerifier 协议检验器
type IPacketVerifier interface {
	// Clear 清除
	Clear()
	// AppendVerifyHandler
	// 追加验证处理函数
	AppendVerifyHandler(handler FuncVerifyPacket)
	// AppendVerifier
	// 追加验证处理器
	AppendVerifier(verify IPacketVerifier)
	// Verify 验证请求入口
	Verify(extName string, pid string, uid string) (rsCode int32)
}
