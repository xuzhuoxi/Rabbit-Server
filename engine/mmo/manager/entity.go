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

type IEntityCreator interface {
	// CreateWorld 构造世界
	CreateWorld(worldId string, worldName string, asRoot bool) (basis.IWorldEntity, error)
	// CreateZoneAt 构造区域
	CreateZoneAt(zoneId string, zoneName string, container basis.IEntityContainer) (basis.IZoneEntity, error)
	// CreateRoomAt 构造房间
	CreateRoomAt(roomId string, roomName string, container basis.IEntityContainer) (basis.IRoomEntity, error)

	// CreateTeam 创建队伍
	CreateTeam(userId string) (basis.ITeamEntity, error)
	// CreateTeamCorps 创建团队
	CreateTeamCorps(teamId string) (basis.ITeamCorpsEntity, error)
	// CreateChannel 构造频道
	CreateChannel(chanId string, chanName string) (basis.IChannelEntity, error)
}

type IEntityIndexSet interface {
	ZoneIndex() basis.IZoneIndex
	RoomIndex() basis.IRoomIndex
	UserIndex() basis.IUserIndex
	TeamIndex() basis.ITeamIndex
	TeamCorpsIndex() basis.ITeamCorpsIndex
	ChannelIndex() basis.IChannelIndex
	GetEntityIndex(entityType basis.EntityType) basis.IEntityIndex
}

type IEntityGetter interface {
	// GetZone 获取区域实例
	GetZone(zoneId string) (basis.IZoneEntity, bool)
	// GetRoom 获取房间实例
	GetRoom(roomId string) (basis.IRoomEntity, bool)
	// GetUser 获取用户实例
	GetUser(userId string) (basis.IUserEntity, bool)
	// GetTeam 获取队伍实例
	GetTeam(teamId string) (basis.ITeamEntity, bool)
	// GetTeamCorps 获取队伍实例
	GetTeamCorps(corpsId string) (basis.ITeamCorpsEntity, bool)
	// GetChannel 获取频道实例
	GetChannel(chanId string) (basis.IChannelEntity, bool)
	// GetEntity 获取实例
	GetEntity(entityType basis.EntityType, entityId string) (basis.IEntity, bool)
}

type IEntityManager interface {
	eventx.IEventDispatcher
	IEntityCreator
	IEntityGetter
	IEntityIndexSet
	basis.IManagerBase

	World() basis.IWorldEntity
	ConstructWorlds(cfg *config.MMOConfig) error
	ConstructWorld(cfg *config.MMOConfig, worldId string) (world basis.IWorldEntity, err error)
	ConstructWorldDefault(cfg *config.MMOConfig) (world basis.IWorldEntity, err error)
}

func NewIEntityManager() IEntityManager {
	return NewEntityManager()
}

func NewEntityManager() IEntityManager {
	rs := &EntityManager{logger: logx.DefaultLogger()}
	rs.worldIndex = index.NewIWorldIndex()
	rs.zoneIndex = index.NewIZoneIndex()
	rs.roomIndex = index.NewIRoomIndex()
	rs.userIndex = index.NewIUserIndex()
	rs.teamIndex = index.NewITeamIndex()
	rs.teamCorpsIndex = index.NewITeamCorpsIndex()
	rs.channelIndex = index.NewIChannelIndex()
	return rs
}

//----------------------------

type EntityManager struct {
	worldIndex       basis.IWorldIndex
	worldIndexMu     sync.RWMutex
	zoneIndex        basis.IZoneIndex
	zoneIndexMu      sync.RWMutex
	roomIndex        basis.IRoomIndex
	roomIndexMu      sync.RWMutex
	userIndex        basis.IUserIndex
	userIndexMu      sync.RWMutex
	teamIndex        basis.ITeamIndex
	teamIndexMu      sync.RWMutex
	teamCorpsIndex   basis.ITeamCorpsIndex
	teamCorpsIndexMu sync.RWMutex
	channelIndex     basis.IChannelIndex
	chanIndexMu      sync.RWMutex

	rootWorld basis.IWorldEntity
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

func (o *EntityManager) ConstructWorlds(cfg *config.MMOConfig) error {
	if cfg.Relations == nil || len(cfg.Relations.Relations) == 0 {
		return nil
	}
	for index := range cfg.Relations.Relations {
		_, err := o.ConstructWorld(cfg, cfg.Relations.Relations[index].WorldId)
		if nil != err {
			return err
		}
	}
	return nil
}

func (o *EntityManager) ConstructWorldDefault(cfg *config.MMOConfig) (world basis.IWorldEntity, err error) {
	return o.ConstructWorld(cfg, cfg.DefaultWorld)
}

func (o *EntityManager) ConstructWorld(cfg *config.MMOConfig, worldId string) (world basis.IWorldEntity, err error) {
	if relation, ok := cfg.Relations.GetWorldRelation(worldId); ok {
		o.logger.Infoln("Start Construct World:", worldId, cfg)
		worldEntity, ok1 := cfg.Entities.FindWorld(worldId)
		if !ok1 {
			err = errors.New("Construct World Fail: " + worldId + " not configured.")
			o.logger.Warnln(err)
			return
		}
		newWorld, err1 := o.CreateWorld(worldEntity.Id, worldEntity.Name, true)
		if nil != err1 {
			err = err1
			o.logger.Warnln(err)
			return
		}
		for _, zoneCfg := range relation.Zones {
			z, ok2 := cfg.Entities.FindZone(zoneCfg.ZoneId)
			if !ok2 {
				err = errors.New("Construct Zone Fail: " + zoneCfg.ZoneId + " not configured.")
				o.logger.Warnln(err)
				return
			}
			zone, err2 := o.CreateZoneAt(z.Id, z.Name, newWorld)
			if nil != err2 {
				err = err2
				o.logger.Warnln(err)
				return
			}
			for _, roomId := range zoneCfg.Rooms {
				r, ok3 := cfg.Entities.FindRoom(roomId)
				if !ok3 {
					err = errors.New("Construct Room Fail: " + roomId + " not configured.")
					o.logger.Warnln(err)
					return
				}
				_, err3 := o.CreateRoomAt(r.Id, r.Name, zone)
				if nil != err3 {
					err = err3
					o.logger.Warnln(err)
					return
				}
			}
		}
		world = newWorld
		o.logger.Infoln("Finish Construct World:", newWorld.UID())
		return
	}
	return nil, errors.New("No World Relation: " + worldId)
}

func (o *EntityManager) CreateWorld(worldId string, worldName string, asRoot bool) (basis.IWorldEntity, error) {
	o.worldIndexMu.Lock()
	defer o.worldIndexMu.Unlock()
	if o.worldIndex.CheckWorld(worldId) {
		return nil, errors.New("EntityManager.CreateWorld Error: WorldId(" + worldId + ") Duplicate!")
	}
	world := entity.CreateWorldEntity(worldId, worldName)
	world.InitEntity()
	o.addEntityEventListener(world)
	if asRoot {
		o.rootWorld = world
	}
	return world, nil
}

func (o *EntityManager) CreateZoneAt(zoneId string, zoneName string, container basis.IEntityContainer) (basis.IZoneEntity, error) {
	o.zoneIndexMu.Lock()
	defer o.zoneIndexMu.Unlock()
	if o.zoneIndex.CheckZone(zoneId) {
		return nil, errors.New("EntityManager.CreateZoneAt Error: ZoneId(" + zoneId + ") Duplicate!")
	}
	zone := entity.NewIZoneEntity(zoneId, zoneName)
	zone.InitEntity()
	o.addEntityEventListener(zone)
	o.zoneIndex.AddZone(zone)
	if nil != container {
		if e, ok := container.(basis.IEntity); ok {
			zone.SetParent(e.UID())
			container.AddChild(zone)
		}
	}
	return zone, nil
}

func (o *EntityManager) CreateRoomAt(roomId string, roomName string, container basis.IEntityContainer) (basis.IRoomEntity, error) {
	o.roomIndexMu.Lock()
	defer o.roomIndexMu.Unlock()
	if o.roomIndex.CheckRoom(roomId) {
		return nil, errors.New("EntityManager.CreateRoomAt Error: RoomId(" + roomId + ") Duplicate")
	}
	room := entity.NewIRoomEntity(roomId, roomName)
	room.InitEntity()
	o.addEntityEventListener(room)
	o.roomIndex.AddRoom(room)
	if nil != container {
		if e, ok := container.(basis.IEntity); ok {
			room.SetParent(e.UID())
			container.AddChild(room)
		}
	}
	return room, nil
}

func (o *EntityManager) CreateTeam(userId string) (basis.ITeamEntity, error) {
	o.teamIndexMu.Lock()
	defer o.teamIndexMu.Unlock()
	if userId == "" || !o.userIndex.CheckUser(userId) {
		return nil, errors.New(fmt.Sprintf("EntityManager.CreateTeam Error: User(%s) does not exist", userId))
	}
	team := entity.NewITeamEntity(basis.GetTeamId(), basis.TeamName, basis.MaxTeamMember)
	team.InitEntity()
	o.addEntityEventListener(team)
	o.teamIndex.AddTeam(team)
	team.AddChild(o.userIndex.GetUser(userId))
	team.SetParent(userId)
	return team, nil
}

func (o *EntityManager) CreateTeamCorps(teamId string) (basis.ITeamCorpsEntity, error) {
	o.teamCorpsIndexMu.Lock()
	defer o.teamCorpsIndexMu.Unlock()
	if teamId == "" || !o.teamIndex.CheckTeam(teamId) {
		return nil, errors.New(fmt.Sprintf("EntityManager.CreateTeamCorps Error: Team(%s) does not exist", teamId))
	}
	teamCorps := entity.NewITeamCorpsEntity(basis.GetTeamCorpsId(), basis.TeamCorpsName)
	teamCorps.InitEntity()
	o.addEntityEventListener(teamCorps)
	o.teamCorpsIndex.AddCorps(teamCorps)
	teamCorps.AddChild(o.teamIndex.GetTeam(teamId))
	teamCorps.SetParent(teamId)
	return teamCorps, nil
}

func (o *EntityManager) CreateChannel(chanId string, chanName string) (basis.IChannelEntity, error) {
	o.chanIndexMu.Lock()
	defer o.chanIndexMu.Unlock()
	if o.channelIndex.CheckChannel(chanId) {
		return nil, errors.New("EntityManager.CreateChannel Error: ChanId(" + chanId + ") Duplicate!")
	}
	channel := entity.NewIChannelEntity(chanId, chanName)
	channel.InitEntity()
	o.addEntityEventListener(channel)
	o.channelIndex.AddChannel(channel)
	return channel, nil
}

func (o *EntityManager) addEntityEventListener(entity basis.IEntity) {
	if dispatcher, ok := entity.(basis.IVariableSupport); ok {
		dispatcher.AddEventListener(basis.EventEntityVarChanged, o.onEntityVar)
	}
}

func (o *EntityManager) removeEntityEventListener(entity basis.IEntity) {
	if dispatcher, ok := entity.(basis.IVariableSupport); ok {
		dispatcher.RemoveEventListener(basis.EventEntityVarChanged, o.onEntityVar)
	}
}

//事件转发
func (o *EntityManager) onEntityVar(evd *eventx.EventData) {
	evd.StopImmediatePropagation()
	o.DispatchEvent(basis.EventManagerVarChanged, o, basis.ManagerVarEventData{
		Entity: evd.CurrentTarget.(basis.IEntity),
		Data:   evd.Data.(encodingx.IKeyValue)}) //[0]为实体目标，[1]为变量
}

//----------------------------

func (o *EntityManager) World() basis.IWorldEntity {
	return o.rootWorld
}

func (o *EntityManager) GetZone(zoneId string) (basis.IZoneEntity, bool) {
	o.zoneIndexMu.RLock()
	defer o.zoneIndexMu.RUnlock()
	if zone := o.zoneIndex.GetZone(zoneId); nil != zone {
		return zone, true
	}
	return nil, false
}

func (o *EntityManager) GetRoom(roomId string) (basis.IRoomEntity, bool) {
	o.roomIndexMu.RLock()
	defer o.roomIndexMu.RUnlock()
	if room := o.roomIndex.GetRoom(roomId); nil != room {
		return room, true
	}
	return nil, false
}

func (o *EntityManager) GetUser(userId string) (basis.IUserEntity, bool) {
	o.userIndexMu.RLock()
	defer o.userIndexMu.RUnlock()
	if user := o.userIndex.GetUser(userId); nil != user {
		return user, true
	}
	return nil, false
}

func (o *EntityManager) GetTeam(teamId string) (basis.ITeamEntity, bool) {
	o.teamIndexMu.RLock()
	defer o.teamIndexMu.RUnlock()
	if team := o.teamIndex.GetTeam(teamId); nil != team {
		return team, true
	}
	return nil, false
}

func (o *EntityManager) GetTeamCorps(corpsId string) (basis.ITeamCorpsEntity, bool) {
	o.teamCorpsIndexMu.RLock()
	defer o.teamCorpsIndexMu.RUnlock()
	if teamCorps := o.teamCorpsIndex.GetCorps(corpsId); nil != teamCorps {
		return teamCorps, true
	}
	return nil, false
}

func (o *EntityManager) GetChannel(chanId string) (basis.IChannelEntity, bool) {
	o.chanIndexMu.RLock()
	defer o.chanIndexMu.RUnlock()
	if channel := o.channelIndex.GetChannel(chanId); nil != channel {
		return channel, true
	}
	return nil, false
}

func (o *EntityManager) GetEntity(entityType basis.EntityType, entityId string) (basis.IEntity, bool) {
	panic("implement me")
}

//-----------------------

func (o *EntityManager) ZoneIndex() basis.IZoneIndex {
	return o.zoneIndex
}

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
	return o.channelIndex
}

func (o *EntityManager) GetEntityIndex(entityType basis.EntityType) basis.IEntityIndex {
	switch entityType {
	case basis.EntityZone:
		return o.zoneIndex
	case basis.EntityRoom:
		return o.roomIndex
	case basis.EntityUser:
		return o.userIndex
	case basis.EntityTeam:
		return o.teamIndex
	case basis.EntityTeamCorps:
		return o.teamCorpsIndex
	case basis.EntityChannel:
		return o.channelIndex
	}
	return nil
}
