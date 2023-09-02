// Package index
// Created by xuzhuoxi
// on 2019-03-09.
// @author xuzhuoxi
package index

import (
	"github.com/xuzhuoxi/Rabbit-Server/engine/mmo/basis"
)

func NewIUserIndex() basis.IUserIndex {
	return NewUserIndex()
}

func NewUserIndex() *UserIndex {
	return &UserIndex{EntityIndex: NewEntityIndex("UserIndex", basis.EntityUser)}
}

type UserIndex struct {
	EntityIndex basis.IEntityIndex
}

func (o *UserIndex) EntityType() basis.EntityType {
	return o.EntityIndex.EntityType()
}

func (o *UserIndex) ForEachEntity(each func(entity basis.IEntity)) {
	o.EntityIndex.ForEachEntity(each)
}

func (o *UserIndex) CheckUser(userId string) bool {
	return o.EntityIndex.Check(userId)
}

func (o *UserIndex) GetUser(userId string) (user basis.IUserEntity, ok bool) {
	user, ok = o.EntityIndex.Get(userId).(basis.IUserEntity)
	return
}

func (o *UserIndex) AddUser(user basis.IUserEntity) (rsCode int32, err error) {
	num, err1 := o.EntityIndex.Add(user)
	if nil == err1 {
		return
	}
	if num == 1 || num == 2 {
		return basis.CodeMMOIndexType, err1
	}
	return basis.CodeMMOUserExist, err1
}

func (o *UserIndex) RemoveUser(userId string) (user basis.IUserEntity, rsCode int32, err error) {
	c, _, err1 := o.EntityIndex.Remove(userId)
	if nil != c {
		return c.(basis.IUserEntity), 0, nil
	}
	return nil, basis.CodeMMOUserNotExist, err1
}

func (o *UserIndex) UpdateUser(user basis.IUserEntity) (rsCode int32, err error) {
	_, err1 := o.EntityIndex.Update(user)
	if nil != err1 {
		return basis.CodeMMOIndexType, err1
	}
	return
}
