// Package cmd
// Create on 2025/2/9
// @author xuzhuoxi
package cmd

import (
	"fmt"
	"github.com/xuzhuoxi/infra-go/cmdx"
	"strings"
)

const (
	infoIdKey  = "id"
	infoAllKey = "all"
)

// OnCmdInfo info -id=ID -all=true
func OnCmdInfo(flagSet *cmdx.FlagSetExtend, args []string) {
	id := flagSet.String(infoIdKey, "", "-id=Id")
	all := flagSet.Bool(infoAllKey, false, "-all=true|false")
	flagSet.Parse(args)

	if *all {
		showInfos()
		return
	}
	idValue := strings.TrimSpace(*id)
	if len(idValue) != 0 {
		showInfo(idValue)
		return
	}
	showInfoDefault()
}

func showInfos() {
	fmt.Println("Info -all=true")
}

func showInfo(id string) {
	fmt.Println(fmt.Sprintf("Info -id=%s", id))
}

func showInfoDefault() {
	fmt.Println("Info")
}
