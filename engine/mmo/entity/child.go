// Package entity
// Created by xuzhuoxi
// on 2019-02-18.
// @author xuzhuoxi
package entity

import (
	"github.com/xuzhuoxi/Rabbit-Server/engine/mmo/basis"
	"sync"
)

func NewIEntityChildSupport() basis.IEntityChild {
	return &EntityChildSupport{}
}

func NewEntityChildSupport() *EntityChildSupport {
	return &EntityChildSupport{}
}

type EntityChildSupport struct {
	Owner string
	lock  sync.RWMutex
}

func (o *EntityChildSupport) GetParent() string {
	o.lock.RLock()
	defer o.lock.RUnlock()
	return o.Owner
}

func (o *EntityChildSupport) IsNoneParent() bool {
	o.lock.RLock()
	defer o.lock.RUnlock()
	return o.Owner == ""
}

func (o *EntityChildSupport) SetParent(parentId string) {
	o.lock.Lock()
	defer o.lock.Unlock()
	o.Owner = parentId
}

func (o *EntityChildSupport) ClearParent() {
	o.lock.Lock()
	defer o.lock.Unlock()
	o.Owner = ""
}
