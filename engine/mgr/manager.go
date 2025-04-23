// Package mgr
// Create on 2023/6/14
// @author xuzhuoxi
package mgr

import (
	"github.com/xuzhuoxi/Rabbit-Server/engine"
	"github.com/xuzhuoxi/Rabbit-Server/engine/clock"
	"github.com/xuzhuoxi/Rabbit-Server/engine/config"
	"github.com/xuzhuoxi/Rabbit-Server/engine/mmo"
	mmoConfig "github.com/xuzhuoxi/Rabbit-Server/engine/mmo/config"
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
	GetConfigs() (root server.CfgRabbitRoot, log config.CfgRabbitLog, clock clock.CfgClock,
		mmo mmoConfig.MMOConfig, server config.CfgRabbitServer, verify config.CfgVerifyRoot)

	// InitLoggerManager 初始化Log管理器
	InitLoggerManager() logx.ILoggerManager
	// InitServers 初始化逻辑服务器
	InitServers()
	// BuildMMOEnv 创建MMO世界
	BuildMMOEnv() mmo.IMMOManager
	// InitClockManager 初始化时钟管理器
	InitClockManager()
}

type IRabbitManager interface {
	logx.ILoggerGetter
	GetInitManager() IRabbitInitManager        // 获取初始化管理器
	GetLogManager() logx.ILoggerManager        // 获取日志管理器
	GetClockManger() clock.IRabbitClockManager // 获取时钟管理器
	GetMMOManger() mmo.IMMOManager             // 获取MMO世界管理器
	GetServerManager() IRabbitServerManager    // 获取逻辑服务器管理器
	GetConnManager() IRabbitConnManager        // 获取连接管理器
}

type RabbitManager struct {
	CfgRoot   *server.CfgRabbitRoot
	CfgLog    *config.CfgRabbitLog    // 日志配置
	CfgClock  *clock.CfgClock         // 时钟配置
	CfgMMO    *mmoConfig.MMOConfig    // MMO世界配置
	CfgServer *config.CfgRabbitServer // 逻辑服务器配置
	CfgVerify *config.CfgVerifyRoot   // 验证配置

	LogManager   logx.ILoggerManager
	ClockManager clock.IRabbitClockManager

	MMOManager    mmo.IMMOManager
	ServerManager *serverMgr
	ConnManager   *RabbitConnManager
}

func (o *RabbitManager) GetInitManager() IRabbitInitManager {
	return o
}

func (o *RabbitManager) GetLogManager() logx.ILoggerManager {
	return o.LogManager
}

func (o *RabbitManager) GetClockManger() clock.IRabbitClockManager {
	return o.ClockManager
}

func (o *RabbitManager) GetMMOManger() mmo.IMMOManager {
	return o.MMOManager
}

func (o *RabbitManager) GetServerManager() IRabbitServerManager {
	return o.ServerManager
}

func (o *RabbitManager) GetConnManager() IRabbitConnManager {
	return o.ConnManager
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
	cfgClock, err3 := cfgRoot.LoadClockConfig()
	if nil != err3 {
		return err3
	}
	cfgMMO, err4 := cfgRoot.LoadMMOConfig()
	if nil != err4 {
		return err4
	}
	cfgMMO.MergeRelationToTags()
	cfgServer, err5 := cfgRoot.LoadServerConfig()
	if nil != err5 {
		return err5
	}
	cfgVerify, err6 := cfgRoot.LoadVerifyConfig()
	if nil != err6 {
		return err6
	}
	o.CfgRoot, o.CfgLog, o.CfgClock, o.CfgMMO, o.CfgServer, o.CfgVerify = cfgRoot, cfgLog, cfgClock, cfgMMO, cfgServer, cfgVerify
	return nil
}

func (o *RabbitManager) GetConfigs() (root server.CfgRabbitRoot, log config.CfgRabbitLog, clock clock.CfgClock,
	mmo mmoConfig.MMOConfig, server config.CfgRabbitServer, verify config.CfgVerifyRoot) {
	return *o.CfgRoot, *o.CfgLog, *o.CfgClock, *o.CfgMMO, *o.CfgServer, *o.CfgVerify
}

func (o *RabbitManager) InitLoggerManager() logx.ILoggerManager {
	o.initLogger()
	return o.LogManager
}

func (o *RabbitManager) InitServers() {
	o.initServers()
}

func (o *RabbitManager) BuildMMOEnv() mmo.IMMOManager {
	o.initMMOEnv()
	return o.MMOManager
}

func (o *RabbitManager) InitClockManager() {
	o.initClockManager()
}
