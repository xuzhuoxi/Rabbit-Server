package server

import (
	"errors"
	"fmt"
	"github.com/xuzhuoxi/Rabbit-Server/engine/config"
	"github.com/xuzhuoxi/infra-go/logx"
	"github.com/xuzhuoxi/infra-go/netx"
	"sync"
)

// Server ---

type IRabbitServer interface {
	IRabbitServerInfo
	IRabbitServerController
	GetConnSet() (set netx.IServerConnSet, ok bool)
	GetExtensionManager() (mgr IRabbitExtensionManager, ok bool)
}

type IRabbitServerController interface {
	Init(cfg config.CfgRabbitServerItem)
	Start()
	Stop()
	Restart()
	Save()
}

type IRabbitServerInfo interface {
	logx.ILoggerSupport
	GetId() string
	GetName() string
}

type FuncNewRabbitExtension func(name string) IProtoExtension
type IRabbitExtensionContainer = IProtoExtensionContainer
type IRabbitExtensionManager interface {
	IExtensionManager
	SetCustomVerify(reqVerify IReqVerify)
}

// Register ---

const NameRabbitServer = "RabbitServer"

type FuncNewRabbitServer = func() IRabbitServer

type metaExtension struct {
	Name string
	C    FuncNewRabbitExtension
}

var (
	serverMap = make(map[string]FuncNewRabbitServer)
	extList   = make([]metaExtension, 0, 0)
	lock      sync.RWMutex
)

// Register Server ---

func NewRabbitServer(name string) (server IRabbitServer, err error) {
	lock.RLock()
	defer lock.RUnlock()
	if f, ok := serverMap[name]; ok {
		return f(), nil
	}
	return nil, errors.New(fmt.Sprintf("No name[%s] at RabbitServer list. ", name))
}

func NewRabbitServerDefault() (server IRabbitServer, err error) {
	return NewRabbitServer(NameRabbitServer)
}

func RegisterRabbitServer(name string, server FuncNewRabbitServer) {
	if len(name) == 0 {
		panic(errors.New(fmt.Sprintf("RegisterRabbitServer Fail: name[%s]", name)))
	}
	if nil == server {
		panic(errors.New(fmt.Sprintf("RegisterRabbitServer Fail: server is ni ")))
	}
	lock.Lock()
	defer lock.Unlock()
	serverMap[name] = server
}

func RegisterRabbitServerDefault(server FuncNewRabbitServer) {
	RegisterRabbitServer(NameRabbitServer, server)
}

// Register Extension ---

func GetAllExtensions() []string {
	if len(extList) == 0 {
		return nil
	}
	rs := make([]string, len(extList), len(extList))
	for index := range extList {
		rs[index] = extList[index].Name
	}
	return rs
}

func NewRabbitExtension(name string) (extension IProtoExtension, err error) {
	lock.RLock()
	defer lock.RUnlock()
	for _, meta := range extList {
		if meta.Name == name {
			return meta.C(name), nil
		}
	}
	return nil, errors.New(fmt.Sprintf("No name[%s] at RabbitExtension list.", name))
}

func RegisterRabbitExtension(name string, extension FuncNewRabbitExtension) {
	if len(name) == 0 {
		panic(errors.New(fmt.Sprintf("RegisterRabbitServer Fail: name[%s]", name)))
	}
	if nil == extension {
		panic(errors.New(fmt.Sprintf("RegisterRabbitExtension Fail: extension is ni ")))
	}
	lock.Lock()
	defer lock.Unlock()
	extList = append(extList, metaExtension{Name: name, C: extension})
}
