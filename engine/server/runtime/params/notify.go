// Package params
// Create on 2023/8/6
// @author xuzhuoxi
package params

import (
	"github.com/xuzhuoxi/Rabbit-Server/engine/server"
	"github.com/xuzhuoxi/Rabbit-Server/engine/server/runtime/params/packets"
	"github.com/xuzhuoxi/infra-go/encodingx"
)

func NewISockNotify() server.IExtensionNotify {
	return NewSockNotify()
}

func NewSockNotify() *SockNotify {
	return &SockNotify{ResponsePacket: *packets.NewResponsePacket()}
}

type SockNotify struct {
	packets.ResponsePacket
}

func (o *SockNotify) SetCodingHandler(codingHandler encodingx.ICodingHandler) {
	o.ParamHandler = packets.NewIPacketCoding(nil, codingHandler)
}
