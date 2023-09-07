// Package entity
// Created by xuzhuoxi
// on 2019-02-18.
// @author xuzhuoxi
package entity

import (
	"github.com/xuzhuoxi/Rabbit-Server/engine/mmo/basis"
	"github.com/xuzhuoxi/Rabbit-Server/engine/mmo/vars"
	"sync"
)

func NewITeamCorpsEntity(corpsId string, corpsName string) basis.ITeamCorpsEntity {
	return &TeamCorpsEntity{CorpsId: corpsId, CorpsName: corpsName}
}

func NewTeamCorpsEntity(corpsId string, corpsName string) *TeamCorpsEntity {
	return &TeamCorpsEntity{CorpsId: corpsId, CorpsName: corpsName}
}

type TeamCorpsEntity struct {
	CorpsId   string
	CorpsName string
	lock      sync.RWMutex

	vars.VariableSupport
	EntityListGroup
	ChildSupport
}

func (o *TeamCorpsEntity) UID() string {
	return o.CorpsId
}

func (o *TeamCorpsEntity) Name() string {
	return o.CorpsName
}

func (o *TeamCorpsEntity) EntityType() basis.EntityType {
	return basis.EntityTeamCorps
}

func (o *TeamCorpsEntity) InitEntity() {
	o.VariableSupport = *vars.NewVariableSupport(o)
	o.EntityListGroup = *NewEntityListGroup(basis.EntityTeam)
	o.ChildSupport = *NewChildEntitySupport()
}

func (o *TeamCorpsEntity) TeamList() []string {
	return o.EntityListGroup.Entities()
}

func (o *TeamCorpsEntity) ContainTeam(corpsId string) bool {
	return o.EntityListGroup.ContainEntity(corpsId)
}

func (o *TeamCorpsEntity) AddTeam(corpsId string) error {
	return o.EntityListGroup.Accept(corpsId)
}

func (o *TeamCorpsEntity) RemoveTeam(corpsId string) error {
	return o.EntityListGroup.Drop(corpsId)
}
