// Package protox
// Created by xuzhuoxi
// on 2019-02-26.
// @author xuzhuoxi
//
package extension

import (
	"github.com/xuzhuoxi/Rabbit-Server/engine/server"
	"github.com/xuzhuoxi/infra-go/encodingx"
)

// ExtensionProtoInfo
// 协议定义
type ExtensionProtoInfo struct {
	ProtoId          string
	ParamType        server.ExtensionParamType
	ExtensionHandler interface{}

	ParamHandler    server.IProtoParamsHandler
	ReqParamFactory server.FuncParamCtor
}

//---------------------------------------

func NewProtoExtensionSupport(Name string) ProtoExtensionSupport {
	return ProtoExtensionSupport{
		Name: Name, ProtoIdToInfo: make(map[string]*ExtensionProtoInfo),
	}
}

type ProtoExtensionSupport struct {
	Name          string
	ProtoIdToInfo map[string]*ExtensionProtoInfo
}

func (s *ProtoExtensionSupport) ExtensionName() string {
	return s.Name
}
func (s *ProtoExtensionSupport) CheckProtocolId(protoId string) bool {
	_, ok := s.ProtoIdToInfo[protoId]
	return ok
}
func (s *ProtoExtensionSupport) GetParamInfo(protoId string) (paramType server.ExtensionParamType, handler server.IProtoParamsHandler) {
	info, _ := s.ProtoIdToInfo[protoId]
	return info.ParamType, info.ParamHandler
}

func (s *ProtoExtensionSupport) SetRequestHandler(protoId string, handler server.ExtensionHandlerNoneParam) {
	s.ProtoIdToInfo[protoId] = &ExtensionProtoInfo{ProtoId: protoId, ParamType: server.None, ExtensionHandler: handler}
}
func (s *ProtoExtensionSupport) SetRequestHandlerBinary(protoId string, handler server.ExtensionHandlerBinaryParam) {
	s.ProtoIdToInfo[protoId] = &ExtensionProtoInfo{ProtoId: protoId, ParamType: server.Binary, ExtensionHandler: handler}
}
func (s *ProtoExtensionSupport) SetRequestHandlerString(protoId string, handler server.ExtensionHandlerStringParam) {
	s.ProtoIdToInfo[protoId] = &ExtensionProtoInfo{ProtoId: protoId, ParamType: server.String, ExtensionHandler: handler}
}
func (s *ProtoExtensionSupport) SetRequestHandlerObject(protoId string, handler server.ExtensionHandlerObjectParam,
	factory server.FuncParamCtor, codingHandler encodingx.ICodingHandler) {
	s.ProtoIdToInfo[protoId] = &ExtensionProtoInfo{ProtoId: protoId, ParamType: server.Object, ExtensionHandler: handler,
		ReqParamFactory: factory, ParamHandler: NewProtoObjectParamsHandler(factory, codingHandler)}
}
func (s *ProtoExtensionSupport) ClearRequestHandler(protoId string) {
	delete(s.ProtoIdToInfo, protoId)
}
func (s *ProtoExtensionSupport) ClearRequestHandlers() {
	s.ProtoIdToInfo = make(map[string]*ExtensionProtoInfo)
}

func (s *ProtoExtensionSupport) OnRequest(resp server.IExtensionResponse, req server.IExtensionRequest) {
	info, _ := s.ProtoIdToInfo[req.ProtoId()]
	switch info.ParamType {
	case server.None:
		handler := info.ExtensionHandler.(server.ExtensionHandlerNoneParam)
		handler(resp.(server.IExtensionResponse), req.(server.IExtensionRequest))
	case server.Binary:
		handler := info.ExtensionHandler.(server.ExtensionHandlerBinaryParam)
		handler(resp.(server.IExtensionResponse), req.(server.IBinaryRequest))
	case server.String:
		handler := info.ExtensionHandler.(server.ExtensionHandlerStringParam)
		handler(resp.(server.IExtensionResponse), req.(server.IStringRequest))
	case server.Object:
		handler := info.ExtensionHandler.(server.ExtensionHandlerObjectParam)
		handler(resp.(server.IExtensionResponse), req.(server.IObjectRequest))
	}
}

//---------------------------------------

func NewGoroutineExtensionSupport(MaxGo int) GoroutineExtensionSupport {
	return GoroutineExtensionSupport{MaxGoroutine: MaxGo}
}

type GoroutineExtensionSupport struct {
	MaxGoroutine int
}

func (s *GoroutineExtensionSupport) MaxGo() int {
	return s.MaxGoroutine
}
