// Package manager
// Created by xuzhuoxi
// on 2019-03-15.
// @author xuzhuoxi
//
package manager

import (
	"errors"
	"fmt"
	"github.com/xuzhuoxi/Rabbit-Server/engine/mmo/basis"
	"github.com/xuzhuoxi/Rabbit-Server/engine/mmo/events"
	"github.com/xuzhuoxi/Rabbit-Server/engine/mmo/vars"
	"github.com/xuzhuoxi/Rabbit-Server/engine/server"
	"github.com/xuzhuoxi/infra-go/eventx"
	"sync"
)

type IPlayerManager interface {
	basis.IManagerBase
	eventx.IEventDispatcher
	// LinkTheWorld 进入MMO世界，玩家实例不存在则进行创建
	LinkTheWorld(playerId string, roomId string) (player basis.IPlayerEntity, room basis.IRoomEntity, rsCode int32, err error)
	// UnlinkTheWorld 离开MMO世界，移除玩家实例
	UnlinkTheWorld(playerId string) (player basis.IPlayerEntity, rsCode int32, err error)

	// EnterRoom 进入房间，要求玩家实例已经存在
	EnterRoom(player basis.IPlayerEntity, roomId string) (room basis.IRoomEntity, rsCode int32, err error)
	// LeaveRoom 离开房间，要求玩家实例已经存在
	LeaveRoom(playerId string) (prevRoomId string, rsCode int32, err error)
	// Transfer 在世界转移，要求玩家实例已经存在
	Transfer(playerId string, toRoomId string, pos basis.XYZ) (room basis.IRoomEntity, rsCode int32, err error)
	// DragIntoRoom 把玩家插入房间
	DragIntoRoom(playerId string, toRoomId string, pos basis.XYZ) (room basis.IRoomEntity, rsCode int32, err error)
}

func NewIPlayerManager(entityMgr IEntityManager) IPlayerManager {
	return NewPlayerManager(entityMgr)
}

func NewPlayerManager(entityMgr IEntityManager) *PlayerManager {
	return &PlayerManager{entityMgr: entityMgr}
}

//----------------------------

type PlayerManager struct {
	eventx.EventDispatcher
	entityMgr IEntityManager
	transLock sync.RWMutex
}

func (o *PlayerManager) InitManager() {
	return
}

func (o *PlayerManager) DisposeManager() {
	return
}

func (o *PlayerManager) LinkTheWorld(playerId string, roomId string) (player basis.IPlayerEntity, room basis.IRoomEntity, rsCode int32, err error) {
	var pos basis.XYZ
	if player1, exist := o.entityMgr.GetPlayer(playerId); exist {
		if roomId == player1.RoomId() {
			room1, ok1 := o.entityMgr.GetRoom(roomId)
			if !ok1 {
				return nil, nil, basis.CodeMMORoomNotExist, nil
			}
			return player1, room1, server.CodeSuc, nil
		}
		player = player1
		pos = player1.Position()
	} else {
		pos = basis.RandomXYZ()
		vs := vars.DefaultVarSetPool.GetInstance()
		defer vars.DefaultVarSetPool.Recycle(vs)
		vs.Set(vars.PlayerPos, pos.Array())
		player, rsCode, err = o.entityMgr.CreatePlayer(playerId, vs)
		if nil != err {
			return nil, nil, rsCode, err
		}
	}
	room, rsCode, err = o.forwardTransfer(player, roomId, true, pos)
	if nil != err {
		return nil, nil, rsCode, err
	}
	return player, room, 0, nil
}

func (o *PlayerManager) UnlinkTheWorld(playerId string) (player basis.IPlayerEntity, rsCode int32, err error) {
	playerIndex := o.entityMgr.PlayerIndex()
	player, rsCode, err = playerIndex.RemovePlayer(playerId)
	if rsCode == server.CodeSuc {
		o.leaveRoom(player, player.RoomId())
		o.DispatchEvent(events.EventPlayerLeaveRoom, o, &events.PlayerEventDataLeaveRoom{RoomId: player.RoomId(), PlayerId: player.UID()})
	}
	return
}

func (o *PlayerManager) EnterRoom(player basis.IPlayerEntity, roomId string) (room basis.IRoomEntity, rsCode int32, err error) {
	if nil == player {
		return nil, basis.CodeMMOPlayerNotExist, errors.New("PlayerManager.EnterRoom Error: player is nil. ")
	}
	pos := basis.RandomXYZ()
	return o.forwardTransfer(player, roomId, true, pos)
}

func (o *PlayerManager) LeaveRoom(playerId string) (prevRoomId string, rsCode int32, err error) {
	playerIndex := o.entityMgr.PlayerIndex()
	player, ok := playerIndex.GetPlayer(playerId)
	if !ok {
		return "", basis.CodeMMOPlayerNotExist, errors.New(fmt.Sprintf("LeaveRoom Error: player(%s) does not exist", playerId))
	}
	return o.backwardTransfer(player, basis.ZeroXYZ)
}

func (o *PlayerManager) Transfer(playerId string, toRoomId string, pos basis.XYZ) (room basis.IRoomEntity, rsCode int32, err error) {
	playerIndex := o.entityMgr.PlayerIndex()
	player, ok := playerIndex.GetPlayer(playerId)
	if !ok {
		return nil, basis.CodeMMOPlayerNotExist, errors.New(fmt.Sprintf("Transfer Error: player(%s) does not exist", playerId))
	}
	return o.forwardTransfer(player, toRoomId, true, pos)
}

func (o *PlayerManager) DragIntoRoom(playerId string, toRoomId string, pos basis.XYZ) (room basis.IRoomEntity, rsCode int32, err error) {
	playerIndex := o.entityMgr.PlayerIndex()
	player, ok := playerIndex.GetPlayer(playerId)
	if !ok {
		return nil, basis.CodeMMOPlayerNotExist, errors.New(fmt.Sprintf("DragInto Error: player(%s) does not exist", playerId))
	}
	return o.forwardTransfer(player, toRoomId, false, pos)
}

func (o *PlayerManager) forwardTransfer(player basis.IPlayerEntity, toRoomId string, active bool, pos basis.XYZ) (room basis.IRoomEntity, rsCode int32, err error) {
	o.transLock.Lock()
	room1, ok := o.getRoom(toRoomId)
	if !ok {
		o.transLock.Unlock()
		return nil, basis.CodeMMORoomNotExist, errors.New("Room is not exist:" + toRoomId)
	}
	if room1.ContainsById(player.UID()) {
		return
	}
	oldRoom := o.leaveRoom(player, toRoomId)
	code2, err2 := o.forwardToRoom(room1, player, pos)
	if err2 != nil {
		if oldRoom != nil {
			oldRoom.UndoRemove(player)
		}
		player.ConfirmNextRoom(false)
		o.transLock.Unlock()
		return nil, code2, err2
	}
	player.ConfirmNextRoom(true)
	o.transLock.Unlock()
	if nil != oldRoom {
		o.DispatchEvent(events.EventPlayerLeaveRoom, o, &events.PlayerEventDataLeaveRoom{RoomId: oldRoom.UID(), PlayerId: player.UID()})
	}
	if active {
		o.DispatchEvent(events.EventPlayerEnterRoomActive, o, player)
	} else {
		o.DispatchEvent(events.EventPlayerEnterRoomPassive, o, player)
	}
	return room1, 0, nil
}

func (o *PlayerManager) backwardTransfer(player basis.IPlayerEntity, pos basis.XYZ) (prevRoomId string, rsCode int32, err error) {
	o.transLock.Lock()
	currentRoomId := player.RoomId()
	prevRoomId, _ = player.GetPrevRoomId()
	room, ok := o.getRoom(prevRoomId)
	if !ok {
		o.transLock.Unlock()
		return "", basis.CodeMMORoomNotExist, errors.New("Prev Room is not exist! ")
	}
	oldRoom := o.leaveRoom(player, prevRoomId)
	code2, err2 := o.backwardToRoom(room, player, pos)
	if err2 != nil {
		if oldRoom != nil {
			oldRoom.UndoRemove(player)
		}
		o.transLock.Unlock()
		return "", code2, err2
	}
	player.BackToPrevRoom()
	o.transLock.Unlock()
	o.DispatchEvent(events.EventPlayerLeaveRoom, o, &events.PlayerEventDataLeaveRoom{RoomId: currentRoomId, PlayerId: player.UID()})
	o.DispatchEvent(events.EventPlayerEnterRoomActive, o, player)
	return prevRoomId, 0, nil
}

func (o *PlayerManager) forwardToRoom(room basis.IRoomEntity, player basis.IPlayerEntity, pos basis.XYZ) (rsCode int32, err error) {
	errNum, err1 := room.AddChild(player)
	if errNum == 3 {
		return basis.CodeMMORoomCapLimit, err1
	}
	player.SetNextRoom(room.UID())
	player.SetPosition(pos, false)
	return
}
func (o *PlayerManager) backwardToRoom(room basis.IRoomEntity, player basis.IPlayerEntity, pos basis.XYZ) (rsCode int32, err error) {
	errNum, err1 := room.AddChild(player)
	if errNum == 3 {
		return basis.CodeMMORoomCapLimit, err1
	}
	player.SetPosition(pos, false)
	return
}

func (o *PlayerManager) leaveRoom(player basis.IPlayerEntity, nextRoomId string) (room basis.IRoomEntity) {
	player.SetNextRoom(nextRoomId)
	var ok bool
	room, ok = o.getRoom(player.RoomId())
	if ok {
		_, _ = room.RemoveChild(player)
	}
	return
}

func (o *PlayerManager) getRoom(roomId string) (room basis.IRoomEntity, ok bool) {
	if len(roomId) == 0 {
		return nil, false
	}
	roomIndex := o.entityMgr.RoomIndex()
	room, ok = roomIndex.GetRoom(roomId)
	return
}
