// Package config
// Create on 2023/6/14
// @author xuzhuoxi
package config

import (
	"fmt"
	"github.com/xuzhuoxi/infra-go/slicex"
)

type CfgZoneRelation struct {
	ZoneId string   `json:"zone" yaml:"zone"`
	Rooms  []string `json:"rooms" yaml:"rooms"`
}

func (o CfgZoneRelation) ExistZone(zoneId string) bool {
	return o.ZoneId == zoneId
}

func (o CfgZoneRelation) ExistRoom(roomId string) bool {
	if len(o.Rooms) == 0 {
		return false
	}
	return slicex.ContainsString(o.Rooms, roomId)
}

type CfgWorldRelation struct {
	WorldId string            `json:"world" yaml:"world"`
	Zones   []CfgZoneRelation `json:"zones" yaml:"zones"`
}

func (o CfgWorldRelation) ExistWorld(worldId string) bool {
	return o.WorldId == worldId
}

func (o CfgWorldRelation) ExistZone(zoneId string) bool {
	if len(o.Zones) == 0 {
		return false
	}
	for index := range o.Zones {
		if o.Zones[index].ExistZone(zoneId) {
			return true
		}
	}
	return false
}

func (o CfgWorldRelation) ExistRoom(roomId string) bool {
	if len(o.Zones) == 0 {
		return false
	}
	for index := range o.Zones {
		if o.Zones[index].ExistRoom(roomId) {
			return true
		}
	}
	return false
}

type CfgRelations struct {
	Relations []CfgWorldRelation `json:"worlds" yaml:"worlds"`
}

func (o *CfgRelations) String() string {
	return fmt.Sprintf("{Relation[%d]}", len(o.Relations))
}

func (o *CfgRelations) GetWorldRelation(worldId string) (relation CfgWorldRelation, ok bool) {
	if len(o.Relations) == 0 {
		return
	}
	for index := range o.Relations {
		if o.Relations[index].WorldId == worldId {
			return o.Relations[index], true
		}
	}
	return
}

func (o *CfgRelations) CheckBelong(childId string, parentId string) bool {
	if len(o.Relations) == 0 || len(childId) == 0 || len(parentId) == 0 {
		return false
	}
	return o.checkZoneBelongWorld(childId, parentId) ||
		o.checkRoomBelongWorld(childId, parentId) ||
		o.checkRoomBelongZone(childId, parentId)
}

func (o *CfgRelations) checkRoomBelongZone(roomId string, zoneId string) bool {
	for rIdx := range o.Relations {
		for zIdx := range o.Relations[rIdx].Zones {
			if o.Relations[rIdx].Zones[zIdx].ZoneId == zoneId &&
				slicex.ContainsString(o.Relations[rIdx].Zones[zIdx].Rooms, roomId) {
				return true
			}
		}
	}
	return false
}

func (o *CfgRelations) checkRoomBelongWorld(roomId string, worldId string) bool {
	for rIdx := range o.Relations {
		if o.Relations[rIdx].WorldId != worldId {
			continue
		}
		for zIdx := range o.Relations[rIdx].Zones {
			if slicex.ContainsString(o.Relations[rIdx].Zones[zIdx].Rooms, roomId) {
				return true
			}
		}
	}
	return false
}

func (o *CfgRelations) checkZoneBelongWorld(zoneId string, worldId string) bool {
	for rIdx := range o.Relations {
		if o.Relations[rIdx].WorldId != worldId {
			continue
		}
		for zIdx := range o.Relations[rIdx].Zones {
			if o.Relations[rIdx].Zones[zIdx].ZoneId == zoneId {
				return true
			}
		}
	}
	return false
}
