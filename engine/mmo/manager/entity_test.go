// Created by xuzhuoxi
// on 2019-06-10.
// @author xuzhuoxi
package manager

import (
	"fmt"
	"github.com/xuzhuoxi/Rabbit-Server/engine/mmo/basis"
	"github.com/xuzhuoxi/Rabbit-Server/engine/mmo/config"
	"github.com/xuzhuoxi/infra-go/filex"
	"github.com/xuzhuoxi/infra-go/logx"
	"github.com/xuzhuoxi/infra-go/osxu"
	"testing"
)

var path = filex.Combine(osxu.GetRunningDir(), "conf/mmo.yaml")

func TestEntityManager_ConstructWorld(t *testing.T) {
	cfg, err := config.ParseByYamlPath(path)
	if nil != err {
		t.Fatal(err)
	}
	fmt.Println(cfg)
	eMgr := NewIEntityManager()
	eMgr.ConstructWorldDefault(cfg)
	eMgr.World().ForEachChild(func(child basis.IEntity) (interruptCurrent bool, interruptRecurse bool) {
		logx.Traceln(child.UID(), child.(basis.IEntityChild).GetParent())
		return
	})
}
