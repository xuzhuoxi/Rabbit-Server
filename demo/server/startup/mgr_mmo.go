// Package startup
// Create on 2023/7/2
// @author xuzhuoxi
package startup

import (
	"github.com/xuzhuoxi/Rabbit-Server/engine/mgr"
	"github.com/xuzhuoxi/Rabbit-Server/engine/mmo"
	"github.com/xuzhuoxi/infra-go/eventx"
	"github.com/xuzhuoxi/infra-go/serialx"
)

type mgrMMO struct {
	eventx.EventDispatcher
	mmoMgr mmo.IMMOManager
}

func (o *mgrMMO) Name() string {
	return "Manage MMO World"
}

func (o *mgrMMO) StartModule() {
	o.mmoMgr = mgr.DefaultManager.GetInitManager().BuildMMOEnv()
	o.DispatchEvent(serialx.EventOnStartupModuleStarted, o, nil)
}

func (o *mgrMMO) StopModule() {
	o.mmoMgr.DisposeManager()
	o.DispatchEvent(serialx.EventOnStartupModuleStopped, o, nil)
}

func (o *mgrMMO) SaveModule() {
	o.DispatchEvent(serialx.EventOnStartupModuleSaved, o, nil)
}
