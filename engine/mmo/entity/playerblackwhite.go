// Package entity
// Created by xuzhuoxi
// on 2019-03-07.
// @author xuzhuoxi
package entity

import "github.com/xuzhuoxi/Rabbit-Server/engine/mmo/basis"

func NewIPlayerSubscriber() basis.IPlayerSubscriber {
	return &PlayerSubscriber{PlayerBlackList: *NewPlayerBlackList(), PlayerWhiteList: *NewPlayerWhiteList()}
}

func NewIPlayerWhiteList() basis.IPlayerWhiteList {
	return NewPlayerWhiteList()
}
func NewIPlayerBlackList() basis.IPlayerBlackList {
	return NewPlayerBlackList()
}

func NewPlayerSubscriber() *PlayerSubscriber {
	return &PlayerSubscriber{PlayerBlackList: *NewPlayerBlackList(), PlayerWhiteList: *NewPlayerWhiteList()}
}

func NewPlayerWhiteList() *PlayerWhiteList {
	return &PlayerWhiteList{whiteGroup: NewEntityListGroup(basis.EntityPlayer)}
}
func NewPlayerBlackList() *PlayerBlackList {
	return &PlayerBlackList{blackGroup: NewEntityListGroup(basis.EntityPlayer)}
}

type PlayerSubscriber struct {
	PlayerBlackList
	PlayerWhiteList
}

func (o *PlayerSubscriber) OnActive(targetId string) bool {
	return o.OnWhite(targetId) && !o.OnBlack(targetId)
}

type PlayerWhiteList struct {
	whiteGroup basis.IEntityGroup
}

func (o *PlayerWhiteList) Whites() []string {
	return o.whiteGroup.Entities()
}

func (o *PlayerWhiteList) AddWhite(targetId string) error {
	return o.whiteGroup.Accept(targetId)
}

func (o *PlayerWhiteList) RemoveWhite(targetId string) error {
	return o.whiteGroup.Drop(targetId)
}

func (o *PlayerWhiteList) OnWhite(targetId string) bool {
	return o.whiteGroup.ContainEntity(targetId)
}

type PlayerBlackList struct {
	blackGroup basis.IEntityGroup
}

func (o *PlayerBlackList) Blacks() []string {
	return o.blackGroup.Entities()
}

func (o *PlayerBlackList) AddBlack(targetId string) error {
	return o.blackGroup.Accept(targetId)
}

func (o *PlayerBlackList) RemoveBlack(targetId string) error {
	return o.blackGroup.Drop(targetId)
}

func (o *PlayerBlackList) OnBlack(targetId string) bool {
	return o.blackGroup.ContainEntity(targetId)
}
