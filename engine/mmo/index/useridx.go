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
	return &UserIndex{EntityIndex: *NewEntityIndex("UserIndex", basis.EntityUser)}
}

type UserIndex struct {
	EntityIndex
}

func (o *UserIndex) CheckUser(userId string) bool {
	return o.EntityIndex.Check(userId)
}

func (o *UserIndex) GetUser(userId string) (user basis.IUserEntity, ok bool) {
	user, ok = o.EntityIndex.Get(userId).(basis.IUserEntity)
	return
}

func (o *UserIndex) AddUser(user basis.IUserEntity) error {
	return o.EntityIndex.Add(user)
}

func (o *UserIndex) RemoveUser(userId string) (basis.IUserEntity, error) {
	c, err := o.EntityIndex.Remove(userId)
	if nil != c {
		return c.(basis.IUserEntity), err
	}
	return nil, err
}

func (o *UserIndex) UpdateUser(user basis.IUserEntity) error {
	return o.EntityIndex.Update(user)
}
