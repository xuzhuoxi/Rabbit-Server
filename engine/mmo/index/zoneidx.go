// Package index
// Created by xuzhuoxi
// on 2019-03-09.
// @author xuzhuoxi
package index

import (
	"github.com/xuzhuoxi/Rabbit-Server/engine/mmo/basis"
)

func NewIZoneIndex() basis.IZoneIndex {
	return NewZoneIndex()
}

func NewZoneIndex() *ZoneIndex {
	return &ZoneIndex{EntityIndex: *NewEntityIndex("ZoneIndex", basis.EntityZone)}
}

type ZoneIndex struct {
	EntityIndex
}

func (o *ZoneIndex) CheckZone(zoneId string) bool {
	return o.EntityIndex.Check(zoneId)
}

func (o *ZoneIndex) GetZone(zoneId string) basis.IZoneEntity {
	entity := o.EntityIndex.Get(zoneId)
	if nil != entity {
		return entity.(basis.IZoneEntity)
	}
	return nil
}

func (o *ZoneIndex) AddZone(zone basis.IZoneEntity) error {
	return o.EntityIndex.Add(zone)
}

func (o *ZoneIndex) RemoveZone(zoneId string) (basis.IZoneEntity, error) {
	c, err := o.EntityIndex.Remove(zoneId)
	if nil != c {
		return c.(basis.IZoneEntity), err
	}
	return nil, err
}

func (o *ZoneIndex) UpdateZone(zone basis.IZoneEntity) error {
	return o.EntityIndex.Update(zone)
}
