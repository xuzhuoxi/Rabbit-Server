// Package entity
// Created by xuzhuoxi
// on 2019-02-18.
// @author xuzhuoxi
package entity

import (
	"github.com/xuzhuoxi/Rabbit-Server/engine/mmo/basis"
	"github.com/xuzhuoxi/Rabbit-Server/engine/mmo/events"
	"github.com/xuzhuoxi/Rabbit-Server/engine/mmo/vars"
	"github.com/xuzhuoxi/Rabbit-Server/engine/server"
	"github.com/xuzhuoxi/infra-go/eventx"
)

func NewIAOBRoomEntity(roomId string, roomRefId string, roomName string, playerCap int) basis.IRoomEntity {
	return NewAOBRoomEntity(roomId, roomRefId, roomName, playerCap)
}

func NewIRoomEntity(roomId string, roomRefId string, roomName string, playerCap int) basis.IRoomEntity {
	return NewRoomEntity(roomId, roomRefId, roomName, playerCap)
}

func NewAOBRoomEntity(roomId string, roomRefId string, roomName string, playerCap int) *AOBRoomEntity {
	room := &AOBRoomEntity{RoomEntity: *NewRoomEntity(roomId, roomRefId, roomName, playerCap)}
	return room
}

func NewRoomEntity(roomId string, roomRefId string, roomName string, playerCap int) *RoomEntity {
	room := &RoomEntity{RoomId: roomId, RefId: roomRefId, RoomName: roomName, _PlayerCap: playerCap}
	return room
}

type RoomConfig struct {
	Id        string
	Name      string
	Private   bool
	MaxMember int
}

// AOBRoomEntity 范围广播房间，适用于mmo大型场景
type AOBRoomEntity struct {
	RoomEntity
}

func (e *AOBRoomEntity) Broadcast(speaker string, handler func(receiver string)) int {
	panic("+++++++++++++++++++")
}

// RoomEntity 常规房间
type RoomEntity struct {
	RoomId     string
	RefId      string
	RoomName   string
	_PlayerCap int
	vars.VariableSupport
	UnitContainer UnitContainer
	ListEntityContainer
	TagsSupport
}

func (o *RoomEntity) UID() string {
	return o.RoomId
}

func (o *RoomEntity) EntityType() basis.EntityType {
	return basis.EntityRoom
}

func (o *RoomEntity) Name() string {
	return o.RoomName
}

func (o *RoomEntity) InitEntity() {
	o.ListEntityContainer = *NewListEntityContainer(o._PlayerCap)
	o.VariableSupport = *vars.NewVariableSupport(o)
	o.UnitContainer = *NewUnitContainer(o.RoomId)
}

func (o *RoomEntity) DestroyEntity() {
	o.UnitContainer.ForEachUnit(func(child basis.IUnitEntity) (interrupt bool) {
		o.removeUnitEventListener(child)
		return false
	})
}

func (o *RoomEntity) RoomRefId() string {
	return o.RefId
}

func (o *RoomEntity) RoomMapId() string {
	return o.RefId
}

func (o *RoomEntity) PlayerCap() int {
	return o._PlayerCap
}

func (o *RoomEntity) SetPlayerCap(cap int) {
	o._PlayerCap = cap
}

func (o *RoomEntity) PlayerCount() int {
	return o.ListEntityContainer.NumChildren()
}

func (o *RoomEntity) Players() []basis.IPlayerEntity {
	num := o.ListEntityContainer.NumChildren()
	if num == 0 {
		return nil
	}
	rs := make([]basis.IPlayerEntity, 0, num)
	o.ListEntityContainer.ForEachChildByType(basis.EntityPlayer, func(child basis.IEntity) {
		rs = append(rs, child.(basis.IPlayerEntity))
	}, false)
	return rs
}

func (o *RoomEntity) UnitCount() int {
	return o.UnitContainer.UnitIndex.Size()
}

func (o *RoomEntity) UnitIndex() basis.IUnitIndex {
	return o.UnitContainer.UnitIndex
}

// IUnitContainer ---------- ---------- ---------- ----------

func (o *RoomEntity) Units() []basis.IUnitEntity {
	return o.UnitContainer.Units()
}

func (o *RoomEntity) CreateUnit(params basis.UnitParams) (unit basis.IUnitEntity, rsCode int32, err error) {
	unit, rsCode, err = o.UnitContainer.CreateUnit(params)
	if rsCode == server.CodeSuc {
		o.addUnitEventListener(unit)
		defer o.DispatchEvent(events.EventUnitInit, o, unit)
	}
	return
}

func (o *RoomEntity) CreateUnits(params []basis.UnitParams, mustAll bool) (units []basis.IUnitEntity, rsCode int32, err error) {
	units, rsCode, err = o.UnitContainer.CreateUnits(params, mustAll)
	if rsCode == server.CodeSuc {
		for index := range units {
			o.addUnitEventListener(units[index])
		}
		defer func() {
			for index := range units {
				o.DispatchEvent(events.EventUnitInit, o, units[index])
			}
		}()
	}
	return
}

func (o *RoomEntity) DestroyUnit(unitId string) (unit basis.IUnitEntity, rsCode int32, err error) {
	unit, rsCode, err = o.UnitContainer.DestroyUnit(unitId)
	if rsCode == server.CodeSuc {
		o.removeUnitEventListener(unit)
		defer o.DispatchEvent(events.EventUnitDestroy, o, unit)
	}
	return
}

func (o *RoomEntity) DestroyUnitsByOwner(owner string) []basis.IUnitEntity {
	rs := o.UnitContainer.DestroyUnitsByOwner(owner)
	if len(rs) == 0 {
		return nil
	}
	for index := range rs {
		o.removeUnitEventListener(rs[index])
	}
	defer func() {
		for index := range rs {
			o.DispatchEvent(events.EventUnitDestroy, o, rs[index])
		}
	}()
	return rs
}

func (o *RoomEntity) ForEachUnit(each func(child basis.IUnitEntity) (interrupt bool)) {
	o.UnitContainer.ForEachUnit(each)
}

func (o *RoomEntity) addUnitEventListener(unit eventx.IEventDispatcher) {
	unit.AddEventListener(events.EventEntityVarMod, o.onEventRedirect)
	unit.AddEventListener(events.EventEntityVarsMod, o.onEventRedirect)
	unit.AddEventListener(events.EventEntityVarDel, o.onEventRedirect)
	unit.AddEventListener(events.EventEntityVarsDel, o.onEventRedirect)
}

func (o *RoomEntity) removeUnitEventListener(unit eventx.IEventDispatcher) {
	unit.RemoveEventListener(events.EventEntityVarsDel, o.onEventRedirect)
	unit.RemoveEventListener(events.EventEntityVarDel, o.onEventRedirect)
	unit.RemoveEventListener(events.EventEntityVarsMod, o.onEventRedirect)
	unit.RemoveEventListener(events.EventEntityVarMod, o.onEventRedirect)
}

// 事件重定向
func (o *RoomEntity) onEventRedirect(evd *eventx.EventData) {
	//fmt.Println("[RoomEntity.onEventRedirect]", evd.EventType)
	evd.StopImmediatePropagation()
	o.DispatchEvent(evd.EventType, o, evd.Data)
}
