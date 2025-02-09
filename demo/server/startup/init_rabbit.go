// Package startup
// Create on 2023/7/2
// @author xuzhuoxi
package startup

import (
	"flag"
	"github.com/xuzhuoxi/Rabbit-Server/engine/mgr"
	"github.com/xuzhuoxi/infra-go/eventx"
	"github.com/xuzhuoxi/infra-go/serialx"
)

type initRabbitManager struct {
	eventx.EventDispatcher
}

func (o *initRabbitManager) Name() string {
	return "Init Rabbit"
}

func (o *initRabbitManager) StartModule() {
	err := mgr.DefaultManager.GetInitManager().LoadRabbitConfig(getConfigPath())
	if nil != err {
		panic(err)
	}
	o.DispatchEvent(serialx.EventOnStartupModuleStarted, o, nil)
}

func (o *initRabbitManager) StopModule() {
	o.DispatchEvent(serialx.EventOnStartupModuleStopped, o, nil)
}

func (o *initRabbitManager) SaveModule() {
	o.DispatchEvent(serialx.EventOnStartupModuleSaved, o, nil)
}

func getConfigPath() string {
	confPath := flag.String("conf", "rabbit.yaml", "RootConfig file for running")
	flag.Parse()
	return *confPath
}
