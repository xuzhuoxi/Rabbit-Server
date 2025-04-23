// Package params
// Created by xuzhuoxi
// on 2019-05-18.
// @author xuzhuoxi
//
package params

import (
	"fmt"
	"github.com/xuzhuoxi/Rabbit-Server/engine/server"
	"github.com/xuzhuoxi/Rabbit-Server/engine/server/runtime/params/packets"
	"github.com/xuzhuoxi/infra-go/netx"
)

func NewISockRequest() server.IExtensionRequest {
	return NewSockRequest()
}

func NewSockRequest() *SockRequest {
	return &SockRequest{}
}

type iExtRequest interface {
	server.IExtensionRequestSettings
	server.IExtensionRequest
}

type SockRequest struct {
	packets.PacketHeader
	ParamType server.ExtensionParamType
	connInfo  netx.IConnInfo
	binData   [][]byte
	strData   []string
	objData   []interface{}
}

func (o *SockRequest) SetConnInfo(connInfo netx.IConnInfo) {
	o.connInfo = connInfo
}

func (o *SockRequest) GetConnInfo() netx.IConnInfo {
	return o.connInfo
}

func (o *SockRequest) String() string {
	return fmt.Sprintf("{Request: %v, %v, %v, %v}",
		o.PacketHeader, o.ParamType, o.binData, o.objData)
}

func (o *SockRequest) DataSize() int {
	switch o.ParamType {
	case server.Binary:
		return len(o.binData)
	case server.String:
		return len(o.strData)
	case server.Object:
		return len(o.objData)
	}
	return 0
}

func (o *SockRequest) SetRequestData(paramType server.ExtensionParamType, paramHandler server.IPacketCoding, data [][]byte) {
	o.ParamType = paramType
	o.binData = data
	switch paramType {
	case server.None:
		o.strData, o.objData = nil, nil
	case server.Binary:
		o.strData, o.objData = nil, nil
	case server.String:
		o.strData, o.objData = o.toStringArray(data), nil
	case server.Object:
		objData := paramHandler.DecodeRequestParams(data)
		o.strData, o.objData = nil, objData
	}
}

func (o *SockRequest) BinaryData() [][]byte {
	return o.binData
}

func (o *SockRequest) FirstBinary() []byte {
	if len(o.binData) == 0 {
		return nil
	} else {
		return o.binData[0]
	}
}

func (o *SockRequest) StringData() []string {
	return o.strData
}

func (o *SockRequest) FirstString() string {
	if len(o.strData) == 0 {
		return ""
	} else {
		return o.strData[0]
	}
}

func (o *SockRequest) ObjectData() []interface{} {
	return o.objData
}

func (o *SockRequest) FirstObject() interface{} {
	if len(o.objData) == 0 {
		return nil
	} else {
		return o.objData[0]
	}
}

func (o *SockRequest) toStringArray(data [][]byte) []string {
	if nil == data || len(data) == 0 {
		return nil
	}
	rs := make([]string, len(data))
	for index := range data {
		rs[index] = string(data[index])
	}
	return rs
}
