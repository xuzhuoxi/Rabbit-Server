// Package entity
// Created by xuzhuoxi
// on 2019-02-18.
// @author xuzhuoxi
package entity

import (
	"github.com/xuzhuoxi/Rabbit-Server/engine/mmo/basis"
	"sync"
)

func NewIUserEntity(userId string, userName string) basis.IUserEntity {
	return NewUserEntity(userId, userName)
}

func NewUserEntity(userId string, userName string) *UserEntity {
	return &UserEntity{Uid: userId, Nick: userName}
}

type UserEntity struct {
	Uid  string //用户标识，唯一，内部使用
	Name string //用户名，唯一
	Nick string //用户昵称
	Addr string //用户历史或当前连接地址

	LocType basis.EntityType
	LocId   string
	locMu   sync.RWMutex

	CorpsId string
	TeamId  string
	teamMu  sync.RWMutex

	Pos   basis.XYZ
	posMu sync.RWMutex

	UserSubscriber
	VariableSupport
}

func (o *UserEntity) UID() string {
	return o.Uid
}

func (o *UserEntity) UserName() string {
	return o.Name
}

func (o *UserEntity) NickName() string {
	return o.Nick
}

func (o *UserEntity) EntityType() basis.EntityType {
	return basis.EntityUser
}

func (o *UserEntity) InitEntity() {
	o.UserSubscriber = *NewUserSubscriber()
	o.VariableSupport = *NewVariableSupport(o)
}

func (o *UserEntity) DestroyEntity() {
}

func (o *UserEntity) GetLocation() (idType basis.EntityType, id string) {
	o.locMu.RLock()
	defer o.locMu.RUnlock()
	return o.LocType, o.LocId
}

func (o *UserEntity) SetLocation(idType basis.EntityType, id string) {
	o.locMu.Lock()
	defer o.locMu.Unlock()
	if idType != o.LocType {
		o.LocType = idType
	}
	if id != o.LocId {
		o.UserSubscriber.RemoveWhite(o.LocId)
		o.LocId = id
		o.UserSubscriber.AddWhite(o.LocId)
	}
}

//---------------------------------

func (o *UserEntity) GetTeamInfo() (teamId string, corpsId string) {
	o.teamMu.RLock()
	defer o.teamMu.RUnlock()
	return o.TeamId, o.CorpsId
}

func (o *UserEntity) SetTeamInfo(teamId string, corpsId string) {
	o.teamMu.Lock()
	defer o.teamMu.Unlock()
	if o.TeamId != teamId {
		o.TeamId = teamId
	}
	if o.CorpsId != corpsId {
		o.CorpsId = corpsId
	}
}

func (o *UserEntity) SetCorps(corpsId string) {
	o.teamMu.Lock()
	defer o.teamMu.Unlock()
	if o.CorpsId != corpsId || o.TeamId == "" {
		o.CorpsId = corpsId
	}
}

func (o *UserEntity) SetTeam(teamId string) {
	o.teamMu.Lock()
	defer o.teamMu.Unlock()
	if o.TeamId != teamId {
		o.TeamId = teamId
		if teamId == "" { //没有队伍，不能回防兵团中
			o.CorpsId = ""
		}
	}
}

//---------------------------------

func (o *UserEntity) GetPosition() basis.XYZ {
	o.posMu.RLock()
	defer o.posMu.RUnlock()
	return o.Pos
}

func (o *UserEntity) SetPosition(pos basis.XYZ) {
	o.posMu.Lock()
	defer o.posMu.Unlock()
	o.Pos = pos
}
