// Package extension
// Created by xuzhuoxi
// on 2019-02-26.
// @author xuzhuoxi
//
package extension

import (
	"github.com/xuzhuoxi/Rabbit-Server/engine/server"
	"github.com/xuzhuoxi/Rabbit-Server/engine/server/packet"
	"github.com/xuzhuoxi/infra-go/encodingx"
)

// RequestMetaInfo
// 协议定义
type RequestMetaInfo struct {
	ProtoId string

	ParamType       server.ExtensionParamType
	ParamHandler    server.IPacketParamsHandler
	ReqParamFactory server.FuncParamObjectCtor

	OnBeforeHandler  server.FuncBeforeRequest
	OnAfterHandler   server.FuncAfterRequest
	OnRequestHandler interface{}
}

//---------------------------------------

func NewOnRequestSupport(Name string) OnRequestSupport {
	return OnRequestSupport{
		Name: Name, ProtoIdToInfo: make(map[string]*RequestMetaInfo),
	}
}

type OnRequestSupport struct {
	Name          string
	ProtoIdToInfo map[string]*RequestMetaInfo
}

func (s *OnRequestSupport) ExtensionName() string {
	return s.Name
}
func (s *OnRequestSupport) CheckProtoId(protoId string) bool {
	_, ok := s.ProtoIdToInfo[protoId]
	return ok
}
func (s *OnRequestSupport) GetParamInfo(protoId string) (paramType server.ExtensionParamType, handler server.IPacketParamsHandler) {
	info, _ := s.ProtoIdToInfo[protoId]
	return info.ParamType, info.ParamHandler
}

func (s *OnRequestSupport) SetBeforeRequestHandler(protoId string, handler server.FuncBeforeRequest) {
	if info, ok := s.ProtoIdToInfo[protoId]; ok {
		info.OnBeforeHandler = handler
		return
	}
	s.ProtoIdToInfo[protoId] = &RequestMetaInfo{ProtoId: protoId, OnBeforeHandler: handler}
}

func (s *OnRequestSupport) SetAfterRequestHandler(protoId string, handler server.FuncAfterRequest) {
	if info, ok := s.ProtoIdToInfo[protoId]; ok {
		info.OnAfterHandler = handler
		return
	}
	s.ProtoIdToInfo[protoId] = &RequestMetaInfo{ProtoId: protoId, OnAfterHandler: handler}
}

func (s *OnRequestSupport) SetOnRequestHandler(protoId string, handler server.FuncOnNoneParamRequest) {
	if info, ok := s.ProtoIdToInfo[protoId]; ok {
		info.ParamType, info.OnRequestHandler = server.None, handler
		return
	}
	s.ProtoIdToInfo[protoId] = &RequestMetaInfo{ProtoId: protoId, ParamType: server.None, OnRequestHandler: handler}
}
func (s *OnRequestSupport) SetOnBinaryRequestHandler(protoId string, handler server.FuncOnBinaryRequest) {
	if info, ok := s.ProtoIdToInfo[protoId]; ok {
		info.ParamType, info.OnRequestHandler = server.Binary, handler
		return
	}
	s.ProtoIdToInfo[protoId] = &RequestMetaInfo{ProtoId: protoId, ParamType: server.Binary, OnRequestHandler: handler}
}
func (s *OnRequestSupport) SetOnStringRequestHandler(protoId string, handler server.FuncOnStringRequest) {
	if info, ok := s.ProtoIdToInfo[protoId]; ok {
		info.ParamType, info.OnRequestHandler = server.String, handler
		return
	}
	s.ProtoIdToInfo[protoId] = &RequestMetaInfo{ProtoId: protoId, ParamType: server.String, OnRequestHandler: handler}
}

func (s *OnRequestSupport) SetOnObjectRequestHandler(protoId string, handler server.FuncOnObjectRequest,
	factory server.FuncParamObjectCtor, codingHandler encodingx.ICodingHandler) {
	if info, ok := s.ProtoIdToInfo[protoId]; ok {
		info.ParamType, info.OnRequestHandler = server.Object, handler
		info.ReqParamFactory, info.ParamHandler = factory, packet.NewPacketParamsHandler(factory, codingHandler)
		return
	}
	s.ProtoIdToInfo[protoId] = &RequestMetaInfo{ProtoId: protoId, ParamType: server.Object, OnRequestHandler: handler,
		ReqParamFactory: factory, ParamHandler: packet.NewPacketParamsHandler(factory, codingHandler)}
}
func (s *OnRequestSupport) ClearRequestHandler(protoId string) {
	delete(s.ProtoIdToInfo, protoId)
}

func (s *OnRequestSupport) ClearRequestHandlers() {
	s.ProtoIdToInfo = make(map[string]*RequestMetaInfo)
}

func (s *OnRequestSupport) OnRequest(resp server.IExtensionResponse, req server.IExtensionRequest) {
	info, _ := s.ProtoIdToInfo[req.ProtoId()]
	switch info.ParamType {
	case server.None:
		handler := info.OnRequestHandler.(server.FuncOnNoneParamRequest)
		handler(resp.(server.IExtensionResponse), req.(server.IExtensionRequest))
	case server.Binary:
		handler := info.OnRequestHandler.(server.FuncOnBinaryRequest)
		handler(resp.(server.IExtensionResponse), req.(server.IBinaryRequest))
	case server.String:
		handler := info.OnRequestHandler.(server.FuncOnStringRequest)
		handler(resp.(server.IExtensionResponse), req.(server.IStringRequest))
	case server.Object:
		handler := info.OnRequestHandler.(server.FuncOnObjectRequest)
		handler(resp.(server.IExtensionResponse), req.(server.IObjectRequest))
	}
}

//---------------------------------------

func NewGoroutineExtensionSupport(MaxGo int) GoroutineOnRequestSupport {
	return GoroutineOnRequestSupport{MaxGoroutine: MaxGo}
}

type GoroutineOnRequestSupport struct {
	MaxGoroutine int
}

func (s *GoroutineOnRequestSupport) MaxGo() int {
	return s.MaxGoroutine
}
