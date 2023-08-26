// Package basis
// Created by xuzhuoxi
// on 2019-03-14.
// @author xuzhuoxi
package basis

type ChannelType uint16

const (
	// None 无效
	None ChannelType = iota
	// StatusChannel 状态
	StatusChannel
	// ChatChannel 聊天
	ChatChannel
	// EventChannel 事件
	EventChannel
)

// IChannelBehavior 频道行为
type IChannelBehavior interface {
	MyChannel() IChannelEntity
	// TouchChannel 订阅频道
	TouchChannel(subscriber string)
	// UnTouchChannel 取消频道订阅
	UnTouchChannel(subscriber string)
	// Broadcast 消息广播
	Broadcast(speaker string, handler func(receiver string)) int
	// BroadcastSome 消息指定目标广播
	BroadcastSome(speaker string, receiver []string, handler func(receiver string)) int
}
