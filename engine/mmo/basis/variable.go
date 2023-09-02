// Package basis
// Created by xuzhuoxi
// on 2019-03-14.
// @author xuzhuoxi
package basis

import (
	"github.com/xuzhuoxi/infra-go/encodingx"
	"github.com/xuzhuoxi/infra-go/eventx"
)

const (
	VarKeyUserNick = "nick"
	VarKeyUserPos  = "pos"
)

func NewVarSet() encodingx.IKeyValue {
	return encodingx.NewCodingMap()
}

// IVariableSupport 变量列表
type IVariableSupport interface {
	eventx.IEventDispatcher
	// GetVar 取Tag数据
	GetVar(key string) (interface{}, bool)
	// CheckVar 检查Tag是否存在
	CheckVar(key string) bool
	// Vars 取Tag数据集合
	Vars() encodingx.IKeyValue

	// SetVar 设置Tag
	SetVar(kv string, value interface{})
	// SetVars 批量设置Tag
	SetVars(kv encodingx.IKeyValue)
}
