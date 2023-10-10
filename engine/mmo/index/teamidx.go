// Package index
// Created by xuzhuoxi
// on 2019-03-08.
// @author xuzhuoxi
package index

import (
	"github.com/xuzhuoxi/Rabbit-Server/engine/mmo/basis"
)

func NewITeamIndex() basis.ITeamIndex {
	return NewTeamIndex()
}

func NewTeamIndex() *TeamIndex {
	return &TeamIndex{EntityIndex: NewEntityIndex("TeamIndex", basis.EntityTeam)}
}

type TeamIndex struct {
	EntityIndex basis.IEntityIndex
}

func (o *TeamIndex) Size() int {
	return o.EntityIndex.Size()
}

func (o *TeamIndex) EntityType() basis.EntityType {
	return o.EntityIndex.EntityType()
}

func (o *TeamIndex) ForEachEntity(each func(entity basis.IEntity) (interrupt bool)) {
	o.EntityIndex.ForEachEntity(each)
}

func (o *TeamIndex) CheckTeam(teamId string) bool {
	return o.EntityIndex.Check(teamId)
}

func (o *TeamIndex) GetTeam(teamId string) (team basis.ITeamEntity, ok bool) {
	team, ok = o.EntityIndex.Get(teamId).(basis.ITeamEntity)
	return
}

func (o *TeamIndex) AddTeam(team basis.ITeamEntity) (rsCode int32, err error) {
	num, err1 := o.EntityIndex.Add(team)
	if nil == err1 {
		return
	}
	if num == 1 || num == 2 {
		return basis.CodeMMOIndexType, err1
	}
	return basis.CodeMMOTeamExist, err1
}

func (o *TeamIndex) RemoveTeam(teamId string) (team basis.ITeamEntity, rsCode int32, err error) {
	c, _, err1 := o.EntityIndex.Remove(teamId)
	if nil != c {
		return c.(basis.ITeamEntity), 0, nil
	}
	return nil, basis.CodeMMOTeamNotExist, err1
}

func (o *TeamIndex) UpdateTeam(team basis.ITeamEntity) (rsCode int32, err error) {
	_, err1 := o.EntityIndex.Update(team)
	if nil != err1 {
		return basis.CodeMMOIndexType, err1
	}
	return
}
