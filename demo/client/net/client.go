// Package net
// Created by xuzhuoxi
// on 2019-03-24.
// @author xuzhuoxi
//
package net

import (
	"github.com/xuzhuoxi/infra-go/mathx"
	"github.com/xuzhuoxi/infra-go/netx"
	"github.com/xuzhuoxi/infra-go/netx/tcpx"
)

const (
	RemoteAddress1 = "127.0.0.1:41000"
	RemoteAddress2 = "127.0.0.1:42000"
	Network        = netx.TcpNetwork
)

var count = 0

func NewUserClient(uId string) *UserClient {
	return &UserClient{SockClient: tcpx.NewTCPClient(), UserId: uId}
}

type UserClient struct {
	UserId     string
	SockClient netx.ISockClient
}

func (uc *UserClient) Open() error {
	var err error = nil
	if count == 0 {
		err = uc.OpenWitAddr(RemoteAddress1)
	} else {
		err = uc.OpenWitAddr(RemoteAddress2)
	}
	count = mathx.AbsInt(count - 1)
	return err
}

func (uc *UserClient) OpenWitAddr(addr string) error {
	return uc.SockClient.OpenClient(netx.SockParams{RemoteAddress: addr, Network: Network})
}
