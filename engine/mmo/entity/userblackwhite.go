// Package entity
// Created by xuzhuoxi
// on 2019-03-07.
// @author xuzhuoxi
package entity

import "github.com/xuzhuoxi/Rabbit-Server/engine/mmo/basis"

func NewIUserSubscriber() basis.IUserSubscriber {
	return &UserSubscriber{UserBlackList: *NewUserBlackList(), UserWhiteList: *NewUserWhiteList()}
}

func NewIUserWhiteList() basis.IUserWhiteList {
	return NewUserWhiteList()
}
func NewIUserBlackList() basis.IUserBlackList {
	return NewUserBlackList()
}

func NewUserSubscriber() *UserSubscriber {
	return &UserSubscriber{UserBlackList: *NewUserBlackList(), UserWhiteList: *NewUserWhiteList()}
}

func NewUserWhiteList() *UserWhiteList {
	return &UserWhiteList{whiteGroup: NewEntityListGroup(basis.EntityUser)}
}
func NewUserBlackList() *UserBlackList {
	return &UserBlackList{blackGroup: NewEntityListGroup(basis.EntityUser)}
}

type UserSubscriber struct {
	UserBlackList
	UserWhiteList
}

func (o *UserSubscriber) OnActive(targetId string) bool {
	return o.OnWhite(targetId) && !o.OnBlack(targetId)
}

type UserWhiteList struct {
	whiteGroup basis.IEntityGroup
}

func (o *UserWhiteList) Whites() []string {
	return o.whiteGroup.Entities()
}

func (o *UserWhiteList) AddWhite(targetId string) error {
	return o.whiteGroup.Accept(targetId)
}

func (o *UserWhiteList) RemoveWhite(targetId string) error {
	return o.whiteGroup.Drop(targetId)
}

func (o *UserWhiteList) OnWhite(targetId string) bool {
	return o.whiteGroup.ContainEntity(targetId)
}

type UserBlackList struct {
	blackGroup basis.IEntityGroup
}

func (o *UserBlackList) Blacks() []string {
	return o.blackGroup.Entities()
}

func (o *UserBlackList) AddBlack(targetId string) error {
	return o.blackGroup.Accept(targetId)
}

func (o *UserBlackList) RemoveBlack(targetId string) error {
	return o.blackGroup.Drop(targetId)
}

func (o *UserBlackList) OnBlack(targetId string) bool {
	return o.blackGroup.ContainEntity(targetId)
}
