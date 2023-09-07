// Package events
// Created by xuzhuoxi
// on 2019-03-08.
// @author xuzhuoxi
package events

import (
	"github.com/xuzhuoxi/Rabbit-Server/engine/mmo/basis"
	"github.com/xuzhuoxi/infra-go/encodingx"
)

const (
	// EventEntityVarChanged 实体变量更新事件
	// 事件数据格式：*basis.VarEventData
	EventEntityVarChanged = "Entity-VariableChanged"
	// EventEntityVarsChanged 实体变量批量更新事件
	// 事件数据格式：*basis.VarsEventData
	EventEntityVarsChanged = "Entity-VariablesChanged"
)

type VarEventData struct {
	Entity basis.IEntity
	Key    string
	Value  interface{}
}

type VarsEventData struct {
	Entity basis.IEntity
	Vars   encodingx.IKeyValue
}
