// Package extension
// Created by xuzhuoxi
// on 2019-05-12.
// @author xuzhuoxi
//
package extension

import (
	"github.com/xuzhuoxi/Rabbit-Server/engine/server"
)

const (
	NameDemo  = "Demo"
	NameLogin = "Login"
)

func init() {
	server.RegisterRabbitExtension(NameDemo, NewRabbitDemoExtension)
	server.RegisterRabbitExtension(NameLogin, func(name string) server.IRabbitExtension { return NewRabbitLoginExtension(name) })
}
