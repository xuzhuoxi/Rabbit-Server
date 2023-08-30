// Package index
// Created by xuzhuoxi
// on 2019-03-09.
// @author xuzhuoxi
package index

import (
	"github.com/xuzhuoxi/Rabbit-Server/engine/mmo/basis"
)

func NewIRoomIndex() basis.IRoomIndex {
	return NewRoomIndex()
}

func NewRoomIndex() *RoomIndex {
	return &RoomIndex{EntityIndex: *NewEntityIndex("RoomIndex", basis.EntityRoom)}
}

type RoomIndex struct {
	EntityIndex
}

func (o *RoomIndex) CheckRoom(roomId string) bool {
	return o.EntityIndex.Check(roomId)
}

func (o *RoomIndex) GetRoom(roomId string) (room basis.IRoomEntity, ok bool) {
	room, ok = o.EntityIndex.Get(roomId).(basis.IRoomEntity)
	return
}

func (o *RoomIndex) AddRoom(room basis.IRoomEntity) error {
	return o.EntityIndex.Add(room)
}

func (o *RoomIndex) RemoveRoom(roomId string) (basis.IRoomEntity, error) {
	c, err := o.EntityIndex.Remove(roomId)
	if nil != c {
		return c.(basis.IRoomEntity), err
	}
	return nil, err
}

func (o *RoomIndex) UpdateRoom(room basis.IRoomEntity) error {
	return o.EntityIndex.Update(room)
}
