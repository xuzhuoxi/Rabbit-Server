// Package rabbit
// Create on 2023/6/14
// @author xuzhuoxi
package rabbit

import (
	"github.com/xuzhuoxi/infra-go/extendx/protox"
	"sync"
)

type IServerExtensionContainer = protox.IProtocolExtensionContainer

func NewServerExtensionContainer() IServerExtensionContainer {
	return protox.NewIProtocolExtensionContainer()
}

type IServerExtension interface {
	protox.IProtocolExtension
}

type FuncServerExtension func() IServerExtension

var (
	extConstructors []FuncServerExtension
	lock            sync.RWMutex
)

func RegisterExtension(constructor FuncServerExtension) {
	lock.Lock()
	defer lock.Unlock()
	extConstructors = append(extConstructors, constructor)
}

func ForeachExtensionConstructor(eachFunc func(constructor FuncServerExtension)) {
	lock.RLock()
	defer lock.RUnlock()
	for _, c := range extConstructors {
		eachFunc(c)
	}
}
