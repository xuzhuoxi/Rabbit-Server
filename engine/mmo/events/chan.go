// Package events
// Create on 2023/9/8
// @author xuzhuoxi
package events

const (
	// EventChanInit 频道创建
	// 事件数据： IChanEntity
	EventChanInit = "Manager.EventChanInit"
	// EventChanDestroy 频道销毁
	// 事件数据： IChanEntity
	EventChanDestroy = "Manager.EventChanDestroy"
)
