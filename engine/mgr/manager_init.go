// Package mgr
// Create on 2023/8/16
// @author xuzhuoxi
package mgr

import (
	"github.com/xuzhuoxi/Rabbit-Server/engine/clock"
	"github.com/xuzhuoxi/Rabbit-Server/engine/mmo"
	"github.com/xuzhuoxi/Rabbit-Server/engine/server"
	"github.com/xuzhuoxi/infra-go/logx"
)

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
	var servers []server.IRabbitServer
	for _, cfgServerItem := range o.CfgServer.Servers {
		s, err := server.NewRabbitServer(cfgServerItem.Name)
		if nil != err {
			panic(err)
		}
		servers = append(servers, s)
		if cfgServerItem.LogRef == "" {
			s.SetLogger(o.LogManager.GetDefaultLogger())
		} else {
			s.SetLogger(o.LogManager.FindLogger(cfgServerItem.LogRef))
		}
		s.Init(cfgServerItem)
		if connSet, ok := s.GetConnSet(); ok {
			o.ConnManager.AddConnSet(s.GetName(), connSet)
		}
	}
	o.ServerManager.Servers = servers
}

func (o *RabbitManager) initMMOEnv() {
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
	o.MMOManager.BuildEnv(o.CfgMMO)
}

func (o *RabbitManager) initClockManager() {
	if o.CfgMMO == nil {
		return
	}
	o.ClockManager = clock.NewIRabbitClockManager()
	err := o.ClockManager.Init(o.CfgClock)
	if nil != err {
		o.LogManager.Warnln("[RabbitManager.initClockManager]", err)
		return
	}
}
