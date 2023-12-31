// Package basis
// Created by xuzhuoxi
// on 2019-03-09.
// @author xuzhuoxi
package basis

import (
	"strconv"
)

var (
	MaxTeamMember = 0
	TeamId        = 200
	TeamName      = "的队伍"
)

var (
	TeamCorpsId   = 1000
	TeamCorpsName = "的军团"
)

func GetTeamId() string {
	defer func() { TeamId++ }()
	return "T_" + strconv.Itoa(TeamId)
}

func GetTeamCorpsId() string {
	defer func() { TeamCorpsId++ }()
	return "TC_" + strconv.Itoa(TeamId)
}

type ITeamControl interface {
	// Leader 队长
	Leader() string
	// MemberList 用户列表
	MemberList() []string
	// ContainMember 检查用户
	ContainMember(memberId string) bool
	// AcceptMember 加入用户,进行唯一性检查
	AcceptMember(memberId string) error
	// DropMember 从组中移除用户
	DropMember(memberId string) error
	// RiseLeader 从组中移除用户
	RiseLeader(memberId string) error
	// DisbandTeam 解散队伍
	DisbandTeam() error
}
