// Package events
// Create on 2023/9/8
// @author xuzhuoxi
package events

const (
	// EventTeamInit 队伍创建
	// 事件数据： ITeamEntity
	EventTeamInit = "MMO.EventTeamInit"
	// EventTeamDestroy 队伍销毁
	// 事件数据： ITeamEntity
	EventTeamDestroy = "MMO.EventTeamDestroy"
)
