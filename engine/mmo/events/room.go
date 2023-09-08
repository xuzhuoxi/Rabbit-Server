// Package events
// Create on 2023/9/8
// @author xuzhuoxi
package events

const (
	// EventRoomInit 房间创建
	// 事件数据： IRoomEntity
	EventRoomInit = "Manager.EventRoomInit"
	// EventRoomDestroy 房间销毁
	// 事件数据： IRoomEntity
	EventRoomDestroy = "Manager.EventRoomDestroy"
)
