// Package extension
// Created by xuzhuoxi
// on 2019-05-18.
// @author xuzhuoxi
//
package extension

import (
	"fmt"
	"github.com/xuzhuoxi/Rabbit-Server/engine/server"
	"github.com/xuzhuoxi/Rabbit-Server/engine/server/packet"
)

func NewSockRequest() *SockRequest {
	return &SockRequest{}
}

type SockRequest struct {
	packet.PacketHeader
	ParamType server.ExtensionParamType
	binData   [][]byte
	strData   []string
	objData   []interface{}
}

func (req *SockRequest) String() string {
	return fmt.Sprintf("{Request: %v, %v, %v, %v}",
		req.PacketHeader, req.ParamType, req.binData, req.objData)
}

func (req *SockRequest) DataSize() int {
	switch req.ParamType {
	case server.Binary:
		return len(req.binData)
	case server.String:
		return len(req.strData)
	case server.Object:
		return len(req.objData)
	}
	return 0
}

func (req *SockRequest) SetRequestData(paramType server.ExtensionParamType, paramHandler server.IPacketParamsHandler, data [][]byte) {
	req.ParamType = paramType
	req.binData = data
	switch paramType {
	case server.None:
		req.strData, req.objData = nil, nil
	case server.Binary:
		req.strData, req.objData = nil, nil
	case server.String:
		req.strData, req.objData = req.toStringArray(data), nil
	case server.Object:
		objData := paramHandler.DecodeRequestParams(data)
		req.strData, req.objData = nil, objData
	}
}

func (req *SockRequest) BinaryData() [][]byte {
	return req.binData
}

func (req *SockRequest) FirstBinary() []byte {
	if len(req.binData) == 0 {
		return nil
	} else {
		return req.binData[0]
	}
}

func (req *SockRequest) StringData() []string {
	return req.strData
}

func (req *SockRequest) FirstString() string {
	if len(req.strData) == 0 {
		return ""
	} else {
		return req.strData[0]
	}
}

func (req *SockRequest) ObjectData() []interface{} {
	return req.objData
}

func (req *SockRequest) FirstObject() interface{} {
	if len(req.objData) == 0 {
		return nil
	} else {
		return req.objData[0]
	}
}

func (req *SockRequest) toStringArray(data [][]byte) []string {
	if nil == data || len(data) == 0 {
		return nil
	}
	rs := make([]string, len(data))
	for index := range data {
		rs[index] = string(data[index])
	}
	return rs
}
