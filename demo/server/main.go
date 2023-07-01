// Create on 2023/6/12
// @author xuzhuoxi
package main

import (
	"fmt"
	"github.com/xuzhuoxi/Rabbit-Server/demo/server/cmd"
	_ "github.com/xuzhuoxi/Rabbit-Server/demo/server/extension"
	"github.com/xuzhuoxi/Rabbit-Server/engine/loader"
	_ "github.com/xuzhuoxi/Rabbit-Server/engine/server/rabbit"
)

func main() {
	fmt.Println()
	fmt.Println("Rabbit-Server:demo Start... ")
	loader := loader.DefaultLoader
	err := loader.LoadConfig("rabbit.yaml")
	if nil != err {
		panic(err)
	}
	loader.InitLoader()
	loader.StartServers()
	cmd.StartCmdListener()
}
