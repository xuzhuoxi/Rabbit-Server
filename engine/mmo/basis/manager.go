// Package basis
// Created by xuzhuoxi
// on 2019-03-16.
// @author xuzhuoxi
package basis

type IManagerBase interface {
	// InitManager 初始化管理器
	InitManager()
	// DisposeManager 销毁管理器
	DisposeManager()
}
