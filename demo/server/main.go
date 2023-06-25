// Create on 2023/6/12
// @author xuzhuoxi
package main

import (
	"github.com/xuzhuoxi/Rabbit-Server/demo/server/cmd"
	_ "github.com/xuzhuoxi/Rabbit-Server/demo/server/extension"
	"github.com/xuzhuoxi/Rabbit-Server/engine/loader"
	_ "github.com/xuzhuoxi/Rabbit-Server/engine/server/rabbit"
)

func main() {
	loader := loader.DefaultLoader
	err := loader.LoadConfig("conf/server.yaml")
	if nil != err {
		panic(err)
	}
	loader.InitServers()
	loader.StartServer()
	cmd.StartCmdListener()
}
