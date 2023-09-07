// Package entity
// Created by xuzhuoxi
// on 2019-03-08.
// @author xuzhuoxi
package entity

import (
	"errors"
	"fmt"
	"github.com/xuzhuoxi/Rabbit-Server/engine/mmo/basis"
	"github.com/xuzhuoxi/Rabbit-Server/engine/mmo/vars"
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
	lock      sync.RWMutex

	vars.VariableSupport
	EntityListGroup
	ChildSupport
}

func (o *TeamEntity) UID() string {
	return o.TeamId
}

func (o *TeamEntity) Name() string {
	return o.TeamName
}

func (o *TeamEntity) EntityType() basis.EntityType {
	return basis.EntityTeam
}

func (o *TeamEntity) InitEntity() {
	o.VariableSupport = *vars.NewVariableSupport(o)
	o.EntityListGroup = *NewEntityListGroup(basis.EntityTeam)
	o.ChildSupport = *NewChildEntitySupport()
}

func (o *TeamEntity) Leader() string {
	return o.Owner
}

func (o *TeamEntity) MemberList() []string {
	return o.Entities()
}

func (o *TeamEntity) ContainMember(memberId string) bool {
	return o.ContainEntity(memberId)
}

func (o *TeamEntity) AcceptMember(memberId string) error {
	return o.Accept(memberId)
}

func (o *TeamEntity) DropMember(memberId string) error {
	o.lock.RLock()
	defer o.lock.RUnlock()
	err := o.Drop(memberId)
	if nil != err {
		return err
	}
	if memberId == o.Owner {
		if 0 == o.Len() {
			return o.disbandTeam()
		}
		o.SetParent(o.Entities()[0])
	}
	return nil
}

func (o *TeamEntity) RiseLeader(memberId string) error {
	o.lock.Lock()
	defer o.lock.Unlock()
	if memberId == o.Owner {
		return errors.New(fmt.Sprintf("%s is already the leader", memberId))
	}
	if !o.ContainEntity(memberId) {
		return errors.New(fmt.Sprintf("%s is not a member", memberId))
	}
	o.SetParent(memberId)
	return nil
}

func (o *TeamEntity) DisbandTeam() error {
	if o.Len() == 0 {
		return nil
	}
	return o.disbandTeam()
}

func (o *TeamEntity) disbandTeam() error {
	o.ClearParent()
	return nil
}
