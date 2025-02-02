// Package proto
// Create on 2025/2/2
// @author xuzhuoxi
package server

// IProtoHeader
// 协议参数头接口
type IProtoHeader interface {
	// ProtoGroup
	// 协议分组
	ProtoGroup() string
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
	GetHeaderInfo() IProtoHeader
}
