// Package protox
// Created by xuzhuoxi
// on 2019-03-22.
// @author xuzhuoxi
//
package extension

import (
	"github.com/xuzhuoxi/Rabbit-Server/engine/server"
	"github.com/xuzhuoxi/infra-go/netx"
)

func NewSockResponse() *SockResponse {
	return &SockResponse{
		ProtoReturnMessage: *NewProtoReturnMessage(),
	}
}

type SockResponse struct {
	ProtoReturnMessage
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

func (resp *SockResponse) SetParamInfo(paramType server.ExtensionParamType, paramHandler server.IProtoParamsHandler) {
	resp.ParamType, resp.ParamHandler = paramType, paramHandler
}

func (resp *SockResponse) SetResultCode(rsCode int32) {
	resp.RsCode = rsCode
}

func (resp *SockResponse) SendResponse() error {
	return resp.sendRedirectMsg(resp.PGroup, resp.PId)
}

func (resp *SockResponse) SendResponseTo(interruptOnErr bool, clientIds ...string) error {
	return resp.sendRedirectMsgTo(resp.PGroup, resp.PId, interruptOnErr, clientIds...)
}

func (resp *SockResponse) SendNotify(eName string, notifyPId string) error {
	return resp.sendRedirectMsg(eName, notifyPId)
}

func (resp *SockResponse) SendNotifyTo(eName string, notifyPId string, interruptOnErr bool, clientIds ...string) error {
	return resp.sendRedirectMsgTo(eName, notifyPId, interruptOnErr, clientIds...)
}

// private

func (resp *SockResponse) sendRedirectMsgTo(eName string, pId string,
	interruptOnErr bool, clientIds ...string) error {
	if len(clientIds) == 0 {
		return nil
	}
	msg, err1 := resp.genMsgBytes(eName, pId)
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
	msg, err := resp.genMsgBytes(eName, pId)
	if nil != err {
		return err
	}
	return resp.SockSender.SendPackTo(msg, resp.CAddress)
}
