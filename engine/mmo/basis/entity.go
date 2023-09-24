// Package basis
// Created by xuzhuoxi
// on 2019-03-14.
// @author xuzhuoxi
package basis

import (
	"strings"
)

type EntityType uint16

const (
	// EntityRoom 房间实体
	EntityRoom EntityType = 1 << iota
	// EntityPlayer 用户实体
	EntityPlayer
	// EntityTeam 队伍实体
	EntityTeam
	// EntityTeamCorps 军团实体
	EntityTeamCorps
	// EntityChannel 频道实体
	EntityChannel

	// EntityNone 不是实体
	EntityNone EntityType = 0
	// EntityAll 全部实体
	EntityAll EntityType = EntityRoom | EntityPlayer | EntityTeamCorps | EntityTeam | EntityChannel
)

var (
	entityNames = make(map[EntityType]string)
	entities    = []EntityType{EntityRoom, EntityPlayer, EntityTeamCorps, EntityTeam, EntityChannel}
)

func init() {
	entityNames[EntityRoom] = "Room"
	entityNames[EntityPlayer] = "Player"
	entityNames[EntityTeam] = "Team"
	entityNames[EntityTeamCorps] = "TeamCorps"
	entityNames[EntityChannel] = "Channel"
}

func (o EntityType) String() string {
	var names []string
	for _, t := range entities {
		if o.Match(t) {
			names = append(names, entityNames[t])
		}
	}
	if len(names) == 0 {
		return "None"
	}
	return strings.Join(names, "|")
}

func (o EntityType) Match(check EntityType) bool {
	return o&check > 0
}

func (o EntityType) Include(check EntityType) bool {
	return o&check == check
}

type IEntity interface {
	// UID 唯一标识
	UID() string
	// Name 昵称，显示使用
	Name() string
	// EntityType 实体类型
	EntityType() EntityType
}

type IInitEntity interface {
	// InitEntity 初始化实体
	InitEntity()
}

type IDestroyEntity interface {
	// DestroyEntity 释放实体
	DestroyEntity()
}

// IPlayerEntity 用户实体
type IPlayerEntity interface {
	IEntity
	IInitEntity
	IDestroyEntity
	IPlayerSubscriber
	IVariableSupport

	// NickName 用户昵称
	NickName() string
	// Position 取坐标
	Position() XYZ
	// SetPosition 设置坐标
	SetPosition(pos XYZ, notify bool)

	// RoomId 取房间Id
	RoomId() string
	// GetPrevRoomId 取上一个房间Id
	GetPrevRoomId() (roomId string, ok bool)
	// SetNextRoom 设置下一个房间Id
	SetNextRoom(roomId string)
	// ConfirmNextRoom 确认下一个房间Id
	ConfirmNextRoom(confirm bool)
	// BackToPrevRoom 回来上一个房间
	BackToPrevRoom()

	// TeamId 队伍Id
	TeamId() string
	// CorpsId 团队Id
	CorpsId() string
	// GetTeamInfo 取队伍相关信息
	GetTeamInfo() (teamId string, corpsId string)
	// SetTeam 设置队伍Id
	SetTeam(teamId string, notify bool)
	// SetCorps 设置团队Id
	SetCorps(corpsId string, notify bool)
}

// IRoomEntity 房间实体
type IRoomEntity interface {
	IEntity
	IInitEntity

	IEntityContainer
	IVariableSupport
	ITagsSupport
	Players() []IPlayerEntity
}

// ITeamEntity 队伍实体
type ITeamEntity interface {
	IEntity
	IInitEntity
	ITeamControl
	IVariableSupport
}

// ITeamCorpsEntity 兵团实体
type ITeamCorpsEntity interface {
	IEntity
	IInitEntity
	IVariableSupport
}

// IChannelEntity 频道实体
type IChannelEntity interface {
	IEntity
	IInitEntity
	IChannelBehavior
	IVariableSupport
}

func EntityEqual(entity1 IEntity, entity2 IEntity) bool {
	return nil != entity1 && nil != entity2 && entity1.UID() == entity2.UID() && entity1.EntityType() == entity2.EntityType() && entity1.Name() == entity2.Name()
}
