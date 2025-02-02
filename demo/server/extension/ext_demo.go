// Package extension
// Created by xuzhuoxi
// on 2019-02-19.
// @author xuzhuoxi
//
package extension

import (
	"github.com/xuzhuoxi/Rabbit-Server/engine/server"
)

func NewRabbitDemoExtension(Name string) server.IProtoExtension {
	return &RabbitDemoExtension{RabbitDemoExtensionSupport: NewRabbitDemoExtensionSupport(Name)}
}

func newOriginObj() interface{} {
	return &originObj{}
}

type originObj struct {
	Data  int
	Data2 string
}

type paramCodingHandler struct {
}

func (c *paramCodingHandler) HandleEncode(data interface{}) (bs []byte, err error) {
	return
}

func (c *paramCodingHandler) HandleDecode(bs []byte, data interface{}) error {
	return nil
}

// RabbitDemoExtension
// Extension至少实现两个接口
// IProtocolExtension(必须)
// IOnNoneRequestExtension、IOnBinaryRequestExtension、IOnObjectRequestExtension(选一)
// IGoroutineExtension、IBatchExtension、IBeforeRequestExtension、IAfterRequestExtension(可选)
type RabbitDemoExtension struct {
	RabbitDemoExtensionSupport
}

func (e *RabbitDemoExtension) InitExtension() error {
	e.GetLogger().Debugln("DemoExtension.InitExtension", e.Name)
	e.SetRequestHandler("N_0", e.onRequestNoneParam)
	e.SetRequestHandlerBinary("B_0", e.onRequestBinary)
	e.SetRequestHandlerString("J_0", e.onRequestJson)
	e.SetRequestHandlerObject("Obj_0", e.onRequestObj, newOriginObj, &paramCodingHandler{})
	return nil
}

func (e *RabbitDemoExtension) DestroyExtension() error {
	e.ClearRequestHandler("J_0")
	e.ClearRequestHandler("B_0")
	e.ClearRequestHandler("N_0")
	e.GetLogger().Debugln("DemoExtension.DestroyExtension", e.Name)
	return nil
}

func (e *RabbitDemoExtension) BeforeRequest(protoId string) {
	e.GetLogger().Debugln("DemoExtension.BeforeRequest!", protoId)
}

func (e *RabbitDemoExtension) onRequestNoneParam(resp server.IExtensionResponse, req server.IExtensionRequest) {
	e.GetLogger().Debugln("DemoExtension.onRequestNoneParam!", req, resp)
}

func (e *RabbitDemoExtension) onRequestBinary(resp server.IExtensionResponse, req server.IBinaryRequest) {
	e.GetLogger().Debugln("DemoExtension.onRequestBinary!", req, resp)
}

func (e *RabbitDemoExtension) onRequestJson(resp server.IExtensionResponse, req server.IStringRequest) {
	e.GetLogger().Debugln("DemoExtension.onRequestJson!", req, resp)
}

func (e *RabbitDemoExtension) onRequestObj(resp server.IExtensionResponse, req server.IObjectRequest) {
	e.GetLogger().Debugln("DemoExtension.onRequestObj!", req, resp)
}

func (e *RabbitDemoExtension) AfterRequest(protoId string) {
	e.GetLogger().Debugln("DemoExtension.AfterRequest!", protoId)
}

func (e *RabbitDemoExtension) SaveExtension() error {
	e.GetLogger().Debugln("DemoExtension.SaveExtension", e.Name)
	return nil
}
