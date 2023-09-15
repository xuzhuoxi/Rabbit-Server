// Package entity
// Created by xuzhuoxi
// on 2019-02-18.
// @author xuzhuoxi
package entity

import (
	"github.com/xuzhuoxi/Rabbit-Server/engine/mmo/basis"
	"github.com/xuzhuoxi/Rabbit-Server/engine/mmo/vars"
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
	room.ListEntityContainer = *NewListEntityContainer(room.MaxMember)
	room.VariableSupport = *vars.NewVariableSupport(room)
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
	ListEntityContainer
	TagsSupport
	vars.VariableSupport
}

func (o *RoomEntity) UID() string {
	return o.RoomId
}

func (o *RoomEntity) Name() string {
	return o.RoomName
}

func (o *RoomEntity) EntityType() basis.EntityType {
	return basis.EntityRoom
}

func (o *RoomEntity) InitEntity() {
	o.ListEntityContainer = *NewListEntityContainer(o.MaxMember)
	o.VariableSupport = *vars.NewVariableSupport(o)
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
