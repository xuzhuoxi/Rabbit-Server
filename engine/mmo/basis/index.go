// Package basis
// Created by xuzhuoxi
// on 2019-03-14.
// @author xuzhuoxi
package basis

type IAbsEntityIndex interface {
	// EntityType 实体类型
	EntityType() EntityType
	// ForEachEntity 遍历全部实体
	ForEachEntity(each func(entity IEntity))
}

type IEntityIndex interface {
	IAbsEntityIndex
	// Check 检查存在
	Check(id string) bool
	// Get 获取one
	Get(id string) (entity IEntity)
	// Add 添加
	// errNum=1: entity=nil
	// errNum=2: EntityType 不匹配
	// errNum=3: entity.UID 重复
	Add(entity IEntity) (errNum int, err error)
	// Remove 从索引中移除
	// errNum=1: 找不到id索引
	Remove(id string) (entity IEntity, errNum int, err error)
	// Update 更新
	// errNum=1: entity=nil
	// errNum=2: EntityType 不匹配
	Update(entity IEntity) (errNum int, err error)
}

// IUserIndex 玩家索引
type IUserIndex interface {
	IAbsEntityIndex
	// CheckUser 检查User是否存在
	CheckUser(userId string) bool
	// GetUser 获取User
	GetUser(userId string) (user IUserEntity, ok bool)
	// AddUser 添加一个新User到索引中
	AddUser(user IUserEntity) (rsCode int32, err error)
	// RemoveUser 从索引中移除一个User
	RemoveUser(userId string) (user IUserEntity, rsCode int32, err error)
	// UpdateUser 从索引中更新一个User
	UpdateUser(user IUserEntity) (rsCode int32, err error)
}

// IRoomIndex 房间索引
type IRoomIndex interface {
	IAbsEntityIndex
	// CheckRoom 检查Room是否存在
	CheckRoom(roomId string) bool
	// GetRoom 获取Room
	GetRoom(roomId string) (room IRoomEntity, ok bool)
	// AddRoom 添加一个新Room到索引中
	AddRoom(room IRoomEntity) (rsCode int32, err error)
	// RemoveRoom 从索引中移除一个Room
	RemoveRoom(roomId string) (room IRoomEntity, rsCode int32, err error)
	// UpdateRoom 从索引中更新一个Room
	UpdateRoom(room IRoomEntity) (rsCode int32, err error)
}

// ITeamIndex 队伍索引
type ITeamIndex interface {
	IAbsEntityIndex
	// CheckTeam 检查Team是否存在
	CheckTeam(teamId string) bool
	// GetTeam 获取Team
	GetTeam(teamId string) (team ITeamEntity, ok bool)
	// AddTeam 添加一个新Team到索引中
	AddTeam(team ITeamEntity) (rsCode int32, err error)
	// RemoveTeam 从索引中移除一个Team
	RemoveTeam(teamId string) (team ITeamEntity, rsCode int32, err error)
	// UpdateTeam 从索引中更新一个Team
	UpdateTeam(team ITeamEntity) (rsCode int32, err error)
}

type ITeamCorpsIndex interface {
	IAbsEntityIndex
	// CheckCorps 检查Corps是否存在
	CheckCorps(corpsId string) bool
	// GetCorps 获取Corps
	GetCorps(corpsId string) (corps ITeamCorpsEntity, ok bool)
	// AddCorps 添加一个新Corps到索引中
	AddCorps(corps ITeamCorpsEntity) (rsCode int32, err error)
	// RemoveCorps 从索引中移除一个Corps
	RemoveCorps(corpsId string) (corps ITeamCorpsEntity, rsCode int32, err error)
	// UpdateCorps 更新一个新Corps到索引中
	UpdateCorps(corps ITeamCorpsEntity) (rsCode int32, err error)
}

// IChannelIndex 频道索引
type IChannelIndex interface {
	IAbsEntityIndex
	// CheckChannel 检查Channel是否存在
	CheckChannel(chanId string) bool
	// GetChannel 获取Channel
	GetChannel(chanId string) (channel IChannelEntity, ok bool)
	// AddChannel 从索引中增加一个Channel
	AddChannel(channel IChannelEntity) (rsCode int32, err error)
	// RemoveChannel 从索引中移除一个Channel
	RemoveChannel(chanId string) (channel IChannelEntity, rsCode int32, err error)
	// UpdateChannel 从索引中更新一个Channel
	UpdateChannel(channel IChannelEntity) (rsCode int32, err error)
}
