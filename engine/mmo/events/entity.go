// Package events
// Create on 2023/9/3
// @author xuzhuoxi
package events

const (
	// EventRoomAwake 房间创建
	EventRoomAwake = "Manager.EventRoomAwake"
	// EventRoomDestroy 房间销毁
	EventRoomDestroy = "Manager.EventRoomDestroy"

	// EventTeamAwake 队伍创建
	EventTeamAwake = "Manager.EventTeamAwake"
	// EventTeamDestroy 队伍销毁
	EventTeamDestroy = "Manager.EventTeamDestroy"

	// EventTeamCorpsAwake 军团创建
	EventTeamCorpsAwake = "Manager.EventTeamCorpsAwake"
	// EventTeamCorpsDestroy 军团销毁
	EventTeamCorpsDestroy = "Manager.EventTeamCorpsDestroy"

	// EventChanAwake 频道创建
	EventChanAwake = "Manager.EventChanAwake"
	// EventChanDestroy 频道销毁
	EventChanDestroy = "Manager.EventChanDestroy"
)
