// Package packet
// Created by xuzhuoxi
// on 2019-05-20.
// @author xuzhuoxi
//
package packet

import (
	"fmt"
	"github.com/xuzhuoxi/Rabbit-Server/engine/server"
)

type PacketHeader struct {
	EName    string
	PId      string
	CId      string
	CAddress string
}

func (h *PacketHeader) String() string {
	return fmt.Sprintf("Header{EName='%s', PId='%s', CID='%s', CAddr='%s'}",
		h.EName, h.PId, h.CId, h.CAddress)
}

func (h *PacketHeader) GetHeaderInfo() server.IPacketHeader {
	return h
}

func (h *PacketHeader) ExtName() string {
	return h.EName
}

func (h *PacketHeader) ProtoId() string {
	return h.PId
}

func (h *PacketHeader) ClientId() string {
	return h.CId
}

func (h *PacketHeader) ClientAddress() string {
	return h.CAddress
}

func (h *PacketHeader) SetHeader(extName string, protoId string, clientId string, clientAddress string) {
	h.EName, h.PId, h.CId, h.CAddress = extName, protoId, clientId, clientAddress
}
