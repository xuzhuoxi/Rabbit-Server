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
	return &RoomIndex{EntityIndex: NewEntityIndex("RoomIndex", basis.EntityRoom)}
}

type RoomIndex struct {
	EntityIndex basis.IEntityIndex
}

func (o *RoomIndex) Size() int {
	return o.EntityIndex.Size()
}

func (o *RoomIndex) EntityType() basis.EntityType {
	return o.EntityIndex.EntityType()
}

func (o *RoomIndex) ForEachEntity(each func(entity basis.IEntity) (interrupt bool)) {
	o.EntityIndex.ForEachEntity(each)
}

func (o *RoomIndex) CheckRoom(roomId string) bool {
	return o.EntityIndex.Check(roomId)
}

func (o *RoomIndex) GetRoom(roomId string) (room basis.IRoomEntity, ok bool) {
	room, ok = o.EntityIndex.Get(roomId).(basis.IRoomEntity)
	return
}

func (o *RoomIndex) AddRoom(room basis.IRoomEntity) (rsCode int32, err error) {
	num, err1 := o.EntityIndex.Add(room)
	if nil == err1 {
		return
	}
	if num == 1 || num == 2 {
		return basis.CodeMMOIndexType, err1
	}
	return basis.CodeMMORoomExist, err1
}

func (o *RoomIndex) RemoveRoom(roomId string) (room basis.IRoomEntity, rsCode int32, err error) {
	c, _, err1 := o.EntityIndex.Remove(roomId)
	if nil != c {
		return c.(basis.IRoomEntity), 0, nil
	}
	return nil, basis.CodeMMORoomNotExist, err1
}

func (o *RoomIndex) UpdateRoom(room basis.IRoomEntity) (rsCode int32, err error) {
	_, err1 := o.EntityIndex.Update(room)
	if nil != err1 {
		return basis.CodeMMOIndexType, err1
	}
	return
}
