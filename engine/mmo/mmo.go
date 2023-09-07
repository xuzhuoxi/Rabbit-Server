// Package mmo
// Created by xuzhuoxi
// on 2019-03-15.
// @author xuzhuoxi
//
package mmo

import (
	"github.com/xuzhuoxi/Rabbit-Server/engine/mmo/basis"
	"github.com/xuzhuoxi/Rabbit-Server/engine/mmo/config"
	"github.com/xuzhuoxi/Rabbit-Server/engine/mmo/manager"
	"github.com/xuzhuoxi/infra-go/logx"
)

type IMMOManager interface {
	basis.IManagerBase
	logx.ILoggerSetter

	GetConfig() config.MMOConfig
	BuildEnv(cfg *config.MMOConfig) error
	GetEntityManager() manager.IEntityManager
	GetPlayerManager() manager.IPlayerManager
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
	playerMgr manager.IPlayerManager
	varMgr    manager.IVariableManager
	cfg       *config.MMOConfig
	logger    logx.ILogger
}

func (o *MMOManager) InitManager() {
	if nil != o.entityMgr {
		return
	}
	o.entityMgr = manager.NewIEntityManager()
	o.entityMgr.InitManager()
	o.playerMgr = manager.NewIPlayerManager(o.entityMgr)
	o.playerMgr.InitManager()
	o.bcMgr = manager.NewIBroadcastManager(o.entityMgr)
	o.bcMgr.InitManager()
	o.varMgr = manager.NewIVariableManager(o.entityMgr, o.bcMgr)
	o.varMgr.InitManager()
	o.SetLogger(o.logger)
}

func (o *MMOManager) DisposeManager() {
	o.varMgr.DisposeManager()
	o.bcMgr.DisposeManager()
	o.playerMgr.DisposeManager()
	o.entityMgr.DisposeManager()
}

func (o *MMOManager) SetLogger(logger logx.ILogger) {
	o.logger = logger
	if nil != o.varMgr {
		o.varMgr.SetLogger(logger)
	}
}

func (o *MMOManager) GetConfig() config.MMOConfig {
	return *o.cfg
}

func (o *MMOManager) BuildEnv(cfg *config.MMOConfig) error {
	o.cfg = cfg
	return o.entityMgr.BuildEnv(cfg)
}

func (o *MMOManager) GetEntityManager() manager.IEntityManager {
	return o.entityMgr
}

func (o *MMOManager) GetPlayerManager() manager.IPlayerManager {
	return o.playerMgr
}

func (o *MMOManager) GetBroadcastManager() manager.IBroadcastManager {
	return o.bcMgr
}

func (o *MMOManager) GetVarManager() manager.IVariableManager {
	return o.varMgr
}
