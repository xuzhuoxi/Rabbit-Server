// Create on 2023/6/12
// @author xuzhuoxi
package main

import (
	"github.com/xuzhuoxi/Rabbit-Server/demo/server/cmd"
	_ "github.com/xuzhuoxi/Rabbit-Server/demo/server/extension"
	"github.com/xuzhuoxi/Rabbit-Server/demo/server/startup"
	"time"
)

func main() {
	startup.StartupManager.StartManager()
	startCmd()
}

func startCmd() {
	cmdLine := cmd.GenCommandLine()
	time.Sleep(2 * time.Second)
	cmdLine.StartListen() //这里会发生阻塞，保证程序不会结束
}
