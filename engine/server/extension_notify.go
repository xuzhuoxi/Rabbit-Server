// Package server
// Create on 2025/2/2
// @author xuzhuoxi
package server

import "github.com/xuzhuoxi/infra-go/encodingx"

type IExtensionNotify interface {
	IResponsePacket
	SetCodingHandler(encodeHandler encodingx.ICodingHandler)
}
