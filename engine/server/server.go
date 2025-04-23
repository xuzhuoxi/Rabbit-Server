package server

import (
	"encoding/base64"
	"errors"
	"fmt"
	"github.com/xuzhuoxi/Rabbit-Server/engine/config"
	"github.com/xuzhuoxi/infra-go/logx"
	"github.com/xuzhuoxi/infra-go/netx"
	"sync"
)

// Server ---

type IRabbitServerInfo interface {
	logx.ILoggerSupport
	// GetId
	// 获取服务器ID
	GetId() string
	// GetPlatformId
	// 获取平台Id
	GetPlatformId() string
	// GetTypeName
	// 获取服务器类型名称
	GetTypeName() string
}

type IRabbitServerController interface {
	// Init
	// 初始化服务器
	Init(cfg config.CfgRabbitServerItem)
	// Start
	// 启动服务器
	Start()
	// Stop
	// 停止服务器
	Stop()
	// Restart
	// 重启服务器
	Restart()
	// Save
	// 保存服务器数据
	Save()
}

type IRabbitServer interface {
	IRabbitServerInfo
	IRabbitServerController
	// GetConnSet
	// 获取服务器连接集合
	GetConnSet() (set netx.IServerConnSet, ok bool)
	// GetExtensionManager
	// 获取服务器扩展管理器
	GetExtensionManager() (mgr IRabbitExtensionManager, ok bool)
}

type FuncNewRabbitExtension = func(extName string) IRabbitExtension

// Register ---

const NameRabbitServer = "Rabbit-Server"

type FuncNewRabbitServer = func() IRabbitServer

type metaExtension struct {
	Name string
	C    FuncNewRabbitExtension
}

var (
	serverMap      = make(map[string]FuncNewRabbitServer)
	extList        = make([]metaExtension, 0, 0)
	lock           sync.RWMutex
	Base64Encoding = base64.RawURLEncoding
)

// Register Server ---

// NewRabbitServer
// 创建服务器
func NewRabbitServer(name string) (server IRabbitServer, err error) {
	lock.RLock()
	defer lock.RUnlock()
	if f, ok := serverMap[name]; ok {
		return f(), nil
	}
	return nil, errors.New(fmt.Sprintf("No name[%s] at RabbitServer list. ", name))
}

// NewRabbitServerDefault
// 创建默认类型服务器
func NewRabbitServerDefault() (server IRabbitServer, err error) {
	return NewRabbitServer(NameRabbitServer)
}

// RegisterRabbitServer
// 注册服务器类型名称与对应创建函数
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

// RegisterRabbitServerDefault
// 注册默认服务器类型名称与对应创建函数
func RegisterRabbitServerDefault(server FuncNewRabbitServer) {
	RegisterRabbitServer(NameRabbitServer, server)
}

// Register Extension ---

// GetAllExtensions
// 获取全部扩展名称
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

// NewRabbitExtension
// 创建一个扩展
func NewRabbitExtension(name string) (extension IRabbitExtension, err error) {
	lock.RLock()
	defer lock.RUnlock()
	for _, meta := range extList {
		if meta.Name == name {
			return meta.C(name), nil
		}
	}
	return nil, errors.New(fmt.Sprintf("No name[%s] at RabbitExtension list.", name))
}

// RegisterRabbitExtension
// 注册扩展创建函数
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
