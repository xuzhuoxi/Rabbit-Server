// Package verify
// Create on 2025/4/24
// @author xuzhuoxi
package verify

import "github.com/xuzhuoxi/Rabbit-Server/engine/server"

func NewExtensionVerifyItem(container server.IRabbitExtensionContainer) *ExtensionVerifyItem {
	return &ExtensionVerifyItem{Container: container}
}

type ExtensionVerifyItem struct {
	Container server.IRabbitExtensionContainer
}

func (o *ExtensionVerifyItem) Verify(extName string, pid string, uid string) (rsCode int32) {
	ext, ok := o.findExtension(extName)
	if !ok {
		return server.CodeExtensionNotExist
	}
	if !ext.CheckProtoId(pid) { //有效性检查
		return server.CodeProtoNotExist
	}
	if e, ok := ext.(server.IEnableExtension); ok && !e.Enable() {
		return server.CodeExtensionDisable
	}
	return server.CodeSuc
}

func (o *ExtensionVerifyItem) findExtension(extName string) (re server.IRabbitExtension, ok bool) {
	if pe, ok := o.Container.GetExtension(extName).(server.IRabbitExtension); ok {
		return pe, true
	}
	return nil, false
}
