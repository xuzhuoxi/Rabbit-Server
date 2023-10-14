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
	"github.com/xuzhuoxi/Rabbit-Server/engine/mmo/events"
	"github.com/xuzhuoxi/Rabbit-Server/engine/mmo/index"
	"github.com/xuzhuoxi/infra-go/encodingx"
	"github.com/xuzhuoxi/infra-go/eventx"
	"github.com/xuzhuoxi/infra-go/logx"
	"sync"
)

type IEntityFactory interface {
	// CreateRoom 构造房间
	CreateRoom(roomId string, roomName string, tags []string, vars encodingx.IKeyValue) (room basis.IRoomEntity, rsCode int32, err error)
	// CreatePlayer 创建玩家实体
	CreatePlayer(playerId string, vars encodingx.IKeyValue) (player basis.IPlayerEntity, rsCode int32, err error)
	// CreateTeam 创建队伍
	CreateTeam(playerId string, vars encodingx.IKeyValue) (team basis.ITeamEntity, rsCode int32, err error)
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
	GetRoomUnitIndex(roomId string) (index basis.IUnitIndex, ok bool)
	RoomIndex() basis.IRoomIndex
	PlayerIndex() basis.IPlayerIndex
	TeamIndex() basis.ITeamIndex
	TeamCorpsIndex() basis.ITeamCorpsIndex
	ChannelIndex() basis.IChannelIndex
	GetEntityIndex(entityType basis.EntityType) basis.IAbsEntityIndex
}

type IEntityGetter interface {
	// GetUnit 获取单位实例
	GetUnit(roomId string, unitId string) (unit basis.IUnitEntity, ok bool)
	// GetRoom 获取房间实例
	GetRoom(roomId string) (room basis.IRoomEntity, ok bool)
	// GetPlayer 获取玩家实例
	GetPlayer(playerId string) (player basis.IPlayerEntity, ok bool)
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
	// ForEachUnit 遍历每个单位实体
	ForEachUnit(roomId string, each func(room basis.IUnitEntity) (interrupt bool))
	// ForEachRoom 遍历每个房间实体
	ForEachRoom(each func(room basis.IRoomEntity) (interrupt bool))
	// ForEachPlayer 遍历每个玩家实体
	ForEachPlayer(each func(player basis.IPlayerEntity) (interrupt bool))
	// ForEachTeam 遍历每个队伍实体
	ForEachTeam(each func(team basis.ITeamEntity) (interrupt bool))
	// ForEachTeamCorps 遍历每个军团实体
	ForEachTeamCorps(each func(corps basis.ITeamCorpsEntity) (interrupt bool))
	// ForEachChannel 遍历每个频道实体
	ForEachChannel(each func(channel basis.IChannelEntity) (interrupt bool))
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
	rs.playerIndex = index.NewIPlayerIndex()
	rs.teamIndex = index.NewITeamIndex()
	rs.teamCorpsIndex = index.NewITeamCorpsIndex()
	rs.chanIndex = index.NewIChannelIndex()
	return rs
}

//----------------------------

type EntityManager struct {
	roomIndex      basis.IRoomIndex      // 协程安全
	playerIndex    basis.IPlayerIndex    // 协程安全
	teamIndex      basis.ITeamIndex      // 协程安全
	teamCorpsIndex basis.ITeamCorpsIndex // 协程安全
	chanIndex      basis.IChannelIndex   // 协程安全

	buildLock sync.RWMutex
	logger    logx.ILogger
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
	o.buildLock.Lock()
	defer o.buildLock.Unlock()
	for _, room := range cfg.Entities.Rooms {
		_, _, err1 := o.CreateRoom(room.Id, room.Name, room.Tags, nil)
		if nil != err1 {
			return err1
		}
	}
	return nil
}

func (o *EntityManager) CreateRoom(roomId string, roomName string, tags []string, vars encodingx.IKeyValue) (room basis.IRoomEntity, rsCode int32, err error) {
	if o.roomIndex.CheckRoom(roomId) {
		return nil, basis.CodeMMORoomExist, errors.New("EntityManager.CreateRoomAt Error: RoomId(" + roomId + ") Duplicate")
	}
	room = entity.NewIRoomEntity(roomId, roomName)
	room.InitEntity()
	room.SetVars(vars, false)
	room.SetTags(tags)
	rsCode, err = o.roomIndex.AddRoom(room)
	if nil != err {
		return nil, rsCode, err
	}
	o.addEntityEventListener(room)
	o.DispatchEvent(events.EventRoomInit, o, room)
	return room, 0, nil
}

func (o *EntityManager) CreatePlayer(playerId string, vars encodingx.IKeyValue) (player basis.IPlayerEntity, rsCode int32, err error) {
	if playerId == "" || o.playerIndex.CheckPlayer(playerId) {
		return nil, basis.CodeMMOPlayerExist, errors.New(fmt.Sprintf("EntityManager.CreatePlayer Error: Player(%s) is nil or exist", playerId))
	}
	player = entity.NewIPlayerEntity(playerId)
	player.InitEntity()
	player.SetVars(vars, false)
	rsCode, err = o.playerIndex.AddPlayer(player)
	if nil != err {
		return nil, rsCode, err
	}
	o.addEntityEventListener(player)
	o.DispatchEvent(events.EventPlayerInit, o, player)
	return player, 0, nil
}

func (o *EntityManager) CreateTeam(playerId string, vars encodingx.IKeyValue) (team basis.ITeamEntity, rsCode int32, err error) {
	_, okPlayer := o.playerIndex.GetPlayer(playerId)
	if !okPlayer {
		return nil, basis.CodeMMOTeamExist, errors.New(fmt.Sprintf("EntityManager.CreateTeam Error: Player(%s) does not exist", playerId))
	}
	team = entity.NewITeamEntity(basis.GetTeamId(), basis.TeamName, basis.MaxTeamMember)
	team.InitEntity()
	team.SetVars(vars, false)
	rsCode, err = o.teamIndex.AddTeam(team)
	if nil != err {
		return nil, rsCode, err
	}
	o.addEntityEventListener(team)
	o.DispatchEvent(events.EventTeamInit, o, team)
	return team, 0, nil
}

func (o *EntityManager) CreateTeamCorps(teamId string, vars encodingx.IKeyValue) (corps basis.ITeamCorpsEntity, rsCode int32, err error) {
	_, okTeam := o.teamIndex.GetTeam(teamId)
	if !okTeam {
		return nil, basis.CodeMMOTeamNotExist, errors.New(fmt.Sprintf("EntityManager.CreateTeamCorps Error: Team(%s) does not exist", teamId))
	}
	teamCorps := entity.NewITeamCorpsEntity(basis.GetTeamCorpsId(), basis.TeamCorpsName)
	teamCorps.InitEntity()
	teamCorps.SetVars(vars, false)
	rsCode, err = o.teamCorpsIndex.AddCorps(teamCorps)
	if nil != err {
		return nil, rsCode, err
	}
	o.addEntityEventListener(teamCorps)
	o.DispatchEvent(events.EventTeamCorpsInit, o, teamCorps)
	return teamCorps, 0, nil
}

func (o *EntityManager) CreateChannel(chanId string, chanName string, vars encodingx.IKeyValue) (channel basis.IChannelEntity, rsCode int32, err error) {
	if o.chanIndex.CheckChannel(chanId) {
		return nil, basis.CodeMMOChanExist, errors.New("EntityManager.CreateChannel Error: ChanId(" + chanId + ") Duplicate!")
	}
	channel = entity.NewIChannelEntity(chanId, chanName)
	channel.InitEntity()
	channel.SetVars(vars, false)
	rsCode, err = o.chanIndex.AddChannel(channel)
	if nil != err {
		return nil, rsCode, err
	}
	o.addEntityEventListener(channel)
	o.DispatchEvent(events.EventChanInit, o, channel)
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
		entity, rsCode, err = o.roomIndex.RemoveRoom(eId)
		if nil != entity {
			defer o.DispatchEvent(events.EventRoomDestroy, o, entity)
		}
	case basis.EntityPlayer:
		entity, rsCode, err = o.playerIndex.RemovePlayer(eId)
		if nil != entity {
			defer o.DispatchEvent(events.EventPlayerDestroy, o, entity)
		}
	case basis.EntityTeam:
		entity, rsCode, err = o.teamIndex.RemoveTeam(eId)
		if nil != entity {
			defer o.DispatchEvent(events.EventTeamDestroy, o, entity)
		}
	case basis.EntityTeamCorps:
		entity, rsCode, err = o.teamCorpsIndex.RemoveCorps(eId)
		if nil != entity {
			defer o.DispatchEvent(events.EventTeamCorpsDestroy, o, entity)
		}
	case basis.EntityChannel:
		entity, rsCode, err = o.chanIndex.RemoveChannel(eId)
		if nil != entity {
			defer o.DispatchEvent(events.EventChanDestroy, o, entity)
		}
	}
	if nil != entity {
		if ed, ok := entity.(eventx.IEventDispatcher); ok {
			o.removeEntityEventListener(ed)
		}
		if initEntity, ok := entity.(basis.IInitEntity); ok {
			initEntity.DestroyEntity()
		}
	}
	return
}

func (o *EntityManager) addEntityEventListener(entity eventx.IEventDispatcher) {
	entity.AddEventListener(events.EventEntityVarChanged, o.onEventRedirect)
	entity.AddEventListener(events.EventEntityVarsChanged, o.onEventRedirect)
	entity.AddEventListener(events.EventUnitInit, o.onEventRedirect)
	entity.AddEventListener(events.EventUnitDestroy, o.onEventRedirect)
}

func (o *EntityManager) removeEntityEventListener(entity eventx.IEventDispatcher) {
	entity.RemoveEventListener(events.EventUnitDestroy, o.onEventRedirect)
	entity.RemoveEventListener(events.EventUnitInit, o.onEventRedirect)
	entity.RemoveEventListener(events.EventEntityVarsChanged, o.onEventRedirect)
	entity.RemoveEventListener(events.EventEntityVarChanged, o.onEventRedirect)
}

// 事件重定向
func (o *EntityManager) onEventRedirect(evd *eventx.EventData) {
	//fmt.Println("[EntityManager.onEventRedirect]", evd.EventType)
	evd.StopImmediatePropagation()
	o.DispatchEvent(evd.EventType, o, evd.Data)
}

//----------------------------

func (o *EntityManager) GetUnit(roomId string, unitId string) (unit basis.IUnitEntity, ok bool) {
	if len(roomId) == 0 || len(unitId) == 0 {
		return
	}
	index, ok1 := o.GetRoomUnitIndex(roomId)
	if !ok1 {
		return
	}
	return index.GetUnit(unitId)
}

func (o *EntityManager) GetRoom(roomId string) (room basis.IRoomEntity, ok bool) {
	if len(roomId) == 0 {
		return
	}
	return o.roomIndex.GetRoom(roomId)
}

func (o *EntityManager) GetPlayer(playerId string) (player basis.IPlayerEntity, ok bool) {
	if len(playerId) == 0 {
		return
	}
	return o.playerIndex.GetPlayer(playerId)
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
	case basis.EntityPlayer:
		entity, ok = o.playerIndex.GetPlayer(eId)
	case basis.EntityTeam:
		entity, ok = o.teamIndex.GetTeam(eId)
	case basis.EntityTeamCorps:
		entity, ok = o.teamCorpsIndex.GetCorps(eId)
	case basis.EntityChannel:
		entity, ok = o.chanIndex.GetChannel(eId)
	}
	return
}

func (o *EntityManager) ForEachUnit(roomId string, each func(room basis.IUnitEntity) (interrupt bool)) {
	index, ok := o.GetRoomUnitIndex(roomId)
	if !ok {
		return
	}
	index.ForEachEntity(func(entity basis.IEntity) (interrupt bool) {
		return each(entity.(basis.IUnitEntity))
	})
}

func (o *EntityManager) ForEachRoom(each func(room basis.IRoomEntity) (interrupt bool)) {
	o.roomIndex.ForEachEntity(func(entity basis.IEntity) (interrupt bool) {
		return each(entity.(basis.IRoomEntity))
	})
}

func (o *EntityManager) ForEachPlayer(each func(player basis.IPlayerEntity) (interrupt bool)) {
	o.playerIndex.ForEachEntity(func(entity basis.IEntity) (interrupt bool) {
		return each(entity.(basis.IPlayerEntity))
	})
}

func (o *EntityManager) ForEachTeam(each func(team basis.ITeamEntity) (interrupt bool)) {
	o.teamIndex.ForEachEntity(func(entity basis.IEntity) (interrupt bool) {
		return each(entity.(basis.ITeamEntity))
	})
}

func (o *EntityManager) ForEachTeamCorps(each func(corps basis.ITeamCorpsEntity) (interrupt bool)) {
	o.teamCorpsIndex.ForEachEntity(func(entity basis.IEntity) (interrupt bool) {
		return each(entity.(basis.ITeamCorpsEntity))
	})
}

func (o *EntityManager) ForEachChannel(each func(channel basis.IChannelEntity) (interrupt bool)) {
	o.chanIndex.ForEachEntity(func(entity basis.IEntity) (interrupt bool) {
		return each(entity.(basis.IChannelEntity))
	})
}

//-----------------------

func (o *EntityManager) GetRoomUnitIndex(roomId string) (index basis.IUnitIndex, ok bool) {
	room, ok1 := o.roomIndex.GetRoom(roomId)
	if !ok1 {
		return
	}
	return room.UnitIndex(), true
}

func (o *EntityManager) RoomIndex() basis.IRoomIndex {
	return o.roomIndex
}

func (o *EntityManager) PlayerIndex() basis.IPlayerIndex {
	return o.playerIndex
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
	case basis.EntityPlayer:
		return o.playerIndex
	case basis.EntityTeam:
		return o.teamIndex
	case basis.EntityTeamCorps:
		return o.teamCorpsIndex
	case basis.EntityChannel:
		return o.chanIndex
	}
	return nil
}
