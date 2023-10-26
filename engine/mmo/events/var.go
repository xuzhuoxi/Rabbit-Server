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
	// EventEntityVarMod 实体变量更新事件
	// 事件数据格式：*VarModEventData
	EventEntityVarMod = "MMO.EventEntityVarMod"
	// EventEntityVarsMod 实体变量批量更新事件
	// 事件数据格式：*VarsModEventData
	EventEntityVarsMod = "MMO.EventEntityVarsMod"
	// EventEntityVarDel 实体变量删除事件
	// 事件数据格式：*VarDelEventData
	EventEntityVarDel = "MMO.EventEntityVarDel"
	// EventEntityVarsDel 实体变量批量删除事件
	// 事件数据格式：*VarsDelEventData
	EventEntityVarsDel = "MMO.EventEntityVarsDel"
)

type VarModEventData struct {
	Entity basis.IEntity
	Key    string
	Value  interface{}
}

type VarsModEventData struct {
	Entity  basis.IEntity
	VarSet  encodingx.IKeyValue
	VarKeys []string
}

type VarDelEventData struct {
	Entity basis.IEntity
	Key    string
}

type VarsDelEventData struct {
	Entity  basis.IEntity
	VarKeys []string
}
