// Package server
// Create on 2025/2/2
// @author xuzhuoxi
package server

import "github.com/xuzhuoxi/infra-go/netx"

type IExtensionRequestSettings interface {
	netx.IConnInfoSetter
}

// IExtensionRequest
// 请求对象参数集合接口
type IExtensionRequest interface {
	IPacketHeader
	netx.IConnInfoGetter
	// DataSize 数据长度
	DataSize() int
	// SetRequestData
	// 设置集合数据信息
	SetRequestData(paramType ExtensionParamType, paramHandler IPacketParamsHandler, data [][]byte)
}

// IBinaryRequest
// 数据参数为二进制的请求对象参数集合接口
type IBinaryRequest interface {
	IExtensionRequest
	// BinaryData
	// RequestBinaryData
	// 请求的参数数据(二进制)
	BinaryData() [][]byte
	// FirstBinary
	// 第一个请求参数
	FirstBinary() []byte
}

// IStringRequest
// 数据参数为Json的请求对象参数集合接口
type IStringRequest interface {
	IExtensionRequest
	// StringData
	// 请求的参数数据(String)
	StringData() []string
	// FirstString
	// 第一个请求参数
	FirstString() string
}

// IObjectRequest
// 数据参数为结构体的请求对象参数集合接口
type IObjectRequest interface {
	IExtensionRequest
	// ObjectData
	// 请求的参数数据(具体数据)
	ObjectData() []interface{}
	// FirstObject
	// 第一个请求参数
	FirstObject() interface{}
}
