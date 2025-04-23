// Package runtime
// Created by xuzhuoxi
// on 2019-02-26.
// @author xuzhuoxi
//
package runtime

import (
	"github.com/xuzhuoxi/Rabbit-Server/engine/server"
	"github.com/xuzhuoxi/Rabbit-Server/engine/server/runtime/params/packets"
	"github.com/xuzhuoxi/infra-go/encodingx"
)

// ExtensionProtoMeta
// 扩展协议元数据定义
type ExtensionProtoMeta struct {
	ProtoId string

	ParamType       server.ExtensionParamType
	ParamHandler    server.IPacketCoding
	ReqParamFactory server.FuncParamObjectCtor

	OnBeforeHandler  server.FuncBeforeRequest
	OnAfterHandler   server.FuncAfterRequest
	OnRequestHandler interface{}
}

//---------------------------------------

func NewExtensionMeta(Name string) ExtensionMeta {
	return ExtensionMeta{
		Name: Name, ProtoIdToInfo: make(map[string]*ExtensionProtoMeta),
	}
}

// ExtensionMeta
// 扩展元数据定义
type ExtensionMeta struct {
	Name          string
	ProtoIdToInfo map[string]*ExtensionProtoMeta
}

func (s *ExtensionMeta) ExtensionName() string {
	return s.Name
}
func (s *ExtensionMeta) CheckProtoId(protoId string) bool {
	_, ok := s.ProtoIdToInfo[protoId]
	return ok
}
func (s *ExtensionMeta) GetParamInfo(protoId string) (paramType server.ExtensionParamType, handler server.IPacketCoding) {
	info, _ := s.ProtoIdToInfo[protoId]
	return info.ParamType, info.ParamHandler
}

func (s *ExtensionMeta) SetBeforeRequestHandler(protoId string, handler server.FuncBeforeRequest) {
	if info, ok := s.ProtoIdToInfo[protoId]; ok {
		info.OnBeforeHandler = handler
		return
	}
	s.ProtoIdToInfo[protoId] = &ExtensionProtoMeta{ProtoId: protoId, OnBeforeHandler: handler}
}

func (s *ExtensionMeta) SetAfterRequestHandler(protoId string, handler server.FuncAfterRequest) {
	if info, ok := s.ProtoIdToInfo[protoId]; ok {
		info.OnAfterHandler = handler
		return
	}
	s.ProtoIdToInfo[protoId] = &ExtensionProtoMeta{ProtoId: protoId, OnAfterHandler: handler}
}

func (s *ExtensionMeta) SetOnRequestHandler(protoId string, handler server.FuncOnNoneParamRequest) {
	if info, ok := s.ProtoIdToInfo[protoId]; ok {
		info.ParamType, info.OnRequestHandler = server.None, handler
		return
	}
	s.ProtoIdToInfo[protoId] = &ExtensionProtoMeta{ProtoId: protoId, ParamType: server.None, OnRequestHandler: handler}
}
func (s *ExtensionMeta) SetOnBinaryRequestHandler(protoId string, handler server.FuncOnBinaryRequest) {
	if info, ok := s.ProtoIdToInfo[protoId]; ok {
		info.ParamType, info.OnRequestHandler = server.Binary, handler
		return
	}
	s.ProtoIdToInfo[protoId] = &ExtensionProtoMeta{ProtoId: protoId, ParamType: server.Binary, OnRequestHandler: handler}
}
func (s *ExtensionMeta) SetOnStringRequestHandler(protoId string, handler server.FuncOnStringRequest) {
	if info, ok := s.ProtoIdToInfo[protoId]; ok {
		info.ParamType, info.OnRequestHandler = server.String, handler
		return
	}
	s.ProtoIdToInfo[protoId] = &ExtensionProtoMeta{ProtoId: protoId, ParamType: server.String, OnRequestHandler: handler}
}

func (s *ExtensionMeta) SetOnObjectRequestHandler(protoId string, handler server.FuncOnObjectRequest,
	factory server.FuncParamObjectCtor, codingHandler encodingx.ICodingHandler) {
	if info, ok := s.ProtoIdToInfo[protoId]; ok {
		info.ParamType, info.OnRequestHandler = server.Object, handler
		info.ReqParamFactory, info.ParamHandler = factory, packets.NewIPacketCoding(factory, codingHandler)
		return
	}
	s.ProtoIdToInfo[protoId] = &ExtensionProtoMeta{ProtoId: protoId, ParamType: server.Object, OnRequestHandler: handler,
		ReqParamFactory: factory, ParamHandler: packets.NewIPacketCoding(factory, codingHandler)}
}
func (s *ExtensionMeta) ClearRequestHandler(protoId string) {
	delete(s.ProtoIdToInfo, protoId)
}

func (s *ExtensionMeta) ClearRequestHandlers() {
	s.ProtoIdToInfo = make(map[string]*ExtensionProtoMeta)
}

func (s *ExtensionMeta) OnRequest(resp server.IExtensionResponse, req server.IExtensionRequest) {
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

func NewGoroutineExtensionMeta(MaxGo int) GoroutineExtensionMeta {
	return GoroutineExtensionMeta{MaxGoroutine: MaxGo}
}

type GoroutineExtensionMeta struct {
	MaxGoroutine int
}

func (s *GoroutineExtensionMeta) MaxGo() int {
	return s.MaxGoroutine
}
