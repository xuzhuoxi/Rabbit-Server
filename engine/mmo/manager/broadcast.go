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
	"math"
	"sync"
)

type IBroadcastManager interface {
	basis.IManagerBase

	// 以下为基础方法------

	// SetNearDistance 设置附近值
	SetNearDistance(distance float64)

	// 以下为广播方法

	// BroadcastCurrent 广播当前用户所在区域
	// source不能为nil
	BroadcastCurrent(sourceUser basis.IUserEntity, excludeBlack bool, handler func(user basis.IUserEntity)) error
	// BroadcastUsers 广播部分用户
	// targets 为用户实体IUserEntity的UID集合
	// sourceUser 可以为nil，当不为nil时会进行黑名单过滤，和本身过滤
	BroadcastUsers(sourceUser basis.IUserEntity, userIds []string, handler func(user basis.IUserEntity)) error

	// BroadcastEntity 广播实体
	// target 为房间实体
	// sourceUser 可以为nil，当不为nil时会进行黑名单过滤，和本身过滤
	BroadcastEntity(sourceUser basis.IUserEntity, entity basis.IEntityContainer, handler func(user basis.IUserEntity)) error
	// BroadcastRoom 广播房间
	// target 为房间实体
	// sourceUser 可以为nil，当不为nil时会进行黑名单过滤，和本身过滤
	BroadcastRoom(sourceUser basis.IUserEntity, room basis.IRoomEntity, handler func(user basis.IUserEntity)) error
	// BroadcastRooms 广播多个房间
	// rooms 指定的房间 Id 集合
	// sourceUser 可以为nil，当不为nil时会进行黑名单过滤，和本身过滤
	BroadcastRooms(sourceUser basis.IUserEntity, rooms []string, handler func(user basis.IUserEntity)) error
	// BroadcastZone 广播区域
	// zoneId 区域Id
	// sourceUser 可以为nil，当不为nil时会进行黑名单过滤，和本身过滤
	BroadcastZone(sourceUser basis.IUserEntity, zoneId string, handler func(user basis.IUserEntity)) error
	// BroadcastZones 广播多个房间
	// zones 指定的区域Id集合
	// sourceUser 可以为nil，当不为nil时会进行黑名单过滤，和本身过滤
	BroadcastZones(sourceUser basis.IUserEntity, zones []string, handler func(user basis.IUserEntity)) error
	// BroadcastWorld 广播世界
	// worldId 世界id
	// sourceUser 可以为nil，当不为nil时会进行黑名单过滤，和本身过滤
	BroadcastWorld(sourceUser basis.IUserEntity, worldId string, handler func(user basis.IUserEntity)) error
	// BroadcastWorlds 广播多个房间
	// worlds 指定的世界 Id 集合
	// sourceUser 可以为nil，当不为nil时会进行黑名单过滤，和本身过滤
	BroadcastWorlds(sourceUser basis.IUserEntity, worlds []string, handler func(user basis.IUserEntity)) error
	// BroadcastTag 广播包含tag的全部房间
	// tags 指定的房间 Id 集合
	// sourceUser 可以为nil，当不为nil时会进行黑名单过滤，和本身过滤
	BroadcastTag(sourceUser basis.IUserEntity, tag string, handler func(user basis.IUserEntity)) error
	// BroadcastTags 根据tags广播房间
	// tags 指定的房间 Id 集合
	// tagAnd true: 要求包含全部Tag: false: 要求包含其中一个tag
	BroadcastTags(sourceUser basis.IUserEntity, tags []string, tagAnd bool, handler func(user basis.IUserEntity)) error
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
				if excludeBlack && checkBlack(source, child) { //黑名单
					return
				}
				if userChild, ok := child.(basis.IUserEntity); ok {
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

func (o *BroadcastManager) BroadcastUsers(source basis.IUserEntity, userIds []string, handler func(entity basis.IUserEntity)) error {
	if len(userIds) == 0 {
		return errors.New("Len of target list is 0 ")
	}
	o.lock.RLock()
	defer o.lock.RUnlock()
	userIndex := o.entityMgr.UserIndex()
	for _, targetId := range userIds {
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

func (o *BroadcastManager) BroadcastRoom(sourceUser basis.IUserEntity, room basis.IRoomEntity, handler func(user basis.IUserEntity)) error {
	if nil == room {
		return errors.New("Target is nil. ")
	}
	o.lock.RLock()
	defer o.lock.RUnlock()
	if c, ok := room.(basis.IEntityContainer); ok {
		if nil == sourceUser {
			c.ForEachChildByType(basis.EntityUser, func(child basis.IEntity) {
				o.toChild(child, handler)
			}, false)
		} else {
			c.ForEachChildByType(basis.EntityUser, func(child basis.IEntity) {
				o.sourceToChild(sourceUser, child, handler)
			}, false)
		}
	}
	return nil
}

func (o *BroadcastManager) BroadcastRooms(sourceUser basis.IUserEntity, rooms []string, handler func(user basis.IUserEntity)) error {
	if len(rooms) == 0 {
		return errors.New("Len of room list is 0 ")
	}
	o.lock.RLock()
	defer o.lock.RUnlock()
	roomIndex := o.entityMgr.RoomIndex()
	for index := range rooms {
		if room, ok1 := roomIndex.GetRoom(rooms[index]); ok1 {
			if c, ok2 := room.(basis.IEntityContainer); ok2 {
				if sourceUser != nil {
					c.ForEachChildByType(basis.EntityUser, func(child basis.IEntity) {
						o.toChild(child, handler)
					}, false)
				} else {
					c.ForEachChildByType(basis.EntityUser, func(child basis.IEntity) {
						o.sourceToChild(sourceUser, child, handler)
					}, false)
				}
			}
		}
	}
	return nil
}

func (o *BroadcastManager) BroadcastZone(sourceUser basis.IUserEntity, zoneId string, handler func(user basis.IUserEntity)) error {
	return o.BroadcastTag(sourceUser, zoneId, handler)
}

func (o *BroadcastManager) BroadcastZones(sourceUser basis.IUserEntity, zones []string, handler func(user basis.IUserEntity)) error {
	return o.BroadcastTags(sourceUser, zones, false, handler)
}

func (o *BroadcastManager) BroadcastWorld(sourceUser basis.IUserEntity, worldId string, handler func(user basis.IUserEntity)) error {
	return o.BroadcastTag(sourceUser, worldId, handler)
}

func (o *BroadcastManager) BroadcastWorlds(sourceUser basis.IUserEntity, worlds []string, handler func(user basis.IUserEntity)) error {
	return o.BroadcastTags(sourceUser, worlds, false, handler)
}

func (o *BroadcastManager) BroadcastTag(sourceUser basis.IUserEntity, tag string, handler func(user basis.IUserEntity)) error {
	if len(tag) == 0 {
		return errors.New("Tag is nil. ")
	}
	o.lock.RLock()
	defer o.lock.RUnlock()
	roomIndex := o.entityMgr.RoomIndex()
	roomIndex.ForEachEntity(func(entity basis.IEntity) {
		room, ok1 := entity.(basis.IRoomEntity)
		if !ok1 || !room.ContainsTag(tag) {
			return
		}
		if c, ok2 := entity.(basis.IEntityContainer); ok2 {
			if sourceUser != nil {
				c.ForEachChildByType(basis.EntityUser, func(child basis.IEntity) {
					o.toChild(child, handler)
				}, false)
			} else {
				c.ForEachChildByType(basis.EntityUser, func(child basis.IEntity) {
					o.sourceToChild(sourceUser, child, handler)
				}, false)
			}
		}
	})
	return nil
}

func (o *BroadcastManager) BroadcastTags(sourceUser basis.IUserEntity, tags []string, tagAnd bool, handler func(user basis.IUserEntity)) error {
	if len(tags) == 0 {
		return errors.New("Len of tag list is 0 ")
	}
	o.lock.RLock()
	defer o.lock.RUnlock()
	roomIndex := o.entityMgr.RoomIndex()
	roomIndex.ForEachEntity(func(entity basis.IEntity) {
		room, ok1 := entity.(basis.IRoomEntity)
		if !ok1 || !room.ContainsTags(tags, tagAnd) {
			return
		}
		if c, ok2 := entity.(basis.IEntityContainer); ok2 {
			if sourceUser != nil {
				c.ForEachChildByType(basis.EntityUser, func(child basis.IEntity) {
					o.toChild(child, handler)
				}, false)
			} else {
				c.ForEachChildByType(basis.EntityUser, func(child basis.IEntity) {
					o.sourceToChild(sourceUser, child, handler)
				}, false)
			}
		}
	})
	return nil
}

func (o *BroadcastManager) BroadcastEntity(sourceUser basis.IUserEntity, entity basis.IEntityContainer, handler func(entity basis.IUserEntity)) error {
	o.lock.RLock()
	defer o.lock.RUnlock()
	if nil == entity {
		return errors.New(fmt.Sprintf("Entity is nil. "))
	}
	if nil == sourceUser {
		entity.ForEachChildByType(basis.EntityUser, func(child basis.IEntity) {
			o.toChild(child, handler)
		}, false)
	} else {
		entity.ForEachChildByType(basis.EntityUser, func(child basis.IEntity) {
			o.sourceToChild(sourceUser, child, handler)
		}, false)
	}
	return nil
}

//-----------------------------

func (o *BroadcastManager) toChild(child basis.IEntity, handler func(user basis.IUserEntity)) {
	if userChild, ok := child.(basis.IUserEntity); ok {
		handler(userChild)
	}
}

func (o *BroadcastManager) sourceToChild(sourceUser basis.IUserEntity, child basis.IEntity, handler func(user basis.IUserEntity)) {
	if checkSame(sourceUser, child) { //本身
		return
	}
	if checkBlack(sourceUser, child) { //黑名单
		return
	}
	if userChild, ok := child.(basis.IUserEntity); ok {
		handler(userChild)
	}
}

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
