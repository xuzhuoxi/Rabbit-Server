// Package entity
// Created by xuzhuoxi
// on 2019-03-07.
// @author xuzhuoxi
package entity

import (
	"errors"
	"github.com/xuzhuoxi/Rabbit-Server/engine/mmo/basis"
	"github.com/xuzhuoxi/infra-go/slicex"
	"sync"
)

func NewIEntityGroup(entityType basis.EntityType, userMap bool) basis.IEntityGroup {
	if userMap {
		return NewEntityMapGroup(entityType)
	} else {
		return NewEntityListGroup(entityType)
	}
}

func NewEntityListGroup(entityType basis.EntityType) *EntityListGroup {
	return &EntityListGroup{entityType: entityType}
}

func NewEntityMapGroup(entityType basis.EntityType) *EntityMapGroup {
	return &EntityMapGroup{entityType: entityType}
}

//------------------------------

type EntityMapGroup struct {
	entityType basis.EntityType
	entityMap  map[string]*struct{}
	max        int
	mapLock    sync.RWMutex
}

func (o *EntityMapGroup) EntityType() basis.EntityType {
	return o.entityType
}

func (o *EntityMapGroup) MaxLen() int {
	return o.max
}

func (o *EntityMapGroup) Len() int {
	o.mapLock.RLock()
	defer o.mapLock.RUnlock()
	return len(o.entityMap)
}

func (o *EntityMapGroup) IsFull() bool {
	o.mapLock.RLock()
	defer o.mapLock.RUnlock()
	return o.isFull()
}

func (o *EntityMapGroup) Entities() []string {
	return o.CopyEntities()
}

func (o *EntityMapGroup) CopyEntities() []string {
	o.mapLock.RLock()
	defer o.mapLock.RUnlock()
	var rs []string
	for key, _ := range o.entityMap {
		rs = append(rs, key)
	}
	return rs
}

func (o *EntityMapGroup) ContainEntity(entityId string) bool {
	o.mapLock.RLock()
	defer o.mapLock.RUnlock()
	_, ok := o.entityMap[entityId]
	return ok
}

func (o *EntityMapGroup) Accept(entityId string) error {
	o.mapLock.Lock()
	defer o.mapLock.Unlock()
	_, ok := o.entityMap[entityId]
	if ok {
		return errors.New("EntityMapGroup.Accept Error: Entity(" + entityId + ") Duplicate")
	}
	if o.isFull() {
		return errors.New("EntityMapGroup.Accept Error: Group is Full")
	}
	o.entityMap[entityId] = nil
	return nil
}

func (o *EntityMapGroup) AcceptMulti(entityId []string) (count int, err error) {
	o.mapLock.Lock()
	defer o.mapLock.Unlock()
	if len(entityId) == 0 {
		return 0, errors.New("EntityMapGroup.AcceptMulti Error: len = 0")
	}
	for _, id := range entityId {
		_, ok := o.entityMap[id]
		if ok && nil != err {
			err = errors.New("EntityMapGroup.AcceptMulti Error: Entity Duplicate")
			continue
		}
		count++
		o.entityMap[id] = nil
	}
	return
}

func (o *EntityMapGroup) Drop(entityId string) error {
	o.mapLock.Lock()
	defer o.mapLock.Unlock()
	_, ok := o.entityMap[entityId]
	if !ok {
		return errors.New("EntityMapGroup.Drop Error: No Entity(" + entityId + ")")
	}
	delete(o.entityMap, entityId)
	return nil
}

func (o *EntityMapGroup) DropMulti(entityId []string) (count int, err error) {
	o.mapLock.Lock()
	defer o.mapLock.Unlock()
	if len(entityId) == 0 {
		return 0, errors.New("EntityMapGroup.DropMulti Error: len = 0")
	}
	if o.max > 0 && len(o.entityMap)+len(entityId) > o.max {
		return 0, errors.New("EntityMapGroup.AcceptMulti Error: Overcrowding")
	}
	for _, id := range entityId {
		_, ok := o.entityMap[id]
		if !ok && nil != err {
			err = errors.New("EntityMapGroup.DropMulti Error: No Entity")
			continue
		}
		count++
		delete(o.entityMap, id)
	}
	return
}

func (o *EntityMapGroup) isFull() bool {
	return o.max > 0 && len(o.entityMap) >= o.max
}

//---------------------------------------

type EntityListGroup struct {
	entityType basis.EntityType
	entityList []string
	max        int
	entityMu   sync.RWMutex
}

func (o *EntityListGroup) EntityType() basis.EntityType {
	return o.entityType
}

func (o *EntityListGroup) MaxLen() int {
	return o.max
}

func (o *EntityListGroup) Len() int {
	o.entityMu.RLock()
	defer o.entityMu.RUnlock()
	return len(o.entityList)
}

func (o *EntityListGroup) IsFull() bool {
	o.entityMu.RLock()
	defer o.entityMu.RUnlock()
	return o.isFull()
}

func (o *EntityListGroup) Entities() []string {
	o.entityMu.RLock()
	defer o.entityMu.RUnlock()
	return o.entityList
}

func (o *EntityListGroup) CopyEntities() []string {
	o.entityMu.RLock()
	defer o.entityMu.RUnlock()
	return slicex.CopyString(o.entityList)
}

func (o *EntityListGroup) ContainEntity(entityId string) bool {
	o.entityMu.RLock()
	defer o.entityMu.RUnlock()
	_, ok := slicex.IndexString(o.entityList, entityId)
	return ok
}

func (o *EntityListGroup) Accept(entityId string) error {
	o.entityMu.Lock()
	defer o.entityMu.Unlock()
	_, ok := slicex.IndexString(o.entityList, entityId)
	if ok {
		return errors.New("EntityListGroup.Accept Error: Entity(" + entityId + ") Duplicate")
	}
	if o.isFull() {
		return errors.New("EntityListGroup.Accept Error: Group is Full")
	}
	o.entityList = append(o.entityList, entityId)
	return nil
}

func (o *EntityListGroup) AcceptMulti(entityId []string) (count int, err error) {
	o.entityMu.Lock()
	defer o.entityMu.Unlock()
	if len(entityId) == 0 {
		return 0, errors.New("EntityListGroup.AcceptMulti Error: len = 0")
	}
	if o.max > 0 && len(o.entityList)+len(entityId) > o.max {
		return 0, errors.New("EntityListGroup.AcceptMulti Error: Overcrowding")
	}
	for _, id := range entityId {
		_, ok := slicex.IndexString(o.entityList, id)
		if ok && nil != err {
			err = errors.New("EntityListGroup.AcceptMulti Error: Entity Duplicate")
			continue
		}
		count++
		o.entityList = append(o.entityList, id)
	}
	return
}

func (o *EntityListGroup) Drop(entityId string) error {
	o.entityMu.Lock()
	defer o.entityMu.Unlock()
	index, ok := slicex.IndexString(o.entityList, entityId)
	if !ok {
		return errors.New("EntityListGroup.Drop Error: No Entity(" + entityId + ")")
	}
	o.entityList = append(o.entityList[:index], o.entityList[index+1:]...)
	return nil
}

func (o *EntityListGroup) DropMulti(entityId []string) (count int, err error) {
	o.entityMu.Lock()
	defer o.entityMu.Unlock()
	if len(entityId) == 0 {
		return 0, errors.New("EntityListGroup.DropMulti Error: len = 0")
	}
	for _, id := range entityId {
		index, ok := slicex.IndexString(o.entityList, id)
		if !ok && nil != err {
			err = errors.New("EntityListGroup.DropMulti Error: No Entity")
			continue
		}
		count++
		o.entityList = append(o.entityList[:index], o.entityList[index+1:]...)
	}
	return
}

func (o *EntityListGroup) isFull() bool {
	return o.max > 0 && len(o.entityList) >= o.max
}
