// Package protox
// Created by xuzhuoxi
// on 2019-05-21.
// @author xuzhuoxi
//
package extension

import (
	"github.com/xuzhuoxi/Rabbit-Server/engine/server"
)

type ExtensionManagerCustomizeSupport struct {
	FuncStartOnPack     server.FuncStartOnPack
	FuncParseMessage    server.FuncParseMessage
	FuncGetExtension    server.FuncGetExtension
	FuncVerify          server.FuncVerify
	ReqVerify           server.IReqVerify
	FuncStartOnRequest  server.FuncStartOnRequest
	FuncFinishOnRequest server.FuncFinishOnRequest
}

func (o *ExtensionManagerCustomizeSupport) SetCustomStartOnPackFunc(funcStartOnPack server.FuncStartOnPack) {
	o.FuncStartOnPack = funcStartOnPack
}
func (o *ExtensionManagerCustomizeSupport) SetCustomParseFunc(funcParse server.FuncParseMessage) {
	o.FuncParseMessage = funcParse
}
func (o *ExtensionManagerCustomizeSupport) SetCustomGetExtensionFunc(funcGet server.FuncGetExtension) {
	o.FuncGetExtension = funcGet
}
func (o *ExtensionManagerCustomizeSupport) SetCustomVerifyFunc(funcVerify server.FuncVerify) {
	o.FuncVerify = funcVerify
}
func (o *ExtensionManagerCustomizeSupport) SetCustomVerify(reqVerify server.IReqVerify) {
	o.ReqVerify = reqVerify
}
func (o *ExtensionManagerCustomizeSupport) SetCustomStartOnRequestFunc(funcStart server.FuncStartOnRequest) {
	o.FuncStartOnRequest = funcStart
}
func (o *ExtensionManagerCustomizeSupport) SetCustomFinishOnRequestFunc(funcFinish server.FuncFinishOnRequest) {
	o.FuncFinishOnRequest = funcFinish
}
func (o *ExtensionManagerCustomizeSupport) SetCustom(funcStartOnPack server.FuncStartOnPack, funcParse server.FuncParseMessage, funcVerify server.FuncVerify, funcStart server.FuncStartOnRequest, funcFinish server.FuncFinishOnRequest) {
	o.FuncStartOnPack, o.FuncParseMessage, o.FuncVerify, o.FuncStartOnRequest, o.FuncFinishOnRequest = funcStartOnPack, funcParse, funcVerify, funcStart, funcFinish
}

func (o *ExtensionManagerCustomizeSupport) CustomStartOnPack(senderAddress string) {
	if nil != o.FuncStartOnPack {
		o.FuncStartOnPack(senderAddress)
	}
}
func (o *ExtensionManagerCustomizeSupport) CustomParseMessage(msgBytes []byte) (name string, pid string, uid string, data [][]byte) {
	if nil != o.FuncParseMessage {
		return o.FuncParseMessage(msgBytes)
	}
	return
}
func (o *ExtensionManagerCustomizeSupport) CustomGetExtension(name string) (extension server.IProtoExtension, rsCode int32) {
	if nil != o.FuncGetExtension {
		return o.FuncGetExtension(name)
	}
	return nil, server.CodeProtoFail
}
func (o *ExtensionManagerCustomizeSupport) CustomVerify(name string, pid string, uid string) (rsCode int32) {
	if nil != o.FuncVerify {
		return o.FuncVerify(name, pid, uid)
	}
	if nil != o.ReqVerify {
		return o.ReqVerify.Verify(name, pid, uid)
	}
	return
}
func (o *ExtensionManagerCustomizeSupport) CustomStartOnRequest(resp server.IExtensionResponse, req server.IExtensionRequest) {
	if nil != o.FuncStartOnRequest {
		o.FuncStartOnRequest(resp, req)
	}
}
func (o *ExtensionManagerCustomizeSupport) CustomFinishOnRequest(resp server.IExtensionResponse, req server.IExtensionRequest) {
	if nil != o.FuncFinishOnRequest {
		o.FuncFinishOnRequest(resp, req)
	}
}
