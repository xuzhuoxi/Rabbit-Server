// Package entity
// Created by xuzhuoxi
// on 2019-02-18.
// @author xuzhuoxi
package entity

import "github.com/xuzhuoxi/Rabbit-Server/engine/mmo/basis"

func NewITeamCorpsEntity(corpsId string, corpsName string) basis.ITeamCorpsEntity {
	return &TeamCorpsEntity{CorpsId: corpsId, CorpsName: corpsName}
}

func NewTeamCorpsEntity(corpsId string, corpsName string) *TeamCorpsEntity {
	return &TeamCorpsEntity{CorpsId: corpsId, CorpsName: corpsName}
}

type TeamCorpsEntity struct {
	CorpsId   string
	CorpsName string
	VariableSupport

	//EntityChildSupport
	//ListEntityContainer
	//TeamGroup *EntityListGroup
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
	//o.ListEntityContainer = *NewListEntityContainer(0)
	//e.TeamGroup = NewEntityListGroup(EntityTeam)
	o.VariableSupport = *NewVariableSupport(o)
}

//func (e *TeamCorpsEntity) TeamList() []string {
//	return e.TeamGroup.Entities()
//}
//
//func (e *TeamCorpsEntity) ContainTeam(corpsId string) bool {
//	return e.TeamGroup.ContainEntity(corpsId)
//}
//
//func (e *TeamCorpsEntity) AddTeam(corpsId string) error {
//	return e.TeamGroup.Accept(corpsId)
//}
//
//func (e *TeamCorpsEntity) RemoveTeam(corpsId string) error {
//	return e.TeamGroup.Drop(corpsId)
//}
