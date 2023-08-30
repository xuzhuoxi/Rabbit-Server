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
	"github.com/xuzhuoxi/infra-go/encodingx"
	"math"
	"sync"
)

type IBroadcastManager interface {
	basis.IManagerBase

	// 以下为基础方法------

	// SetNearDistance 设置附近值
	SetNearDistance(distance float64)

	// 以下为广播方法

	// BroadcastToUsers 广播部分用户
	// targets 为用户实体IUserEntity的UID集合
	// sourceUser 可以为nil，当不为nil时会进行黑名单过滤，和本身过滤
	BroadcastToUsers(sourceUser basis.IUserEntity, targets []string, handler func(user basis.IUserEntity)) error
	// BroadcastToRoom 广播房间
	// target 为房间实体
	// sourceUser 可以为nil，当不为nil时会进行黑名单过滤，和本身过滤
	BroadcastToRoom(sourceUser basis.IUserEntity, target basis.IRoomEntity, handler func(user basis.IUserEntity)) error
	// BroadcastToRooms 广播多个房间
	// rooms 指定的房间 Id 集合
	BroadcastToRooms(sourceUser basis.IUserEntity, rooms []string, handler func(user basis.IUserEntity)) error
	// BroadcastToTag 广播包含tag的全部房间
	// tags 指定的房间 Id 集合
	BroadcastToTag(sourceUser basis.IUserEntity, tag string, handler func(user basis.IUserEntity)) error
	// BroadcastToTags 根据tags广播房间
	// tags 指定的房间 Id 集合
	// tagAnd true: 要求包含全部Tag: false: 要求包含其中一个tag
	BroadcastToTags(sourceUser basis.IUserEntity, tags []string, tagAnd bool, handler func(user basis.IUserEntity)) error

	// BroadcastEntity 广播整个实体
	// target 为环境实体或用户实体
	// source可以为nil，当不为nil时会进行黑名单过滤，和本身过滤
	BroadcastEntity(source basis.IUserEntity, target basis.IEntity, handler func(user basis.IUserEntity)) error
	// BroadcastCurrent 广播当前用户所在区域
	// source不能为nil
	BroadcastCurrent(source basis.IUserEntity, excludeBlack bool, handler func(user basis.IUserEntity)) error

	//以下为业务型方法------

	// NotifyEnvVar 环境实体变量批量更新
	NotifyEnvVar(varTarget basis.IEntity, key string, value interface{})
	// NotifyEnvVars 环境实体变量批量更新
	NotifyEnvVars(varTarget basis.IEntity, vars encodingx.IKeyValue)
	// NotifyUserVar 用户实体变量更新
	NotifyUserVar(source basis.IUserEntity, key string, value interface{})
	// NotifyUserVars 用户实体变量批量更新
	NotifyUserVars(source basis.IUserEntity, vars encodingx.IKeyValue)
	// NotifyUserVarCurrent 用户实体变量更新
	NotifyUserVarCurrent(source basis.IUserEntity, key string, value interface{})
	// NotifyUserVarsCurrent 用户实体变量批量更新
	NotifyUserVarsCurrent(source basis.IUserEntity, vars encodingx.IKeyValue)
}

func NewIBroadcastManager(entityMgr IEntityManager) IBroadcastManager {
	return NewBroadcastManager(entityMgr)
}

func NewBroadcastManager(entityMgr IEntityManager) *BroadcastManager {
	return &BroadcastManager{entityMgr: entityMgr, distance: math.MaxFloat64}
}

//----------------------------------

type BroadcastManager struct {
	entityMgr IEntityManager
	lock      sync.RWMutex
	distance  float64
}

func (o *BroadcastManager) InitManager() {
	return
}

func (o *BroadcastManager) DisposeManager() {
	return
}

func (o *BroadcastManager) SetNearDistance(distance float64) {
	o.lock.Lock()
	defer o.lock.Unlock()
	o.distance = distance
}

func (o *BroadcastManager) BroadcastToRoom(sourceUser basis.IUserEntity, target basis.IRoomEntity, handler func(user basis.IUserEntity)) error {
	//TODO implement me
	panic("implement me")
}

func (o *BroadcastManager) BroadcastToRooms(sourceUser basis.IUserEntity, roomTags []string, handler func(user basis.IUserEntity)) error {
	//TODO implement me
	panic("implement me")
}

func (o *BroadcastManager) BroadcastToTag(sourceUser basis.IUserEntity, tag string, handler func(user basis.IUserEntity)) error {
	//TODO implement me
	panic("implement me")
}

func (o *BroadcastManager) BroadcastToTags(sourceUser basis.IUserEntity, tags []string, tagAnd bool, handler func(user basis.IUserEntity)) error {
	//TODO implement me
	panic("implement me")
}

func (o *BroadcastManager) BroadcastEntity(source basis.IUserEntity, target basis.IEntity, handler func(entity basis.IUserEntity)) error {
	o.lock.RLock()
	defer o.lock.RUnlock()
	if nil == target {
		return errors.New(fmt.Sprintf("Target is nil. "))
	}
	if userTarget, ok := target.(basis.IUserEntity); ok {
		if nil != source && (checkSame(source, userTarget) || checkBlack(source, userTarget)) { //本身 或 黑名单
			return nil
		}
		handler(userTarget)
		return nil
	}
	if entityContainer, ok := target.(basis.IEntityContainer); ok { //容器判断
		if nil == source {
			entityContainer.ForEachChildByType(basis.EntityUser, func(child basis.IEntity) {
				handler(child.(basis.IUserEntity))
			}, true)
		} else {
			entityContainer.ForEachChild(func(child basis.IEntity) (interruptCurrent bool, interruptRecurse bool) {
				if basis.EntityUser != child.EntityType() || checkSame(source, child) { //不是用户实体 或 是自己本身
					return
				}
				if checkBlack(source, child) { //黑名单
					return false, true
				}
				handler(child.(basis.IUserEntity))
				return
			})
		}
	}
	return nil
}

func (o *BroadcastManager) BroadcastToUsers(source basis.IUserEntity, targets []string, handler func(entity basis.IUserEntity)) error {
	o.lock.RLock()
	defer o.lock.RUnlock()
	if len(targets) == 0 {
		return errors.New("Targets's len is 0 ")
	}
	userIndex := o.entityMgr.UserIndex()
	for _, targetId := range targets {
		if targetUser, ok := userIndex.GetUser(targetId); ok { //目标用户存在
			if nil != source {
				if checkSame(source, targetUser) { //本身
					continue
				}
				if checkBlack(source, targetUser) { //黑名单
					continue
				}
			}
			handler(targetUser)
		}
	}
	return nil
}

func (o *BroadcastManager) BroadcastCurrent(source basis.IUserEntity, excludeBlack bool, handler func(entity basis.IUserEntity)) error {
	o.lock.RLock()
	defer o.lock.RUnlock()
	if nil == source {
		return errors.New(fmt.Sprintf("Source is nil. "))
	}
	roomId := source.GetRoomId()
	if room, ok1 := o.entityMgr.GetRoom(roomId); ok1 {
		if ec, ok2 := room.(basis.IEntityContainer); ok2 {
			ec.ForEachChildByType(basis.EntityUser, func(child basis.IEntity) {
				if checkSame(source, child) { //本身
					return
				}
				if userChild, ok := child.(basis.IUserEntity); ok {
					if excludeBlack && checkBlack(source, userChild) { //黑名单
						return
					}
					if !basis.NearXYZ(source.GetPosition(), userChild.GetPosition(), o.distance) { //位置不相近
						return
					}
					handler(userChild)
				}
			}, false)
		}
	}
	return nil
}

//-----------------------------

func (o *BroadcastManager) NotifyEnvVars(varTarget basis.IEntity, vars encodingx.IKeyValue) {
	if varTarget == nil {
		return
	}
	o.BroadcastEntity(nil, varTarget, func(entity basis.IUserEntity) {
	})
}

func (o *BroadcastManager) NotifyEnvVar(varTarget basis.IEntity, key string, value interface{}) {
}

func (o *BroadcastManager) NotifyUserVar(source basis.IUserEntity, key string, value interface{}) {
}

func (o *BroadcastManager) NotifyUserVars(source basis.IUserEntity, vars encodingx.IKeyValue) {
}

func (o *BroadcastManager) NotifyUserVarCurrent(source basis.IUserEntity, key string, value interface{}) {
}

func (o *BroadcastManager) NotifyUserVarsCurrent(source basis.IUserEntity, vars encodingx.IKeyValue) {
}

//-----------------------------

func checkSame(source basis.IUserEntity, target basis.IEntity) bool {
	return source.UID() == target.UID()
}

func checkBlack(source basis.IUserEntity, target basis.IEntity) bool {
	if source.OnBlack(target.UID()) { //source黑名单
		return true
	}
	if userTarget, ok := target.(basis.IUserEntity); ok { //target黑名单
		return userTarget.OnBlack(source.UID())
	}
	return false
}
