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
	//// CreateWorld 构造世界
	//CreateWorld(worldId string, worldName string, asRoot bool, vars encodingx.IKeyValue) (world basis.IWorldEntity, rsCode int32, err error)
	//// CreateZoneAt 构造区域
	//CreateZoneAt(zoneId string, zoneName string, container basis.IEntityContainer, vars encodingx.IKeyValue) (basis.IZoneEntity, error)

	// CreateRoom 构造房间
	CreateRoom(roomId string, roomName string, tags []string, vars encodingx.IKeyValue) (basis.IRoomEntity, error)

	// CreateUser 创建用户实体
	CreateUser(userId string, vars encodingx.IKeyValue) (basis.IUserEntity, error)
	// CreateTeam 创建队伍
	CreateTeam(userId string, vars encodingx.IKeyValue) (basis.ITeamEntity, error)
	// CreateTeamCorps 创建团队
	CreateTeamCorps(teamId string, vars encodingx.IKeyValue) (basis.ITeamCorpsEntity, error)
	// CreateChannel 构造频道
	CreateChannel(chanId string, chanName string, vars encodingx.IKeyValue) (basis.IChannelEntity, error)

	// DestroyEntity 删除实体
	DestroyEntity(entity basis.IEntity) error
	// DestroyEntityBy 通过类型和Id删除实体
	DestroyEntityBy(entityType basis.EntityType, eId string) (basis.IEntity, error)
}

type IEntityIndexSet interface {
	//ZoneIndex() basis.ITagIndex
	RoomIndex() basis.IRoomIndex
	UserIndex() basis.IUserIndex
	TeamIndex() basis.ITeamIndex
	TeamCorpsIndex() basis.ITeamCorpsIndex
	ChannelIndex() basis.IChannelIndex
	GetEntityIndex(entityType basis.EntityType) basis.IEntityIndex
}

type IEntityGetter interface {
	//// GetWorld 获取区域实例
	//GetWorld(worldId string) (world basis.IWorldEntity, ok bool)
	//// GetZone 获取区域实例
	//GetZone(zoneId string) (zone basis.IZoneEntity, ok bool)

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

type IEntityManager interface {
	eventx.IEventDispatcher
	IEntityFactory
	IEntityGetter
	IEntityIndexSet
	basis.IManagerBase

	// BuildEnv 构建MMO环境
	BuildEnv(cfg *config.MMOConfig) error
	//ConstructWorld(cfg *config.MMOConfig, worldId string) (world basis.IWorldEntity, err error)
	//ConstructWorldDefault(cfg *config.MMOConfig) (world basis.IWorldEntity, err error)
}

func NewIEntityManager() IEntityManager {
	return NewEntityManager()
}

func NewEntityManager() IEntityManager {
	rs := &EntityManager{logger: logx.DefaultLogger()}
	//rs.worldIndex = index.NewIWorldIndex()
	//rs.zoneIndex = index.NewIZoneIndex()

	rs.roomIndex = index.NewIRoomIndex()
	rs.userIndex = index.NewIUserIndex()
	rs.teamIndex = index.NewITeamIndex()
	rs.teamCorpsIndex = index.NewITeamCorpsIndex()
	rs.chanIndex = index.NewIChannelIndex()
	return rs
}

//----------------------------

type EntityManager struct {
	//worldIndex     basis.IWorldIndex // 协程安全
	//worldLock      sync.RWMutex
	//zoneIndex      basis.ITagIndex // 协程安全
	//zoneLock       sync.RWMutex
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

	//rootWorld basis.IWorldEntity
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
		_, err1 := o.createRoom(room.Id, room.Name, room.Tags, nil)
		if nil != err1 {
			return err1
		}
	}
	return nil
}

//func (o *EntityManager) ConstructWorldDefault(cfg *config.MMOConfig) (world basis.IWorldEntity, err error) {
//	return o.ConstructWorld(cfg, cfg.DefaultWorld)
//}
//
//func (o *EntityManager) ConstructWorld(cfg *config.MMOConfig, worldId string) (world basis.IWorldEntity, err error) {
//	if relation, ok := cfg.Relations.GetWorldRelation(worldId); ok {
//		o.logger.Infoln("Start Construct World:", worldId, cfg)
//		worldEntity, ok1 := cfg.Entities.FindWorld(worldId)
//		if !ok1 {
//			err = errors.New("Construct World Fail: " + worldId + " not configured.")
//			o.logger.Warnln(err)
//			return
//		}
//		newWorld, err1 := o.CreateWorld(worldEntity.Id, worldEntity.Name, true, nil)
//		if nil != err1 {
//			err = err1
//			o.logger.Warnln(err)
//			return
//		}
//		for _, zoneCfg := range relation.Zones {
//			z, ok2 := cfg.Entities.FindZone(zoneCfg.ZoneId)
//			if !ok2 {
//				err = errors.New("Construct Zone Fail: " + zoneCfg.ZoneId + " not configured.")
//				o.logger.Warnln(err)
//				return
//			}
//			zone, err2 := o.CreateZoneAt(z.Id, z.Name, newWorld, nil)
//			if nil != err2 {
//				err = err2
//				o.logger.Warnln(err)
//				return
//			}
//			for _, roomId := range zoneCfg.Rooms {
//				r, ok3 := cfg.Entities.FindRoom(roomId)
//				if !ok3 {
//					err = errors.New("Construct Room Fail: " + roomId + " not configured.")
//					o.logger.Warnln(err)
//					return
//				}
//				_, err3 := o.CreateRoom(r.Id, r.Name, zone, nil)
//				if nil != err3 {
//					err = err3
//					o.logger.Warnln(err)
//					return
//				}
//			}
//		}
//		world = newWorld
//		o.logger.Infoln("Finish Construct World:", newWorld.UID())
//		return
//	}
//	return nil, errors.New("No World Relation: " + worldId)
//}

//func (o *EntityManager) CreateWorld(worldId string, worldName string, asRoot bool, vars encodingx.IKeyValue) (world basis.IWorldEntity, rsCode int32, err error) {
//	o.worldLock.Lock()
//	defer o.worldLock.Unlock()
//	if o.worldIndex.CheckWorld(worldId) {
//		return nil, basis.CodeMMOWorldExist, errors.New("EntityManager.CreateWorld Error: WorldId(" + worldId + ") Duplicate!")
//	}
//	world = entity.CreateWorldEntity(worldId, worldName)
//	world.InitEntity()
//	world.SetVars(vars)
//	err = o.worldIndex.AddWorld(world)
//	if nil != err {
//		return nil, err
//	}
//	o.addEntityEventListener(world)
//	if asRoot {
//		o.rootWorld = world
//	}
//	return world, nil
//}
//
//func (o *EntityManager) CreateZoneAt(zoneId string, zoneName string, container basis.IEntityContainer,
//	vars encodingx.IKeyValue) (basis.IZoneEntity, error) {
//	o.zoneLock.Lock()
//	defer o.zoneLock.Unlock()
//	if o.zoneIndex.CheckTag(zoneId) {
//		return nil, errors.New("EntityManager.CreateZoneAt Error: ZoneId(" + zoneId + ") Duplicate!")
//	}
//	zone := entity.NewIZoneEntity(zoneId, zoneName)
//	zone.InitEntity()
//	zone.SetVars(vars)
//	err := o.zoneIndex.AddZone(zone)
//	if nil != err {
//		return nil, err
//	}
//	o.addEntityEventListener(zone)
//	if nil != container {
//		if e, ok := container.(basis.IEntity); ok {
//			zone.SetParent(e.UID())
//			err1 := container.AddChild(zone)
//			if err1 != nil {
//				return nil, err1
//			}
//		}
//	}
//	return zone, nil
//}

func (o *EntityManager) CreateRoom(roomId string, roomName string, tags []string, vars encodingx.IKeyValue) (basis.IRoomEntity, error) {
	o.roomLock.Lock()
	defer o.roomLock.Unlock()
	return o.createRoom(roomId, roomName, tags, vars)
}

func (o *EntityManager) createRoom(roomId string, roomName string, tags []string, vars encodingx.IKeyValue) (basis.IRoomEntity, error) {
	if o.roomIndex.CheckRoom(roomId) {
		return nil, errors.New("EntityManager.CreateRoomAt Error: RoomId(" + roomId + ") Duplicate")
	}
	room := entity.NewIRoomEntity(roomId, roomName)
	room.InitEntity()
	room.SetVars(vars)
	room.SetTags(tags)
	err := o.roomIndex.AddRoom(room)
	if nil != err {
		return nil, err
	}
	o.addEntityEventListener(room)
	return room, nil
}

func (o *EntityManager) CreateUser(userId string, vars encodingx.IKeyValue) (basis.IUserEntity, error) {
	o.userLock.Lock()
	defer o.userLock.Unlock()
	if userId == "" || o.userIndex.CheckUser(userId) {
		return nil, errors.New(fmt.Sprintf("EntityManager.CreateUser Error: User(%s) is nil or exist", userId))
	}
	user := entity.NewIUserEntity(userId)
	user.SetVars(vars)
	err := o.userIndex.AddUser(user)
	if nil != err {
		return nil, err
	}
	o.addEntityEventListener(user)
	return user, nil
}

func (o *EntityManager) CreateTeam(userId string, vars encodingx.IKeyValue) (basis.ITeamEntity, error) {
	o.teamLock.Lock()
	defer o.teamLock.Unlock()
	_, okUser := o.userIndex.GetUser(userId)
	if !okUser {
		return nil, errors.New(fmt.Sprintf("EntityManager.CreateTeam Error: User(%s) does not exist", userId))
	}
	team := entity.NewITeamEntity(basis.GetTeamId(), basis.TeamName, basis.MaxTeamMember)
	team.InitEntity()
	team.SetVars(vars)
	err := o.teamIndex.AddTeam(team)
	if nil != err {
		return nil, err
	}
	o.addEntityEventListener(team)
	//_ = team.AddChild(user)
	//team.SetParent(userId)
	return team, nil
}

func (o *EntityManager) CreateTeamCorps(teamId string, vars encodingx.IKeyValue) (basis.ITeamCorpsEntity, error) {
	o.teamCorpsLock.Lock()
	defer o.teamCorpsLock.Unlock()
	_, okTeam := o.teamIndex.GetTeam(teamId)
	if !okTeam {
		return nil, errors.New(fmt.Sprintf("EntityManager.CreateTeamCorps Error: Team(%s) does not exist", teamId))
	}
	teamCorps := entity.NewITeamCorpsEntity(basis.GetTeamCorpsId(), basis.TeamCorpsName)
	teamCorps.InitEntity()
	teamCorps.SetVars(vars)
	err := o.teamCorpsIndex.AddCorps(teamCorps)
	if nil != err {
		return nil, err
	}
	o.addEntityEventListener(teamCorps)
	//_ = teamCorps.AddChild(team)
	//teamCorps.SetParent(teamId)
	return teamCorps, nil
}

func (o *EntityManager) CreateChannel(chanId string, chanName string, vars encodingx.IKeyValue) (basis.IChannelEntity, error) {
	o.chanLock.Lock()
	defer o.chanLock.Unlock()
	if o.chanIndex.CheckChannel(chanId) {
		return nil, errors.New("EntityManager.CreateChannel Error: ChanId(" + chanId + ") Duplicate!")
	}
	channel := entity.NewIChannelEntity(chanId, chanName)
	channel.InitEntity()
	channel.SetVars(vars)
	err := o.chanIndex.AddChannel(channel)
	if nil != err {
		return nil, err
	}
	o.addEntityEventListener(channel)
	return channel, nil
}

func (o *EntityManager) DestroyEntity(entity basis.IEntity) error {
	if nil == entity {
		return errors.New("DestroyEntity Error at: entity is nil. ")
	}
	_, err := o.DestroyEntityBy(entity.EntityType(), entity.UID())
	return err
}

func (o *EntityManager) DestroyEntityBy(entityType basis.EntityType, eId string) (entity basis.IEntity, err error) {
	switch entityType {
	//case basis.EntityWorld:
	//	o.worldLock.Lock()
	//	defer o.worldLock.Unlock()
	//	entity, err = o.worldIndex.RemoveWorld(eId)
	//case basis.EntityZone:
	//	o.zoneLock.Lock()
	//	defer o.zoneLock.Unlock()
	//	entity, err = o.zoneIndex.RemoveZone(eId)
	case basis.EntityRoom:
		o.roomLock.Lock()
		defer o.roomLock.Unlock()
		entity, err = o.roomIndex.RemoveRoom(eId)
	case basis.EntityUser:
		o.userLock.Lock()
		defer o.userLock.Unlock()
		entity, err = o.userIndex.RemoveUser(eId)
	case basis.EntityTeamCorps:
		o.teamCorpsLock.Lock()
		defer o.teamCorpsLock.Unlock()
		entity, err = o.teamCorpsIndex.RemoveCorps(eId)
	case basis.EntityTeam:
		o.teamLock.Lock()
		defer o.teamLock.Unlock()
		entity, err = o.teamIndex.RemoveTeam(eId)
	case basis.EntityChannel:
		o.chanLock.Lock()
		defer o.chanLock.Unlock()
		entity, err = o.chanIndex.RemoveChannel(eId)
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

//func (o *EntityManager) DefaultWorld() basis.IWorldEntity {
//	return o.rootWorld
//}
//
//func (o *EntityManager) GetWorld(worldId string) (world basis.IWorldEntity, ok bool) {
//	if len(worldId) == 0 {
//		return
//	}
//	return o.worldIndex.GetWorld(worldId)
//}
//
//func (o *EntityManager) GetZone(zoneId string) (zone basis.IZoneEntity, ok bool) {
//	if len(zoneId) == 0 {
//		return
//	}
//	return o.zoneIndex.GetZone(zoneId)
//}

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
	//case basis.EntityWorld:
	//	entity = o.worldIndex.Get(eId)
	//case basis.EntityZone:
	//	entity = o.zoneIndex.Get(eId)
	case basis.EntityRoom:
		entity = o.roomIndex.Get(eId)
	case basis.EntityUser:
		entity = o.userIndex.Get(eId)
	case basis.EntityTeamCorps:
		entity = o.teamCorpsIndex.Get(eId)
	case basis.EntityTeam:
		entity = o.teamIndex.Get(eId)
	case basis.EntityChannel:
		entity = o.chanIndex.Get(eId)
	}
	return entity, entity != nil
}

//-----------------------

//func (o *EntityManager) ZoneIndex() basis.ITagIndex {
//	return o.zoneIndex
//}

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

func (o *EntityManager) GetEntityIndex(entityType basis.EntityType) basis.IEntityIndex {
	switch entityType {
	//case basis.EntityZone:
	//	return o.zoneIndex
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
