// Package extension
// Created by xuzhuoxi
// on 2019-05-19.
// @author xuzhuoxi
//
package extension

import (
	"github.com/xuzhuoxi/Rabbit-Server/engine/server"
	"github.com/xuzhuoxi/infra-go/lang"
)

var (
	DefaultRequestPool  = NewRequestPool()
	DefaultResponsePool = NewResponsePool()
	DefaultNotifyPool   = NewNotifyPool()
)

func init() {
	DefaultRequestPool.Register(func() server.IExtensionRequest {
		return NewSockRequest()
	})
	DefaultResponsePool.Register(func() server.IExtensionResponse {
		return NewSockResponse()
	})
	DefaultNotifyPool.Register(func() server.IExtensionNotify {
		return NewSockNotify()
	})
}

//--------------------------------------------

func NewRequestPool() server.IRequestPool { return &reqPool{pool: lang.NewObjectPoolSync()} }

func NewResponsePool() server.IResponsePool {
	return &respPool{pool: lang.NewObjectPoolSync()}
}

func NewNotifyPool() server.INotifyPool {
	return &notifyPool{pool: lang.NewObjectPoolSync()}
}

type reqPool struct {
	pool lang.IObjectPool
}

func (p *reqPool) Register(newFunc func() server.IExtensionRequest) {
	p.pool.Register(func() interface{} {
		return newFunc()
	}, func(instance interface{}) bool {
		return nil != instance
	})
}

func (p *reqPool) GetInstance() server.IExtensionRequest {
	return p.pool.GetInstance().(server.IExtensionRequest)
}

func (p *reqPool) Recycle(instance server.IExtensionRequest) bool {
	return p.pool.Recycle(instance)
}

type respPool struct {
	pool lang.IObjectPool
}

func (p *respPool) Register(newFunc func() server.IExtensionResponse) {
	p.pool.Register(func() interface{} {
		return newFunc()
	}, func(instance interface{}) bool {
		return nil != instance
	})
}

func (p *respPool) GetInstance() server.IExtensionResponse {
	return p.pool.GetInstance().(server.IExtensionResponse)
}

func (p *respPool) Recycle(instance server.IExtensionResponse) bool {
	return p.pool.Recycle(instance)
}

type notifyPool struct {
	pool lang.IObjectPool
}

func (p *notifyPool) Register(newFunc func() server.IExtensionNotify) {
	p.pool.Register(func() interface{} {
		return newFunc()
	}, func(instance interface{}) bool {
		return nil != instance
	})
}

func (p *notifyPool) GetInstance() server.IExtensionNotify {
	return p.pool.GetInstance().(server.IExtensionNotify)
}

func (p *notifyPool) Recycle(instance server.IExtensionNotify) bool {
	return p.pool.Recycle(instance)
}
