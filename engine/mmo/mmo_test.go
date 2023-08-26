// Package mmo
// Create on 2023/8/26
// @author xuzhuoxi
package mmo

import (
	"fmt"
	"github.com/xuzhuoxi/Rabbit-Server/engine/mmo/basis"
	"github.com/xuzhuoxi/Rabbit-Server/engine/mmo/config"
	"github.com/xuzhuoxi/Rabbit-Server/engine/mmo/manager"
	"github.com/xuzhuoxi/infra-go/eventx"
	"github.com/xuzhuoxi/infra-go/filex"
	"github.com/xuzhuoxi/infra-go/osxu"
	"math/rand"
	"testing"
	"time"
)

var (
	cfgPath = filex.Combine(osxu.GetRunningDir(), "conf/mmo.yaml")
	count   = 20
	ch      = make(chan struct{})
)

func TestNewMMOManager(t *testing.T) {
	cfg, err := config.ParseByYamlPath(cfgPath)
	if nil != err {
		t.Fatal(err)
	}
	fmt.Println(cfg)
	eMgr := manager.NewIEntityManager()
	world, err1 := eMgr.ConstructWorldDefault(cfg)
	if nil != err1 {
		t.Fatal(err1)
	}
	eMgr.AddEventListener(basis.EventManagerVarChanged, onVarChanged)
	world.ForEachChild(func(child basis.IEntity) (interruptCurrent bool, interruptRecurse bool) {
		go setVar(child, time.Duration(rand.Int63n(5)))
		return
	})
	time.Sleep(time.Second * 5)
}

func onVarChanged(evd *eventx.EventData) {
	varData := evd.Data.(basis.ManagerVarEventData)
	fmt.Println(varData.Entity.UID(), varData.Entity.EntityType(), varData.Data.Len())
}

func setVar(entity basis.IEntity, interval time.Duration) {
	if v, ok := entity.(basis.IVariableSupport); ok {
		ran := rand.Intn(10)
		if ran > 5 {
			v.SetVar("temp", rand.Intn(1000))
		} else {
			vs := basis.NewVarSet()
			vs.Set("temp1", rand.Intn(300))
			vs.Set("temp2", rand.Intn(300))
			v.SetVars(vs)
		}
	}
	time.Sleep(interval + time.Duration(rand.Int63n(int64(time.Second)*5)))
	setVar(entity, interval)
}
