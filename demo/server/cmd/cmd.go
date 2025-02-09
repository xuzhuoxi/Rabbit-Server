// Package cmd
// Create on 2025/2/9
// @author xuzhuoxi
package cmd

import (
	"github.com/xuzhuoxi/infra-go/cmdx"
)

const (
	// Info 查看信息
	Info = "info"
)

func GenCommandLine() cmdx.ICommandLineListener {
	cmdLine := cmdx.CreateCommandLineListener("Rabbit-Server —— 请输入命令：", 0)
	cmdLine.MapCommand(Info, OnCmdInfo)
	return cmdLine
}
