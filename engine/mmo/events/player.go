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
	// EventPlayerEnterRoomActive 玩家进入房间(主动)
	// 事件数据：IPlayerEntity
	EventPlayerEnterRoomActive = "MMO.EventPlayerEnterRoomActive"
	// EventPlayerEnterRoomPassive 玩家进入房间(被动)
	// 事件数据：IPlayerEntity
	EventPlayerEnterRoomPassive = "MMO.EventPlayerEnterRoomPassive"
)

type PlayerEventDataLeaveRoom struct {
	RoomId   string
	PlayerId string
}
