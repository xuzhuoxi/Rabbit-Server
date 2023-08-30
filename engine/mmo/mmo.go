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
)

type IMMOManager interface {
	basis.IManagerBase
	logx.ILoggerSetter

	GetEntityManager() manager.IEntityManager
	GetUserManager() manager.IUserManager
	GetBroadcastManager() manager.IBroadcastManager
	GetVarManager() manager.IVariableManager
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
	bcMgr     manager.IBroadcastManager
	userMgr   manager.IUserManager
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
	o.bcMgr = manager.NewIBroadcastManager(o.entityMgr)
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

func (o *MMOManager) SetLogger(logger logx.ILogger) {
	o.logger = logger
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

func (o *MMOManager) GetVarManager() manager.IVariableManager {
	return o.varMgr
}
