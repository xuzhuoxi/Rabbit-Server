// Package loader
// Create on 2023/6/14
// @author xuzhuoxi
package loader

import (
	"github.com/xuzhuoxi/Rabbit-Server/engine/mmo"
	"github.com/xuzhuoxi/Rabbit-Server/engine/mmo/config"
	"github.com/xuzhuoxi/Rabbit-Server/engine/server"
	"github.com/xuzhuoxi/infra-go/logx"
)

var (
	DefaultLoader = NewRabbitLoader()
)

func NewRabbitLoader() IRabbitLoader {
	return &RabbitLoader{}
}

type IRabbitLoader interface {
	logx.ILoggerGetter
	LoadConfig(rootPath string) error
	InitLoader() (logManager logx.ILoggerManager, mmoManager mmo.IMMOManager)

	StartServers()
	StopServers()
	SaveServers()

	StartServer(id string)
	StopServer(id string)
	SaveServer(id string)
}

type RabbitLoader struct {
	CfgLog    *server.CfgRabbitLog
	CfgServer *server.CfgRabbitServer
	CfgMMO    *config.MMOConfig

	LogManager logx.ILoggerManager
	MMOManager mmo.IMMOManager
	Servers    []server.IRabbitServer
}

func (o *RabbitLoader) GetLogger() logx.ILogger {
	if nil == o.CfgLog {
		return o.LogManager.FindLogger(logx.DefaultLoggerName)
	}
	return o.LogManager.FindLogger(o.CfgLog.Default)
}

func (o *RabbitLoader) LoadConfig(rootPath string) error {
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
	o.CfgLog, o.CfgServer, o.CfgMMO = cfgLog, cfgServer, cfgMMO
	return nil
}

func (o *RabbitLoader) InitLoader() (logManager logx.ILoggerManager, mmoManager mmo.IMMOManager) {
	o.initLogger()
	o.initMMO()
	o.initServers()
	return o.LogManager, o.MMOManager
}

func (o *RabbitLoader) initLogger() {
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

func (o *RabbitLoader) initMMO() {
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

func (o *RabbitLoader) initServers() {
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

func (o *RabbitLoader) StartServer(id string) {
	for _, s := range o.Servers {
		if s.GetId() == id {
			go s.Start()
		}
	}
}

func (o *RabbitLoader) StopServer(id string) {
	for _, s := range o.Servers {
		if s.GetId() == id {
			s.Stop()
		}
	}
}

func (o *RabbitLoader) SaveServer(id string) {
	for _, s := range o.Servers {
		if s.GetId() == id {
			s.Save()
		}
	}
}

func (o *RabbitLoader) StartServers() {
	for _, s := range o.Servers {
		go s.Start()
	}
}

func (o *RabbitLoader) StopServers() {
	for index := len(o.Servers) - 1; index >= 0; index -= 1 {
		o.Servers[index].Stop()
	}
}

func (o *RabbitLoader) SaveServers() {
	for _, s := range o.Servers {
		s.Save()
	}
}
