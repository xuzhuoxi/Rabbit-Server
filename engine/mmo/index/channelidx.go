// Package index
// Created by xuzhuoxi
// on 2019-03-09.
// @author xuzhuoxi
package index

import (
	"github.com/xuzhuoxi/Rabbit-Server/engine/mmo/basis"
)

func NewIChannelIndex() basis.IChannelIndex {
	return NewChannelIndex()
}

func NewChannelIndex() *ChannelIndex {
	return &ChannelIndex{EntityIndex: *NewEntityIndex("ChannelIndex", basis.EntityChannel)}
}

type ChannelIndex struct {
	EntityIndex
}

func (o *ChannelIndex) CheckChannel(chanId string) bool {
	return o.EntityIndex.Check(chanId)
}

func (o *ChannelIndex) GetChannel(chanId string) basis.IChannelEntity {
	entity := o.EntityIndex.Get(chanId)
	if nil != entity {
		return entity.(basis.IChannelEntity)
	}
	return nil
}

func (o *ChannelIndex) AddChannel(channel basis.IChannelEntity) error {
	return o.EntityIndex.Add(channel)
}

func (o *ChannelIndex) RemoveChannel(chanId string) (basis.IChannelEntity, error) {
	c, err := o.EntityIndex.Remove(chanId)
	if nil != c {
		return c.(basis.IChannelEntity), err
	}
	return nil, err
}

func (o *ChannelIndex) UpdateChannel(channel basis.IChannelEntity) error {
	return o.EntityIndex.Update(channel)
}
