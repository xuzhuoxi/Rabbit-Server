// Package basis
// Created by xuzhuoxi
// on 2019-03-08.
// @author xuzhuoxi
package basis

import "github.com/xuzhuoxi/infra-go/encodingx"

const (
	EventEntityVarChanged = "Entity-VariableChanged"
	EventEntityVarDeleted = "Entity-VariableDeleted"
)

const (
	EventManagerVarChanged = "Manager-VariableChanged"
	EventManagerVarDeleted = "Manager-VariableDeleted"
)

type ManagerVarEventData struct {
	Entity IEntity
	Data   encodingx.IKeyValue
}
