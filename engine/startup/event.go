// Package startup
// Create on 2023/7/2
// @author xuzhuoxi
package startup

const (
	EventOnModuleSaved   = "startup-Module:SaveFinish"
	EventOnModuleStarted = "startup-Module:StartFinish"
	EventOnModuleStopped = "startup-Module:StopFinish"

	EventOnManagerSaveFinish     = "startup-Manager:SaveFinish"
	EventOnManagerStartupFinish  = "startup-Manager:StartupFinish"
	EventOnManagerShutdownFinish = "startup-Manager:ShutdownFinish"
	EventOnManagerRebootFinish   = "startup-Manager:RebootFinish"
)
