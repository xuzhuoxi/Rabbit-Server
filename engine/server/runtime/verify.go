// Package runtime
// Create on 2025/4/23
// @author xuzhuoxi
package runtime

import (
	"github.com/xuzhuoxi/Rabbit-Server/engine/server"
)

type VerifySupport struct {
}

func (o *VerifySupport) VerifyExtension(container server.IRabbitExtensionContainer, extName string, pid string) (match server.IRabbitExtension, rsCode int32) {
	ext, ok := o.findExtension(container, extName)
	if !ok {
		return nil, server.CodeExtensionNotExist
	}
	if !ext.CheckProtoId(pid) { //有效性检查
		return nil, server.CodeProtoNotExist
	}
	if e, ok := ext.(server.IEnableExtension); ok && !e.Enable() {
		return nil, server.CodeExtensionDisable
	}
	return ext, server.CodeSuc
}

func (o *VerifySupport) findExtension(container server.IRabbitExtensionContainer, extName string) (pe server.IRabbitExtension, ok bool) {
	if pe, ok := container.GetExtension(extName).(server.IRabbitExtension); ok {
		return pe, true
	}
	return nil, false
}
