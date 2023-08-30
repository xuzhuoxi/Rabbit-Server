// Package basis
// Created by xuzhuoxi
// on 2019-03-08.
// @author xuzhuoxi
package basis

import "github.com/xuzhuoxi/infra-go/encodingx"

const (
	EventEntityVarChanged  = "Entity-VariableChanged"
	EventEntityVarsChanged = "Entity-VariablesChanged"
)

const (
	EventManagerVarChanged  = "Manager-VariableChanged"
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
