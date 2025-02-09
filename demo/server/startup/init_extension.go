// Package startup
// Create on 2023/7/2
// @author xuzhuoxi
package startup

import (
	"github.com/xuzhuoxi/Rabbit-Server/demo/server/extension"
	"github.com/xuzhuoxi/infra-go/eventx"
	"github.com/xuzhuoxi/infra-go/serialx"
)

type initExtensions struct {
	eventx.EventDispatcher
}

func (o *initExtensions) Name() string {
	return "Register Extensions"
}

func (o *initExtensions) StartModule() {
	extension.Register()
	o.DispatchEvent(serialx.EventOnStartupModuleStarted, o, nil)
}

func (o *initExtensions) StopModule() {
	o.DispatchEvent(serialx.EventOnStartupModuleStopped, o, nil)
}

func (o *initExtensions) SaveModule() {
	o.DispatchEvent(serialx.EventOnStartupModuleSaved, o, nil)
}
