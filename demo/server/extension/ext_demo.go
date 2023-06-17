// Package server
// Created by xuzhuoxi
// on 2019-02-19.
// @author xuzhuoxi
//
package extension

import (
	"github.com/xuzhuoxi/Rabbit-Server/engine/server"
	"github.com/xuzhuoxi/infra-go/encodingx"
	"github.com/xuzhuoxi/infra-go/extendx/protox"
)

func NewRabbitDemoExtension(Name string) server.IRabbitExtension {
	return &RabbitDemoExtension{RabbitDemoExtensionSupport: NewRabbitDemoExtensionSupport(Name)}
}

type originObj struct {
	Data  int
	Data2 string
}

type paramHandler struct {
}

func (h *paramHandler) SetCodingHandler(handler encodingx.ICodingHandler) {
	panic("implement me")
}

func (h *paramHandler) HandleRequestParam(data []byte) interface{} {
	panic("implement me")
}

func (h *paramHandler) HandleRequestParams(data [][]byte) []interface{} {
	panic("implement me")
}

func (h *paramHandler) HandleResponseParam(data interface{}) []byte {
	panic("implement me")
}

func (h *paramHandler) HandleResponseParams(data []interface{}) [][]byte {
	panic("implement me")
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
	e.SetRequestHandlerJson("J_0", e.onRequestJson)
	e.SetRequestHandlerObject("Obj_0", e.onRequestObj, originObj{}, &paramHandler{})
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

func (e *RabbitDemoExtension) onRequestNoneParam(resp protox.IExtensionResponse, req protox.IExtensionRequest) {
	e.GetLogger().Debugln("DemoExtension.onRequestNoneParam!", req, resp)
}

func (e *RabbitDemoExtension) onRequestBinary(resp protox.IExtensionBinaryResponse, req protox.IExtensionBinaryRequest) {
	e.GetLogger().Debugln("DemoExtension.onRequestBinary!", req, resp)
}

func (e *RabbitDemoExtension) onRequestJson(resp protox.IExtensionJsonResponse, req protox.IExtensionJsonRequest) {
	e.GetLogger().Debugln("DemoExtension.onRequestJson!", req, resp)
}

func (e *RabbitDemoExtension) onRequestObj(resp protox.IExtensionObjectResponse, req protox.IExtensionObjectRequest) {
	e.GetLogger().Debugln("DemoExtension.onRequestObj!", req, resp)
}

func (e *RabbitDemoExtension) AfterRequest(protoId string) {
	e.GetLogger().Debugln("DemoExtension.AfterRequest!", protoId)
}

func (e *RabbitDemoExtension) SaveExtension() error {
	e.GetLogger().Debugln("DemoExtension.SaveExtension", e.Name)
	return nil
}
