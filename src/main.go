// Create on 2023/6/12
// @author xuzhuoxi
package main

import (
	"github.com/xuzhuoxi/Rabbit-Server/engine/loader"
	"github.com/xuzhuoxi/Rabbit-Server/src/cmd"
)

func main() {
	loader := loader.DefaultLoader
	err := loader.LoadConfig("conf/server.yaml")
	if nil != err {
		panic(err)
	}
	loader.InitServer()
	loader.StartServer()
	cmd.StartCmdListener()
}
