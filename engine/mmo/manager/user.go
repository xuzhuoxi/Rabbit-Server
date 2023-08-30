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
	EnterRoomAuto(userId string, roomId string, pos basis.XYZ) (user basis.IUserEntity, err error)
	// EnterRoom 进入房间，要求玩家实例已经存在
	EnterRoom(user basis.IUserEntity, roomId string, pos basis.XYZ) error
	// LeaveRoom 离开房间，要求玩家实例已经存在
	LeaveRoom(userId string) error
	// Transfer 在世界转移，要求玩家实例已经存在
	Transfer(userId string, toRoomId string, pos basis.XYZ) error
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

func (o *UserManager) EnterRoomAuto(userId string, roomId string, pos basis.XYZ) (user basis.IUserEntity, err error) {
	o.transLock.Lock()
	defer o.transLock.Unlock()
	if _, exist := o.entityMgr.GetUser(userId); exist {
		return nil, errors.New("UserManager.UserBorn Error: User " + userId + " already exist. ")
	}
	vars := basis.NewVarSet()
	vars.Set(basis.VarKeyUserPos, pos)
	userEntity, err1 := o.entityMgr.CreateUser(userId, vars)
	if nil != err1 {
		return nil, err1
	}
	err2 := o.transfer(userEntity, roomId, pos)
	if nil != err2 {
		return nil, err2
	}
	return userEntity, nil
}

func (o *UserManager) EnterRoom(user basis.IUserEntity, roomId string, pos basis.XYZ) error {
	if nil == user {
		return errors.New("EnterRoom Error: user is nil. ")
	}
	o.transLock.Lock()
	defer o.transLock.Unlock()
	return o.transfer(user, roomId, pos)
}

func (o *UserManager) LeaveRoom(userId string) error {
	o.transLock.Lock()
	defer o.transLock.Unlock()
	userIndex := o.entityMgr.UserIndex()
	user, ok := userIndex.GetUser(userId)
	if !ok {
		return errors.New(fmt.Sprintf("LeaveRoom Error: user(%s) does not exist", userId))
	}
	leave, err := o.leaveRoom(user)
	if nil != err {
		if nil != leave {
			leave.UndoRemove(user)
		}
		return err
	}
	user.ConfirmNextRoom(true)
	return nil
}

func (o *UserManager) Transfer(userId string, toRoomId string, pos basis.XYZ) error {
	o.transLock.Lock()
	defer o.transLock.Unlock()
	userIndex := o.entityMgr.UserIndex()
	user, ok := userIndex.GetUser(userId)
	if !ok {
		return errors.New(fmt.Sprintf("Transfer Error: user(%s) does not exist", userId))
	}
	return o.transfer(user, toRoomId, pos)
}

func (o *UserManager) transfer(user basis.IUserEntity, toRoomId string, pos basis.XYZ) error {
	oldRoom, err1 := o.leaveRoom(user)
	if nil != err1 {
		if oldRoom != nil {
			oldRoom.UndoRemove(user)
		}
		return err1
	}
	newRoom, err2 := o.enterRoom(user, toRoomId, pos)
	if err2 != nil {
		if oldRoom != nil {
			oldRoom.UndoRemove(user)
		}
		if nil != newRoom {
			newRoom.UndoAdd(user)
		}
		user.ConfirmNextRoom(false)
		return err2
	}
	user.ConfirmNextRoom(true)
	return nil
}

func (o *UserManager) enterRoom(user basis.IUserEntity, roomId string, pos basis.XYZ) (room basis.IRoomEntity, err error) {
	roomIndex := o.entityMgr.RoomIndex()
	r, ok := roomIndex.GetRoom(roomId)
	if !ok {
		return nil, errors.New("WorldManager.enterRoom Error: Room(" + roomId + ") does not exist")
	}
	err = r.AddChild(user)
	if nil != err {
		return
	}
	user.SetNextRoom(roomId)
	user.SetPosition(pos)
	return r, nil
}

func (o *UserManager) leaveRoom(user basis.IUserEntity) (room1 basis.IRoomEntity, err error) {
	roomId := user.GetRoomId()
	if len(roomId) == 0 {
		// 不在房间
		return nil, nil
	}
	roomIndex := o.entityMgr.RoomIndex()
	if "" == roomId || !roomIndex.CheckRoom(roomId) {
		return nil, nil
	}
	room, _ := roomIndex.GetRoom(roomId)
	err = room.RemoveChild(user)
	if nil != err {
		return
	}
	user.SetNextRoom("")
	return
}
