// Package protox
// Created by xuzhuoxi
// on 2019-05-19.
// @author xuzhuoxi
//
package extension

import (
	"github.com/xuzhuoxi/Rabbit-Server/engine/server"
	"github.com/xuzhuoxi/infra-go/encodingx"
)

func NewProtoObjectParamsHandler(factory server.FuncParamCtor, handler encodingx.ICodingHandler) server.IProtoParamsHandler {
	rs := &ProtoObjectParamsHandler{ParamFactory: factory}
	rs.SetCodingHandler(handler)
	return rs
}

type ProtoObjectParamsHandler struct {
	ParamFactory  server.FuncParamCtor
	ReqHandler    encodingx.ICodingHandler
	ReturnHandler encodingx.ICodingHandler
}

func (o *ProtoObjectParamsHandler) SetCodingHandler(handler encodingx.ICodingHandler) {
	o.ReqHandler, o.ReturnHandler = handler, handler
}

func (o *ProtoObjectParamsHandler) SetCodingHandlers(reqHandler encodingx.ICodingHandler, returnHandler encodingx.ICodingHandler) {
	o.ReqHandler, o.ReturnHandler = reqHandler, returnHandler
}

func (o *ProtoObjectParamsHandler) HandleRequestParam(data []byte) interface{} {
	rs := o.ParamFactory()
	err := o.ReqHandler.HandleDecode(data, rs)
	if nil != err {
		return nil
	}
	return rs
}

func (o *ProtoObjectParamsHandler) HandleRequestParams(data [][]byte) []interface{} {
	var objData []interface{}
	for index := range data {
		objData = append(objData, o.HandleRequestParam(data[index]))
	}
	return objData
}

func (o *ProtoObjectParamsHandler) HandleReturnParam(data interface{}) []byte {
	bs, err := o.ReturnHandler.HandleEncode(data)
	if nil != err {
		return nil
	}
	return bs
}

func (o *ProtoObjectParamsHandler) HandleReturnParams(data []interface{}) [][]byte {
	var byteData [][]byte
	for index := range data {
		byteData = append(byteData, o.HandleReturnParam(data[index]))
	}
	return byteData
}
