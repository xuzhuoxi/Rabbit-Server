// Package index
// Created by xuzhuoxi
// on 2019-03-09.
// @author xuzhuoxi
package index

import (
	"github.com/xuzhuoxi/Rabbit-Server/engine/mmo/basis"
)

func NewITeamCorpsIndex() basis.ITeamCorpsIndex {
	return NewTeamCorpsIndex()
}

func NewTeamCorpsIndex() *TeamCorpsIndex {
	return &TeamCorpsIndex{EntityIndex: *NewEntityIndex("TeamCorpsIndex", basis.EntityTeamCorps)}
}

type TeamCorpsIndex struct {
	EntityIndex
}

func (o *TeamCorpsIndex) CheckCorps(corpsId string) bool {
	return o.EntityIndex.Check(corpsId)
}

func (o *TeamCorpsIndex) GetCorps(corpsId string) (corps basis.ITeamCorpsEntity, ok bool) {
	corps, ok = o.EntityIndex.Get(corpsId).(basis.ITeamCorpsEntity)
	return
}

func (o *TeamCorpsIndex) AddCorps(corps basis.ITeamCorpsEntity) error {
	return o.EntityIndex.Add(corps)
}

func (o *TeamCorpsIndex) RemoveCorps(corpsId string) (basis.ITeamCorpsEntity, error) {
	c, err := o.EntityIndex.Remove(corpsId)
	if nil != c {
		return c.(basis.ITeamCorpsEntity), err
	}
	return nil, err
}

func (o *TeamCorpsIndex) UpdateCorps(corps basis.ITeamCorpsEntity) error {
	return o.EntityIndex.Update(corps)
}
