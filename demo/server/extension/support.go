// Package extension
// Created by xuzhuoxi
// on 2019-02-18.
// @author xuzhuoxi
//
package extension

import (
	"github.com/xuzhuoxi/infra-go/extendx/protox"
	"github.com/xuzhuoxi/infra-go/logx"
)

func NewRabbitDemoExtensionSupport(Name string) RabbitDemoExtensionSupport {
	support := protox.NewProtoExtensionSupport(Name)
	return RabbitDemoExtensionSupport{ProtoExtensionSupport: support}
}

type RabbitDemoExtensionSupport struct {
	protox.ProtoExtensionSupport
	logx.LoggerSupport
}
