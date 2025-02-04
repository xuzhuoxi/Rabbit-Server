// Package extension
// Create on 2023/8/6
// @author xuzhuoxi
package extension

import (
	"github.com/xuzhuoxi/Rabbit-Server/engine/server/packet"
	"github.com/xuzhuoxi/infra-go/encodingx"
)

func NewSockNotify() *SockNotify {
	return &SockNotify{ResponsePacket: *packet.NewResponsePacket()}
}

type SockNotify struct {
	packet.ResponsePacket
}

func (o *SockNotify) SetCodingHandler(codingHandler encodingx.ICodingHandler) {
	o.ParamHandler = packet.NewPacketParamsHandler(nil, codingHandler)
}
