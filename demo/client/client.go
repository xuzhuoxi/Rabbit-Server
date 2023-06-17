// Package client
// Created by xuzhuoxi
// on 2019-03-24.
// @author xuzhuoxi
//
package client

import (
	"github.com/xuzhuoxi/infra-go/netx"
	"github.com/xuzhuoxi/infra-go/netx/tcpx"
)

const (
	RemoteAddress0 = "127.0.0.1:41000"
	RemoteAddress1 = "127.0.0.1:42000"
	Network        = netx.TcpNetwork
)

var ClientCreator = tcpx.NewTCP4Client
var count = 0

func NewUserClient(uId string) *UserClient {
	client := ClientCreator()
	return &UserClient{SockClient: client, UserId: uId}
}

type UserClient struct {
	UserId     string
	SockClient netx.ISockClient
}

func (uc *UserClient) Open() error {
	var remoteAddr string
	if count == 0 {
		remoteAddr = RemoteAddress0
		count = 1
	} else {
		remoteAddr = RemoteAddress1
		count = 0
	}
	return uc.SockClient.OpenClient(netx.SockParams{RemoteAddress: remoteAddr, Network: Network})
}
