// Package extension
// Created by xuzhuoxi
// on 2019-02-18.
// @author xuzhuoxi
//
package extension

import (
	"github.com/xuzhuoxi/Rabbit-Server/engine/server/runtime"
	"github.com/xuzhuoxi/infra-go/logx"
)

func NewRabbitDemoExtensionSupport(Name string) RabbitDemoExtensionSupport {
	support := runtime.NewExtensionMeta(Name)
	return RabbitDemoExtensionSupport{ExtensionMeta: support}
}

type RabbitDemoExtensionSupport struct {
	runtime.ExtensionMeta
	logx.LoggerSupport
}
