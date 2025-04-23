// Package delete
// Created by xuzhuoxi
// on 2019-05-21.
// @author xuzhuoxi
//
package delete

//// FuncStartOnPack
//// 响应入口
//type FuncStartOnPack = func(connInfo netx.IConnInfo)
//
//// FuncParsePacket
//// 解释二进制数据
//type FuncParsePacket = func(msgBytes []byte) (extName string, pid string, uid string, data [][]byte)
//
//// FuncGetExtension
//// 消息处理入口，这里是并发方法
//type FuncGetExtension = func(extName string) (extension IRabbitExtension, rsCode int32)
//
//// FuncStartOnRequest
//// 响应开始
//type FuncStartOnRequest = func(resp IExtensionResponse, req IExtensionRequest)
//
//// FuncFinishOnRequest
//// 响应完成
//type FuncFinishOnRequest = func(resp IExtensionResponse, req IExtensionRequest)
//
//type ICustomManagerSetting interface {
//	// SetCustomStartOnPackFunc
//	// 设置自定义响应开始行为
//	SetCustomStartOnPackFunc(funcStartOnPack FuncStartOnPack)
//	// SetCustomParsePacketFunc
//	// 设置自定义数据解释行为
//	SetCustomParsePacketFunc(funcParse FuncParsePacket)
//	// SetCustomGetExtensionFunc
//	// 设置自定义扩展获取
//	SetCustomGetExtensionFunc(funcVerify FuncGetExtension)
//	// SetCustomVerifyFunc
//	// 设置自定义验证
//	SetCustomVerifyFunc(funcVerify FuncVerifyPacket)
//	// SetCustomPacketVerifier
//	// 设置自定义的消息包校验器
//	SetCustomPacketVerifier(reqVerify IPacketVerifier)
//	// SetCustomStartOnRequestFunc
//	// 设置自定义响应前置行为
//	SetCustomStartOnRequestFunc(funcStart FuncStartOnRequest)
//	// SetCustomFinishOnRequestFunc
//	// 设置自定义响应完成行为
//	SetCustomFinishOnRequestFunc(funcFinish FuncFinishOnRequest)
//	// SetCustom
//	// 设置自定义行为
//	SetCustom(funcStartOnPack FuncStartOnPack, funcParse FuncParsePacket, funcVerify FuncVerifyPacket, funcStart FuncStartOnRequest, funcFinish FuncFinishOnRequest)
//}
//
//type ICustomManagerSupport interface {
//	CustomStartOnPack(connInfo netx.IConnInfo)
//	CustomParsePacket(msgBytes []byte) (extName string, pid string, uid string, data [][]byte)
//	CustomGetExtension(extName string) (extension IRabbitExtension, rsCode int32)
//	CustomVerify(extName string, pid string, uid string) (rsCode int32)
//	CustomStartOnRequest(resp IExtensionResponse, req IExtensionRequest)
//	CustomFinishOnRequest(resp IExtensionResponse, req IExtensionRequest)
//}
//
//
//type CustomManagerSupport struct {
//	FuncStartOnPack     server.FuncStartOnPack
//	FuncParsePacket     server.FuncParsePacket
//	FuncGetExtension    server.FuncGetExtension
//	FuncVerifyPacket    server.FuncVerifyPacket
//	PacketVerifier      server.IPacketVerifier
//	FuncStartOnRequest  server.FuncStartOnRequest
//	FuncFinishOnRequest server.FuncFinishOnRequest
//
//	SupportMutex sync.RWMutex
//}
//
//func (o *CustomManagerSupport) SetCustomStartOnPackFunc(funcStartOnPack server.FuncStartOnPack) {
//	o.SupportMutex.Lock()
//	defer o.SupportMutex.Unlock()
//	o.FuncStartOnPack = funcStartOnPack
//}
//func (o *CustomManagerSupport) SetCustomParsePacketFunc(funcParse server.FuncParsePacket) {
//	o.SupportMutex.Lock()
//	defer o.SupportMutex.Unlock()
//	o.FuncParsePacket = funcParse
//}
//func (o *CustomManagerSupport) SetCustomGetExtensionFunc(funcGet server.FuncGetExtension) {
//	o.SupportMutex.Lock()
//	defer o.SupportMutex.Unlock()
//	o.FuncGetExtension = funcGet
//}
//func (o *CustomManagerSupport) SetCustomVerifyFunc(funcVerify server.FuncVerifyPacket) {
//	o.SupportMutex.Lock()
//	defer o.SupportMutex.Unlock()
//	o.FuncVerifyPacket = funcVerify
//}
//func (o *CustomManagerSupport) SetCustomPacketVerifier(reqVerify server.IPacketVerifier) {
//	o.SupportMutex.Lock()
//	defer o.SupportMutex.Unlock()
//	o.PacketVerifier = reqVerify
//}
//func (o *CustomManagerSupport) SetCustomStartOnRequestFunc(funcStart server.FuncStartOnRequest) {
//	o.SupportMutex.Lock()
//	defer o.SupportMutex.Unlock()
//	o.FuncStartOnRequest = funcStart
//}
//func (o *CustomManagerSupport) SetCustomFinishOnRequestFunc(funcFinish server.FuncFinishOnRequest) {
//	o.SupportMutex.Lock()
//	defer o.SupportMutex.Unlock()
//	o.FuncFinishOnRequest = funcFinish
//}
//func (o *CustomManagerSupport) SetCustom(funcStartOnPack server.FuncStartOnPack, funcParse server.FuncParsePacket, funcVerify server.FuncVerifyPacket, funcStart server.FuncStartOnRequest, funcFinish server.FuncFinishOnRequest) {
//	o.SupportMutex.Lock()
//	defer o.SupportMutex.Unlock()
//	o.FuncStartOnPack, o.FuncParsePacket, o.FuncVerifyPacket, o.FuncStartOnRequest, o.FuncFinishOnRequest = funcStartOnPack, funcParse, funcVerify, funcStart, funcFinish
//}
//
//// Custom
//
//func (o *CustomManagerSupport) CustomStartOnPack(connInfo netx.IConnInfo) {
//	if nil != o.FuncStartOnPack {
//		o.FuncStartOnPack(connInfo)
//	}
//}
//func (o *CustomManagerSupport) CustomParsePacket(msgBytes []byte) (name string, pid string, uid string, data [][]byte) {
//	if nil != o.FuncParsePacket {
//		return o.FuncParsePacket(msgBytes)
//	}
//	return
//}
//func (o *CustomManagerSupport) CustomGetExtension(name string) (extension server.IRabbitExtension, rsCode int32) {
//	if nil != o.FuncGetExtension {
//		return o.FuncGetExtension(name)
//	}
//	return nil, server.CodeProtoFail
//}
//func (o *CustomManagerSupport) CustomVerify(name string, pid string, uid string) (rsCode int32) {
//	if nil != o.FuncVerifyPacket {
//		return o.FuncVerifyPacket(name, pid, uid)
//	}
//	if nil != o.PacketVerifier {
//		return o.PacketVerifier.VerifyExtension(name, pid, uid)
//	}
//	return
//}
//func (o *CustomManagerSupport) CustomStartOnRequest(resp server.IExtensionResponse, req server.IExtensionRequest) {
//	if nil != o.FuncStartOnRequest {
//		o.FuncStartOnRequest(resp, req)
//	}
//}
//func (o *CustomManagerSupport) CustomFinishOnRequest(resp server.IExtensionResponse, req server.IExtensionRequest) {
//	if nil != o.FuncFinishOnRequest {
//		o.FuncFinishOnRequest(resp, req)
//	}
//}
