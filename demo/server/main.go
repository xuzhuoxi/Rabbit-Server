// Create on 2023/6/12
// @author xuzhuoxi
package main

import (
	"fmt"
	"github.com/xuzhuoxi/Rabbit-Server/demo/server/cmd"
	_ "github.com/xuzhuoxi/Rabbit-Server/demo/server/extension"
	"github.com/xuzhuoxi/Rabbit-Server/engine/mgr"
)

func main() {
	fmt.Println()
	fmt.Println("Rabbit-Server:demo Start... ")
	mgr := mgr.DefaultManager
	err := mgr.GetInitManager().LoadRabbitConfig("rabbit.yaml")
	if nil != err {
		panic(err)
	}
	mgr.GetInitManager().InitLoggerManager()
	mgr.GetServerManager().StartServers()
	cmd.StartCmdListener()
}
