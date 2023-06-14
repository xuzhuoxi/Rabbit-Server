// Package loader
// Create on 2023/6/14
// @author xuzhuoxi
package loader

import (
	"github.com/xuzhuoxi/Rabbit-Server/engine/mmo"
	"github.com/xuzhuoxi/Rabbit-Server/engine/mmo/config"
	"github.com/xuzhuoxi/Rabbit-Server/engine/server"
	"github.com/xuzhuoxi/Rabbit-Server/engine/server/rabbit"
	"github.com/xuzhuoxi/infra-go/logx"
)

var (
	DefaultLoader = NewRabbitLoader()
)

func NewRabbitLoader() IRabbitLoader {
	return &RabbitLoader{}
}

type IRabbitLoader interface {
	LoadConfig(cfgPath string) error
	InitServer()
	StartServer()
	StopServer()
}

type RabbitLoader struct {
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

func (o *RabbitLoader) InitServer() {
	o.initMMO()
	o.initServers()
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
		logger.SetConfig(logx.LogConfig{Type: cfgLog.LogType, Level: cfgLog.LogLevel,
			FilePath: cfgLog.GetLogPath(), MaxSize: cfgLog.MaxSize()})
		o.MMOManager.SetLogger(logger)
	}
	o.MMOManager.GetEntityManager().ConstructWorldDefault(o.ConfigMMO)
}

func (o *RabbitLoader) initServers() {
	if o.ConfigServer == nil || len(o.ConfigServer.Servers) == 0 {
		return
	}
	for _, cfgServer := range o.ConfigServer.Servers {
		s := rabbit.NewRabbitServer(cfgServer)
		o.Servers = append(o.Servers, s)
		s.Init()
	}
}

func (o *RabbitLoader) StartServer() {
	for _, s := range o.Servers {
		go s.Start()
	}
}

func (o *RabbitLoader) StopServer() {
	for index := len(o.Servers) - 1; index >= 0; index -= 1 {
		o.Servers[index].Stop()
	}
}
