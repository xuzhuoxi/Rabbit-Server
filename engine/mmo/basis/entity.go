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
	// EntityUnit 单位实体
	EntityUnit EntityType = 1 << iota
	// EntityPlayer 用户实体
	EntityPlayer
	// EntityRoom 房间实体
	EntityRoom
	// EntityTeam 队伍实体
	EntityTeam
	// EntityTeamCorps 军团实体
	EntityTeamCorps
	// EntityChannel 频道实体
	EntityChannel

	// EntityNone 不是实体
	EntityNone EntityType = 0
	// EntityAll 全部实体
	EntityAll EntityType = EntityUnit | EntityRoom | EntityPlayer | EntityTeamCorps | EntityTeam | EntityChannel
)

var (
	entityNames = make(map[EntityType]string)
	entities    = []EntityType{EntityRoom, EntityPlayer, EntityTeamCorps, EntityTeam, EntityChannel}
)

func init() {
	entityNames[EntityUnit] = "Unit"
	entityNames[EntityPlayer] = "Player"
	entityNames[EntityRoom] = "Room"
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
	// EntityType 实体类型
	EntityType() EntityType
}

type INameEntity interface {
	// Name 昵称，显示使用
	Name() string
}

type IInitEntity interface {
	// InitEntity 初始化实体
	InitEntity()
	// DestroyEntity 释放实体
	DestroyEntity()
}

type IUnitEntity interface {
	IEntity
	IInitEntity
	IVariableSupport
	// Position 取坐标
	Position() (pos XYZ)
	// SetPosition 设置坐标
	SetPosition(pos XYZ, notify bool)
	// Owner 拥有者
	Owner() string
	// SetOwner 设置拥有者
	SetOwner(owner string, notify bool)
	// RoomId 房间Id
	RoomId() string
}

// IPlayerEntity 用户实体
type IPlayerEntity interface {
	IEntity
	INameEntity
	IInitEntity
	IPlayerSubscriber
	IVariableSupport
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
	INameEntity
	IInitEntity
	IVariableSupport
	ITagsSupport

	IEntityContainer
	IUnitContainer
	// PlayerCount 玩家数量
	PlayerCount() int
	// Players 全部玩家
	Players() []IPlayerEntity
	// UnitCount 单位数量
	UnitCount() int
	// UnitIndex 单位索引
	UnitIndex() IUnitIndex
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
	INameEntity
	IInitEntity
	IChannelBehavior
	IVariableSupport
}

func EntityEqual(entity1 IEntity, entity2 IEntity) bool {
	return nil != entity1 && nil != entity2 && entity1.UID() == entity2.UID() && entity1.EntityType() == entity2.EntityType()
}
