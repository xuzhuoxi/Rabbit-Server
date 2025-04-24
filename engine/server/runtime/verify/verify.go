// Package verify
// Create on 2023/8/23
// @author xuzhuoxi
package verify

import "github.com/xuzhuoxi/Rabbit-Server/engine/server"

func NewIPacketVerifier() server.IPacketVerifier {
	return &PacketVerifier{}
}

type packetVerifyItem struct {
	handler  server.FuncVerifyPacket
	verifier server.IPacketVerifier
}

type PacketVerifier struct {
	handlers []*packetVerifyItem
}

func (o *PacketVerifier) Clear() {
	o.handlers = nil
}

func (o *PacketVerifier) AppendVerifyHandler(handler server.FuncVerifyPacket) {
	o.handlers = append(o.handlers, &packetVerifyItem{handler: handler})
}

func (o *PacketVerifier) AppendVerifier(verifier server.IPacketVerifier) {
	o.handlers = append(o.handlers, &packetVerifyItem{verifier: verifier})
}

func (o *PacketVerifier) Verify(extName string, pid string, uid string) (rsCode int32) {
	if len(o.handlers) == 0 {
		return
	}
	for _, item := range o.handlers {
		if item == nil {
			continue
		}
		if item.handler != nil {
			code := item.handler(extName, pid, uid)
			if code != server.CodeSuc {
				return code
			}
		} else if item.verifier != nil {
			code := item.verifier.Verify(extName, pid, uid)
			if code != server.CodeSuc {
				return code
			}
		}
	}
	return
}
