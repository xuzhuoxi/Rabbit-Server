// Package extension
// Created by xuzhuoxi
// on 2019-02-18.
// @author xuzhuoxi
//
package extension

import (
	"github.com/xuzhuoxi/Rabbit-Server/engine/server/extension"
	"github.com/xuzhuoxi/infra-go/logx"
)

func NewRabbitDemoExtensionSupport(Name string) RabbitDemoExtensionSupport {
	support := extension.NewProtoExtensionSupport(Name)
	return RabbitDemoExtensionSupport{ProtoExtensionSupport: support}
}

type RabbitDemoExtensionSupport struct {
	extension.ProtoExtensionSupport
	logx.LoggerSupport
}
