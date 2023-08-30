// Package basis
// Created by xuzhuoxi
// on 2019-03-14.
// @author xuzhuoxi
package basis

type IEntityChild interface {
	// GetParent 父节点Id
	GetParent() string
	// IsNoneParent 判断父节点是否存在
	IsNoneParent() bool

	// SetParent 设置父节点Id
	SetParent(ownerId string)
	// ClearParent 清除父节点Id
	ClearParent()
}
