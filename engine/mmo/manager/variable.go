// Package manager
// Created by xuzhuoxi
// on 2019-03-16.
// @author xuzhuoxi
//
package manager

import (
	"github.com/xuzhuoxi/Rabbit-Server/engine/mmo/basis"
	"github.com/xuzhuoxi/infra-go/encodingx"
	"github.com/xuzhuoxi/infra-go/eventx"
	"github.com/xuzhuoxi/infra-go/logx"
)

type IVariableManager interface {
	basis.IManagerBase
}

func NewIVariableManager(entityManager IEntityManager, broadcastManager IBroadcastManager) IVariableManager {
	return NewVariableManager(entityManager, broadcastManager)
}

func NewVariableManager(entityManager IEntityManager, broadcastManager IBroadcastManager) *VariableManager {
	return &VariableManager{entityMgr: entityManager, bcMgr: broadcastManager, logger: logx.DefaultLogger()}
}

//--------------------------------

type VariableManager struct {
	entityMgr IEntityManager
	bcMgr     IBroadcastManager
	logger    logx.ILogger
}

func (o *VariableManager) InitManager() {
	o.entityMgr.AddEventListener(basis.EventVariableChanged, o.onEntityVar)
}

func (o *VariableManager) DisposeManager() {
	o.entityMgr.RemoveEventListener(basis.EventVariableChanged, o.onEntityVar)
}

func (o *VariableManager) SetLogger(logger logx.ILogger) {
	o.logger = logger
}

func (o *VariableManager) onEntityVar(evd *eventx.EventData) {
	data := evd.Data.([]interface{})
	currentTarget := data[0].(basis.IEntity)
	varSet := data[1].(encodingx.IKeyValue)
	if nil != o.logger {
		o.logger.Traceln("onEntityVar", currentTarget.UID(), varSet)
	}
	if currentTarget.EntityType() == basis.EntityUser {
		o.bcMgr.NotifyUserVars(currentTarget.(basis.IUserEntity), varSet)
	} else {
		o.bcMgr.NotifyEnvVars(currentTarget, varSet)
	}
}
