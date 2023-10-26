// Package mmo
// Create on 2023/8/26
// @author xuzhuoxi
package mmo

import (
	"fmt"
	"github.com/xuzhuoxi/Rabbit-Server/engine/mmo/basis"
	"github.com/xuzhuoxi/Rabbit-Server/engine/mmo/config"
	"github.com/xuzhuoxi/Rabbit-Server/engine/mmo/events"
	"github.com/xuzhuoxi/Rabbit-Server/engine/mmo/manager"
	"github.com/xuzhuoxi/Rabbit-Server/engine/mmo/vars"
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
	err1 := eMgr.BuildEnv(cfg)
	if nil != err1 {
		t.Fatal(err1)
	}
	eMgr.AddEventListener(events.EventEntityVarMod, onVarChanged)
	eMgr.AddEventListener(events.EventEntityVarsMod, onVarsChanged)
	eMgr.ForEachRoom(func(room basis.IRoomEntity) (interrupt bool) {
		go setVar(room, time.Duration(rand.Int63n(5)))
		return false
	})
	time.Sleep(time.Second * 5)
}

func onVarChanged(evd *eventx.EventData) {
	varData := evd.Data.(events.VarModEventData)
	fmt.Println(varData.Entity.UID(), varData.Entity.EntityType(), varData.Key, varData.Value)
}

func onVarsChanged(evd *eventx.EventData) {
	varData := evd.Data.(events.VarsModEventData)
	fmt.Println(varData.Entity.UID(), varData.Entity.EntityType(), varData.VarKeys)
}

func setVar(entity basis.IEntity, interval time.Duration) {
	if v, ok := entity.(basis.IVariableSupport); ok {
		ran := rand.Intn(10)
		if ran > 5 {
			v.SetVar("temp", rand.Intn(1000), true)
		} else {
			vs := vars.DefaultVarSetPool.GetInstance()
			vs.Set("temp1", rand.Intn(300))
			vs.Set("temp2", rand.Intn(300))
			v.SetVars(vs, true)
			vars.DefaultVarSetPool.Recycle(vs)
		}
	}
	time.Sleep(interval + time.Duration(rand.Int63n(int64(time.Second)*5)))
	setVar(entity, interval)
}
