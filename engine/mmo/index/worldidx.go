// Package index
// Created by xuzhuoxi
// on 2019-03-09.
// @author xuzhuoxi
package index

import (
	"github.com/xuzhuoxi/Rabbit-Server/engine/mmo/basis"
)

func NewIWorldIndex() basis.IWorldIndex {
	return NewWorldIndex()
}

func NewWorldIndex() *WorldIndex {
	return &WorldIndex{EntityIndex: *NewEntityIndex("WorldIndex", basis.EntityWorld)}
}

type WorldIndex struct {
	EntityIndex
}

func (o *WorldIndex) CheckWorld(worldId string) bool {
	return o.EntityIndex.Check(worldId)
}

func (o *WorldIndex) GetWorld(worldId string) basis.IWorldEntity {
	entity := o.EntityIndex.Get(worldId)
	if nil != entity {
		return entity.(basis.IWorldEntity)
	}
	return nil
}

func (o *WorldIndex) AddWorld(world basis.IWorldEntity) error {
	return o.EntityIndex.Add(world)
}

func (o *WorldIndex) RemoveWorld(world string) (basis.IWorldEntity, error) {
	c, err := o.EntityIndex.Remove(world)
	if nil != c {
		return c.(basis.IWorldEntity), err
	}
	return nil, err
}

func (o *WorldIndex) UpdateWorld(world basis.IWorldEntity) error {
	return o.EntityIndex.Update(world)
}
