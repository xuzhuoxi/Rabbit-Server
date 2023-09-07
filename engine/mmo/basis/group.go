// Package basis
// Created by xuzhuoxi
// on 2019-03-14.
// @author xuzhuoxi
package basis

type ITeamGroup interface {
	// TeamList 队伍列表
	TeamList() []string
	// ContainTeam 检查队伍存在性
	ContainTeam(roomId string) bool
	// AddTeam 添加房间
	AddTeam(roomId string) error
	// RemoveTeam 移除房间
	RemoveTeam(roomId string) error
}

type IRoomGroup interface {
	// RoomList 房间列表
	RoomList() []string
	// ContainRoom 检查房间存在性
	ContainRoom(roomId string) bool
	// AddRoom 添加房间
	AddRoom(roomId string) error
	// RemoveRoom 移除房间
	RemoveRoom(roomId string) error
}

type IPlayerGroup interface {
	// PlayerList 玩家列表
	PlayerList() []string
	// ContainsPlayer 检查玩家
	ContainsPlayer(playerId string) bool
	// AcceptPlayer 加入玩家,进行唯一性检查
	AcceptPlayer(playerId string) error
	// DropPlayer 从组中移除玩家
	DropPlayer(playerId string) error
}

// IEntityGroup 组
type IEntityGroup interface {
	// EntityType 接纳实体的类型
	EntityType() EntityType
	// MaxLen 最大实例数
	MaxLen() int
	// Len 实体数量
	Len() int
	// IsFull 实体已满
	IsFull() bool

	// Entities 包含实体id
	Entities() []string
	// CopyEntities 包含实体id
	CopyEntities() []string
	// ContainEntity 检查实体是否属于当前组
	ContainEntity(entityId string) bool

	// Accept 加入实体到组,进行唯一性检查
	Accept(entity string) error
	// AcceptMulti 加入实体到组,进行唯一性检查
	AcceptMulti(entityId []string) (count int, err error)
	// Drop 从组中移除实体
	Drop(entityId string) error
	// DropMulti 从组中移除实体
	DropMulti(entityId []string) (count int, err error)
}
