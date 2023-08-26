// Package entity
// Created by xuzhuoxi
// on 2019-03-03.
// @author xuzhuoxi
package entity

import (
	"github.com/xuzhuoxi/Rabbit-Server/engine/mmo/basis"
	"github.com/xuzhuoxi/infra-go/encodingx"
	"github.com/xuzhuoxi/infra-go/eventx"
	"sync"
)

func NewIVariableSupport(currentTarget basis.IEntity) basis.IVariableSupport {
	return NewVariableSupport(currentTarget)
}

func NewVariableSupport(currentTarget basis.IEntity) *VariableSupport {
	return &VariableSupport{currentTarget: currentTarget, vars: basis.NewVarSet()}
}

//---------------------------------------------

type VariableSupport struct {
	currentTarget basis.IEntity
	eventx.EventDispatcher
	vars encodingx.IKeyValue
	mu   sync.RWMutex
}

func (o *VariableSupport) Vars() encodingx.IKeyValue {
	o.mu.RLock()
	defer o.mu.RUnlock()
	return o.vars
}

func (o *VariableSupport) SetVar(key string, value interface{}) {
	o.mu.Lock()
	defer o.mu.Unlock()
	if diff, ok := o.vars.Set(key, value); ok {
		o.DispatchEvent(basis.EventEntityVarChanged, o.currentTarget, diff)
	}
}

func (o *VariableSupport) SetVars(kv encodingx.IKeyValue) {
	o.mu.Lock()
	defer o.mu.Unlock()
	diff := o.vars.Merge(kv)
	if nil != diff {
		o.DispatchEvent(basis.EventEntityVarChanged, o.currentTarget, diff)
	}
}

func (o *VariableSupport) GetVar(key string) (interface{}, bool) {
	o.mu.RLock()
	defer o.mu.RUnlock()
	return o.vars.Get(key)
}

func (o *VariableSupport) CheckVar(key string) bool {
	o.mu.RLock()
	defer o.mu.RUnlock()
	return o.vars.Check(key)
}

func (o *VariableSupport) RemoveVar(key string) {
	o.mu.Lock()
	defer o.mu.Unlock()
	o.vars.Delete(key)
}
