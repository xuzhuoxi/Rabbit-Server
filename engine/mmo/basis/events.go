// Package basis
// Created by xuzhuoxi
// on 2019-03-08.
// @author xuzhuoxi
package basis

import "github.com/xuzhuoxi/infra-go/encodingx"

const (
	EventEntityVarChanged  = "Entity-VariableChanged"
	EventManagerVarChanged = "Manager-VariableChanged"
)

type ManagerVarEventData struct {
	Entity IEntity
	Data   encodingx.IKeyValue
}
