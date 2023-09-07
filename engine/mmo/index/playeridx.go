// Package index
// Created by xuzhuoxi
// on 2019-03-09.
// @author xuzhuoxi
package index

import (
	"github.com/xuzhuoxi/Rabbit-Server/engine/mmo/basis"
)

func NewIPlayerIndex() basis.IPlayerIndex {
	return NewPlayerIndex()
}

func NewPlayerIndex() *PlayerIndex {
	return &PlayerIndex{EntityIndex: NewEntityIndex("PlayerIndex", basis.EntityPlayer)}
}

type PlayerIndex struct {
	EntityIndex basis.IEntityIndex
}

func (o *PlayerIndex) EntityType() basis.EntityType {
	return o.EntityIndex.EntityType()
}

func (o *PlayerIndex) ForEachEntity(each func(entity basis.IEntity)) {
	o.EntityIndex.ForEachEntity(each)
}

func (o *PlayerIndex) CheckPlayer(playerId string) bool {
	return o.EntityIndex.Check(playerId)
}

func (o *PlayerIndex) GetPlayer(playerId string) (player basis.IPlayerEntity, ok bool) {
	player, ok = o.EntityIndex.Get(playerId).(basis.IPlayerEntity)
	return
}

func (o *PlayerIndex) AddPlayer(player basis.IPlayerEntity) (rsCode int32, err error) {
	num, err1 := o.EntityIndex.Add(player)
	if nil == err1 {
		return
	}
	if num == 1 || num == 2 {
		return basis.CodeMMOIndexType, err1
	}
	return basis.CodeMMOPlayerExist, err1
}

func (o *PlayerIndex) RemovePlayer(playerId string) (player basis.IPlayerEntity, rsCode int32, err error) {
	c, _, err1 := o.EntityIndex.Remove(playerId)
	if nil != c {
		return c.(basis.IPlayerEntity), 0, nil
	}
	return nil, basis.CodeMMOPlayerNotExist, err1
}

func (o *PlayerIndex) UpdatePlayer(player basis.IPlayerEntity) (rsCode int32, err error) {
	_, err1 := o.EntityIndex.Update(player)
	if nil != err1 {
		return basis.CodeMMOIndexType, err1
	}
	return
}
