// Package entity
// Created by xuzhuoxi
// on 2019-02-18.
// @author xuzhuoxi
package entity

import "github.com/xuzhuoxi/Rabbit-Server/engine/mmo/basis"

func CreateWorldEntity(worldId string, worldName string) basis.IWorldEntity {
	return &WorldEntity{WorldId: worldId, WorldName: worldName}
}

type WorldEntity struct {
	WorldId   string
	WorldName string
	ListEntityContainer

	//ZoneGroup *EntityListGroup
	VariableSupport
}

func (o *WorldEntity) UID() string {
	return o.WorldId
}

func (o *WorldEntity) NickName() string {
	return o.WorldName
}

func (o *WorldEntity) EntityType() basis.EntityType {
	return basis.EntityWorld
}

func (o *WorldEntity) InitEntity() {
	o.ListEntityContainer = *NewListEntityContainer(0)
	//w.ZoneGroup = NewEntityListGroup(EntityZone)
	o.VariableSupport = *NewVariableSupport(o)
}

//func (w *WorldEntity) ZoneList() []string {
//	return w.ZoneGroup.Entities()
//}
//
//func (w *WorldEntity) ContainZone(zoneId string) bool {
//	return w.ZoneGroup.ContainEntity(zoneId)
//}
//
//func (w *WorldEntity) AddZone(zoneId string) error {
//	return w.ZoneGroup.Accept(zoneId)
//}
//
//func (w *WorldEntity) RemoveZone(zoneId string) error {
//	return w.ZoneGroup.Drop(zoneId)
//}
