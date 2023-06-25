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
	LoadConfig(cfgPath string) error

	InitServers()
	StartServers()
	StopServers()
	SaveServers()

	StartServer(id string)
	StopServer(id string)
	SaveServer(id string)
}

type RabbitLoader struct {
	logx.LoggerSupport
	ConfigServer *server.CfgRabbitServerConfig
	ConfigMMO    *config.MMOConfig

	MMOManager mmo.IMMOManager
	Servers    []server.IRabbitServer
}

func (o *RabbitLoader) LoadConfig(cfgPath string) error {
	cfgServer, err := server.PauseServerConfig(cfgPath)
	if nil != err {
		return err
	}
	var cfgMMO *config.MMOConfig
	if cfgServer.MMO != "" {
		cfgMMO, err = config.ParseByYamlPath(cfgServer.MMO)
		if nil != err {
			return err
		}
	}
	o.ConfigServer, o.ConfigMMO = cfgServer, cfgMMO
	return nil
}

func (o *RabbitLoader) InitServers() {
	o.initLogger()
	o.initMMO()
	o.initServers()
}

func (o *RabbitLoader) initLogger() {
	logger := logx.NewLogger()
	o.SetLogger(logger)
	if o.ConfigServer == nil || o.ConfigServer.Logger == nil {
		return
	}
	logger.SetConfig(o.ConfigServer.Logger.ToLogConfig())
}

func (o *RabbitLoader) initMMO() {
	if o.ConfigMMO == nil {
		return
	}
	o.MMOManager = mmo.NewMMOManager()
	o.MMOManager.InitManager()
	cfgLog := o.ConfigMMO.Log
	if nil == cfgLog {
		o.MMOManager.SetLogger(logx.DefaultLogger())
	} else {
		logger := logx.NewLogger()
		logger.SetConfig(cfgLog.ToLogConfig())
		o.MMOManager.SetLogger(logger)
	}
	o.MMOManager.GetEntityManager().ConstructWorldDefault(o.ConfigMMO)
}

func (o *RabbitLoader) initServers() {
	if o.ConfigServer == nil || len(o.ConfigServer.Servers) == 0 {
		return
	}
	for _, cfgServer := range o.ConfigServer.Servers {
		s, err := server.NewRabbitServer(cfgServer.Name)
		if nil != err {
			panic(err)
		}
		o.Servers = append(o.Servers, s)
		s.Init(cfgServer)
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
