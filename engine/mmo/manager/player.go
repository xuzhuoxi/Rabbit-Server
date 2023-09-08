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
	"github.com/xuzhuoxi/infra-go/eventx"
	"sync"
)

type IPlayerManager interface {
	basis.IManagerBase
	eventx.IEventDispatcher
	// EnterRoomAuto 进入房间，玩家实例不存在则进行创建
	EnterRoomAuto(playerId string, roomId string, pos basis.XYZ) (player basis.IPlayerEntity, rsCode int32, err error)
	// EnterRoom 进入房间，要求玩家实例已经存在
	EnterRoom(player basis.IPlayerEntity, roomId string, pos basis.XYZ) (rsCode int32, err error)
	// LeaveRoom 离开房间，要求玩家实例已经存在
	LeaveRoom(playerId string) (roomId string, rsCode int32, err error)
	// Transfer 在世界转移，要求玩家实例已经存在
	Transfer(playerId string, toRoomId string, pos basis.XYZ) (rsCode int32, err error)
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

func (o *PlayerManager) EnterRoomAuto(playerId string, roomId string, pos basis.XYZ) (player basis.IPlayerEntity, rsCode int32, err error) {
	o.transLock.Lock()
	defer o.transLock.Unlock()
	if _, exist := o.entityMgr.GetPlayer(playerId); exist {
		return nil, basis.CodeMMOPlayerExist, errors.New("PlayerManager.EnterRoomAuto Error: Player " + playerId + " already exist. ")
	}
	vs := vars.NewVarSet()
	vs.Set(vars.PlayerPos, pos)
	player, rsCode, err = o.entityMgr.CreatePlayer(playerId, vs)
	if nil != err {
		return nil, rsCode, err
	}
	rsCode, err = o.forwardTransfer(player, roomId, pos)
	if nil != err {
		return nil, rsCode, err
	}
	return player, 0, nil
}

func (o *PlayerManager) EnterRoom(player basis.IPlayerEntity, roomId string, pos basis.XYZ) (rsCode int32, err error) {
	if nil == player {
		return basis.CodeMMOPlayerNotExist, errors.New("PlayerManager.EnterRoom Error: player is nil. ")
	}
	o.transLock.Lock()
	defer o.transLock.Unlock()
	return o.forwardTransfer(player, roomId, pos)
}

func (o *PlayerManager) LeaveRoom(playerId string) (roomId string, rsCode int32, err error) {
	o.transLock.Lock()
	defer o.transLock.Unlock()
	playerIndex := o.entityMgr.PlayerIndex()
	player, ok := playerIndex.GetPlayer(playerId)
	if !ok {
		return "", basis.CodeMMOPlayerNotExist, errors.New(fmt.Sprintf("LeaveRoom Error: player(%s) does not exist", playerId))
	}
	return o.backwardTransfer(player, basis.ZeroXYZ)
}

func (o *PlayerManager) Transfer(playerId string, toRoomId string, pos basis.XYZ) (rsCode int32, err error) {
	o.transLock.Lock()
	defer o.transLock.Unlock()
	playerIndex := o.entityMgr.PlayerIndex()
	player, ok := playerIndex.GetPlayer(playerId)
	if !ok {
		return basis.CodeMMOPlayerNotExist, errors.New(fmt.Sprintf("Transfer Error: player(%s) does not exist", playerId))
	}
	return o.forwardTransfer(player, toRoomId, pos)
}

func (o *PlayerManager) forwardTransfer(player basis.IPlayerEntity, toRoomId string, pos basis.XYZ) (rsCode int32, err error) {
	room, ok := o.getRoom(toRoomId)
	if !ok {
		return basis.CodeMMORoomNotExist, errors.New("Room is not exist:" + toRoomId)
	}
	if room.ContainsById(player.UID()) {
		return
	}
	oldRoom := o.leaveRoom(player, toRoomId)
	code2, err2 := o.forwardToRoom(room, player, pos)
	if err2 != nil {
		if oldRoom != nil {
			oldRoom.UndoRemove(player)
		}
		player.ConfirmNextRoom(false)
		return code2, err2
	}
	player.ConfirmNextRoom(true)
	oldRoomId := ""
	if nil != oldRoom {
		oldRoomId = oldRoom.UID()
	}
	o.DispatchEvent(events.EventPlayerLeaveRoom, o, &events.PlayerEventDataLeaveRoom{RoomId: oldRoomId, PlayerId: player.UID()})
	o.DispatchEvent(events.EventPlayerEnterRoom, o, player)
	return 0, nil
}

func (o *PlayerManager) backwardTransfer(player basis.IPlayerEntity, pos basis.XYZ) (roomId string, rsCode int32, err error) {
	oldRoomId := player.RoomId()
	roomId, _ = player.GetPrevRoomId()
	room, ok := o.getRoom(roomId)
	if !ok {
		return "", basis.CodeMMORoomNotExist, errors.New("Prev Room is not exist! ")
	}
	oldRoom := o.leaveRoom(player, roomId)
	code2, err2 := o.backwardToRoom(room, player, pos)
	if err2 != nil {
		if oldRoom != nil {
			oldRoom.UndoRemove(player)
		}
		return "", code2, err2
	}
	player.BackToPrevRoom()
	o.DispatchEvent(events.EventPlayerLeaveRoom, o, &events.PlayerEventDataLeaveRoom{RoomId: oldRoomId, PlayerId: player.UID()})
	o.DispatchEvent(events.EventPlayerEnterRoom, o, player)
	return roomId, 0, nil
}

func (o *PlayerManager) forwardToRoom(room basis.IRoomEntity, player basis.IPlayerEntity, pos basis.XYZ) (rsCode int32, err error) {
	errNum, err1 := room.AddChild(player)
	if errNum == 3 {
		return basis.CodeMMORoomCapLimit, err1
	}
	player.SetNextRoom(room.UID())
	player.SetPosition(pos)
	return
}
func (o *PlayerManager) backwardToRoom(room basis.IRoomEntity, player basis.IPlayerEntity, pos basis.XYZ) (rsCode int32, err error) {
	errNum, err1 := room.AddChild(player)
	if errNum == 3 {
		return basis.CodeMMORoomCapLimit, err1
	}
	player.SetPosition(pos)
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
