package server

import (
	"github.com/xuzhuoxi/infra-go/logx"
)

//
//var (
//	cMap     = make(map[string]FuncNewGameServer)
//	cDefault FuncNewGameServer
//)
//
//type FuncNewGameServer = func() IRabbitServer

type IRabbitServer interface {
	IRabbitServerInfo
	IRabbitServerController
}

type IRabbitServerController interface {
	Init()
	Start()
	Stop()
	Restart()
	Save()
}

type IRabbitServerInfo interface {
	GetId() string
	GetName() string
	GetLogger() logx.ILogger
}

//func NewRabbitServerDefault() IRabbitServer {
//	return cDefault()
//}
//
//func NewRabbitServer(name string) IRabbitServer {
//	if f, ok := cMap[name]; ok {
//		return f()
//	}
//	panic(errors.New(fmt.Sprintf("Duplicate name[%s] at RegisterRabbitServer", name)))
//}
//
//func RegisterRabbitServer(name string, server FuncNewGameServer) {
//	if _, ok := cMap[name]; ok {
//		panic(errors.New(fmt.Sprintf("Duplicate name[%s] at RegisterRabbitServer", name)))
//	}
//	cMap[name] = server
//}
//func RegisterRabbitServerDefault(server FuncNewGameServer) {
//	if cDefault != nil {
//		panic(errors.New("Duplicate Register Default! "))
//	}
//	cDefault = server
//}
