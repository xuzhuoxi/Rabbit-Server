// Package extension
// Created by xuzhuoxi
// on 2019-03-03.
// @author xuzhuoxi
//
package extension

import (
	"github.com/xuzhuoxi/Rabbit-Server/engine/server"
	"github.com/xuzhuoxi/Rabbit-Server/engine/server/core"
	"time"
)

const (
	LoginId   = "LI"
	ReLoginId = "RLI"
)

func NewRabbitLoginExtension(Name string) *RabbitLoginExtension {
	return &RabbitLoginExtension{RabbitDemoExtensionSupport: NewRabbitDemoExtensionSupport(Name)}
}

type RabbitLoginExtension struct {
	RabbitDemoExtensionSupport
}

func (e *RabbitLoginExtension) InitExtension() error {
	e.GetLogger().Debugln("LoginExtension.InitExtension", e.Name)
	e.SetRequestHandlerString(LoginId, e.onRequestLogin)
	e.SetRequestHandlerString(ReLoginId, e.onRequestReLogin)
	return nil
}

func (e *RabbitLoginExtension) DestroyExtension() error {
	e.ClearRequestHandler(ReLoginId)
	e.ClearRequestHandler(LoginId)
	e.GetLogger().Debugln("[RabbitLoginExtension.DestroyExtension]", e.Name)
	return nil
}

func (e *RabbitLoginExtension) onRequestLogin(resp server.IExtensionResponse, req server.IStringRequest) {
	password := string(req.StringData()[0])
	if e.check(req.ClientId(), password) {
		core.RabbitAddressProxy.MapIdAddress(req.ClientId(), req.ClientAddress())
		time.Sleep(time.Millisecond * 20)
		resp.SendStringResponse("ok", "200")
		e.GetLogger().Traceln("[RabbitLoginExtension.onRequestLogin]", "Check Suc!", req.ProtoId(), req.ClientId(), password)
	} else {
		e.GetLogger().Warnln("[RabbitLoginExtension.onRequestLogin]", "Check Fail!", req.ProtoId(), req.ClientId(), password)
	}
}

func (e *RabbitLoginExtension) onRequestReLogin(resp server.IExtensionResponse, req server.IStringRequest) {
	password := req.StringData()[0]
	if e.check(req.ClientId(), password) {
		core.RabbitAddressProxy.MapIdAddress(req.ClientId(), req.ClientAddress())
		time.Sleep(time.Millisecond * 20)
		resp.SendStringResponse("ok")
		e.GetLogger().Traceln("[RabbitLoginExtension.onRequestReLogin]", "Check Suc!", req.ProtoId(), req.ClientId(), password)
	} else {
		e.GetLogger().Warnln("[RabbitLoginExtension.onRequestReLogin]", "Check Fail!", req.ProtoId(), req.ClientId(), password)
	}
}

func (e *RabbitLoginExtension) check(uid string, password string) bool {
	return uid == password
}
