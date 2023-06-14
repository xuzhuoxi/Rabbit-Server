// Package cmd
// Create on 2023/6/3
// @author xuzhuoxi
package cmd

import (
	"github.com/xuzhuoxi/infra-go/cmdx"
)

func StartCmdListener() {
	cmdLine := cmdx.CreateCommandLineListener("请输入命令：", 0)

	cmdLine.StartListen() //这里会发生阻塞，保证程序不会结束
}
