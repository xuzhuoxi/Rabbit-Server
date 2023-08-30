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
	lock sync.RWMutex
}

func (o *VariableSupport) GetVar(key string) (interface{}, bool) {
	o.lock.RLock()
	defer o.lock.RUnlock()
	return o.vars.Get(key)
}

func (o *VariableSupport) CheckVar(key string) bool {
	o.lock.RLock()
	defer o.lock.RUnlock()
	return o.vars.Check(key)
}

func (o *VariableSupport) Vars() encodingx.IKeyValue {
	o.lock.RLock()
	defer o.lock.RUnlock()
	return o.vars
}

func (o *VariableSupport) SetVar(kv string, value interface{}) {
	if len(kv) == 0 {
		return
	}
	o.lock.Lock()
	defer o.lock.Unlock()
	var ok bool
	if value == nil {
		_, ok = o.vars.Delete(kv)
	} else {
		_, ok = o.vars.Set(kv, value)
	}
	if ok {
		o.DispatchEvent(basis.EventEntityVarChanged, o.currentTarget,
			&basis.VarEventData{Entity: o.currentTarget, Key: kv, Value: value})
	}
}

func (o *VariableSupport) SetVars(kv encodingx.IKeyValue) {
	if nil == kv {
		return
	}
	o.lock.Lock()
	defer o.lock.Unlock()
	diff, _ := o.vars.Merge(kv)
	if nil != diff {
		o.DispatchEvent(basis.EventEntityVarsChanged, o.currentTarget,
			&basis.VarsEventData{Entity: o.currentTarget, Vars: diff})
	}
}

//func (o *VariableSupport) RemoveVar(key string) {
//	o.lock.Lock()
//	defer o.lock.Unlock()
//	_, ok := o.vars.Delete(key)
//	if ok {
//		o.DispatchEvent(basis.EventEntityVarDeleted, o.currentTarget, key)
//	}
//}
//
//func (o *VariableSupport) RemoveVars() (count int) {
//	o.lock.Lock()
//	defer o.lock.Unlock()
//	count = o.vars.Len()
//	o.vars = basis.NewVarSet()
//	if count > 0 {
//		o.DispatchEvent(basis.EventEntityVarsDeleted, o.currentTarget, count)
//	}
//	return
//}
