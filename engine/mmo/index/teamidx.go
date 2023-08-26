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
	return &TeamIndex{EntityIndex: *NewEntityIndex("TeamIndex", basis.EntityTeam)}
}

type TeamIndex struct {
	EntityIndex
}

func (o *TeamIndex) CheckTeam(teamId string) bool {
	return o.EntityIndex.Check(teamId)
}

func (o *TeamIndex) GetTeam(teamId string) basis.ITeamEntity {
	entity := o.EntityIndex.Get(teamId)
	if nil != entity {
		return entity.(basis.ITeamEntity)
	}
	return nil
}

func (o *TeamIndex) AddTeam(team basis.ITeamEntity) error {
	return o.EntityIndex.Add(team)
}

func (o *TeamIndex) RemoveTeam(teamId string) (basis.ITeamEntity, error) {
	c, err := o.EntityIndex.Remove(teamId)
	if nil != c {
		return c.(basis.ITeamEntity), err
	}
	return nil, err
}

func (o *TeamIndex) UpdateTeam(team basis.ITeamEntity) error {
	return o.EntityIndex.Update(team)
}
