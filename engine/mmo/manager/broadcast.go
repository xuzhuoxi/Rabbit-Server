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
	SetNearDistance(distance int32)

	// 以下为广播方法

	// BroadcastCurrent 广播当前用户所在区域
	// source不能为nil
	BroadcastCurrent(sourcePlayer basis.IPlayerEntity, excludeBlack bool, handler func(player basis.IPlayerEntity)) error
	// BroadcastPlayers 广播部分用户
	// targets 为用户实体 IPlayerEntity 的UID集合
	// sourcePlayer 可以为nil，当不为nil时会进行黑名单过滤，和本身过滤
	BroadcastPlayers(sourcePlayer basis.IPlayerEntity, players []string, handler func(player basis.IPlayerEntity)) error

	// BroadcastEntity 广播实体
	// target 为房间实体
	// sourcePlayer 可以为nil，当不为nil时会进行黑名单过滤，和本身过滤
	BroadcastEntity(sourcePlayer basis.IPlayerEntity, entity basis.IEntityContainer, handler func(player basis.IPlayerEntity)) error
	// BroadcastRoom 广播房间
	// target 为房间实体
	// sourcePlayer 可以为nil，当不为nil时会进行黑名单过滤，和本身过滤
	BroadcastRoom(sourcePlayer basis.IPlayerEntity, room basis.IRoomEntity, handler func(player basis.IPlayerEntity)) error
	// BroadcastRooms 广播多个房间
	// rooms 指定的房间 Id 集合
	// sourcePlayer 可以为nil，当不为nil时会进行黑名单过滤，和本身过滤
	BroadcastRooms(sourcePlayer basis.IPlayerEntity, rooms []string, handler func(player basis.IPlayerEntity)) error
	// BroadcastZone 广播区域
	// zoneId 区域Id
	// sourcePlayer 可以为nil，当不为nil时会进行黑名单过滤，和本身过滤
	BroadcastZone(sourcePlayer basis.IPlayerEntity, zoneId string, handler func(player basis.IPlayerEntity)) error
	// BroadcastZones 广播多个房间
	// zones 指定的区域Id集合
	// sourcePlayer 可以为nil，当不为nil时会进行黑名单过滤，和本身过滤
	BroadcastZones(sourcePlayer basis.IPlayerEntity, zones []string, handler func(player basis.IPlayerEntity)) error
	// BroadcastWorld 广播世界
	// worldId 世界id
	// sourcePlayer 可以为nil，当不为nil时会进行黑名单过滤，和本身过滤
	BroadcastWorld(sourcePlayer basis.IPlayerEntity, worldId string, handler func(player basis.IPlayerEntity)) error
	// BroadcastWorlds 广播多个房间
	// worlds 指定的世界 Id 集合
	// sourcePlayer 可以为nil，当不为nil时会进行黑名单过滤，和本身过滤
	BroadcastWorlds(sourcePlayer basis.IPlayerEntity, worlds []string, handler func(player basis.IPlayerEntity)) error
	// BroadcastTag 广播包含tag的全部房间
	// tags 指定的房间 Id 集合
	// sourcePlayer 可以为nil，当不为nil时会进行黑名单过滤，和本身过滤
	BroadcastTag(sourcePlayer basis.IPlayerEntity, tag string, handler func(player basis.IPlayerEntity)) error
	// BroadcastTags 根据tags广播房间
	// tags 指定的房间 Id 集合
	// tagAnd true: 要求包含全部Tag: false: 要求包含其中一个tag
	BroadcastTags(sourcePlayer basis.IPlayerEntity, tags []string, tagAnd bool, handler func(player basis.IPlayerEntity)) error
}

func NewIBroadcastManager(entityMgr IEntityManager) IBroadcastManager {
	return NewBroadcastManager(entityMgr)
}

func NewBroadcastManager(entityMgr IEntityManager) *BroadcastManager {
	return &BroadcastManager{entityMgr: entityMgr, distance: math.MaxInt32}
}

//----------------------------------

type BroadcastManager struct {
	entityMgr IEntityManager
	lock      sync.RWMutex
	distance  int32
}

func (o *BroadcastManager) InitManager() {
	return
}

func (o *BroadcastManager) DisposeManager() {
	return
}

func (o *BroadcastManager) SetNearDistance(distance int32) {
	o.lock.Lock()
	defer o.lock.Unlock()
	o.distance = distance
}

func (o *BroadcastManager) BroadcastCurrent(sourcePlayer basis.IPlayerEntity, excludeBlack bool, handler func(player basis.IPlayerEntity)) error {
	o.lock.RLock()
	defer o.lock.RUnlock()
	if nil == sourcePlayer {
		return errors.New(fmt.Sprintf("Source is nil. "))
	}
	roomId := sourcePlayer.RoomId()
	if room, ok1 := o.entityMgr.GetRoom(roomId); ok1 {
		if ec, ok2 := room.(basis.IEntityContainer); ok2 {
			ec.ForEachChildByType(basis.EntityPlayer, func(child basis.IEntity) {
				if checkSame(sourcePlayer, child) { //本身
					return
				}
				if excludeBlack && checkBlack(sourcePlayer, child) { //黑名单
					return
				}
				if playerChild, ok := child.(basis.IPlayerEntity); ok {
					if !basis.NearXYZ(sourcePlayer.Position(), playerChild.Position(), o.distance) { //位置不相近
						return
					}
					handler(playerChild)
				}
			}, false)
		}
	}
	return nil
}

func (o *BroadcastManager) BroadcastPlayers(sourcePlayer basis.IPlayerEntity, players []string, handler func(player basis.IPlayerEntity)) error {
	if len(players) == 0 {
		return errors.New("Len of target list is 0 ")
	}
	o.lock.RLock()
	defer o.lock.RUnlock()
	playerIndex := o.entityMgr.PlayerIndex()
	for _, targetId := range players {
		if targetPlayer, ok := playerIndex.GetPlayer(targetId); ok { //目标用户存在
			if nil != sourcePlayer {
				if checkSame(sourcePlayer, targetPlayer) { //本身
					continue
				}
				if checkBlack(sourcePlayer, targetPlayer) { //黑名单
					continue
				}
			}
			handler(targetPlayer)
		}
	}
	return nil
}

func (o *BroadcastManager) BroadcastRoom(sourcePlayer basis.IPlayerEntity, room basis.IRoomEntity, handler func(player basis.IPlayerEntity)) error {
	if nil == room {
		return errors.New("Target is nil. ")
	}
	o.lock.RLock()
	defer o.lock.RUnlock()
	if c, ok := room.(basis.IEntityContainer); ok {
		if nil == sourcePlayer {
			c.ForEachChildByType(basis.EntityPlayer, func(child basis.IEntity) {
				o.toChild(child, handler)
			}, false)
		} else {
			c.ForEachChildByType(basis.EntityPlayer, func(child basis.IEntity) {
				o.sourceToChild(sourcePlayer, child, handler)
			}, false)
		}
	}
	return nil
}

func (o *BroadcastManager) BroadcastRooms(sourcePlayer basis.IPlayerEntity, rooms []string, handler func(player basis.IPlayerEntity)) error {
	if len(rooms) == 0 {
		return errors.New("Len of room list is 0 ")
	}
	o.lock.RLock()
	defer o.lock.RUnlock()
	roomIndex := o.entityMgr.RoomIndex()
	for index := range rooms {
		if room, ok1 := roomIndex.GetRoom(rooms[index]); ok1 {
			if c, ok2 := room.(basis.IEntityContainer); ok2 {
				if sourcePlayer != nil {
					c.ForEachChildByType(basis.EntityPlayer, func(child basis.IEntity) {
						o.toChild(child, handler)
					}, false)
				} else {
					c.ForEachChildByType(basis.EntityPlayer, func(child basis.IEntity) {
						o.sourceToChild(sourcePlayer, child, handler)
					}, false)
				}
			}
		}
	}
	return nil
}

func (o *BroadcastManager) BroadcastZone(sourcePlayer basis.IPlayerEntity, zoneId string, handler func(player basis.IPlayerEntity)) error {
	return o.BroadcastTag(sourcePlayer, zoneId, handler)
}

func (o *BroadcastManager) BroadcastZones(sourcePlayer basis.IPlayerEntity, zones []string, handler func(player basis.IPlayerEntity)) error {
	return o.BroadcastTags(sourcePlayer, zones, false, handler)
}

func (o *BroadcastManager) BroadcastWorld(sourcePlayer basis.IPlayerEntity, worldId string, handler func(player basis.IPlayerEntity)) error {
	return o.BroadcastTag(sourcePlayer, worldId, handler)
}

func (o *BroadcastManager) BroadcastWorlds(sourcePlayer basis.IPlayerEntity, worlds []string, handler func(player basis.IPlayerEntity)) error {
	return o.BroadcastTags(sourcePlayer, worlds, false, handler)
}

func (o *BroadcastManager) BroadcastTag(sourcePlayer basis.IPlayerEntity, tag string, handler func(player basis.IPlayerEntity)) error {
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
			if sourcePlayer != nil {
				c.ForEachChildByType(basis.EntityPlayer, func(child basis.IEntity) {
					o.toChild(child, handler)
				}, false)
			} else {
				c.ForEachChildByType(basis.EntityPlayer, func(child basis.IEntity) {
					o.sourceToChild(sourcePlayer, child, handler)
				}, false)
			}
		}
	})
	return nil
}

func (o *BroadcastManager) BroadcastTags(sourcePlayer basis.IPlayerEntity, tags []string, tagAnd bool, handler func(player basis.IPlayerEntity)) error {
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
			if sourcePlayer != nil {
				c.ForEachChildByType(basis.EntityPlayer, func(child basis.IEntity) {
					o.toChild(child, handler)
				}, false)
			} else {
				c.ForEachChildByType(basis.EntityPlayer, func(child basis.IEntity) {
					o.sourceToChild(sourcePlayer, child, handler)
				}, false)
			}
		}
	})
	return nil
}

func (o *BroadcastManager) BroadcastEntity(sourcePlayer basis.IPlayerEntity, entity basis.IEntityContainer, handler func(entity basis.IPlayerEntity)) error {
	o.lock.RLock()
	defer o.lock.RUnlock()
	if nil == entity {
		return errors.New(fmt.Sprintf("Entity is nil. "))
	}
	if nil == sourcePlayer {
		entity.ForEachChildByType(basis.EntityPlayer, func(child basis.IEntity) {
			o.toChild(child, handler)
		}, false)
	} else {
		entity.ForEachChildByType(basis.EntityPlayer, func(child basis.IEntity) {
			o.sourceToChild(sourcePlayer, child, handler)
		}, false)
	}
	return nil
}

//-----------------------------

func (o *BroadcastManager) toChild(child basis.IEntity, handler func(player basis.IPlayerEntity)) {
	if playerChild, ok := child.(basis.IPlayerEntity); ok {
		handler(playerChild)
	}
}

func (o *BroadcastManager) sourceToChild(sourcePlayer basis.IPlayerEntity, child basis.IEntity, handler func(player basis.IPlayerEntity)) {
	if checkSame(sourcePlayer, child) { //本身
		return
	}
	if checkBlack(sourcePlayer, child) { //黑名单
		return
	}
	if playerChild, ok := child.(basis.IPlayerEntity); ok {
		handler(playerChild)
	}
}

func checkSame(source basis.IPlayerEntity, target basis.IEntity) bool {
	return source.UID() == target.UID()
}

func checkBlack(source basis.IPlayerEntity, target basis.IEntity) bool {
	if source.OnBlack(target.UID()) { //source黑名单
		return true
	}
	if playerTarget, ok := target.(basis.IPlayerEntity); ok { //target黑名单
		return playerTarget.OnBlack(source.UID())
	}
	return false
}
