// Package packets
// Created by xuzhuoxi
// on 2019-05-19.
// @author xuzhuoxi
//
package packets

import (
	"github.com/xuzhuoxi/Rabbit-Server/engine/server"
	"github.com/xuzhuoxi/infra-go/encodingx"
)

func NewIPacketCoding(factory server.FuncParamObjectCtor, handler encodingx.ICodingHandler) server.IPacketCoding {
	rs := &PacketCoding{ParamFactory: factory}
	rs.SetCodingHandler(handler)
	return rs
}

type PacketCoding struct {
	ParamFactory server.FuncParamObjectCtor
	ReqHandler   encodingx.ICodingHandler
	RespHandler  encodingx.ICodingHandler
}

func (o *PacketCoding) SetCodingHandler(handler encodingx.ICodingHandler) {
	o.ReqHandler, o.RespHandler = handler, handler
}

func (o *PacketCoding) SetCodingHandlers(reqHandler encodingx.ICodingHandler, respHandler encodingx.ICodingHandler) {
	o.ReqHandler, o.RespHandler = reqHandler, respHandler
}

func (o *PacketCoding) DecodeRequestParam(data []byte) interface{} {
	rs := o.ParamFactory()
	err := o.ReqHandler.HandleDecode(data, rs)
	if nil != err {
		return nil
	}
	return rs
}

func (o *PacketCoding) DecodeRequestParams(data [][]byte) []interface{} {
	var objData []interface{}
	for index := range data {
		objData = append(objData, o.DecodeRequestParam(data[index]))
	}
	return objData
}

func (o *PacketCoding) EncodeResponseParam(data interface{}) []byte {
	bs, err := o.RespHandler.HandleEncode(data)
	if nil != err {
		return nil
	}
	return bs
}

func (o *PacketCoding) EncodeResponseParams(data []interface{}) [][]byte {
	var byteData [][]byte
	for index := range data {
		byteData = append(byteData, o.EncodeResponseParam(data[index]))
	}
	return byteData
}
