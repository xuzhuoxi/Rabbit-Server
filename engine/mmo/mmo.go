// Package mmo
// Created by xuzhuoxi
// on 2019-03-15.
// @author xuzhuoxi
//
package mmo

import (
	"github.com/xuzhuoxi/Rabbit-Server/engine/mmo/basis"
	"github.com/xuzhuoxi/Rabbit-Server/engine/mmo/manager"
	"github.com/xuzhuoxi/infra-go/logx"
	"github.com/xuzhuoxi/infra-go/netx"
)

type IMMOManager interface {
	basis.IManagerBase
	netx.ISockServerSetter
	netx.IAddressProxySetter

	GetEntityManager() manager.IEntityManager
	GetUserManager() manager.IUserManager
	GetBroadcastManager() manager.IBroadcastManager
}

func NewIMMOManager() IMMOManager {
	return NewMMOManager()
}

func NewMMOManager() *MMOManager {
	return &MMOManager{logger: logx.DefaultLogger()}
}

//----------------------------

type MMOManager struct {
	entityMgr manager.IEntityManager
	userMgr   manager.IUserManager
	bcMgr     manager.IBroadcastManager
	varMgr    manager.IVariableManager
	logger    logx.ILogger
}

func (o *MMOManager) InitManager() {
	if nil != o.entityMgr {
		return
	}
	o.entityMgr = manager.NewIEntityManager()
	o.entityMgr.InitManager()
	o.userMgr = manager.NewIUserManager(o.entityMgr)
	o.userMgr.InitManager()
	o.bcMgr = manager.NewIBroadcastManager(o.entityMgr, nil, nil)
	o.bcMgr.InitManager()
	o.varMgr = manager.NewIVariableManager(o.entityMgr, o.bcMgr)
	o.varMgr.InitManager()
	o.SetLogger(o.logger)
}

func (o *MMOManager) DisposeManager() {
	o.varMgr.DisposeManager()
	o.bcMgr.DisposeManager()
	o.userMgr.DisposeManager()
	o.entityMgr.DisposeManager()
}

func (o *MMOManager) SetSockServer(server netx.ISockServer) {
	if nil != o.bcMgr {
		o.bcMgr.SetSockServer(server)
	}
}

func (o *MMOManager) SetAddressProxy(proxy netx.IAddressProxy) {
	if nil != o.bcMgr {
		o.bcMgr.SetAddressProxy(proxy)
	}
}

func (o *MMOManager) SetLogger(logger logx.ILogger) {
	o.logger = logger
	if nil != o.entityMgr {
		o.entityMgr.SetLogger(logger)
	}
	if nil != o.userMgr {
		o.userMgr.SetLogger(logger)
	}
	if nil != o.bcMgr {
		o.bcMgr.SetLogger(logger)
	}
	if nil != o.varMgr {
		o.varMgr.SetLogger(logger)
	}
}

func (o *MMOManager) GetEntityManager() manager.IEntityManager {
	return o.entityMgr
}

func (o *MMOManager) GetUserManager() manager.IUserManager {
	return o.userMgr
}

func (o *MMOManager) GetBroadcastManager() manager.IBroadcastManager {
	return o.bcMgr
}
