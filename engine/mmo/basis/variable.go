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
	GetVar(key string) (interface{}, bool)
	CheckVar(key string) bool
	Vars() encodingx.IKeyValue

	SetVar(kv string, value interface{})
	SetVars(kv encodingx.IKeyValue)
}
