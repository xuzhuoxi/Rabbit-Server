// Package startup
// Create on 2023/7/2
// @author xuzhuoxi
package startup

import (
	"github.com/xuzhuoxi/Rabbit-Server/engine/mgr"
	"github.com/xuzhuoxi/infra-go/eventx"
	"github.com/xuzhuoxi/infra-go/serialx"
)

type initServer struct {
	eventx.EventDispatcher
}

func (o *initServer) Name() string {
	return "Init Servers"
}

func (o *initServer) StartModule() {
	mgr.DefaultManager.GetInitManager().InitServers()
	o.DispatchEvent(serialx.EventOnStartupModuleStarted, o, nil)
}

func (o *initServer) StopModule() {
	o.DispatchEvent(serialx.EventOnStartupModuleStopped, o, nil)
}

func (o *initServer) SaveModule() {
	o.DispatchEvent(serialx.EventOnStartupModuleSaved, o, nil)
}
