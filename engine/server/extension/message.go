// Package protox
// Create on 2023/8/6
// @author xuzhuoxi
package extension

import (
	"errors"
	"github.com/xuzhuoxi/Rabbit-Server/engine/server"
	"github.com/xuzhuoxi/Rabbit-Server/engine/server/proto"
	"github.com/xuzhuoxi/infra-go/binaryx"
	"github.com/xuzhuoxi/infra-go/bytex"
)

func NewProtoReturnMessage() *ProtoReturnMessage {
	return &ProtoReturnMessage{
		DataBuff: bytex.NewDefaultBuffToBlock(),
		MsgBuff:  bytex.NewDefaultBuffToBlock(),
	}
}

type ProtoReturnMessage struct {
	proto.ProtoHeader
	RsCode       int32
	ParamHandler server.IProtoParamsHandler
	MsgBuff      bytex.IBuffToBlock
	DataBuff     bytex.IBuffToBlock
	dataBytes    []byte
}

func (o *ProtoReturnMessage) PrepareData() {
	o.dataBytes = nil
	o.DataBuff.Reset()
}

func (o *ProtoReturnMessage) AppendLen(ln int) error {
	order := o.DataBuff.GetOrder()
	return binaryx.Write(o.DataBuff, order, uint16(ln))
}

func (o *ProtoReturnMessage) AppendBinary(data ...[]byte) error {
	if len(data) == 0 {
		return nil
	}
	for index := range data {
		o.DataBuff.WriteData(data[index])
	}
	return nil
}

func (o *ProtoReturnMessage) AppendCommon(data ...interface{}) error {
	if len(data) == 0 {
		return nil
	}
	order := o.DataBuff.GetOrder()
	for index := range data {
		err := binaryx.Write(o.DataBuff, order, data[index])
		if nil != err {
			return err
		}
	}
	return nil
}

func (o *ProtoReturnMessage) AppendString(data ...string) error {
	if len(data) == 0 {
		return nil
	}
	for index := range data {
		o.DataBuff.WriteString(data[index])
	}
	return nil
}

func (o *ProtoReturnMessage) AppendJson(data ...interface{}) error {
	if len(data) == 0 {
		return nil
	}
	for index := range data {
		jsonStr, err1 := toJson(data[index])
		if nil != err1 {
			return err1
		}
		err2 := o.AppendString(jsonStr)
		if nil != err2 {
			return err2
		}
	}
	return nil
}

func (o *ProtoReturnMessage) AppendObject(data ...interface{}) error {
	if len(data) == 0 {
		return nil
	}
	if o.ParamHandler == nil {
		return errors.New("AppendObject Error: ParamHandler is nil! ")
	}
	for index := range data {
		bs := o.ParamHandler.HandleReturnParam(data[index])
		o.DataBuff.WriteData(bs)
	}
	return nil
}

func (o *ProtoReturnMessage) GenMsgBytes() (msg []byte, err error) {
	return o.genMsgBytes(o.PGroup, o.PId)
}

func (o *ProtoReturnMessage) genMsgBytes(eName string, pId string) (bytes []byte, err error) {
	err1 := o.writeHeaderToMsg(eName, pId)
	if nil != err1 {
		return nil, err1
	}
	err2 := o.writeDataToMsg()
	if nil != err2 {
		return nil, err2
	}
	return o.MsgBuff.ReadBytes(), nil
}

func (o *ProtoReturnMessage) writeHeaderToMsg(eName string, pId string) error {
	o.MsgBuff.Reset()
	o.MsgBuff.WriteString(eName)
	o.MsgBuff.WriteString(pId)
	o.MsgBuff.WriteString(o.CId)
	return binaryx.Write(o.MsgBuff, o.MsgBuff.GetOrder(), o.RsCode)
}

func (o *ProtoReturnMessage) writeDataToMsg() error {
	if nil == o.dataBytes {
		o.dataBytes = o.DataBuff.ReadBytesCopy()
		if nil == o.dataBytes {
			o.dataBytes = []byte{}
		}
	}
	_, err1 := o.MsgBuff.Write(o.dataBytes)
	return err1
}
