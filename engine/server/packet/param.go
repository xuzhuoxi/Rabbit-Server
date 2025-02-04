// Package packet
// Created by xuzhuoxi
// on 2019-05-19.
// @author xuzhuoxi
//
package packet

import (
	"github.com/xuzhuoxi/Rabbit-Server/engine/server"
	"github.com/xuzhuoxi/infra-go/encodingx"
)

func NewPacketParamsHandler(factory server.FuncParamObjectCtor, handler encodingx.ICodingHandler) server.IPacketParamsHandler {
	rs := &PacketParamsHandler{ParamFactory: factory}
	rs.SetCodingHandler(handler)
	return rs
}

type PacketParamsHandler struct {
	ParamFactory server.FuncParamObjectCtor
	ReqHandler   encodingx.ICodingHandler
	RespHandler  encodingx.ICodingHandler
}

func (o *PacketParamsHandler) SetCodingHandler(handler encodingx.ICodingHandler) {
	o.ReqHandler, o.RespHandler = handler, handler
}

func (o *PacketParamsHandler) SetCodingHandlers(reqHandler encodingx.ICodingHandler, respHandler encodingx.ICodingHandler) {
	o.ReqHandler, o.RespHandler = reqHandler, respHandler
}

func (o *PacketParamsHandler) DecodeRequestParam(data []byte) interface{} {
	rs := o.ParamFactory()
	err := o.ReqHandler.HandleDecode(data, rs)
	if nil != err {
		return nil
	}
	return rs
}

func (o *PacketParamsHandler) DecodeRequestParams(data [][]byte) []interface{} {
	var objData []interface{}
	for index := range data {
		objData = append(objData, o.DecodeRequestParam(data[index]))
	}
	return objData
}

func (o *PacketParamsHandler) EncodeResponseParam(data interface{}) []byte {
	bs, err := o.RespHandler.HandleEncode(data)
	if nil != err {
		return nil
	}
	return bs
}

func (o *PacketParamsHandler) EncodeResponseParams(data []interface{}) [][]byte {
	var byteData [][]byte
	for index := range data {
		byteData = append(byteData, o.EncodeResponseParam(data[index]))
	}
	return byteData
}
