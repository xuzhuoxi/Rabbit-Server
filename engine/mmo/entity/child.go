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
	oMu   sync.RWMutex
}

func (o *EntityChildSupport) GetParent() string {
	o.oMu.RLock()
	defer o.oMu.RUnlock()
	return o.Owner
}

func (o *EntityChildSupport) NoneParent() bool {
	o.oMu.RLock()
	defer o.oMu.RUnlock()
	return o.Owner == ""
}

func (o *EntityChildSupport) SetParent(parentId string) {
	o.oMu.Lock()
	defer o.oMu.Unlock()
	o.Owner = parentId
}

func (o *EntityChildSupport) ClearParent() {
	o.oMu.Lock()
	defer o.oMu.Unlock()
	o.Owner = ""
}
