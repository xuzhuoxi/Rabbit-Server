// Package entity
// Created by xuzhuoxi
// on 2019-03-08.
// @author xuzhuoxi
package entity

import (
	"github.com/xuzhuoxi/Rabbit-Server/engine/mmo/basis"
	"sync"
)

func NewITeamEntity(teamId string, teamName string, maxMember int) basis.ITeamEntity {
	return &TeamEntity{TeamId: teamId, TeamName: teamName, MaxMember: maxMember}
}

func NewTeamEntity(teamId string, teamName string, maxMember int) *TeamEntity {
	return &TeamEntity{TeamId: teamId, TeamName: teamName, MaxMember: maxMember}
}

// TeamEntity 常规房间
type TeamEntity struct {
	TeamId    string
	TeamName  string
	MaxMember int
	EntityChildSupport
	ListEntityContainer

	//UserGroup *EntityListGroup
	VariableSupport

	teamMu sync.RWMutex
}

func (o *TeamEntity) UID() string {
	return o.TeamId
}

func (o *TeamEntity) NickName() string {
	return o.TeamName
}

func (o *TeamEntity) EntityType() basis.EntityType {
	return basis.EntityTeam
}

func (o *TeamEntity) InitEntity() {
	o.EntityChildSupport = *NewEntityChildSupport()
	o.ListEntityContainer = *NewListEntityContainer(o.MaxMember)
	//e.UserGroup = NewEntityListGroup(EntityUser)
	o.VariableSupport = *NewVariableSupport(o)
}

//func (e *TeamEntity) Leader() string {
//	e.teamMu.RLock()
//	defer e.teamMu.RUnlock()
//	return e.Owner
//}
//
//func (e *TeamEntity) MemberList() []string {
//	return e.UserGroup.Entities()
//}
//
//func (e *TeamEntity) ContainMember(memberId string) bool {
//	return e.UserGroup.ContainEntity(memberId)
//}
//
//func (e *TeamEntity) AcceptMember(memberId string) error {
//	return e.UserGroup.Accept(memberId)
//}
//
//func (e *TeamEntity) DropMember(memberId string) error {
//	e.teamMu.RLock()
//	defer e.teamMu.RUnlock()
//	err := e.UserGroup.Drop(memberId)
//	if nil != err {
//		return err
//	}
//	if memberId == e.Owner {
//		if 0 == e.UserGroup.Len() {
//			return e.disbandTeam()
//		}
//		e.SetParent(e.UserGroup.Entities()[0])
//	}
//	return nil
//}
//
//func (e *TeamEntity) RiseLeader(memberId string) error {
//	e.teamMu.Lock()
//	defer e.teamMu.Unlock()
//	if memberId == e.Owner {
//		return errors.New(fmt.Sprintf("%s is already the leader", memberId))
//	}
//	if !e.UserGroup.ContainEntity(memberId) {
//		return errors.New(fmt.Sprintf("%s is not a member", memberId))
//	}
//	e.SetParent(memberId)
//	return nil
//}
//
//func (e *TeamEntity) DisbandTeam() error {
//	if e.UserGroup.Len() == 0 {
//		return nil
//	}
//	return e.disbandTeam()
//}
//
//func (e *TeamEntity) disbandTeam() error {
//	return nil
//}
