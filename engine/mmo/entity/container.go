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

func NewIMapEntityContainer(cap int) basis.IEntityContainer {
	return NewMapEntityContainer(cap)
}

func NewIListEntityContainer(cap int) basis.IEntityContainer {
	return NewListEntityContainer(cap)
}

func NewMapEntityContainer(cap int) *MapEntityContainer {
	return &MapEntityContainer{cap: cap, entities: make(map[string]basis.IEntity)}
}

func NewListEntityContainer(cap int) *ListEntityContainer {
	return &ListEntityContainer{cap: cap}
}

//--------------------

type MapEntityContainer struct {
	cap      int
	entities map[string]basis.IEntity
	lock     sync.RWMutex
}

func (o *MapEntityContainer) NumChildren() int {
	o.lock.RLock()
	defer o.lock.RUnlock()
	return len(o.entities)
}

func (o *MapEntityContainer) Full() bool {
	o.lock.RLock()
	defer o.lock.RUnlock()
	return o.isFull()
}

func (o *MapEntityContainer) Contains(entity basis.IEntity) (isContains bool) {
	o.lock.RLock()
	defer o.lock.RUnlock()
	_, isContains = o.entities[entity.UID()]
	return
}

func (o *MapEntityContainer) ContainsById(entityId string) (isContains bool) {
	o.lock.RLock()
	defer o.lock.RUnlock()
	_, isContains = o.entities[entityId]
	return
}

func (o *MapEntityContainer) GetChildById(entityId string) (entity basis.IEntity, ok bool) {
	o.lock.RLock()
	defer o.lock.RUnlock()
	entity, ok = o.entities[entityId]
	return
}

func (o *MapEntityContainer) UpdateChild(entity basis.IEntity) (old basis.IEntity, err error) {
	if nil == entity {
		return nil, errors.New("Entity is nil. ")
	}
	o.lock.Lock()
	defer o.lock.Unlock()
	if o.isFull() {
		return nil, errors.New("Container is full ")
	}
	if v, ok := o.entities[entity.UID()]; ok {
		old = v
	}
	o.entities[entity.UID()] = entity
	return
}

func (o *MapEntityContainer) AddChild(entity basis.IEntity) error {
	if nil == entity {
		return errors.New("Entity is nil. ")
	}
	o.lock.Lock()
	defer o.lock.Unlock()
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
	if nil == entity {
		return errors.New("Entity is nil. ")
	}
	o.lock.Lock()
	defer o.lock.Unlock()
	id := entity.UID()
	_, isContains := o.entities[id]
	if !isContains {
		return errors.New(fmt.Sprintf("Entity(%s) does not exist in the container", id))
	}
	delete(o.entities, id)
	return nil
}

func (o *MapEntityContainer) RemoveChildById(entityId string) (entity basis.IEntity, ok bool) {
	o.lock.Lock()
	defer o.lock.Unlock()
	entity, ok = o.entities[entityId]
	if ok {
		delete(o.entities, entityId)
	}
	return
}

func (o *MapEntityContainer) UndoUpdate(old basis.IEntity, new basis.IEntity) {
	o.lock.Lock()
	defer o.lock.Unlock()
	if old == nil {
		delete(o.entities, new.UID())
	} else {
		o.entities[old.UID()] = old
	}
}

func (o *MapEntityContainer) UndoAdd(added basis.IEntity) {
	o.lock.Lock()
	defer o.lock.Unlock()
	delete(o.entities, added.UID())
}

func (o *MapEntityContainer) UndoRemove(removed basis.IEntity) {
	o.lock.Lock()
	defer o.lock.Unlock()
	o.entities[removed.UID()] = removed
}

func (o *MapEntityContainer) ForEachChild(each func(child basis.IEntity) (interruptCurrent bool, interruptRecurse bool)) {
	o.lock.RLock()
	defer o.lock.RUnlock()
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
	o.lock.RLock()
	defer o.lock.RUnlock()
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
	return o.cap > 0 && o.cap <= len(o.entities)
}

//--------------------

type ListEntityContainer struct {
	cap         int
	entities    []basis.IEntity
	lock        sync.RWMutex
	nextAdds    []basis.IEntity
	nextRemoves []basis.IEntity
}

func (o *ListEntityContainer) NumChildren() int {
	o.lock.RLock()
	defer o.lock.RUnlock()
	return len(o.entities)
}

func (o *ListEntityContainer) Full() bool {
	o.lock.RLock()
	defer o.lock.RUnlock()
	return o.isFull()
}

func (o *ListEntityContainer) Contains(entity basis.IEntity) (isContains bool) {
	o.lock.RLock()
	defer o.lock.RUnlock()
	if nil == entity {
		return false
	}
	e, _, ok := o.firstContains(entity.UID())
	if ok {
		return basis.EntityEqual(e, entity)
	}
	return false
}

func (o *ListEntityContainer) ContainsById(entityId string) (isContains bool) {
	o.lock.RLock()
	defer o.lock.RUnlock()
	_, _, isContains = o.firstContains(entityId)
	return
}

func (o *ListEntityContainer) GetChildById(childId string) (entity basis.IEntity, ok bool) {
	o.lock.RLock()
	defer o.lock.RUnlock()
	entity, _, ok = o.firstContains(childId)
	return
}

func (o *ListEntityContainer) UpdateChild(child basis.IEntity) (old basis.IEntity, err error) {
	if nil == child {
		return nil, errors.New("Entity is nil. ")
	}
	o.lock.Lock()
	defer o.lock.Unlock()
	id := child.UID()
	_, index, isContains := o.firstContains(id)
	if isContains {
		old = o.entities[index]
		o.entities[index] = child
	} else {
		if o.isFull() {
			return nil, errors.New("Container is full ")
		}
		o.entities = append(o.entities, child)
	}
	return
}

func (o *ListEntityContainer) AddChild(child basis.IEntity) error {
	if nil == child {
		return errors.New("Entity is nil. ")
	}
	o.lock.Lock()
	defer o.lock.Unlock()
	id := child.UID()
	_, _, isContains := o.firstContains(id)
	if isContains {
		return errors.New(fmt.Sprintf("Entity(%s) is already in the container", id))
	}
	if o.isFull() {
		return errors.New("Container is full ")
	}
	o.entities = append(o.entities, child)
	return nil
}

func (o *ListEntityContainer) RemoveChild(child basis.IEntity) error {
	if nil == child {
		return errors.New("Entity is nil. ")
	}
	o.lock.Lock()
	defer o.lock.Unlock()
	id := child.UID()
	_, index, isContains := o.firstContains(id)
	if !isContains {
		return errors.New(fmt.Sprintf("Entity(%s) does not exist in the container", id))
	}
	o.entities = append(o.entities[:index], o.entities[index+1:]...)
	return nil
}

func (o *ListEntityContainer) RemoveChildById(childId string) (entity basis.IEntity, ok bool) {
	o.lock.Lock()
	defer o.lock.Unlock()
	var index int
	entity, index, ok = o.firstContains(childId)
	if ok {
		o.entities = append(o.entities[:index], o.entities[index+1:]...)
	}
	return
}

func (o *ListEntityContainer) UndoUpdate(old basis.IEntity, new basis.IEntity) {
	if old == new {
		return
	}
	o.lock.Lock()
	defer o.lock.Unlock()
	_, index, isContains := o.firstContains(new.UID())
	if !isContains {
		return
	}
	if old == nil {
		o.entities = append(o.entities[:index], o.entities[index+1:]...)
	} else {
		o.entities[index] = old
	}
}

func (o *ListEntityContainer) UndoAdd(added basis.IEntity) {
	if nil == added {
		return
	}
	o.lock.Lock()
	defer o.lock.Unlock()
	id := added.UID()
	_, index, isContains := o.lastContains(id)
	if isContains {
		o.entities = append(o.entities[:index], o.entities[index+1:]...)
	}
}

func (o *ListEntityContainer) UndoRemove(removed basis.IEntity) {
	if nil == removed {
		return
	}
	o.lock.Lock()
	defer o.lock.Unlock()
	o.entities = append(o.entities, removed)
}

func (o *ListEntityContainer) ForEachChild(each func(child basis.IEntity) (interruptCurrent bool, interruptRecurse bool)) {
	o.lock.RLock()
	defer o.lock.RUnlock()
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
	o.lock.RLock()
	defer o.lock.RUnlock()
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

func (o *ListEntityContainer) firstContains(entityId string) (entity basis.IEntity, index int, isContains bool) {
	for index = 0; index < len(o.entities); index++ {
		if o.entities[index].UID() == entityId {
			entity = o.entities[index]
			isContains = true
			return
		}
	}
	return nil, -1, false
}

func (o *ListEntityContainer) lastContains(entityId string) (entity basis.IEntity, index int, isContains bool) {
	for index = len(o.entities) - 1; index >= 0; index-- {
		if o.entities[index].UID() == entityId {
			entity = o.entities[index]
			isContains = true
			return
		}
	}
	return nil, -1, false
}

func (o *ListEntityContainer) isFull() bool {
	return o.cap > 0 && o.cap <= len(o.entities)
}
