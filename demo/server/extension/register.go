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
	server.RegisterRabbitExtension(NameDemo, func() server.IRabbitExtension { return NewRabbitDemoExtension(NameDemo) })
	server.RegisterRabbitExtension(NameLogin, func() server.IRabbitExtension { return NewRabbitLoginExtension(NameLogin) })
}
