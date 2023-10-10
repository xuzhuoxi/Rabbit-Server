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
	return &ChannelIndex{EntityIndex: NewEntityIndex("ChannelIndex", basis.EntityChannel)}
}

type ChannelIndex struct {
	EntityIndex basis.IEntityIndex
}

func (o *ChannelIndex) Size() int {
	return o.EntityIndex.Size()
}

func (o *ChannelIndex) EntityType() basis.EntityType {
	return o.EntityIndex.EntityType()
}

func (o *ChannelIndex) ForEachEntity(each func(entity basis.IEntity) (interrupt bool)) {
	o.EntityIndex.ForEachEntity(each)
}

func (o *ChannelIndex) CheckChannel(chanId string) bool {
	return o.EntityIndex.Check(chanId)
}

func (o *ChannelIndex) GetChannel(chanId string) (channel basis.IChannelEntity, ok bool) {
	channel, ok = o.EntityIndex.Get(chanId).(basis.IChannelEntity)
	return
}

func (o *ChannelIndex) AddChannel(channel basis.IChannelEntity) (rsCode int32, err error) {
	num, err1 := o.EntityIndex.Add(channel)
	if nil == err1 {
		return
	}
	if num == 1 || num == 2 {
		return basis.CodeMMOIndexType, err1
	}
	return basis.CodeMMOChanExist, err1
}

func (o *ChannelIndex) RemoveChannel(chanId string) (channel basis.IChannelEntity, rsCode int32, err error) {
	c, _, err1 := o.EntityIndex.Remove(chanId)
	if nil != c {
		return c.(basis.IChannelEntity), 0, nil
	}
	return nil, basis.CodeMMOChanNotExist, err1
}

func (o *ChannelIndex) UpdateChannel(channel basis.IChannelEntity) (rsCode int32, err error) {
	_, err1 := o.EntityIndex.Update(channel)
	if nil != err1 {
		return basis.CodeMMOIndexType, err1
	}
	return
}
