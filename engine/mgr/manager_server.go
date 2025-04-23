// Package mgr
// Create on 2023/8/16
// @author xuzhuoxi
package mgr

import (
	"github.com/xuzhuoxi/Rabbit-Server/engine/server"
)

type IRabbitServerManager interface {
	// StartServers
	// 启用全部服务器
	StartServers()
	// StopServers
	// 停用全部服务器
	StopServers()
	// SaveServers
	// 保存全部服务器数据
	SaveServers()

	// StartServer
	// 启用指定服务器
	// id 服务器ID
	StartServer(id string)
	// StopServer
	// 停用指定服务器
	// id 服务器ID
	StopServer(id string)
	// SaveServer
	// 保存指定服务器数据
	SaveServer(id string)

	// ForEachServer
	// 遍历服务器
	// reverse 是否倒序遍历
	// each 遍历回调
	ForEachServer(reverse bool, each func(index int, server server.IRabbitServer))
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
