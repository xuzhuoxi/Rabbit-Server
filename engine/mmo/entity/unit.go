// Package entity
// Create on 2023/10/7
// @author xuzhuoxi
package entity

import (
	"github.com/xuzhuoxi/Rabbit-Server/engine/mmo/basis"
	"github.com/xuzhuoxi/Rabbit-Server/engine/mmo/events"
	"github.com/xuzhuoxi/Rabbit-Server/engine/mmo/vars"
)

func NewIUnitEntity(unitId string, roomId string) basis.IUnitEntity {
	return NewUnitEntity(unitId, roomId)
}

func NewUnitEntity(unitId string, roomId string) *UnitEntity {
	unit := &UnitEntity{Uid: unitId, roomId: roomId}
	return unit
}

type UnitEntity struct {
	Uid    string //用户标识，唯一，内部使用
	roomId string
	vars.VariableSupport
}

func (o *UnitEntity) UID() string {
	return o.Uid
}

func (o *UnitEntity) EntityType() basis.EntityType {
	return basis.EntityUnit
}

func (o *UnitEntity) InitEntity() {
	o.VariableSupport = *vars.NewVariableSupport(o)
}

func (o *UnitEntity) DestroyEntity() {
}

func (o *UnitEntity) Position() (pos basis.XYZ) {
	val, ok := o.GetVar(vars.UnitPos)
	if !ok {
		return
	}
	arr := val.([]int32)
	if len(arr) >= 3 {
		pos.X, pos.Y, pos.Z = arr[0], arr[1], arr[2]
		return
	}
	if len(arr) == 2 {
		pos.X, pos.Y = arr[0], arr[1]
		return
	}
	if len(arr) == 1 {
		pos.X = arr[0]
		return
	}
	return
}

func (o *UnitEntity) SetPosition(pos basis.XYZ, notify bool) {
	posArr := pos.Array()
	ok := o.SetVar(vars.UnitPos, posArr, false)
	if notify && ok {
		o.VariableSupport.DispatchEvent(events.EventEntityVarChanged, o,
			&events.VarEventData{Entity: o, Key: vars.PlayerPos, Value: posArr})
	}
}

func (o *UnitEntity) RoomId() string {
	return o.roomId
}
