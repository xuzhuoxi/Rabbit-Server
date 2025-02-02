// Package protox
// Created by xuzhuoxi
// on 2019-05-20.
// @author xuzhuoxi
//
package proto

import (
	"fmt"
	"github.com/xuzhuoxi/Rabbit-Server/engine/server"
)

type ProtoHeader struct {
	PGroup   string
	PId      string
	CId      string
	CAddress string
}

func (h *ProtoHeader) String() string {
	return fmt.Sprintf("Header{ENmae='%s', PId='%s', CID='%s', CAddr='%s'}",
		h.PGroup, h.PId, h.CId, h.CAddress)
}

func (h *ProtoHeader) GetHeaderInfo() server.IProtoHeader {
	return h
}

func (h *ProtoHeader) ProtoGroup() string {
	return h.PGroup
}

func (h *ProtoHeader) ProtoId() string {
	return h.PId
}

func (h *ProtoHeader) ClientId() string {
	return h.CId
}

func (h *ProtoHeader) ClientAddress() string {
	return h.CAddress
}

func (h *ProtoHeader) SetHeader(extensionName string, protoId string, clientId string, clientAddress string) {
	h.PGroup, h.PId, h.CId, h.CAddress = extensionName, protoId, clientId, clientAddress
}
