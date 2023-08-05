// Package mgr
// Create on 2023/6/14
// @author xuzhuoxi
package mgr

import (
	"github.com/xuzhuoxi/Rabbit-Server/engine/mmo"
	"github.com/xuzhuoxi/Rabbit-Server/engine/mmo/config"
	"github.com/xuzhuoxi/Rabbit-Server/engine/server"
	"github.com/xuzhuoxi/infra-go/logx"
)

var (
	DefaultManager = NewRabbitManager()
)

func NewRabbitManager() IRabbitManager {
	return &RabbitManager{}
}

type IRabbitManager interface {
	logx.ILoggerGetter
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

	StartServers()
	StopServers()
	SaveServers()

	StartServer(id string)
	StopServer(id string)
	SaveServer(id string)
}

type RabbitManager struct {
	CfgRoot   *server.CfgRabbitRoot
	CfgLog    *server.CfgRabbitLog
	CfgServer *server.CfgRabbitServer
	CfgMMO    *config.MMOConfig

	LogManager logx.ILoggerManager
	MMOManager mmo.IMMOManager
	Servers    []server.IRabbitServer
}

func (o *RabbitManager) GetLogger() logx.ILogger {
	if nil == o.CfgLog {
		return o.LogManager.FindLogger(logx.DefaultLoggerName)
	}
	return o.LogManager.FindLogger(o.CfgLog.Default)
}

func (o *RabbitManager) LoadRabbitConfig(rootPath string) error {
	cfgRoot, err := server.LoadRabbitRootConfig(rootPath)
	if nil != err {
		return err
	}
	cfgLog, err := cfgRoot.LoadLogConfig()
	if nil != err {
		return err
	}
	cfgServer, err := cfgRoot.LoadServerConfig()
	if nil != err {
		return err
	}
	cfgMMO, err := cfgRoot.LoadMMOConfig()
	if nil != err {
		return err
	}
	o.CfgRoot, o.CfgLog, o.CfgServer, o.CfgMMO = cfgRoot, cfgLog, cfgServer, cfgMMO
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

func (o *RabbitManager) initLogger() {
	o.LogManager = logx.DefaultLoggerManager
	if nil == o.CfgLog {
		logger := logx.NewLogger()
		o.LogManager.RegisterLogger(logx.DefaultLoggerName, logger)
		o.LogManager.SetDefault(logx.DefaultLoggerName)
	} else {
		for _, log := range o.CfgLog.Logs {
			o.LogManager.GenLogger(log.Name, log.Conf)
		}
		o.LogManager.SetDefault(o.CfgLog.Default)
	}
}

func (o *RabbitManager) initServers() {
	if o.CfgServer == nil || len(o.CfgServer.Servers) == 0 {
		return
	}
	for _, cfgServerItem := range o.CfgServer.Servers {
		s, err := server.NewRabbitServer(cfgServerItem.Name)
		if nil != err {
			panic(err)
		}
		o.Servers = append(o.Servers, s)
		if cfgServerItem.LogRef == "" {
			s.SetLogger(o.LogManager.GetDefaultLogger())
		} else {
			s.SetLogger(o.LogManager.FindLogger(cfgServerItem.LogRef))
		}
		s.Init(cfgServerItem)
	}
}

func (o *RabbitManager) initMMOWorld() {
	if o.CfgMMO == nil {
		return
	}
	o.MMOManager = mmo.NewMMOManager()
	o.MMOManager.InitManager()
	if nil == o.CfgMMO.Log && len(o.CfgMMO.LogRef) == 0 {
		o.MMOManager.SetLogger(o.LogManager.GetDefaultLogger())
	} else {
		if o.CfgMMO.Log != nil {
			logger := logx.NewLogger()
			logger.SetConfig(o.CfgMMO.Log.ToLogConfig())
			o.MMOManager.SetLogger(logger)
		} else {
			o.MMOManager.SetLogger(o.LogManager.FindLogger(o.CfgMMO.LogRef))
		}
	}
	o.MMOManager.GetEntityManager().ConstructWorldDefault(o.CfgMMO)
}

func (o *RabbitManager) StartServer(id string) {
	for _, s := range o.Servers {
		if s.GetId() == id {
			go s.Start()
		}
	}
}

func (o *RabbitManager) StopServer(id string) {
	for _, s := range o.Servers {
		if s.GetId() == id {
			s.Stop()
		}
	}
}

func (o *RabbitManager) SaveServer(id string) {
	for _, s := range o.Servers {
		if s.GetId() == id {
			s.Save()
		}
	}
}

func (o *RabbitManager) StartServers() {
	for _, s := range o.Servers {
		go s.Start()
	}
}

func (o *RabbitManager) StopServers() {
	for index := len(o.Servers) - 1; index >= 0; index -= 1 {
		o.Servers[index].Stop()
	}
}

func (o *RabbitManager) SaveServers() {
	for _, s := range o.Servers {
		s.Save()
	}
}
