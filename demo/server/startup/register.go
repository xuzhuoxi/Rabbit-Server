// Package startup
// Create on 2023/7/2
// @author xuzhuoxi
package startup

import (
	"github.com/xuzhuoxi/infra-go/serialx"
)

var (
	StartupManager = serialx.NewStartupManager()
)

func init() {
	StartupManager.AppendModule(&initDefault{})
	StartupManager.AppendModule(&initRabbitManager{})
	StartupManager.AppendModule(&initRabbitLogger{})

	StartupManager.AppendModule(&initExtensions{})
	StartupManager.AppendModule(&initServer{})

	StartupManager.AppendModule(&mgrServer{})
	StartupManager.AppendModule(&mgrMMO{})
}
