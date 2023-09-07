// Package basis
// Created by xuzhuoxi
// on 2019-03-14.
// @author xuzhuoxi
package basis

// IPlayerBlackList 黑名单
type IPlayerBlackList interface {
	// Blacks 通信黑名单，返回原始切片，如果要修改的，请先copy
	Blacks() []string
	// AddBlack 增加黑名单
	AddBlack(targetId string) error
	// RemoveBlack 移除黑名单
	RemoveBlack(targetId string) error
	// OnBlack 处于
	OnBlack(targetId string) bool
}

// IPlayerWhiteList 黑名单
type IPlayerWhiteList interface {
	// Whites 通信白名单，返回原始切片，如果要修改的，请先copy
	Whites() []string
	// AddWhite 增加白名单
	AddWhite(targetId string) error
	// RemoveWhite 移除白名单
	RemoveWhite(targetId string) error
	// OnWhite 处于
	OnWhite(targetId string) bool
}

// IPlayerSubscriber 参与者
type IPlayerSubscriber interface {
	IPlayerWhiteList
	IPlayerBlackList
	// OnActive 处于激活
	OnActive(targetId string) bool
}
