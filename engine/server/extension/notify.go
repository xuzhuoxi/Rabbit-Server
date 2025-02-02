// Package protox
// Create on 2023/8/6
// @author xuzhuoxi
package extension

import (
	"github.com/xuzhuoxi/infra-go/encodingx"
)

func NewSockNotify() *SockNotify {
	return &SockNotify{ProtoReturnMessage: *NewProtoReturnMessage()}
}

type SockNotify struct {
	ProtoReturnMessage
}

func (o *SockNotify) SetCodingHandler(codingHandler encodingx.ICodingHandler) {
	o.ParamHandler = NewProtoObjectParamsHandler(nil, codingHandler)
}
