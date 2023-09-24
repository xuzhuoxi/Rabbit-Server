// Package entity
// Created by xuzhuoxi
// on 2019-02-18.
// @author xuzhuoxi
package entity

import (
	"github.com/xuzhuoxi/Rabbit-Server/engine/mmo/basis"
	"github.com/xuzhuoxi/Rabbit-Server/engine/mmo/events"
	"github.com/xuzhuoxi/Rabbit-Server/engine/mmo/vars"
	"sync"
)

func NewIPlayerEntity(playerId string) basis.IPlayerEntity {
	return NewPlayerEntity(playerId)
}

func NewPlayerEntity(playerId string) *PlayerEntity {
	player := &PlayerEntity{Uid: playerId}
	player.PlayerSubscriber = *NewPlayerSubscriber()
	player.VariableSupport = *vars.NewVariableSupport(player)
	return player
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
	x, okX := o.GetVar(vars.PlayerPosX)
	if !okX {
		x = 0
	}
	y, okY := o.GetVar(vars.PlayerPosY)
	if !okY {
		y = 0
	}
	z, okZ := o.GetVar(vars.PlayerPosZ)
	if !okZ {
		z = 0
	}
	return basis.XYZ{X: x.(int32), Y: y.(int32), Z: z.(int32)}
}

func (o *PlayerEntity) SetPosition(pos basis.XYZ, notify bool) {
	okX := o.SetVar(vars.PlayerPosX, pos.X, false)
	okY := o.SetVar(vars.PlayerPosY, pos.Y, false)
	okZ := o.SetVar(vars.PlayerPosZ, pos.Z, false)
	if notify {
		var diff []string
		if okX {
			diff = append(diff, vars.PlayerPosY)
		}
		if okY {
			diff = append(diff, vars.PlayerPosY)
		}
		if okZ {
			diff = append(diff, vars.PlayerPosZ)
		}
		if len(diff) > 0 {
			o.VariableSupport.DispatchEvent(events.EventEntityVarsChanged, o,
				&events.VarsEventData{Entity: o, VarSet: o.Vars(), VarKeys: diff})

		}
	}
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

func (o *PlayerEntity) SetCorps(corpsId string, notify bool) {
	if len(corpsId) == 0 {
		o.SetVar(vars.PlayerTeamCorps, nil, notify)
	} else {
		o.SetVar(vars.PlayerTeamCorps, corpsId, notify)
	}
}

func (o *PlayerEntity) SetTeam(teamId string, notify bool) {
	if len(teamId) == 0 {
		o.SetVar(vars.PlayerTeam, nil, notify)
	} else {
		o.SetVar(vars.PlayerTeam, teamId, notify)
	}
}
