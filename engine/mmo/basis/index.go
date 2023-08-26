// Package basis
// Created by xuzhuoxi
// on 2019-03-14.
// @author xuzhuoxi
package basis

type IEntityIndex interface {
	EntityType() EntityType
	// Check 检查存在
	Check(id string) bool
	// Get 获取one
	Get(id string) IEntity
	// Add 添加
	Add(entity IEntity) error
	// Remove 从索引中移除
	Remove(id string) (IEntity, error)
	// Update 更新
	Update(entity IEntity) error
}

type IWorldIndex interface {
	IEntityIndex
	// CheckWorld 检查World是否存在
	CheckWorld(worldId string) bool
	// GetWorld 获取World
	GetWorld(worldId string) IWorldEntity
	// AddWorld 添加一个新World到索引中
	AddWorld(world IWorldEntity) error
	// RemoveWorld 从索引中移除一个World
	RemoveWorld(worldId string) (IWorldEntity, error)
	// UpdateWorld 更新一个新World到索引中
	UpdateWorld(zone IWorldEntity) error
}

type IZoneIndex interface {
	IEntityIndex
	// CheckZone 检查Zone是否存在
	CheckZone(zoneId string) bool
	// GetZone 获取Zone
	GetZone(zoneId string) IZoneEntity
	// AddZone 添加一个新Zone到索引中
	AddZone(zone IZoneEntity) error
	// RemoveZone 从索引中移除一个Zone
	RemoveZone(zoneId string) (IZoneEntity, error)
	// UpdateZone 更新一个新Zone到索引中
	UpdateZone(zone IZoneEntity) error
}

// IRoomIndex 房间索引
type IRoomIndex interface {
	IEntityIndex
	// CheckRoom 检查Room是否存在
	CheckRoom(roomId string) bool
	// GetRoom 获取Room
	GetRoom(roomId string) IRoomEntity
	// AddRoom 添加一个新Room到索引中
	AddRoom(room IRoomEntity) error
	// RemoveRoom 从索引中移除一个Room
	RemoveRoom(roomId string) (IRoomEntity, error)
	// UpdateRoom 从索引中更新一个Room
	UpdateRoom(room IRoomEntity) error
}

type ITeamCorpsIndex interface {
	IEntityIndex
	// CheckCorps 检查Corps是否存在
	CheckCorps(corpsId string) bool
	// GetCorps 获取Corps
	GetCorps(corpsId string) ITeamCorpsEntity
	// AddCorps 添加一个新Corps到索引中
	AddCorps(corps ITeamCorpsEntity) error
	// RemoveCorps 从索引中移除一个Corps
	RemoveCorps(corpsId string) (ITeamCorpsEntity, error)
	// UpdateCorps 更新一个新Corps到索引中
	UpdateCorps(corps ITeamCorpsEntity) error
}

// ITeamIndex 队伍索引
type ITeamIndex interface {
	IEntityIndex
	// CheckTeam 检查Team是否存在
	CheckTeam(teamId string) bool
	// GetTeam 获取Team
	GetTeam(teamId string) ITeamEntity
	// AddTeam 添加一个新Team到索引中
	AddTeam(team ITeamEntity) error
	// RemoveTeam 从索引中移除一个Team
	RemoveTeam(teamId string) (ITeamEntity, error)
	// UpdateTeam 从索引中更新一个Team
	UpdateTeam(team ITeamEntity) error
}

// IUserIndex 玩家索引
type IUserIndex interface {
	IEntityIndex
	// CheckUser 检查User是否存在
	CheckUser(userId string) bool
	// GetUser 获取User
	GetUser(userId string) IUserEntity
	// AddUser 添加一个新User到索引中
	AddUser(user IUserEntity) error
	// RemoveUser 从索引中移除一个User
	RemoveUser(userId string) (IUserEntity, error)
	// UpdateUser 从索引中更新一个User
	UpdateUser(user IUserEntity) error
}

// IChannelIndex 频道索引
type IChannelIndex interface {
	IEntityIndex
	// CheckChannel 检查Channel是否存在
	CheckChannel(chanId string) bool
	// GetChannel 获取Channel
	GetChannel(chanId string) IChannelEntity
	// AddChannel 从索引中增加一个Channel
	AddChannel(channel IChannelEntity) error
	// RemoveChannel 从索引中移除一个Channel
	RemoveChannel(chanId string) (IChannelEntity, error)
	// UpdateChannel 从索引中更新一个Channel
	UpdateChannel(channel IChannelEntity) error
}
