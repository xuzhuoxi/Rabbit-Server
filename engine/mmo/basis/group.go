// Package basis
// Created by xuzhuoxi
// on 2019-03-14.
// @author xuzhuoxi
package basis

//type IZoneGroup interface {
//	// ZoneList 区域列表
//	ZoneList() []string
//	// ContainZone 检查区域存在性
//	ContainZone(zoneId string) bool
//	// AddZone 添加区域
//	AddZone(zoneId string) error
//	// RemoveZone 移除区域
//	RemoveZone(zoneId string) error
//}

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

type IUserGroup interface {
	// UserList 用户列表
	UserList() []string
	// ContainUser 检查用户
	ContainUser(userId string) bool
	// AcceptUser 加入用户,进行唯一性检查
	AcceptUser(userId string) error
	// DropUser 从组中移除用户
	DropUser(userId string) error
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
