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
	"sync"
)

type IUserManager interface {
	basis.IManagerBase
	// EnterRoomAuto 进入房间，玩家实例不存在则进行创建
	EnterRoomAuto(userId string, roomId string, pos basis.XYZ) (user basis.IUserEntity, rsCode int32, err error)
	// EnterRoom 进入房间，要求玩家实例已经存在
	EnterRoom(user basis.IUserEntity, roomId string, pos basis.XYZ) (rsCode int32, err error)
	// LeaveRoom 离开房间，要求玩家实例已经存在
	LeaveRoom(userId string) (roomId string, rsCode int32, err error)
	// Transfer 在世界转移，要求玩家实例已经存在
	Transfer(userId string, toRoomId string, pos basis.XYZ) (rsCode int32, err error)
}

func NewIUserManager(entityMgr IEntityManager) IUserManager {
	return NewUserManager(entityMgr)
}

func NewUserManager(entityMgr IEntityManager) *UserManager {
	return &UserManager{entityMgr: entityMgr}
}

//----------------------------

type UserManager struct {
	entityMgr IEntityManager
	transLock sync.RWMutex
}

func (o *UserManager) InitManager() {
	return
}

func (o *UserManager) DisposeManager() {
	return
}

func (o *UserManager) EnterRoomAuto(userId string, roomId string, pos basis.XYZ) (user basis.IUserEntity, rsCode int32, err error) {
	o.transLock.Lock()
	defer o.transLock.Unlock()
	if _, exist := o.entityMgr.GetUser(userId); exist {
		return nil, basis.CodeMMOUserExist, errors.New("UserManager.UserBorn Error: User " + userId + " already exist. ")
	}
	vars := basis.NewVarSet()
	vars.Set(basis.VarKeyUserPos, pos)
	user, rsCode, err = o.entityMgr.CreateUser(userId, vars)
	if nil != err {
		return nil, rsCode, err
	}
	rsCode, err = o.forwardTransfer(user, roomId, pos)
	if nil != err {
		return nil, rsCode, err
	}
	return user, 0, nil
}

func (o *UserManager) EnterRoom(user basis.IUserEntity, roomId string, pos basis.XYZ) (rsCode int32, err error) {
	if nil == user {
		return basis.CodeMMOUserNotExist, errors.New("EnterRoom Error: user is nil. ")
	}
	o.transLock.Lock()
	defer o.transLock.Unlock()
	return o.forwardTransfer(user, roomId, pos)
}

func (o *UserManager) LeaveRoom(userId string) (roomId string, rsCode int32, err error) {
	o.transLock.Lock()
	defer o.transLock.Unlock()
	userIndex := o.entityMgr.UserIndex()
	user, ok := userIndex.GetUser(userId)
	if !ok {
		return "", basis.CodeMMOUserNotExist, errors.New(fmt.Sprintf("LeaveRoom Error: user(%s) does not exist", userId))
	}
	return o.backwardTransfer(user, basis.ZeroXYZ)
}

func (o *UserManager) Transfer(userId string, toRoomId string, pos basis.XYZ) (rsCode int32, err error) {
	o.transLock.Lock()
	defer o.transLock.Unlock()
	userIndex := o.entityMgr.UserIndex()
	user, ok := userIndex.GetUser(userId)
	if !ok {
		return basis.CodeMMOUserNotExist, errors.New(fmt.Sprintf("Transfer Error: user(%s) does not exist", userId))
	}
	return o.forwardTransfer(user, toRoomId, pos)
}

func (o *UserManager) forwardTransfer(user basis.IUserEntity, toRoomId string, pos basis.XYZ) (rsCode int32, err error) {
	room, ok := o.getRoom(toRoomId)
	if !ok {
		return basis.CodeMMORoomNotExist, errors.New("Room is not exist:" + toRoomId)
	}
	if room.ContainsById(user.UID()) {
		return
	}
	oldRoom := o.leaveRoom(user, toRoomId)
	code2, err2 := o.forwardToRoom(room, user, pos)
	if err2 != nil {
		if oldRoom != nil {
			oldRoom.UndoRemove(user)
		}
		user.ConfirmNextRoom(false)
		return code2, err2
	}
	user.ConfirmNextRoom(true)
	return 0, nil
}

func (o *UserManager) backwardTransfer(user basis.IUserEntity, pos basis.XYZ) (roomId string, rsCode int32, err error) {
	roomId, _ = user.GetPrevRoomId()
	room, ok := o.getRoom(roomId)
	if !ok {
		return "", basis.CodeMMORoomNotExist, errors.New("Prev Room is not exist! ")
	}
	oldRoom := o.leaveRoom(user, roomId)
	code2, err2 := o.backwardToRoom(room, user, pos)
	if err2 != nil {
		if oldRoom != nil {
			oldRoom.UndoRemove(user)
		}
		return "", code2, err2
	}
	user.BackToPrevRoom()
	return roomId, 0, nil
}

func (o *UserManager) forwardToRoom(room basis.IRoomEntity, user basis.IUserEntity, pos basis.XYZ) (rsCode int32, err error) {
	errNum, err1 := room.AddChild(user)
	if errNum == 3 {
		return basis.CodeMMORoomCapLimit, err1
	}
	user.SetNextRoom(room.UID())
	user.SetPosition(pos)
	return
}
func (o *UserManager) backwardToRoom(room basis.IRoomEntity, user basis.IUserEntity, pos basis.XYZ) (rsCode int32, err error) {
	errNum, err1 := room.AddChild(user)
	if errNum == 3 {
		return basis.CodeMMORoomCapLimit, err1
	}
	user.SetPosition(pos)
	return
}

func (o *UserManager) leaveRoom(user basis.IUserEntity, nextRoomId string) (room basis.IRoomEntity) {
	user.SetNextRoom(nextRoomId)
	var ok bool
	room, ok = o.getRoom(user.GetRoomId())
	if ok {
		_, _ = room.RemoveChild(user)
	}
	return
}

func (o *UserManager) getRoom(roomId string) (room basis.IRoomEntity, ok bool) {
	if len(roomId) == 0 {
		return nil, false
	}
	roomIndex := o.entityMgr.RoomIndex()
	room, ok = roomIndex.GetRoom(roomId)
	return
}
