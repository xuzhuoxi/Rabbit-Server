// Package index
// Created by xuzhuoxi
// on 2019-03-17.
// @author xuzhuoxi
package index

import (
	"errors"
	"fmt"
	"github.com/xuzhuoxi/Rabbit-Server/engine/mmo/basis"
	"sync"
)

func NewIEntityIndex(indexName string, entityType basis.EntityType) basis.IEntityIndex {
	return NewEntityIndex(indexName, entityType)
}

func NewEntityIndex(indexName string, entityType basis.EntityType) *EntityIndex {
	return &EntityIndex{indexName: indexName, entityType: entityType, entityMap: make(map[string]basis.IEntity)}
}

type EntityIndex struct {
	indexName  string
	entityType basis.EntityType
	entityMap  map[string]basis.IEntity
	indexLock  sync.RWMutex
}

func (o *EntityIndex) EntityType() basis.EntityType {
	return o.entityType
}

func (o *EntityIndex) Check(id string) bool {
	o.indexLock.RLock()
	defer o.indexLock.RUnlock()
	return o.check(id)
}

func (o *EntityIndex) check(id string) bool {
	_, ok := o.entityMap[id]
	return ok
}

func (o *EntityIndex) Get(id string) basis.IEntity {
	o.indexLock.RLock()
	defer o.indexLock.RUnlock()
	if !o.check(id) {
		return nil
	}
	return o.entityMap[id]
}

func (o *EntityIndex) Add(entity basis.IEntity) (errNum int, err error) {
	o.indexLock.Lock()
	defer o.indexLock.Unlock()
	if nil == entity {
		//return 1, errors.New(o.indexName + ".Add Error: entity is nil")
		return 1, errors.New(fmt.Sprintf("%s.Add Error: entity is nil", o.indexName))
	}
	if !o.entityType.Include(entity.EntityType()) {
		//return 2, errors.New(o.indexName + ".Add Error: Type is not included")
		return 2, errors.New(fmt.Sprintf("%s.Add Error: Type is not included", o.indexName))
	}
	id := entity.UID()
	if o.check(id) {
		//return 3, errors.New(o.indexName + ".Add Error: Id(" + id + ") Duplicate")
		return 3, errors.New(fmt.Sprintf("%s.Add Error: Id(%s) Duplicate", o.indexName, id))
	}
	o.entityMap[id] = entity
	return 0, nil
}

func (o *EntityIndex) Remove(id string) (entity basis.IEntity, errNum int, err error) {
	o.indexLock.Lock()
	defer o.indexLock.Unlock()
	e, ok := o.entityMap[id]
	if ok {
		delete(o.entityMap, id)
		return e, 0, nil
	}
	//return nil, 1, errors.New(o.indexName + ".Remove Error: No Entity(" + id + ")")
	return nil, 1, errors.New(fmt.Sprintf("%s.Remove Error: No Entity(%s)", o.indexName, id))
}

func (o *EntityIndex) Update(entity basis.IEntity) (errNum int, err error) {
	o.indexLock.Lock()
	defer o.indexLock.Unlock()
	if nil == entity {
		//return 1, errors.New(o.indexName + ".Update Error: entity is nil")
		return 1, errors.New(fmt.Sprintf("%s.Update Error: entity is nil", o.indexName))
	}
	if !o.entityType.Include(entity.EntityType()) {
		//return 2, errors.New(o.indexName + ".Update Error: Type is not included")
		return 2, errors.New(fmt.Sprintf("%s.Update Error: Type is not included", o.indexName))
	}
	o.entityMap[entity.UID()] = entity
	return 0, nil
}

func (o *EntityIndex) ForEachEntity(each func(entity basis.IEntity)) {
	if nil == each {
		return
	}
	o.indexLock.Lock()
	defer o.indexLock.Unlock()
	for _, entity := range o.entityMap {
		each(entity)
	}
}
