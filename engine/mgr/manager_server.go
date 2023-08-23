// Package mgr
// Create on 2023/8/16
// @author xuzhuoxi
package mgr

import "github.com/xuzhuoxi/Rabbit-Server/engine/server"

type IRabbitServerManager interface {
	StartServers()
	StopServers()
	SaveServers()

	StartServer(id string)
	StopServer(id string)
	SaveServer(id string)
}

type serverMgr struct {
	Servers []server.IRabbitServer
}

func (o *serverMgr) StartServer(id string) {
	for _, s := range o.Servers {
		if s.GetId() == id {
			go s.Start()
		}
	}
}

func (o *serverMgr) StopServer(id string) {
	for _, s := range o.Servers {
		if s.GetId() == id {
			s.Stop()
		}
	}
}

func (o *serverMgr) SaveServer(id string) {
	for _, s := range o.Servers {
		if s.GetId() == id {
			s.Save()
		}
	}
}

func (o *serverMgr) StartServers() {
	for _, s := range o.Servers {
		go s.Start()
	}
}

func (o *serverMgr) StopServers() {
	for index := len(o.Servers) - 1; index >= 0; index -= 1 {
		o.Servers[index].Stop()
	}
}

func (o *serverMgr) SaveServers() {
	for _, s := range o.Servers {
		s.Save()
	}
}
