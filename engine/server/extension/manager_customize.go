// Package extension
// Created by xuzhuoxi
// on 2019-05-21.
// @author xuzhuoxi
//
package extension

import (
	"github.com/xuzhuoxi/Rabbit-Server/engine/server"
)

type CustomManagerSupport struct {
	FuncStartOnPack     server.FuncStartOnPack
	FuncParsePacket     server.FuncParsePacket
	FuncGetExtension    server.FuncGetExtension
	FuncVerifyPacket    server.FuncVerifyPacket
	PacketVerifier      server.IPacketVerifier
	FuncStartOnRequest  server.FuncStartOnRequest
	FuncFinishOnRequest server.FuncFinishOnRequest
}

func (o *CustomManagerSupport) SetCustomStartOnPackFunc(funcStartOnPack server.FuncStartOnPack) {
	o.FuncStartOnPack = funcStartOnPack
}
func (o *CustomManagerSupport) SetCustomParsePacketFunc(funcParse server.FuncParsePacket) {
	o.FuncParsePacket = funcParse
}
func (o *CustomManagerSupport) SetCustomGetExtensionFunc(funcGet server.FuncGetExtension) {
	o.FuncGetExtension = funcGet
}
func (o *CustomManagerSupport) SetCustomVerifyFunc(funcVerify server.FuncVerifyPacket) {
	o.FuncVerifyPacket = funcVerify
}
func (o *CustomManagerSupport) SetCustomPacketVerifier(reqVerify server.IPacketVerifier) {
	o.PacketVerifier = reqVerify
}
func (o *CustomManagerSupport) SetCustomStartOnRequestFunc(funcStart server.FuncStartOnRequest) {
	o.FuncStartOnRequest = funcStart
}
func (o *CustomManagerSupport) SetCustomFinishOnRequestFunc(funcFinish server.FuncFinishOnRequest) {
	o.FuncFinishOnRequest = funcFinish
}
func (o *CustomManagerSupport) SetCustom(funcStartOnPack server.FuncStartOnPack, funcParse server.FuncParsePacket, funcVerify server.FuncVerifyPacket, funcStart server.FuncStartOnRequest, funcFinish server.FuncFinishOnRequest) {
	o.FuncStartOnPack, o.FuncParsePacket, o.FuncVerifyPacket, o.FuncStartOnRequest, o.FuncFinishOnRequest = funcStartOnPack, funcParse, funcVerify, funcStart, funcFinish
}

// Custom

func (o *CustomManagerSupport) CustomStartOnPack(senderAddress string) {
	if nil != o.FuncStartOnPack {
		o.FuncStartOnPack(senderAddress)
	}
}
func (o *CustomManagerSupport) CustomParsePacket(msgBytes []byte) (name string, pid string, uid string, data [][]byte) {
	if nil != o.FuncParsePacket {
		return o.FuncParsePacket(msgBytes)
	}
	return
}
func (o *CustomManagerSupport) CustomGetExtension(name string) (extension server.IRabbitExtension, rsCode int32) {
	if nil != o.FuncGetExtension {
		return o.FuncGetExtension(name)
	}
	return nil, server.CodeProtoFail
}
func (o *CustomManagerSupport) CustomVerify(name string, pid string, uid string) (rsCode int32) {
	if nil != o.FuncVerifyPacket {
		return o.FuncVerifyPacket(name, pid, uid)
	}
	if nil != o.PacketVerifier {
		return o.PacketVerifier.Verify(name, pid, uid)
	}
	return
}
func (o *CustomManagerSupport) CustomStartOnRequest(resp server.IExtensionResponse, req server.IExtensionRequest) {
	if nil != o.FuncStartOnRequest {
		o.FuncStartOnRequest(resp, req)
	}
}
func (o *CustomManagerSupport) CustomFinishOnRequest(resp server.IExtensionResponse, req server.IExtensionRequest) {
	if nil != o.FuncFinishOnRequest {
		o.FuncFinishOnRequest(resp, req)
	}
}
