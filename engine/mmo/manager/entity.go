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
	"github.com/xuzhuoxi/Rabbit-Server/engine/mmo/config"
	"github.com/xuzhuoxi/Rabbit-Server/engine/mmo/entity"
	"github.com/xuzhuoxi/Rabbit-Server/engine/mmo/index"
	"github.com/xuzhuoxi/infra-go/encodingx"
	"github.com/xuzhuoxi/infra-go/eventx"
	"github.com/xuzhuoxi/infra-go/logx"
	"sync"
)

type IEntityFactory interface {
	// CreateRoom 构造房间
	CreateRoom(roomId string, roomName string, tags []string, vars encodingx.IKeyValue) (room basis.IRoomEntity, rsCode int32, err error)

	// CreateUser 创建用户实体
	CreateUser(userId string, vars encodingx.IKeyValue) (user basis.IUserEntity, rsCode int32, err error)
	// CreateTeam 创建队伍
	CreateTeam(userId string, vars encodingx.IKeyValue) (team basis.ITeamEntity, rsCode int32, err error)
	// CreateTeamCorps 创建团队
	CreateTeamCorps(teamId string, vars encodingx.IKeyValue) (corps basis.ITeamCorpsEntity, rsCode int32, err error)
	// CreateChannel 构造频道
	CreateChannel(chanId string, chanName string, vars encodingx.IKeyValue) (channel basis.IChannelEntity, rsCode int32, err error)

	// DestroyEntity 删除实体
	DestroyEntity(entity basis.IEntity) (rsCode int32, err error)
	// DestroyEntityBy 通过类型和Id删除实体
	DestroyEntityBy(entityType basis.EntityType, eId string) (entity basis.IEntity, rsCode int32, err error)
}

type IEntityIndexSet interface {
	RoomIndex() basis.IRoomIndex
	UserIndex() basis.IUserIndex
	TeamIndex() basis.ITeamIndex
	TeamCorpsIndex() basis.ITeamCorpsIndex
	ChannelIndex() basis.IChannelIndex
	GetEntityIndex(entityType basis.EntityType) basis.IAbsEntityIndex
}

type IEntityGetter interface {
	// GetRoom 获取房间实例
	GetRoom(roomId string) (room basis.IRoomEntity, ok bool)
	// GetUser 获取用户实例
	GetUser(userId string) (user basis.IUserEntity, ok bool)
	// GetTeam 获取队伍实例
	GetTeam(teamId string) (team basis.ITeamEntity, ok bool)
	// GetTeamCorps 获取队伍实例
	GetTeamCorps(corpsId string) (corps basis.ITeamCorpsEntity, ok bool)
	// GetChannel 获取频道实例
	GetChannel(chanId string) (channel basis.IChannelEntity, ok bool)
	// GetEntity 获取实例
	GetEntity(entityType basis.EntityType, entityId string) (entity basis.IEntity, ok bool)
}

type IEntityIterator interface {
	// ForEachRoom 遍历每个房间实体
	ForEachRoom(each func(room basis.IRoomEntity))
	// ForEachUser 遍历每个用户实体
	ForEachUser(each func(user basis.IUserEntity))
	// ForEachTeam 遍历每个队伍实体
	ForEachTeam(each func(team basis.ITeamEntity))
	// ForEachTeamCorps 遍历每个军团实体
	ForEachTeamCorps(each func(corps basis.ITeamCorpsEntity))
	// ForEachChannel 遍历每个频道实体
	ForEachChannel(each func(channel basis.IChannelEntity))
}

type IEntityManager interface {
	eventx.IEventDispatcher
	IEntityFactory
	IEntityGetter
	IEntityIterator
	IEntityIndexSet
	basis.IManagerBase

	// BuildEnv 构建MMO环境
	BuildEnv(cfg *config.MMOConfig) error
}

func NewIEntityManager() IEntityManager {
	return NewEntityManager()
}

func NewEntityManager() IEntityManager {
	rs := &EntityManager{logger: logx.DefaultLogger()}
	rs.roomIndex = index.NewIRoomIndex()
	rs.userIndex = index.NewIUserIndex()
	rs.teamIndex = index.NewITeamIndex()
	rs.teamCorpsIndex = index.NewITeamCorpsIndex()
	rs.chanIndex = index.NewIChannelIndex()
	return rs
}

//----------------------------

type EntityManager struct {
	roomIndex      basis.IRoomIndex // 协程安全
	roomLock       sync.RWMutex
	userIndex      basis.IUserIndex // 协程安全
	userLock       sync.RWMutex
	teamIndex      basis.ITeamIndex // 协程安全
	teamLock       sync.RWMutex
	teamCorpsIndex basis.ITeamCorpsIndex // 协程安全
	teamCorpsLock  sync.RWMutex
	chanIndex      basis.IChannelIndex // 协程安全
	chanLock       sync.RWMutex

	logger logx.ILogger
	eventx.EventDispatcher
}

func (o *EntityManager) InitManager() {
	return
}

func (o *EntityManager) DisposeManager() {
	return
}

func (o *EntityManager) SetLogger(logger logx.ILogger) {
	o.logger = logger
}

func (o *EntityManager) BuildEnv(cfg *config.MMOConfig) error {
	o.roomLock.Lock()
	defer o.roomLock.Unlock()
	for _, room := range cfg.Entities.Rooms {
		_, _, err1 := o.createRoom(room.Id, room.Name, room.Tags, nil)
		if nil != err1 {
			return err1
		}
	}
	return nil
}

func (o *EntityManager) CreateRoom(roomId string, roomName string, tags []string, vars encodingx.IKeyValue) (room basis.IRoomEntity, rsCode int32, err error) {
	o.roomLock.Lock()
	defer o.roomLock.Unlock()
	return o.createRoom(roomId, roomName, tags, vars)
}

func (o *EntityManager) createRoom(roomId string, roomName string, tags []string, vars encodingx.IKeyValue) (room basis.IRoomEntity, rsCode int32, err error) {
	if o.roomIndex.CheckRoom(roomId) {
		return nil, basis.CodeMMORoomExist, errors.New("EntityManager.CreateRoomAt Error: RoomId(" + roomId + ") Duplicate")
	}
	room = entity.NewIRoomEntity(roomId, roomName)
	room.InitEntity()
	room.SetVars(vars)
	room.SetTags(tags)
	rsCode, err = o.roomIndex.AddRoom(room)
	if nil != err {
		return nil, rsCode, err
	}
	o.addEntityEventListener(room)
	return room, 0, nil
}

func (o *EntityManager) CreateUser(userId string, vars encodingx.IKeyValue) (user basis.IUserEntity, rsCode int32, err error) {
	o.userLock.Lock()
	defer o.userLock.Unlock()
	if userId == "" || o.userIndex.CheckUser(userId) {
		return nil, basis.CodeMMOUserExist, errors.New(fmt.Sprintf("EntityManager.CreateUser Error: User(%s) is nil or exist", userId))
	}
	user = entity.NewIUserEntity(userId)
	user.SetVars(vars)
	rsCode, err = o.userIndex.AddUser(user)
	if nil != err {
		return nil, rsCode, err
	}
	o.addEntityEventListener(user)
	return user, 0, nil
}

func (o *EntityManager) CreateTeam(userId string, vars encodingx.IKeyValue) (team basis.ITeamEntity, rsCode int32, err error) {
	o.teamLock.Lock()
	defer o.teamLock.Unlock()
	_, okUser := o.userIndex.GetUser(userId)
	if !okUser {
		return nil, basis.CodeMMOTeamExist, errors.New(fmt.Sprintf("EntityManager.CreateTeam Error: User(%s) does not exist", userId))
	}
	team = entity.NewITeamEntity(basis.GetTeamId(), basis.TeamName, basis.MaxTeamMember)
	team.InitEntity()
	team.SetVars(vars)
	rsCode, err = o.teamIndex.AddTeam(team)
	if nil != err {
		return nil, rsCode, err
	}
	o.addEntityEventListener(team)
	return team, 0, nil
}

func (o *EntityManager) CreateTeamCorps(teamId string, vars encodingx.IKeyValue) (corps basis.ITeamCorpsEntity, rsCode int32, err error) {
	o.teamCorpsLock.Lock()
	defer o.teamCorpsLock.Unlock()
	_, okTeam := o.teamIndex.GetTeam(teamId)
	if !okTeam {
		return nil, basis.CodeMMOTeamNotExist, errors.New(fmt.Sprintf("EntityManager.CreateTeamCorps Error: Team(%s) does not exist", teamId))
	}
	teamCorps := entity.NewITeamCorpsEntity(basis.GetTeamCorpsId(), basis.TeamCorpsName)
	teamCorps.InitEntity()
	teamCorps.SetVars(vars)
	rsCode, err = o.teamCorpsIndex.AddCorps(teamCorps)
	if nil != err {
		return nil, rsCode, err
	}
	o.addEntityEventListener(teamCorps)
	return teamCorps, 0, nil
}

func (o *EntityManager) CreateChannel(chanId string, chanName string, vars encodingx.IKeyValue) (channel basis.IChannelEntity, rsCode int32, err error) {
	o.chanLock.Lock()
	defer o.chanLock.Unlock()
	if o.chanIndex.CheckChannel(chanId) {
		return nil, basis.CodeMMOChanExist, errors.New("EntityManager.CreateChannel Error: ChanId(" + chanId + ") Duplicate!")
	}
	channel = entity.NewIChannelEntity(chanId, chanName)
	channel.InitEntity()
	channel.SetVars(vars)
	rsCode, err = o.chanIndex.AddChannel(channel)
	if nil != err {
		return nil, rsCode, err
	}
	o.addEntityEventListener(channel)
	return channel, 0, nil
}

func (o *EntityManager) DestroyEntity(entity basis.IEntity) (rsCode int32, err error) {
	if nil == entity {
		return basis.CodeMMOOther, errors.New("DestroyEntity Error at: entity is nil. ")
	}
	_, rsCode, err = o.DestroyEntityBy(entity.EntityType(), entity.UID())
	return
}

func (o *EntityManager) DestroyEntityBy(entityType basis.EntityType, eId string) (entity basis.IEntity, rsCode int32, err error) {
	switch entityType {
	case basis.EntityRoom:
		o.roomLock.Lock()
		defer o.roomLock.Unlock()
		entity, rsCode, err = o.roomIndex.RemoveRoom(eId)
	case basis.EntityUser:
		o.userLock.Lock()
		defer o.userLock.Unlock()
		entity, rsCode, err = o.userIndex.RemoveUser(eId)
	case basis.EntityTeamCorps:
		o.teamCorpsLock.Lock()
		defer o.teamCorpsLock.Unlock()
		entity, rsCode, err = o.teamCorpsIndex.RemoveCorps(eId)
	case basis.EntityTeam:
		o.teamLock.Lock()
		defer o.teamLock.Unlock()
		entity, rsCode, err = o.teamIndex.RemoveTeam(eId)
	case basis.EntityChannel:
		o.chanLock.Lock()
		defer o.chanLock.Unlock()
		entity, rsCode, err = o.chanIndex.RemoveChannel(eId)
	}
	if nil != entity {
		o.removeEntityEventListener(entity)
	}
	return
}

func (o *EntityManager) addEntityEventListener(entity basis.IEntity) {
	if dispatcher, ok := entity.(basis.IVariableSupport); ok {
		dispatcher.AddEventListener(basis.EventEntityVarChanged, o.onEntityVar)
		dispatcher.AddEventListener(basis.EventEntityVarsChanged, o.onEntityVars)
	}
}

func (o *EntityManager) removeEntityEventListener(entity basis.IEntity) {
	if dispatcher, ok := entity.(basis.IVariableSupport); ok {
		dispatcher.RemoveEventListener(basis.EventEntityVarsChanged, o.onEntityVars)
		dispatcher.RemoveEventListener(basis.EventEntityVarChanged, o.onEntityVar)
	}
}

// 事件转发: 单个量变量更新
func (o *EntityManager) onEntityVar(evd *eventx.EventData) {
	evd.StopImmediatePropagation()
	o.DispatchEvent(basis.EventManagerVarChanged, o, evd.Data)
}

// 事件转发: 批量变量更新
func (o *EntityManager) onEntityVars(evd *eventx.EventData) {
	evd.StopImmediatePropagation()
	o.DispatchEvent(basis.EventManagerVarsChanged, o, evd.Data)
}

//----------------------------

func (o *EntityManager) GetRoom(roomId string) (room basis.IRoomEntity, ok bool) {
	if len(roomId) == 0 {
		return
	}
	return o.roomIndex.GetRoom(roomId)
}

func (o *EntityManager) GetUser(userId string) (user basis.IUserEntity, ok bool) {
	if len(userId) == 0 {
		return
	}
	return o.userIndex.GetUser(userId)
}

func (o *EntityManager) GetTeam(teamId string) (team basis.ITeamEntity, ok bool) {
	if len(teamId) == 0 {
		return
	}
	return o.teamIndex.GetTeam(teamId)
}

func (o *EntityManager) GetTeamCorps(corpsId string) (corps basis.ITeamCorpsEntity, ok bool) {
	if len(corpsId) == 0 {
		return
	}
	return o.teamCorpsIndex.GetCorps(corpsId)
}

func (o *EntityManager) GetChannel(chanId string) (channel basis.IChannelEntity, ok bool) {
	if len(chanId) == 0 {
		return
	}
	return o.chanIndex.GetChannel(chanId)
}

func (o *EntityManager) GetEntity(entityType basis.EntityType, eId string) (entity basis.IEntity, ok bool) {
	switch entityType {
	case basis.EntityRoom:
		entity, ok = o.roomIndex.GetRoom(eId)
	case basis.EntityUser:
		entity, ok = o.userIndex.GetUser(eId)
	case basis.EntityTeam:
		entity, ok = o.teamIndex.GetTeam(eId)
	case basis.EntityTeamCorps:
		entity, ok = o.teamCorpsIndex.GetCorps(eId)
	case basis.EntityChannel:
		entity, ok = o.chanIndex.GetChannel(eId)
	}
	return
}

func (o *EntityManager) ForEachRoom(each func(room basis.IRoomEntity)) {
	o.roomIndex.ForEachEntity(func(entity basis.IEntity) {
		each(entity.(basis.IRoomEntity))
	})
}

func (o *EntityManager) ForEachUser(each func(user basis.IUserEntity)) {
	o.userIndex.ForEachEntity(func(entity basis.IEntity) {
		each(entity.(basis.IUserEntity))
	})
}

func (o *EntityManager) ForEachTeam(each func(team basis.ITeamEntity)) {
	o.teamIndex.ForEachEntity(func(entity basis.IEntity) {
		each(entity.(basis.ITeamEntity))
	})
}

func (o *EntityManager) ForEachTeamCorps(each func(corps basis.ITeamCorpsEntity)) {
	o.teamCorpsIndex.ForEachEntity(func(entity basis.IEntity) {
		each(entity.(basis.ITeamCorpsEntity))
	})
}

func (o *EntityManager) ForEachChannel(each func(channel basis.IChannelEntity)) {
	o.chanIndex.ForEachEntity(func(entity basis.IEntity) {
		each(entity.(basis.IChannelEntity))
	})
}

//-----------------------

func (o *EntityManager) RoomIndex() basis.IRoomIndex {
	return o.roomIndex
}

func (o *EntityManager) UserIndex() basis.IUserIndex {
	return o.userIndex
}

func (o *EntityManager) TeamIndex() basis.ITeamIndex {
	return o.teamIndex
}

func (o *EntityManager) TeamCorpsIndex() basis.ITeamCorpsIndex {
	return o.teamCorpsIndex
}

func (o *EntityManager) ChannelIndex() basis.IChannelIndex {
	return o.chanIndex
}

func (o *EntityManager) GetEntityIndex(entityType basis.EntityType) basis.IAbsEntityIndex {
	switch entityType {
	case basis.EntityRoom:
		return o.roomIndex
	case basis.EntityUser:
		return o.userIndex
	case basis.EntityTeam:
		return o.teamIndex
	case basis.EntityTeamCorps:
		return o.teamCorpsIndex
	case basis.EntityChannel:
		return o.chanIndex
	}
	return nil
}
