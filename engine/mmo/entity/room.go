// Package entity
// Created by xuzhuoxi
// on 2019-02-18.
// @author xuzhuoxi
package entity

import "github.com/xuzhuoxi/Rabbit-Server/engine/mmo/basis"

func NewIAOBRoomEntity(roomId string, roomName string) basis.IRoomEntity {
	return NewAOBRoomEntity(roomId, roomName)
}

func NewIRoomEntity(roomId string, roomName string) basis.IRoomEntity {
	return NewRoomEntity(roomId, roomName)
}

func NewAOBRoomEntity(roomId string, roomName string) *AOBRoomEntity {
	return &AOBRoomEntity{RoomEntity: *NewRoomEntity(roomId, roomName)}
}

func NewRoomEntity(roomId string, roomName string) *RoomEntity {
	return &RoomEntity{RoomId: roomId, RoomName: roomName, MaxMember: 0}
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
	VariableSupport
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
	//e.UserGroup = NewEntityListGroup(EntityUser)
	o.VariableSupport = *NewVariableSupport(o)
}

//func (e *RoomEntity) UserList() []string {
//	return e.UserGroup.Entities()
//}
//
//func (e *RoomEntity) ContainUser(userId string) bool {
//	return e.UserGroup.ContainEntity(userId)
//}
//
//func (e *RoomEntity) AcceptUser(userId string) error {
//	return e.UserGroup.Accept(userId)
//}
//
//func (e *RoomEntity) DropUser(userId string) error {
//	return e.UserGroup.Drop(userId)
//}
