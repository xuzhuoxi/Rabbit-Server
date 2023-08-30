// Package config
// Create on 2023/6/14
// @author xuzhuoxi
package config

import "github.com/xuzhuoxi/infra-go/slicex"

type CfgRelationItem struct {
	Id   string   `json:"id" yaml:"id"`
	Name string   `json:"name" yaml:"name"`
	List []string `json:"list" yaml:"list"`
}

type CfgRelations struct {
	Zones  []CfgRelationItem `json:"zones" yaml:"zones"`
	Worlds []CfgRelationItem `json:"worlds" yaml:"worlds"`
}

func (o CfgRelations) ExistZone(zoneId string) bool {
	return o.existItem(o.Zones, zoneId)
}

func (o CfgRelations) ExistWorld(worldId string) bool {
	return o.existItem(o.Worlds, worldId)
}

func (o CfgRelations) ExistItem(items []CfgRelationItem, itemId string) bool {
	return o.existItem(items, itemId)
}

func (o CfgRelations) FindMyZones(listId string) []string {
	var rs []string
	for _, zone := range o.Zones {
		if slicex.ContainsString(zone.List, listId) {
			rs = append(rs, zone.Id)
		}
	}
	return rs
}

func (o CfgRelations) FindMyWorlds(listId string) []string {
	zones := o.FindMyZones(listId)
	var rs []string
	for _, world := range o.Worlds {
		for _, worldListItem := range world.List {
			if listId == worldListItem || slicex.ContainsString(zones, worldListItem) {
				rs = append(rs, world.Id)
				break
			}
		}
	}
	return rs
}

func (o CfgRelations) existItem(items []CfgRelationItem, itemId string) bool {
	if len(items) == 0 {
		return false
	}
	for _, item := range items {
		if item.Id == itemId {
			return true
		}
	}
	return false
}
