// Package entity
// Created by xuzhuoxi
// on 2019-02-18.
// @author xuzhuoxi
package entity

import (
	"github.com/xuzhuoxi/Rabbit-Server/engine/mmo/basis"
	"sync"
)

func NewIUserEntity(userId string) basis.IUserEntity {
	return NewUserEntity(userId)
}

func NewUserEntity(userId string) *UserEntity {
	return &UserEntity{Uid: userId}
}

type UserEntity struct {
	Uid string //用户标识，唯一，内部使用

	RoomId, nextRoomId string
	roomLock           sync.RWMutex

	CorpsId string
	TeamId  string
	teamMu  sync.RWMutex

	UserSubscriber
	VariableSupport
	rooms []string
}

func (o *UserEntity) UID() string {
	return o.Uid
}

func (o *UserEntity) Name() string {
	return o.NickName()
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

func (o *UserEntity) NickName() string {
	nick, ok := o.GetVar(basis.VarKeyUserNick)
	if !ok {
		return ""
	}
	return nick.(string)
}

func (o *UserEntity) GetRoomId() string {
	o.roomLock.RLock()
	defer o.roomLock.RUnlock()
	return o.RoomId
}

func (o *UserEntity) GetPrevRoomId() (roomId string, ok bool) {
	o.roomLock.RLock()
	defer o.roomLock.RUnlock()
	if len(o.rooms) == 0 {
		return
	}
	return o.rooms[len(o.rooms)-1], true
}

func (o *UserEntity) SetNextRoom(roomId string) {
	o.nextRoomId = roomId
}

func (o *UserEntity) ConfirmNextRoom(confirm bool) {
	defer func() {
		o.nextRoomId = ""
	}()
	if !confirm {
		return
	}
	o.roomLock.Lock()
	defer o.roomLock.Unlock()
	if o.RoomId != o.nextRoomId {
		_ = o.UserSubscriber.RemoveWhite(o.RoomId)
		o.rooms = append(o.rooms, o.nextRoomId)
		o.RoomId = o.nextRoomId
		_ = o.UserSubscriber.AddWhite(o.RoomId)
	}
}

func (o *UserEntity) BackToPrevRoom() {
	o.RoomId = o.rooms[len(o.rooms)-1]
	o.rooms = o.rooms[:len(o.rooms)-1]
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

func (o *UserEntity) GetPosition() (pos basis.XYZ) {
	val, ok := o.GetVar(basis.VarKeyUserPos)
	if !ok {
		return
	}
	return val.(basis.XYZ)
}

func (o *UserEntity) SetPosition(pos basis.XYZ) {
	o.SetVar(basis.VarKeyUserPos, pos)
}
