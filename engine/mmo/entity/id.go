// Package entity
// Create on 2023/12/27
// @author xuzhuoxi
package entity

import (
	"github.com/xuzhuoxi/Rabbit-Server/engine/mmo/basis"
	"github.com/xuzhuoxi/infra-go/timex"
	"strconv"
	"sync"
)

type IdGenerator struct {
	Prev  string
	Index int64
	lock  sync.RWMutex
}

func (o *IdGenerator) NextIndex() int64 {
	o.lock.Lock()
	defer o.lock.Unlock()
	o.Index += 1
	return o.Index
}

func (o *IdGenerator) GenId() string {
	o.lock.Lock()
	defer o.lock.Unlock()
	o.Index += 1
	id := o.Prev + "_" + strconv.FormatInt(o.Index, 36)
	return id
}

type IdGeneratorFactory struct {
	generators map[basis.EntityType]*IdGenerator
	mapLock    sync.RWMutex
}

func (o *IdGeneratorFactory) SetGenerator(t basis.EntityType, gen *IdGenerator) {
	o.mapLock.Lock()
	defer o.mapLock.Unlock()
	o.generators[t] = gen
}

func (o *IdGeneratorFactory) NextIndex(t basis.EntityType) int64 {
	o.mapLock.RLock()
	defer o.mapLock.RUnlock()
	return o.generators[t].NextIndex()
}

func (o *IdGeneratorFactory) GenId(t basis.EntityType) string {
	o.mapLock.RLock()
	defer o.mapLock.RUnlock()
	return o.generators[t].GenId()
}

var (
	DefaultIdGeneratorFactory *IdGeneratorFactory
)

func init() {
	DefaultIdGeneratorFactory = &IdGeneratorFactory{generators: make(map[basis.EntityType]*IdGenerator)}
	index := timex.NowMilliseconds1970()
	DefaultIdGeneratorFactory.SetGenerator(basis.EntityUnit, &IdGenerator{Prev: "Unit", Index: index})
	DefaultIdGeneratorFactory.SetGenerator(basis.EntityPlayer, &IdGenerator{Prev: "Player", Index: index})
	DefaultIdGeneratorFactory.SetGenerator(basis.EntityRoom, &IdGenerator{Prev: "Room", Index: index})
	DefaultIdGeneratorFactory.SetGenerator(basis.EntityTeam, &IdGenerator{Prev: "Team", Index: index})
	DefaultIdGeneratorFactory.SetGenerator(basis.EntityTeamCorps, &IdGenerator{Prev: "Corps", Index: index})
	DefaultIdGeneratorFactory.SetGenerator(basis.EntityChannel, &IdGenerator{Prev: "Chan", Index: index})
}

func NextIndex(t basis.EntityType) int64 {
	return DefaultIdGeneratorFactory.NextIndex(t)
}

func GenId(t basis.EntityType) string {
	return DefaultIdGeneratorFactory.GenId(t)
}
