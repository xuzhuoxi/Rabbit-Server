// Package config
// Create on 2023/6/14
// @author xuzhuoxi
package config

import (
	"errors"
	"fmt"
	"github.com/xuzhuoxi/infra-go/slicex"
)

type CfgMMOEntity struct {
	Id   string   `json:"id" yaml:"id"`
	Name string   `json:"name" yaml:"name"`
	Cap  int      `json:"cap" yaml:"cap"`
	Tags []string `json:"tags" yaml:"tags"`
}

func (o *CfgMMOEntity) AppendTag(tag string) bool {
	if len(tag) == 0 || slicex.ContainsString(o.Tags, tag) {
		return false
	}
	o.Tags = append(o.Tags, tag)
	return true
}

func (o *CfgMMOEntity) AppendTags(tags []string) {
	if len(tags) == 0 {
		return
	}
	for index := range tags {
		_ = o.AppendTag(tags[index])
	}
}

type CfgMMOEntities struct {
	Rooms []CfgMMOEntity `json:"rooms" yaml:"rooms"`
}

func (o *CfgMMOEntities) String() string {
	return fmt.Sprintf("{Room[%d]}", len(o.Rooms))
}

func (o *CfgMMOEntities) ExistRoom(roomId string) bool {
	return o.checkEntities(o.Rooms, roomId)
}

func (o *CfgMMOEntities) FindRoom(roomId string) (room CfgMMOEntity, ok bool) {
	if index := o.findEntityIndex(o.Rooms, roomId); index != -1 {
		return o.Rooms[index], true
	}
	return
}

func (o *CfgMMOEntities) FindRoomIndex(roomId string) int {
	return o.findEntityIndex(o.Rooms, roomId)
}

func (o *CfgMMOEntities) CheckDuplicate() error {
	set := make(map[string]struct{})
	for _, e := range o.Rooms {
		if _, exist := set[e.Id]; exist {
			return errors.New("Duplicate at Id:" + e.Id)
		}
		set[e.Id] = struct{}{}
	}
	return nil
}

func (o *CfgMMOEntities) checkEntities(entities []CfgMMOEntity, entityId string) bool {
	return o.findEntityIndex(entities, entityId) != -1
}

func (o *CfgMMOEntities) findEntityIndex(entities []CfgMMOEntity, entityId string) (index int) {
	if len(entities) == 0 {
		return -1
	}
	for index := range entities {
		if entities[index].Id == entityId {
			return index
		}
	}
	return
}
