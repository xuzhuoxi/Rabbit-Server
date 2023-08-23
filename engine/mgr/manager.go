// Package mgr
// Create on 2023/6/14
// @author xuzhuoxi
package mgr

import (
	"github.com/xuzhuoxi/Rabbit-Server/engine/mmo"
	"github.com/xuzhuoxi/Rabbit-Server/engine/mmo/config"
	"github.com/xuzhuoxi/Rabbit-Server/engine/server"
	"github.com/xuzhuoxi/Rabbit-Server/engine/server/clock"
	"github.com/xuzhuoxi/infra-go/logx"
)

var (
	DefaultManager = NewRabbitManager()
)

func NewRabbitManager() IRabbitManager {
	return &RabbitManager{
		ConnManager:   &RabbitConnManager{},
		ServerManager: &serverMgr{},
	}
}

type IRabbitInitManager interface {
	// LoadRabbitConfig 加载配置
	LoadRabbitConfig(rootPath string) error
	// GetConfigs 取加载好的配置信息
	GetConfigs() (root server.CfgRabbitRoot, log server.CfgRabbitLog, server server.CfgRabbitServer, mmo config.MMOConfig)

	// InitLoggerManager 初始化Log管理器
	InitLoggerManager() logx.ILoggerManager
	// InitServers 初始化逻辑服务器
	InitServers()
	// CreateMMOWorld 创建MMO世界
	CreateMMOWorld() mmo.IMMOManager
	// InitClockManager 初始化时钟管理器
	InitClockManager()
}

type IRabbitManager interface {
	logx.ILoggerGetter
	GetInitManager() IRabbitInitManager
	GetLogManager() logx.ILoggerManager
	GetServerManager() IRabbitServerManager
	GetConnManager() IRabbitConnManager
	GetMMOManger() mmo.IMMOManager
	GetClockManger() clock.IRabbitClockManager
}

type RabbitManager struct {
	CfgRoot   *server.CfgRabbitRoot
	CfgLog    *server.CfgRabbitLog
	CfgServer *server.CfgRabbitServer
	CfgMMO    *config.MMOConfig
	CfgClock  *server.CfgClock

	LogManager    logx.ILoggerManager
	ServerManager *serverMgr
	ConnManager   *RabbitConnManager
	MMOManager    mmo.IMMOManager
	ClockManager  clock.IRabbitClockManager
}

func (o *RabbitManager) GetInitManager() IRabbitInitManager {
	return o
}

func (o *RabbitManager) GetLogManager() logx.ILoggerManager {
	return o.LogManager
}

func (o *RabbitManager) GetServerManager() IRabbitServerManager {
	return o.ServerManager
}

func (o *RabbitManager) GetMMOManger() mmo.IMMOManager {
	return o.MMOManager
}

func (o *RabbitManager) GetConnManager() IRabbitConnManager {
	return o.ConnManager
}

func (o *RabbitManager) GetClockManger() clock.IRabbitClockManager {
	return o.ClockManager
}

func (o *RabbitManager) GetLogger() logx.ILogger {
	if nil == o.CfgLog {
		return o.LogManager.FindLogger(logx.DefaultLoggerName)
	}
	return o.LogManager.FindLogger(o.CfgLog.Default)
}

// IRabbitInitManager ---------- ---------- ---------- ----------

func (o *RabbitManager) LoadRabbitConfig(rootPath string) error {
	cfgRoot, err1 := server.LoadRabbitRootConfig(rootPath)
	if nil != err1 {
		return err1
	}
	cfgLog, err2 := cfgRoot.LoadLogConfig()
	if nil != err2 {
		return err2
	}
	cfgServer, err3 := cfgRoot.LoadServerConfig()
	if nil != err3 {
		return err3
	}
	cfgMMO, err4 := cfgRoot.LoadMMOConfig()
	if nil != err4 {
		return err4
	}
	cfgClock, err5 := cfgRoot.LoadClockConfig()
	if nil != err5 {
		return err5
	}
	o.CfgRoot, o.CfgLog, o.CfgServer, o.CfgMMO, o.CfgClock = cfgRoot, cfgLog, cfgServer, cfgMMO, cfgClock
	return nil
}

func (o *RabbitManager) GetConfigs() (root server.CfgRabbitRoot, log server.CfgRabbitLog, server server.CfgRabbitServer, mmo config.MMOConfig) {
	return *o.CfgRoot, *o.CfgLog, *o.CfgServer, *o.CfgMMO
}

func (o *RabbitManager) InitLoggerManager() logx.ILoggerManager {
	o.initLogger()
	return o.LogManager
}

func (o *RabbitManager) InitServers() {
	o.initServers()
}

func (o *RabbitManager) CreateMMOWorld() mmo.IMMOManager {
	o.initMMOWorld()
	return o.MMOManager
}

func (o *RabbitManager) InitClockManager() {
	o.initClockManager()
}
