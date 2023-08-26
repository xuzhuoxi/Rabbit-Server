// Package entity
// Created by xuzhuoxi
// on 2019-02-18.
// @author xuzhuoxi
package entity

import (
	"github.com/xuzhuoxi/Rabbit-Server/engine/mmo/basis"
	"github.com/xuzhuoxi/infra-go/slicex"
	"sync"
)

//-----------------------------------------------

func NewIChannelEntity(chanId string, chanName string) basis.IChannelEntity {
	return &ChannelEntity{ChanId: chanId, ChanName: chanName}
}

func NewChannelEntity(chanId string, chanName string) *ChannelEntity {
	return &ChannelEntity{ChanId: chanId, ChanName: chanName}
}

type ChannelEntity struct {
	ChanId     string
	ChanName   string
	Subscriber []string
	Mu         sync.RWMutex
}

func (o *ChannelEntity) UID() string {
	return o.ChanId
}

func (o *ChannelEntity) NickName() string {
	return o.ChanName
}

func (o *ChannelEntity) EntityType() basis.EntityType {
	return basis.EntityChannel
}

func (o *ChannelEntity) InitEntity() {
}

func (o *ChannelEntity) MyChannel() basis.IChannelEntity {
	return o
}

func (o *ChannelEntity) TouchChannel(subscriber string) {
	o.Mu.Lock()
	defer o.Mu.Unlock()
	if o.hasSubscriber(subscriber) {
		return
	}
	o.Subscriber = append(o.Subscriber, subscriber)
}

func (o *ChannelEntity) UnTouchChannel(subscriber string) {
	o.Mu.Lock()
	defer o.Mu.Unlock()
	index, ok := slicex.IndexString(o.Subscriber, subscriber)
	if !ok {
		return
	}
	o.Subscriber = append(o.Subscriber[:index], o.Subscriber[index+1:]...)
}

func (o *ChannelEntity) Broadcast(speaker string, handler func(receiver string)) int {
	o.Mu.RLock()
	defer o.Mu.RUnlock()
	rs := len(o.Subscriber)
	for _, r := range o.Subscriber {
		if r == speaker {
			continue
		}
		handler(r)
	}
	return rs - 1
}

func (o *ChannelEntity) BroadcastSome(speaker string, receiver []string, handler func(receiver string)) int {
	o.Mu.RLock()
	defer o.Mu.RUnlock()
	count := 0
	for _, v := range o.Subscriber {
		if _, ok := slicex.IndexString(receiver, v); ok && speaker != v {
			handler(v)
			count++
		}
	}
	return count
}

func (o *ChannelEntity) hasSubscriber(subscriber string) bool {
	_, ok := slicex.IndexString(o.Subscriber, subscriber)
	return ok
}
