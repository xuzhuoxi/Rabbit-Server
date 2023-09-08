// Package events
// Create on 2023/9/3
// @author xuzhuoxi
package events

const (
	// EventPlayerInit 玩家创建
	// 事件数据： IPlayerEntity
	EventPlayerInit = "MMO.EventPlayerInit"
	// EventPlayerDestroy 玩家销毁
	// 事件数据： IPlayerEntity
	EventPlayerDestroy = "MMO.EventPlayerDestroy"
)

const (
	// EventPlayerLeaveRoom 玩家离开房间
	// 事件数据：*PlayerEventDataLeaveRoom
	EventPlayerLeaveRoom = "MMO.EventPlayerLeaveRoom"
	// EventPlayerEnterRoom 玩家进入房间
	// 事件数据：IPlayerEntity
	EventPlayerEnterRoom = "MMO.EventPlayerEnterRoom"
)

type PlayerEventDataLeaveRoom struct {
	RoomId   string
	PlayerId string
}
