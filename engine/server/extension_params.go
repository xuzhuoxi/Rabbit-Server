// Package server
// Create on 2025/2/2
// @author xuzhuoxi
package server

import "github.com/xuzhuoxi/infra-go/encodingx"

// ExtensionParamType
// Extension响应参数类型
type ExtensionParamType int

const (
	None ExtensionParamType = iota
	Binary
	String
	Object
)

// FuncParamCtor 通用参数构造器
type FuncParamCtor = func() interface{}

// IProtoParamsHandler
// 协议参数处理器接口
// 要求：并发安全
type IProtoParamsHandler interface {
	// SetCodingHandler
	// 设置编解码器
	SetCodingHandler(handler encodingx.ICodingHandler)
	// SetCodingHandlers
	// 设置编解码器
	SetCodingHandlers(reqHandler encodingx.ICodingHandler, returnHandler encodingx.ICodingHandler)
	// HandleRequestParam
	// 处理请求参数转换：二进制->具体数据
	HandleRequestParam(data []byte) interface{}
	// HandleRequestParams
	// 处理请求参数转换：二进制->具体数据(批量)
	HandleRequestParams(data [][]byte) []interface{}
	// HandleReturnParam
	// 处理响应参数转换：具体数据->二进制
	HandleReturnParam(data interface{}) []byte
	// HandleReturnParams
	// 处理响应参数转换：具体数据->二进制(批量)
	HandleReturnParams(data []interface{}) [][]byte
}
