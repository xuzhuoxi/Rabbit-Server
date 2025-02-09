// Package startup
// Create on 2023/7/2
// @author xuzhuoxi
package startup

import (
	"github.com/xuzhuoxi/Rabbit-Server/engine/mgr"
	"github.com/xuzhuoxi/infra-go/eventx"
	"github.com/xuzhuoxi/infra-go/serialx"
)

type mgrServer struct {
	eventx.EventDispatcher
}

func (o *mgrServer) Name() string {
	return "Manage Servers"
}

func (o *mgrServer) StartModule() {
	mgr.DefaultManager.GetServerManager().StartServers()
	o.DispatchEvent(serialx.EventOnStartupModuleStarted, o, nil)
}

func (o *mgrServer) StopModule() {
	mgr.DefaultManager.GetServerManager().StopServers()
	o.DispatchEvent(serialx.EventOnStartupModuleStopped, o, nil)
}

func (o *mgrServer) SaveModule() {
	o.DispatchEvent(serialx.EventOnStartupModuleSaved, o, nil)
}
