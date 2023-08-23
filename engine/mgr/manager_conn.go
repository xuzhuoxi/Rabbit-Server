// Package mgr
// Create on 2023/8/6
// @author xuzhuoxi
package mgr

import (
	"github.com/xuzhuoxi/infra-go/netx"
	"sync"
)

type IRabbitConnManager interface {
	// CloseConn
	// 关闭指定连接
	CloseConn(address string) (err error, ok bool)
	// FindConn
	// 查找连接
	FindConn(address string) (conn netx.IServerConn, ok bool)
}

type IRabbitConnManagerMod interface {
	// AddConnSet
	// 添加一个连接集合
	AddConnSet(named string, set netx.IServerConnSet) bool
	// RemoveConnSet
	// 移除一个连接集合
	RemoveConnSet(named string) bool
}

type connSetItem struct {
	Name    string
	ConnSet netx.IServerConnSet
}

type RabbitConnManager struct {
	SetItems []*connSetItem
	Lock     sync.RWMutex
}

func (o *RabbitConnManager) CloseConn(address string) (err error, ok bool) {
	o.Lock.RLock()
	defer o.Lock.RUnlock()
	if len(o.SetItems) == 0 || len(address) == 0 {
		return nil, false
	}
	for _, set := range o.SetItems {
		_, ok1 := set.ConnSet.FindConnection(address)
		if ok1 {
			return set.ConnSet.CloseConnection(address)
		}
	}
	return nil, false
}

func (o *RabbitConnManager) FindConn(address string) (conn netx.IServerConn, ok bool) {
	o.Lock.RLock()
	defer o.Lock.RUnlock()
	if len(o.SetItems) == 0 || len(address) == 0 {
		return nil, false
	}
	for _, set := range o.SetItems {
		rs, ok1 := set.ConnSet.FindConnection(address)
		if ok1 {
			return rs, true
		}
	}
	return nil, false
}

func (o *RabbitConnManager) AddConnSet(named string, set netx.IServerConnSet) bool {
	if len(named) == 0 || nil == set {
		return false
	}
	o.Lock.Lock()
	defer o.Lock.Unlock()
	for index := range o.SetItems {
		if o.SetItems[index].Name == named {
			return false
		}
	}
	o.SetItems = append(o.SetItems, &connSetItem{Name: named, ConnSet: set})
	return true
}

func (o *RabbitConnManager) RemoveConnSet(named string) bool {
	if len(named) == 0 {
		return false
	}
	o.Lock.Lock()
	defer o.Lock.Unlock()
	for index := range o.SetItems {
		if o.SetItems[index].Name == named {
			o.SetItems = append(o.SetItems[:index], o.SetItems[index+1:]...)
			return true
		}
	}
	return false
}
