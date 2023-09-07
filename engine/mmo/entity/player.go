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

func NewIPlayerEntity(playerId string) basis.IPlayerEntity {
	return NewPlayerEntity(playerId)
}

func NewPlayerEntity(playerId string) *PlayerEntity {
	return &PlayerEntity{Uid: playerId}
}

type PlayerEntity struct {
	Uid string //用户标识，唯一，内部使用

	roomId, nextRoomId string
	roomLock           sync.RWMutex

	PlayerSubscriber
	vars.VariableSupport
	rooms []string
}

func (o *PlayerEntity) UID() string {
	return o.Uid
}

func (o *PlayerEntity) Name() string {
	return o.NickName()
}

func (o *PlayerEntity) EntityType() basis.EntityType {
	return basis.EntityPlayer
}

func (o *PlayerEntity) InitEntity() {
	o.PlayerSubscriber = *NewPlayerSubscriber()
	o.VariableSupport = *vars.NewVariableSupport(o)
}

func (o *PlayerEntity) DestroyEntity() {
}

// 扩展属性

func (o *PlayerEntity) NickName() string {
	nick, ok := o.GetVar(vars.PlayerNick)
	if !ok {
		return ""
	}
	return nick.(string)
}

func (o *PlayerEntity) Position() (pos basis.XYZ) {
	val, ok := o.GetVar(vars.PlayerPos)
	if !ok {
		return
	}
	return val.(basis.XYZ)
}

func (o *PlayerEntity) SetPosition(pos basis.XYZ) {
	o.SetVar(vars.PlayerPos, pos)
}

func (o *PlayerEntity) RoomId() string {
	o.roomLock.RLock()
	defer o.roomLock.RUnlock()
	return o.roomId
}

func (o *PlayerEntity) GetPrevRoomId() (roomId string, ok bool) {
	o.roomLock.RLock()
	defer o.roomLock.RUnlock()
	if len(o.rooms) == 0 {
		return
	}
	return o.rooms[len(o.rooms)-1], true
}

func (o *PlayerEntity) SetNextRoom(roomId string) {
	o.nextRoomId = roomId
}

func (o *PlayerEntity) ConfirmNextRoom(confirm bool) {
	defer func() {
		o.nextRoomId = ""
	}()
	if !confirm {
		return
	}
	o.roomLock.Lock()
	defer o.roomLock.Unlock()
	if o.roomId != o.nextRoomId {
		_ = o.PlayerSubscriber.RemoveWhite(o.roomId)
		o.rooms = append(o.rooms, o.nextRoomId)
		o.roomId = o.nextRoomId
		_ = o.PlayerSubscriber.AddWhite(o.roomId)
	}
}

func (o *PlayerEntity) BackToPrevRoom() {
	o.roomId = o.rooms[len(o.rooms)-1]
	o.rooms = o.rooms[:len(o.rooms)-1]
}

//---------------------------------

func (o *PlayerEntity) TeamId() string {
	val, ok := o.GetVar(vars.PlayerTeam)
	if !ok {
		return val.(string)
	}
	return ""
}

func (o *PlayerEntity) CorpsId() string {
	val, ok := o.GetVar(vars.PlayerTeamCorps)
	if !ok {
		return val.(string)
	}
	return ""
}

func (o *PlayerEntity) GetTeamInfo() (teamId string, corpsId string) {
	teamId, corpsId = o.TeamId(), o.CorpsId()
	return
}

func (o *PlayerEntity) SetTeamInfo(teamId string, corpsId string) {
	o.SetTeam(teamId)
	o.SetCorps(corpsId)
}

func (o *PlayerEntity) SetCorps(corpsId string) {
	if len(corpsId) == 0 {
		o.SetVar(vars.PlayerTeamCorps, nil)
	} else {
		o.SetVar(vars.PlayerTeamCorps, corpsId)
	}
}

func (o *PlayerEntity) SetTeam(teamId string) {
	if len(teamId) == 0 {
		o.SetVar(vars.PlayerTeam, nil)
	} else {
		o.SetVar(vars.PlayerTeam, teamId)
	}
}
