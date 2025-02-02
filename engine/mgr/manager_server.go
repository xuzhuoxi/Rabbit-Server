// Package mgr
// Create on 2023/8/16
// @author xuzhuoxi
package mgr

import (
	"github.com/xuzhuoxi/Rabbit-Server/engine/server"
)

type IRabbitServerManager interface {
	SetReqVerify(reqVerify server.IReqVerify)
	SetReqVerifyNew(newReqVerify server.FuncNewIReqVerify)

	StartServers()
	StopServers()
	SaveServers()

	StartServer(id string)
	StopServer(id string)
	SaveServer(id string)

	ForEachServer(reverse bool, each func(index int, server server.IRabbitServer))
}

type serverMgr struct {
	Servers []server.IRabbitServer
}

func (o *serverMgr) SetReqVerify(reqVerify server.IReqVerify) {
	for index := range o.Servers {
		mgr, ok := o.Servers[index].GetExtensionManager()
		if !ok {
			return
		}
		mgr.SetCustomVerify(reqVerify)
	}
}

func (o *serverMgr) SetReqVerifyNew(newReqVerify server.FuncNewIReqVerify) {
	for index := range o.Servers {
		mgr, ok := o.Servers[index].GetExtensionManager()
		if !ok {
			return
		}
		mgr.SetCustomVerify(newReqVerify())
	}
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

func (o *serverMgr) ForEachServer(reverse bool, each func(index int, server server.IRabbitServer)) {
	if reverse {
		for index := len(o.Servers) - 1; index >= 0; index-- {
			each(index, o.Servers[index])
		}
	} else {
		for index := range o.Servers {
			each(index, o.Servers[index])
		}
	}
}
