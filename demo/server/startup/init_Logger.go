// Package startup
// Create on 2023/7/2
// @author xuzhuoxi
package startup

import (
	"github.com/xuzhuoxi/Rabbit-Server/engine/mgr"
	"github.com/xuzhuoxi/infra-go/eventx"
	"github.com/xuzhuoxi/infra-go/logx"
	"github.com/xuzhuoxi/infra-go/serialx"
)

type initRabbitLogger struct {
	eventx.EventDispatcher
	logMgr logx.ILoggerManager
}

func (o *initRabbitLogger) Name() string {
	return "Init Rabbit Logger Manager"
}

func (o *initRabbitLogger) StartModule() {
	o.logMgr = mgr.DefaultManager.GetInitManager().InitLoggerManager()
	o.DispatchEvent(serialx.EventOnStartupModuleStarted, o, nil)
}

func (o *initRabbitLogger) StopModule() {
	o.DispatchEvent(serialx.EventOnStartupModuleStopped, o, nil)
}

func (o *initRabbitLogger) SaveModule() {
	o.DispatchEvent(serialx.EventOnStartupModuleSaved, o, nil)
}
