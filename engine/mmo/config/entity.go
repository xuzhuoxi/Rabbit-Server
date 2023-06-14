// Package config
// Create on 2023/6/14
// @author xuzhuoxi
package config

import (
	"errors"
	"fmt"
	"github.com/xuzhuoxi/infra-go/slicex"
)

type EntityType = string

const (
	None  EntityType = "none"
	World EntityType = "world"
	Zone  EntityType = "zone"
	Room  EntityType = "room"
)

type CfgMMOEntity struct {
	Id      string `json:"id" yaml:"id"`
	Name    string `json:"name" yaml:"name"`
	MaxUser int    `json:"max" yaml:"max"`
}

type CfgMMOEntities struct {
	Worlds []CfgMMOEntity `json:"worlds" yaml:"worlds"`
	Zones  []CfgMMOEntity `json:"zones" yaml:"zones"`
	Rooms  []CfgMMOEntity `json:"rooms" yaml:"rooms"`
}

func (o *CfgMMOEntities) String() string {
	return fmt.Sprintf("{World[%d], Zone[%d], Room[%d]}", len(o.Worlds), len(o.Zones), len(o.Rooms))
}

func (o *CfgMMOEntities) ExistWorld(worldId string) bool {
	return o.checkEntities(o.Worlds, worldId)
}

func (o *CfgMMOEntities) ExistZone(zoneId string) bool {
	return o.checkEntities(o.Zones, zoneId)
}

func (o *CfgMMOEntities) ExistRoom(roomId string) bool {
	return o.checkEntities(o.Rooms, roomId)
}

func (o *CfgMMOEntities) FindWorld(worldId string) (world CfgMMOEntity, ok bool) {
	if entity, ok := o.findEntity(o.Worlds, worldId); ok {
		return entity, true
	}
	return
}

func (o *CfgMMOEntities) FindZone(zoneId string) (zone CfgMMOEntity, ok bool) {
	if entity, ok := o.findEntity(o.Zones, zoneId); ok {
		return entity, true
	}
	return
}

func (o *CfgMMOEntities) FindRoom(roomId string) (room CfgMMOEntity, ok bool) {
	if entity, ok := o.findEntity(o.Rooms, roomId); ok {
		return entity, true
	}
	return
}

func (o *CfgMMOEntities) FindEntity(entityId string) (entity CfgMMOEntity, t EntityType, ok bool) {
	if entity, ok := o.findEntity(o.Worlds, entityId); ok {
		return entity, World, true
	}
	if entity, ok := o.findEntity(o.Zones, entityId); ok {
		return entity, Zone, true
	}
	if entity, ok := o.findEntity(o.Rooms, entityId); ok {
		return entity, Room, true
	}
	t = None
	return
}

func (o *CfgMMOEntities) CheckEntities() error {
	var idArr []string
	for _, world := range o.Worlds {
		if slicex.ContainsString(idArr, world.Id) {
			return errors.New("CfgMMOEntity duplicate definition at Id:" + world.Id)
		}
		idArr = append(idArr, world.Id)
	}
	for _, zone := range o.Zones {
		if slicex.ContainsString(idArr, zone.Id) {
			return errors.New("CfgMMOEntity duplicate definition at Id:" + zone.Id)
		}
		idArr = append(idArr, zone.Id)
	}
	for _, room := range o.Rooms {
		if slicex.ContainsString(idArr, room.Id) {
			return errors.New("CfgMMOEntity duplicate definition at Id:" + room.Id)
		}
		idArr = append(idArr, room.Id)
	}
	return nil
}

func (o *CfgMMOEntities) checkEntities(entities []CfgMMOEntity, entityId string) bool {
	if len(entities) == 0 {
		return false
	}
	for index := range entities {
		if entities[index].Id == entityId {
			return true
		}
	}
	return false
}

func (o *CfgMMOEntities) findEntity(entities []CfgMMOEntity, entityId string) (entity CfgMMOEntity, ok bool) {
	if len(entities) == 0 {
		return
	}
	for index := range entities {
		if entities[index].Id == entityId {
			return entities[index], true
		}
	}
	return
}
