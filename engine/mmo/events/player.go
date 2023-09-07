// Package events
// Create on 2023/9/3
// @author xuzhuoxi
package events

const (
	// EventPlayerBorn 玩家创建
	EventPlayerBorn = "Manager.EventPlayerBorn"
	// EventPlayerDestroy 玩家销毁
	EventPlayerDestroy = "Manager.EventPlayerDestroy"
)

const (
	// EventPlayerLeaveRoom 玩家离开房间
	// 事件数据：*PlayerEventDataLeaveRoom
	EventPlayerLeaveRoom = "Manager.EventPlayerLeaveRoom"
	// EventPlayerEnterRoom 玩家进入房间
	// 事件数据：IPlayerEntity
	EventPlayerEnterRoom = "Manager.EventPlayerEnterRoom"
)

type PlayerEventDataLeaveRoom struct {
	RoomId   string
	PlayerId string
}
