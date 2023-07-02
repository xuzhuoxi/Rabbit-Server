// Package startup
// Create on 2023/7/2
// @author xuzhuoxi
package startup

import (
	"github.com/xuzhuoxi/infra-go/eventx"
)

func NewStartupManager() IStartupManager {
	return &StartupManager{}
}

type IStartupModule interface {
	eventx.IEventDispatcher
	Name() string
	Start()
	Stop()
	Save()
}

type IStartupManager interface {
	eventx.IEventDispatcher
	// AppendModule 添加模块
	AppendModule(module IStartupModule)
	// StartUp 启动
	StartUp() bool
	// Save 保存状态
	Save() bool
	// ShutDown 停止
	ShutDown() bool
	// Reboot 重启
	Reboot() bool
}

type StartupManager struct {
	eventx.EventDispatcher
	modules     []IStartupModule
	index       int
	running     bool
	saving      bool
	rebooting   bool
	shutdowning bool
}

func (o *StartupManager) AppendModule(module IStartupModule) {
	if nil == module || o.running || o.saving || o.rebooting || o.shutdowning {
		return
	}
	o.modules = append(o.modules, module)
}

func (o *StartupManager) StartUp() bool {
	if o.running || o.rebooting {
		return false
	}
	o.running, o.index = true, 0
	o.startModule()
	return true
}

func (o *StartupManager) startModule() {
	if o.index >= len(o.modules) {
		o.DispatchEvent(EventOnManagerStartupFinish, o, nil)
		return
	}
	m := o.modules[o.index]
	m.OnceEventListener(EventOnModuleStarted, o.onModuleStarted)
	m.Start()
}

func (o *StartupManager) onModuleStarted(evd *eventx.EventData) {
	o.index += 1
	o.startModule()
}

func (o *StartupManager) Save() bool {
	if !o.running || o.saving {
		return false
	}
	o.saving, o.index = true, 0
	o.saveModule()
	return true
}

func (o *StartupManager) saveModule() {
	if o.index == len(o.modules) {
		o.saving = false
		o.DispatchEvent(EventOnManagerSaveFinish, o, nil)
		return
	}
	m := o.modules[o.index]
	m.OnceEventListener(EventOnModuleSaved, o.onModuleSaved)
	m.Save()
}

func (o *StartupManager) onModuleSaved(evd *eventx.EventData) {
	o.index += 1
	o.saveModule()
}

func (o *StartupManager) ShutDown() bool {
	if !o.running {
		return false
	}
	o.shutdowning = true
	o.prepareShutdown()
	return true
}

func (o *StartupManager) prepareShutdown() {
	o.OnceEventListener(EventOnManagerSaveFinish, o.onManagerShutdownSaved)
	if !o.saving {
		o.saving, o.index = true, 0
		o.saveModule()
	}
}

func (o *StartupManager) onManagerShutdownSaved(evd *eventx.EventData) {
	o.index = len(o.modules) - 1
	o.shutdownModule()
}

func (o *StartupManager) shutdownModule() {
	if o.index < 0 {
		o.running, o.shutdowning = false, false
		o.DispatchEvent(EventOnManagerShutdownFinish, o, nil)
		return
	}
	m := o.modules[o.index]
	m.OnceEventListener(EventOnModuleStopped, o.onModuleStopped)
	m.Stop()
}

func (o *StartupManager) onModuleStopped(evd *eventx.EventData) {
	o.index -= 1
	o.shutdownModule()
}

func (o *StartupManager) Reboot() bool {
	if !o.running || o.rebooting || o.shutdowning {
		return false
	}
	o.rebooting = true
	o.reboot()
	return true
}

func (o *StartupManager) reboot() {
	o.OnceEventListener(EventOnManagerShutdownFinish, o.onManagerShutdownFinish)
	if !o.shutdowning {
		o.shutdowning = true
		o.prepareShutdown()
	}
}

func (o *StartupManager) onManagerShutdownFinish(evd *eventx.EventData) {
	o.OnceEventListener(EventOnManagerStartupFinish, o.onManagerStartupFinish)
	o.running, o.index = true, 0
	o.startModule()
}

func (o *StartupManager) onManagerStartupFinish(evd *eventx.EventData) {
	o.rebooting = false
	o.DispatchEvent(EventOnManagerRebootFinish, o, nil)
}
