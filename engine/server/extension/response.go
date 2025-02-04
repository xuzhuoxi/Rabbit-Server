// Package extension
// Created by xuzhuoxi
// on 2019-03-22.
// @author xuzhuoxi
//
package extension

import (
	"github.com/xuzhuoxi/Rabbit-Server/engine/server"
	"github.com/xuzhuoxi/Rabbit-Server/engine/server/packet"
	"github.com/xuzhuoxi/infra-go/netx"
)

func NewSockResponse() *SockResponse {
	return &SockResponse{
		ResponsePacket: *packet.NewResponsePacket(),
	}
}

type SockResponse struct {
	packet.ResponsePacket
	SockSender   netx.ISockSender
	AddressProxy netx.IAddressProxy
	ParamType    server.ExtensionParamType
}

func (resp *SockResponse) SetAddressProxy(proxy netx.IAddressProxy) {
	resp.AddressProxy = proxy
}

func (resp *SockResponse) SetSockSender(sockSender netx.ISockSender) {
	resp.SockSender = sockSender
}

func (resp *SockResponse) SetParamInfo(paramType server.ExtensionParamType, paramHandler server.IPacketParamsHandler) {
	resp.ParamType, resp.ParamHandler = paramType, paramHandler
}

func (resp *SockResponse) SetResultCode(rsCode int32) {
	resp.RsCode = rsCode
}

func (resp *SockResponse) SendResponse() error {
	return resp.sendRedirectMsg(resp.EName, resp.PId)
}

func (resp *SockResponse) SendResponseTo(interruptOnErr bool, clientIds ...string) error {
	return resp.sendRedirectMsgTo(resp.EName, resp.PId, interruptOnErr, clientIds...)
}

func (resp *SockResponse) SendNotify(eName string, notifyPId string) error {
	return resp.sendRedirectMsg(eName, notifyPId)
}

func (resp *SockResponse) SendNotifyTo(eName string, notifyPId string, interruptOnErr bool, clientIds ...string) error {
	return resp.sendRedirectMsgTo(eName, notifyPId, interruptOnErr, clientIds...)
}

// extend

func (resp *SockResponse) ResponseNone() error {
	resp.PrepareData()
	return resp.SendResponse()
}

func (resp *SockResponse) ResponseNoneToClient(interruptOnErr bool, clientIds ...string) error {
	if len(clientIds) == 0 {
		return nil
	}
	resp.PrepareData()
	return resp.SendResponseTo(interruptOnErr, clientIds...)
}

func (resp *SockResponse) ResponseBinary(data ...[]byte) error {
	resp.PrepareData()
	resp.AppendBinary(data...)
	return resp.SendResponse()
}

func (resp *SockResponse) ResponseCommon(data ...interface{}) error {
	resp.PrepareData()
	err := resp.AppendCommon(data...)
	if nil != err {
		return err
	}
	return resp.SendResponse()
}

func (resp *SockResponse) ResponseString(data ...string) error {
	resp.PrepareData()
	err := resp.AppendString(data...)
	if nil != err {
		return err
	}
	return resp.SendResponse()
}

func (resp *SockResponse) ResponseJson(data ...interface{}) error {
	resp.PrepareData()
	err := resp.AppendJson(data...)
	if nil != err {
		return err
	}
	return resp.SendResponse()
}

func (resp *SockResponse) ResponseObject(data ...interface{}) error {
	resp.PrepareData()
	err := resp.AppendObject(data...)
	if nil != err {
		return err
	}
	return resp.SendResponse()
}

// private

func (resp *SockResponse) sendRedirectMsgTo(eName string, pId string,
	interruptOnErr bool, clientIds ...string) error {
	if len(clientIds) == 0 {
		return nil
	}
	msg, err1 := resp.GenMsgBytes(eName, pId)
	if nil != err1 {
		return err1
	}
	for _, clientId := range clientIds {
		if address, ok := resp.AddressProxy.GetAddress(clientId); ok {
			err := resp.SockSender.SendPackTo(msg, address)
			if nil != err && interruptOnErr {
				return err
			}
		}
	}
	return nil
}

func (resp *SockResponse) sendRedirectMsg(eName string, pId string) error {
	msg, err := resp.GenMsgBytes(eName, pId)
	if nil != err {
		return err
	}
	return resp.SockSender.SendPackTo(msg, resp.CAddress)
}
