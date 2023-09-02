// Package basis
// Created by xuzhuoxi
// on 2019-03-08.
// @author xuzhuoxi
package basis

import "github.com/xuzhuoxi/infra-go/encodingx"

const (
	// EventEntityVarChanged 实体变量更新事件
	// 事件数据格式：*basis.VarEventData
	EventEntityVarChanged = "Entity-VariableChanged"
	// EventEntityVarsChanged 实体变量批量更新事件
	// 事件数据格式：*basis.VarsEventData
	EventEntityVarsChanged = "Entity-VariablesChanged"
)

const (
	// EventManagerVarChanged 管理器抛出的实体变量更新事件
	// 事件数据格式：*basis.VarEventData
	EventManagerVarChanged = "Manager-VariableChanged"
	// EventManagerVarsChanged 管理器抛出的实体变量批量更新事件
	// 事件数据格式：*basis.VarsEventData
	EventManagerVarsChanged = "Manager-VariablesChanged"
)

type VarEventData struct {
	Entity IEntity
	Key    string
	Value  interface{}
}

type VarsEventData struct {
	Entity IEntity
	Vars   encodingx.IKeyValue
}
