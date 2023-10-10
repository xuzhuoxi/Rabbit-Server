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
	return &TeamCorpsIndex{EntityIndex: NewEntityIndex("TeamCorpsIndex", basis.EntityTeamCorps)}
}

type TeamCorpsIndex struct {
	EntityIndex basis.IEntityIndex
}

func (o *TeamCorpsIndex) Size() int {
	return o.EntityIndex.Size()
}

func (o *TeamCorpsIndex) EntityType() basis.EntityType {
	return o.EntityIndex.EntityType()
}

func (o *TeamCorpsIndex) ForEachEntity(each func(entity basis.IEntity) (interrupt bool)) {
	o.EntityIndex.ForEachEntity(each)
}

func (o *TeamCorpsIndex) CheckCorps(corpsId string) bool {
	return o.EntityIndex.Check(corpsId)
}

func (o *TeamCorpsIndex) GetCorps(corpsId string) (corps basis.ITeamCorpsEntity, ok bool) {
	corps, ok = o.EntityIndex.Get(corpsId).(basis.ITeamCorpsEntity)
	return
}

func (o *TeamCorpsIndex) AddCorps(corps basis.ITeamCorpsEntity) (rsCode int32, err error) {
	num, err1 := o.EntityIndex.Add(corps)
	if nil == err1 {
		return
	}
	if num == 1 || num == 2 {
		return basis.CodeMMOIndexType, err1
	}
	return basis.CodeMMOTeamCorpsExist, err1
}

func (o *TeamCorpsIndex) RemoveCorps(corpsId string) (corps basis.ITeamCorpsEntity, rsCode int32, err error) {
	c, _, err1 := o.EntityIndex.Remove(corpsId)
	if nil != c {
		return c.(basis.ITeamCorpsEntity), 0, nil
	}
	return nil, basis.CodeMMOTeamCorpsNotExist, err1
}

func (o *TeamCorpsIndex) UpdateCorps(corps basis.ITeamCorpsEntity) (rsCode int32, err error) {
	_, err1 := o.EntityIndex.Update(corps)
	if nil != err1 {
		return basis.CodeMMOIndexType, err1
	}
	return
}
