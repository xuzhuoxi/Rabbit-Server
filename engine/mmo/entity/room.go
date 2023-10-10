// Package entity
// Created by xuzhuoxi
// on 2019-02-18.
// @author xuzhuoxi
package entity

import (
	"fmt"
	"github.com/xuzhuoxi/Rabbit-Server/engine/mmo/basis"
	"github.com/xuzhuoxi/Rabbit-Server/engine/mmo/events"
	"github.com/xuzhuoxi/Rabbit-Server/engine/mmo/vars"
	"github.com/xuzhuoxi/infra-go/eventx"
	"github.com/xuzhuoxi/infra-go/extendx/protox"
)

func NewIAOBRoomEntity(roomId string, roomName string) basis.IRoomEntity {
	return NewAOBRoomEntity(roomId, roomName)
}

func NewIRoomEntity(roomId string, roomName string) basis.IRoomEntity {
	return NewRoomEntity(roomId, roomName)
}

func NewAOBRoomEntity(roomId string, roomName string) *AOBRoomEntity {
	room := &AOBRoomEntity{RoomEntity: *NewRoomEntity(roomId, roomName)}
	return room
}

func NewRoomEntity(roomId string, roomName string) *RoomEntity {
	room := &RoomEntity{RoomId: roomId, RoomName: roomName, MaxMember: 0}
	room.MaxMember = 100
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
	RoomId    string
	RoomName  string
	MaxMember int
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
	o.ListEntityContainer = *NewListEntityContainer(o.MaxMember)
	o.VariableSupport = *vars.NewVariableSupport(o)
	o.UnitContainer = *NewUnitContainer(o.RoomId, 1000)
}

func (o *RoomEntity) DestroyEntity() {
	o.UnitContainer.ForEachUnit(func(child basis.IUnitEntity) (interrupt bool) {
		child.RemoveEventListener(events.EventEntityVarsChanged, o.onEventRedirect)
		child.RemoveEventListener(events.EventEntityVarChanged, o.onEventRedirect)
		return false
	})
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
	if rsCode == protox.CodeSuc {
		unit.AddEventListener(events.EventEntityVarChanged, o.onEventRedirect)
		unit.AddEventListener(events.EventEntityVarsChanged, o.onEventRedirect)
		defer o.DispatchEvent(events.EventUnitInit, o, unit)
		fmt.Println("[RoomEntity.CreateUnit]")
	}
	return
}

func (o *RoomEntity) CreateUnits(params []basis.UnitParams, mustAll bool) (units []basis.IUnitEntity, rsCode int32, err error) {
	units, rsCode, err = o.UnitContainer.CreateUnits(params, mustAll)
	if rsCode == protox.CodeSuc {
		for index := range units {
			units[index].AddEventListener(events.EventEntityVarChanged, o.onEventRedirect)
			units[index].AddEventListener(events.EventEntityVarsChanged, o.onEventRedirect)
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
	if rsCode == protox.CodeSuc {
		unit.RemoveEventListener(events.EventEntityVarsChanged, o.onEventRedirect)
		unit.RemoveEventListener(events.EventEntityVarChanged, o.onEventRedirect)
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
		rs[index].RemoveEventListener(events.EventEntityVarsChanged, o.onEventRedirect)
		rs[index].RemoveEventListener(events.EventEntityVarChanged, o.onEventRedirect)
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

// 事件重定向
func (o *RoomEntity) onEventRedirect(evd *eventx.EventData) {
	evd.StopImmediatePropagation()
	o.DispatchEvent(evd.EventType, o, evd.Data)
}
