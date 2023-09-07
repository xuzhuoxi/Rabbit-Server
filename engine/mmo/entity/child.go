// Package entity
// Created by xuzhuoxi
// on 2019-02-18.
// @author xuzhuoxi
package entity

import (
	"github.com/xuzhuoxi/Rabbit-Server/engine/mmo/basis"
	"sync"
)

func NewIChildEntitySupport() basis.IChildSupport {
	return &ChildSupport{}
}

func NewChildEntitySupport() *ChildSupport {
	return &ChildSupport{}
}

type ChildSupport struct {
	Owner string
	lock  sync.RWMutex
}

func (o *ChildSupport) GetParent() string {
	o.lock.RLock()
	defer o.lock.RUnlock()
	return o.Owner
}

func (o *ChildSupport) IsNoneParent() bool {
	o.lock.RLock()
	defer o.lock.RUnlock()
	return o.Owner == ""
}

func (o *ChildSupport) SetParent(parentId string) {
	o.lock.Lock()
	defer o.lock.Unlock()
	o.Owner = parentId
}

func (o *ChildSupport) ClearParent() {
	o.lock.Lock()
	defer o.lock.Unlock()
	o.Owner = ""
}
