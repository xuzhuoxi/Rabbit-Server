// Package events
// Create on 2023/9/8
// @author xuzhuoxi
package events

const (
	// EventTeamCorpsInit 军团创建
	// 事件数据： ITeamCorpsEntity
	EventTeamCorpsInit = "Manager.EventTeamCorpsInit"
	// EventTeamCorpsDestroy 军团销毁
	// 事件数据： ITeamCorpsEntity
	EventTeamCorpsDestroy = "Manager.EventTeamCorpsDestroy"
)
