// Package params
// Created by xuzhuoxi
// on 2019-03-22.
// @author xuzhuoxi
//
package params

import (
	"github.com/xuzhuoxi/Rabbit-Server/engine/server"
	"github.com/xuzhuoxi/Rabbit-Server/engine/server/runtime/params/packets"
	"github.com/xuzhuoxi/infra-go/netx"
)

func NewISockResponse() server.IExtensionResponse {
	return NewSockResponse()
}

func NewSockResponse() *SockResponse {
	return &SockResponse{
		ResponsePacket: *packets.NewResponsePacket(),
	}
}

type iExtResponse interface {
	server.IExtensionResponseSettings
	server.IExtensionResponse
}

type SockResponse struct {
	packets.ResponsePacket
	SockSender     netx.ISockSender
	UserConnMapper netx.IUserConnMapper
	ConnInfo       netx.IConnInfo
	ParamType      server.ExtensionParamType
}

func (o *SockResponse) SetConnInfo(connInfo netx.IConnInfo) {
	o.ConnInfo = connInfo
}

func (o *SockResponse) GetConnInfo() netx.IConnInfo {
	return o.ConnInfo
}

func (o *SockResponse) SetUserConnMapper(mapper netx.IUserConnMapper) {
	o.UserConnMapper = mapper
}

func (o *SockResponse) SetSockSender(sockSender netx.ISockSender) {
	o.SockSender = sockSender
}

func (o *SockResponse) SetParamInfo(paramType server.ExtensionParamType, paramHandler server.IPacketCoding) {
	o.ParamType, o.ParamHandler = paramType, paramHandler
}

func (o *SockResponse) SetResultCode(rsCode int32) {
	o.RsCode = rsCode
}

func (o *SockResponse) SendResponse() error {
	return o.sendRedirectMsg(o.EName, o.PId)
}

func (o *SockResponse) SendResponseTo(interruptOnErr bool, clientIds ...string) error {
	return o.sendRedirectMsgTo(o.EName, o.PId, interruptOnErr, clientIds...)
}

func (o *SockResponse) SendNotify(eName string, notifyPId string) error {
	return o.sendRedirectMsg(eName, notifyPId)
}

func (o *SockResponse) SendNotifyTo(eName string, notifyPId string, interruptOnErr bool, clientIds ...string) error {
	return o.sendRedirectMsgTo(eName, notifyPId, interruptOnErr, clientIds...)
}

// extend

func (o *SockResponse) ResponseNone() error {
	o.PrepareData()
	return o.SendResponse()
}

func (o *SockResponse) ResponseNoneToClient(interruptOnErr bool, clientIds ...string) error {
	if len(clientIds) == 0 {
		return nil
	}
	o.PrepareData()
	return o.SendResponseTo(interruptOnErr, clientIds...)
}

func (o *SockResponse) ResponseBinary(data ...[]byte) error {
	o.PrepareData()
	o.AppendBinary(data...)
	return o.SendResponse()
}

func (o *SockResponse) ResponseCommon(data ...interface{}) error {
	o.PrepareData()
	err := o.AppendCommon(data...)
	if nil != err {
		return err
	}
	return o.SendResponse()
}

func (o *SockResponse) ResponseString(data ...string) error {
	o.PrepareData()
	err := o.AppendString(data...)
	if nil != err {
		return err
	}
	return o.SendResponse()
}

func (o *SockResponse) ResponseJson(data ...interface{}) error {
	o.PrepareData()
	err := o.AppendJson(data...)
	if nil != err {
		return err
	}
	return o.SendResponse()
}

func (o *SockResponse) ResponseObject(data ...interface{}) error {
	o.PrepareData()
	err := o.AppendObject(data...)
	if nil != err {
		return err
	}
	return o.SendResponse()
}

// private

func (o *SockResponse) sendRedirectMsgTo(eName string, pId string,
	interruptOnErr bool, userIds ...string) error {
	if len(userIds) == 0 {
		return nil
	}
	msg, err1 := o.GenMsgBytes(eName, pId)
	if nil != err1 {
		return err1
	}
	for _, userId := range userIds {
		if connId, ok := o.UserConnMapper.GetConnId(userId); ok {
			err := o.SockSender.SendPackTo(msg, connId)
			if nil != err && interruptOnErr {
				return err
			}
		}
	}
	return nil
}

func (o *SockResponse) sendRedirectMsg(eName string, pId string) error {
	msg, err := o.GenMsgBytes(eName, pId)
	if nil != err {
		return err
	}
	return o.SockSender.SendPackTo(msg, o.GetConnInfo().GetConnId()) // TODO: 未修改
}
