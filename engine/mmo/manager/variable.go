// Package manager
// Created by xuzhuoxi
// on 2019-03-16.
// @author xuzhuoxi
//
package manager

import (
	"github.com/xuzhuoxi/Rabbit-Server/engine/mmo/basis"
	"github.com/xuzhuoxi/Rabbit-Server/engine/mmo/events"
	"github.com/xuzhuoxi/infra-go/encodingx"
	"github.com/xuzhuoxi/infra-go/eventx"
	"github.com/xuzhuoxi/infra-go/logx"
)

type IVariableManager interface {
	logx.ILoggerSetter
	basis.IManagerBase
	SetVar(eType basis.EntityType, eId string, key string, value interface{})
	SetVars(eType basis.EntityType, eId string, vars encodingx.IKeyValue)
}

func NewIVariableManager(entityManager IEntityManager, broadcastManager IBroadcastManager) IVariableManager {
	return NewVariableManager(entityManager, broadcastManager)
}

func NewVariableManager(entityManager IEntityManager, broadcastManager IBroadcastManager) *VariableManager {
	return &VariableManager{entityMgr: entityManager, bcMgr: broadcastManager}
}

//--------------------------------

type VariableManager struct {
	logx.LoggerSupport
	entityMgr IEntityManager
	bcMgr     IBroadcastManager
}

func (o *VariableManager) SetVar(eType basis.EntityType, eId string, key string, value interface{}) {
	entity, ok1 := o.entityMgr.GetEntity(eType, eId)
	if !ok1 {
		o.GetLogger().Warnln("Entity is not exist: ", eType, eId)
		return
	}
	if ve, ok2 := entity.(basis.IVariableSupport); ok2 {
		ve.SetVar(key, value)
		return
	}
	o.GetLogger().Warnln("Entity does not support variable settings: ", eType, eId)
}

func (o *VariableManager) SetVars(eType basis.EntityType, eId string, vars encodingx.IKeyValue) {
	entity, ok1 := o.entityMgr.GetEntity(eType, eId)
	if !ok1 {
		o.GetLogger().Warnln("Entity is not exist: ", eType, eId)
		return
	}
	if ve, ok2 := entity.(basis.IVariableSupport); ok2 {
		ve.SetVars(vars)
		return
	}
	o.GetLogger().Warnln("Entity does not support variable settings: ", eType, eId)
}

func (o *VariableManager) InitManager() {
	//o.entityMgr.AddEventListener(basis.EventManagerVarChanged, o.onEntityVar)
	//o.entityMgr.AddEventListener(basis.EventManagerVarsChanged, o.onEntityVars)
}

func (o *VariableManager) DisposeManager() {
	//o.entityMgr.RemoveEventListener(basis.EventManagerVarsChanged, o.onEntityVars)
	//o.entityMgr.RemoveEventListener(basis.EventManagerVarChanged, o.onEntityVar)
}

func (o *VariableManager) onEntityVar(evd *eventx.EventData) {
	data := evd.Data.(*events.VarEventData)
	eventEntity := data.Entity
	key := data.Key
	value := data.Value
	logger := o.GetLogger()
	if nil != logger {
		logger.Traceln("onEntityVar", eventEntity.UID(), key, value)
	}
}

func (o *VariableManager) onEntityVars(evd *eventx.EventData) {
	data := evd.Data.(*events.VarsEventData)
	eventEntity := data.Entity
	varSet := data.Vars
	logger := o.GetLogger()
	if nil != logger {
		logger.Traceln("onEntityVars", eventEntity.UID(), varSet)
	}
}
