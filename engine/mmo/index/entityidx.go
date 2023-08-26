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
	entityMu   sync.RWMutex
}

func (o *EntityIndex) EntityType() basis.EntityType {
	return o.entityType
}

func (o *EntityIndex) Check(id string) bool {
	o.entityMu.RLock()
	defer o.entityMu.RUnlock()
	return o.check(id)
}

func (o *EntityIndex) check(id string) bool {
	_, ok := o.entityMap[id]
	return ok
}

func (o *EntityIndex) Get(id string) basis.IEntity {
	o.entityMu.RLock()
	defer o.entityMu.RUnlock()
	return o.entityMap[id]
}

func (o *EntityIndex) Add(entity basis.IEntity) error {
	o.entityMu.Lock()
	defer o.entityMu.Unlock()
	if nil == entity {
		//return errors.New(i.indexName + ".Add Error: entity is nil")
		return errors.New(fmt.Sprintf("%s.Add Error: entity is nil", o.indexName))
	}
	if !o.entityType.Include(entity.EntityType()) {
		//return errors.New(i.indexName + ".Add Error: Type is not included")
		return errors.New(fmt.Sprintf("%s.Add Error: Type is not included", o.indexName))
	}
	id := entity.UID()
	if o.check(id) {
		//return errors.New(i.indexName + ".Add Error: Id(" + id + ") Duplicate")
		return errors.New(fmt.Sprintf("%s.Add Error: Id(%s) Duplicate", o.indexName, id))
	}
	o.entityMap[id] = entity
	return nil
}

func (o *EntityIndex) Remove(id string) (basis.IEntity, error) {
	o.entityMu.Lock()
	defer o.entityMu.Unlock()
	e, ok := o.entityMap[id]
	if ok {
		delete(o.entityMap, id)
		return e, nil
	}
	//return nil, errors.New(i.indexName + ".Remove Error: No Entity(" + id + ")")
	return nil, errors.New(fmt.Sprintf("%s.Remove Error: No Entity(%s)", o.indexName, id))
}

func (o *EntityIndex) Update(entity basis.IEntity) error {
	o.entityMu.Lock()
	defer o.entityMu.Unlock()
	if nil == entity {
		//return errors.New(i.indexName + ".Update Error: entity is nil")
		return errors.New(fmt.Sprintf("%s.Update Error: entity is nil", o.indexName))
	}
	if !o.entityType.Include(entity.EntityType()) {
		//return errors.New(i.indexName + ".Update Error: Type is not included")
		return errors.New(fmt.Sprintf("%s.Update Error: Type is not included", o.indexName))
	}
	o.entityMap[entity.UID()] = entity
	return nil
}
