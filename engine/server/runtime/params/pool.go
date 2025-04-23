// Package extension
// Created by xuzhuoxi
// on 2019-05-19.
// @author xuzhuoxi
//
package params

import (
	"github.com/xuzhuoxi/Rabbit-Server/engine/server"
	"sync"
)

var (
	DefaultRequestPool  = NewRequestPool(NewISockRequest)
	DefaultResponsePool = NewResponsePool(NewISockResponse)
	DefaultNotifyPool   = NewNotifyPool(NewISockNotify)
)

func NewRequestPool(newFunc func() server.IExtensionRequest) server.IRequestPool {
	return &reqPool{pool: &sync.Pool{
		New: func() interface{} { return newFunc() },
	}}
}

type reqPool struct {
	pool *sync.Pool
}

func (o *reqPool) GetInstance() server.IExtensionRequest {
	return o.pool.Get().(server.IExtensionRequest)
}

func (o *reqPool) Recycle(instance server.IExtensionRequest) {
	if nil == instance {
		return
	}
	o.pool.Put(instance)
}

func NewResponsePool(newFunc func() server.IExtensionResponse) server.IResponsePool {
	return &respPool{pool: &sync.Pool{
		New: func() interface{} { return newFunc() },
	}}
}

type respPool struct {
	pool *sync.Pool
}

func (o *respPool) GetInstance() server.IExtensionResponse {
	return o.pool.Get().(server.IExtensionResponse)
}

func (o *respPool) Recycle(instance server.IExtensionResponse) {
	if nil == instance {
		return
	}
	o.pool.Put(instance)
}

func NewNotifyPool(newFunc func() server.IExtensionNotify) server.INotifyPool {
	return &notifyPool{pool: &sync.Pool{
		New: func() interface{} { return newFunc() },
	}}
}

type notifyPool struct {
	pool *sync.Pool
}

func (o *notifyPool) GetInstance() server.IExtensionNotify {
	return o.pool.Get().(server.IExtensionNotify)
}

func (o *notifyPool) Recycle(instance server.IExtensionNotify) {
	if nil == instance {
		return
	}
	o.pool.Put(instance)
}
