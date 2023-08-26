// Package entity
// Created by xuzhuoxi
// on 2019-03-14.
// @author xuzhuoxi
package entity

import (
	"errors"
	"fmt"
	"github.com/xuzhuoxi/Rabbit-Server/engine/mmo/basis"
	"sync"
)

func NewIMapEntityContainer(maxCount int) basis.IEntityContainer {
	return NewMapEntityContainer(maxCount)
}

func NewIListEntityContainer(maxCount int) basis.IEntityContainer {
	return NewListEntityContainer(maxCount)
}

func NewMapEntityContainer(maxCount int) *MapEntityContainer {
	return &MapEntityContainer{maxCount: maxCount, entities: make(map[string]basis.IEntity)}
}

func NewListEntityContainer(maxCount int) *ListEntityContainer {
	return &ListEntityContainer{maxCount: maxCount}
}

//--------------------

type MapEntityContainer struct {
	maxCount    int
	entities    map[string]basis.IEntity
	containerMu sync.RWMutex
}

func (o *MapEntityContainer) NumChildren() int {
	o.containerMu.RLock()
	defer o.containerMu.RUnlock()
	return len(o.entities)
}

func (o *MapEntityContainer) Full() bool {
	o.containerMu.RLock()
	defer o.containerMu.RUnlock()
	return o.isFull()
}

func (o *MapEntityContainer) Contains(entity basis.IEntity) (isContains bool) {
	o.containerMu.RLock()
	defer o.containerMu.RUnlock()
	_, isContains = o.entities[entity.UID()]
	return
}

func (o *MapEntityContainer) ContainsById(entityId string) (isContains bool) {
	o.containerMu.RLock()
	defer o.containerMu.RUnlock()
	_, isContains = o.entities[entityId]
	return
}

func (o *MapEntityContainer) GetChildById(entityId string) (entity basis.IEntity, ok bool) {
	o.containerMu.RLock()
	defer o.containerMu.RUnlock()
	entity, ok = o.entities[entityId]
	return
}

func (o *MapEntityContainer) ReplaceChildInto(entity basis.IEntity) error {
	o.containerMu.Lock()
	defer o.containerMu.Unlock()
	if nil == entity {
		return errors.New("Entity is nil. ")
	}
	if o.isFull() {
		return errors.New("Container is full ")
	}
	o.entities[entity.UID()] = entity
	return nil
}

func (o *MapEntityContainer) AddChild(entity basis.IEntity) error {
	o.containerMu.Lock()
	defer o.containerMu.Unlock()
	if nil == entity {
		return errors.New("Entity is nil. ")
	}
	id := entity.UID()
	_, isContains := o.entities[id]
	if isContains {
		return errors.New(fmt.Sprintf("Entity(%s) is already in the container", id))
	}
	if o.isFull() {
		return errors.New("Container is full ")
	}
	o.entities[id] = entity
	return nil
}

func (o *MapEntityContainer) RemoveChild(entity basis.IEntity) error {
	o.containerMu.Lock()
	defer o.containerMu.Unlock()
	if nil == entity {
		return errors.New("Entity is nil. ")
	}
	id := entity.UID()
	_, isContains := o.entities[id]
	if !isContains {
		return errors.New(fmt.Sprintf("Entity(%s) does not exist in the container", id))
	}
	delete(o.entities, id)
	return nil
}

func (o *MapEntityContainer) RemoveChildById(entityId string) (entity basis.IEntity, ok bool) {
	o.containerMu.Lock()
	defer o.containerMu.Unlock()
	entity, ok = o.entities[entityId]
	if ok {
		delete(o.entities, entityId)
	}
	return
}

func (o *MapEntityContainer) ForEachChild(each func(child basis.IEntity) (interruptCurrent bool, interruptRecurse bool)) {
	o.containerMu.RLock()
	defer o.containerMu.RUnlock()
	if 0 == len(o.entities) {
		return
	}
	for _, entity := range o.entities {
		child := entity
		interruptCurrent, interruptRecurse := each(child)
		if interruptCurrent {
			return
		}
		if interruptRecurse {
			continue
		}
		if container, ok := entity.(basis.IEntityContainer); ok {
			container.ForEachChild(each)
		}
	}
}

func (o *MapEntityContainer) ForEachChildByType(entityType basis.EntityType, each func(child basis.IEntity), recurse bool) {
	o.containerMu.RLock()
	defer o.containerMu.RUnlock()
	if 0 == len(o.entities) {
		return
	}
	if recurse {
		for _, entity := range o.entities {
			if entity.EntityType() != entityType {
				continue
			}
			child := entity
			each(child)
			if container, ok := entity.(basis.IEntityContainer); ok {
				container.ForEachChildByType(entityType, each, true)
			}
		}
	} else {
		for _, entity := range o.entities {
			if entity.EntityType() != entityType {
				continue
			}
			child := entity
			each(child)
		}
	}
}

func (o *MapEntityContainer) isFull() bool {
	return o.maxCount > 0 && o.maxCount <= len(o.entities)
}

//--------------------

type ListEntityContainer struct {
	maxCount    int
	entities    []basis.IEntity
	containerMu sync.RWMutex
}

func (o *ListEntityContainer) NumChildren() int {
	o.containerMu.RLock()
	defer o.containerMu.RUnlock()
	return len(o.entities)
}

func (o *ListEntityContainer) Full() bool {
	o.containerMu.RLock()
	defer o.containerMu.RUnlock()
	return o.isFull()
}

func (o *ListEntityContainer) Contains(entity basis.IEntity) (isContains bool) {
	o.containerMu.RLock()
	defer o.containerMu.RUnlock()
	if nil == entity {
		return false
	}
	e, _, ok := o.contains(entity.UID())
	if ok {
		return basis.EntityEqual(e, entity)
	}
	return false
}

func (o *ListEntityContainer) ContainsById(entityId string) (isContains bool) {
	o.containerMu.RLock()
	defer o.containerMu.RUnlock()
	_, _, isContains = o.contains(entityId)
	return
}

func (o *ListEntityContainer) GetChildById(entityId string) (entity basis.IEntity, ok bool) {
	o.containerMu.RLock()
	defer o.containerMu.RUnlock()
	entity, _, ok = o.contains(entityId)
	return
}

func (o *ListEntityContainer) ReplaceChildInto(entity basis.IEntity) error {
	o.containerMu.Lock()
	defer o.containerMu.Unlock()
	if nil == entity {
		return errors.New("Entity is nil. ")
	}
	id := entity.UID()
	_, index, isContains := o.contains(id)
	if isContains {
		o.entities[index] = entity
	} else {
		if o.isFull() {
			return errors.New("Container is full ")
		}
		o.entities = append(o.entities, entity)
	}
	return nil
}

func (o *ListEntityContainer) AddChild(entity basis.IEntity) error {
	o.containerMu.Lock()
	defer o.containerMu.Unlock()
	if nil == entity {
		return errors.New("Entity is nil. ")
	}
	id := entity.UID()
	_, _, isContains := o.contains(id)
	if isContains {
		return errors.New(fmt.Sprintf("Entity(%s) is already in the container", id))
	}
	if o.isFull() {
		return errors.New("Container is full ")
	}
	o.entities = append(o.entities, entity)
	return nil
}

func (o *ListEntityContainer) RemoveChild(entity basis.IEntity) error {
	o.containerMu.Lock()
	defer o.containerMu.Unlock()
	if nil == entity {
		return errors.New("Entity is nil. ")
	}
	id := entity.UID()
	_, index, isContains := o.contains(id)
	if !isContains {
		return errors.New(fmt.Sprintf("Entity(%s) does not exist in the container", id))
	}
	o.entities = append(o.entities[:index], o.entities[index+1:]...)
	return nil
}

func (o *ListEntityContainer) RemoveChildById(entityId string) (entity basis.IEntity, ok bool) {
	var index int
	entity, index, ok = o.contains(entityId)
	if ok {
		o.entities = append(o.entities[:index], o.entities[index+1:]...)
	}
	return
}

func (o *ListEntityContainer) ForEachChild(each func(child basis.IEntity) (interruptCurrent bool, interruptRecurse bool)) {
	o.containerMu.RLock()
	defer o.containerMu.RUnlock()
	if 0 == len(o.entities) {
		return
	}
	for _, entity := range o.entities {
		child := entity
		interruptCurrent, interruptRecurse := each(child)
		if interruptCurrent {
			return
		}
		if interruptRecurse {
			continue
		}
		if container, ok := entity.(basis.IEntityContainer); ok {
			container.ForEachChild(each)
		}
	}
}

func (o *ListEntityContainer) ForEachChildByType(entityType basis.EntityType, each func(child basis.IEntity), recurse bool) {
	o.containerMu.RLock()
	defer o.containerMu.RUnlock()
	if 0 == len(o.entities) {
		return
	}
	if recurse {
		for _, entity := range o.entities {
			if entity.EntityType() != entityType {
				continue
			}
			child := entity
			each(child)
			if container, ok := entity.(basis.IEntityContainer); ok {
				container.ForEachChildByType(entityType, each, true)
			}
		}
	} else {
		for _, entity := range o.entities {
			if entity.EntityType() != entityType {
				continue
			}
			child := entity
			each(child)
		}
	}
}

func (o *ListEntityContainer) contains(entityId string) (entity basis.IEntity, index int, isContains bool) {
	for index = 0; index < len(o.entities); index++ {
		if o.entities[index].UID() == entityId {
			entity = o.entities[index]
			isContains = true
			return
		}
	}
	return nil, -1, false
}

func (o *ListEntityContainer) isFull() bool {
	return o.maxCount > 0 && o.maxCount <= len(o.entities)
}
