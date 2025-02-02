// Package index
// Created by xuzhuoxi
// on 2019-03-09.
// @author xuzhuoxi
package index

import (
	"github.com/xuzhuoxi/Rabbit-Server/engine/mmo/basis"
	"github.com/xuzhuoxi/Rabbit-Server/engine/server"
	"sync"
)

func NewIUnitIndex() basis.IUnitIndex {
	return NewUnitIndex()
}

func NewUnitIndex() *UnitIndex {
	return &UnitIndex{EntityIndex: NewEntityIndex("UnitIndex", basis.EntityUnit)}
}

type UnitIndex struct {
	EntityIndex basis.IEntityIndex
	lock        sync.RWMutex
}

func (o *UnitIndex) Size() int {
	return o.EntityIndex.Size()
}

func (o *UnitIndex) EntityType() basis.EntityType {
	return o.EntityIndex.EntityType()
}

func (o *UnitIndex) ForEachEntity(each func(entity basis.IEntity) (interrupt bool)) {
	o.EntityIndex.ForEachEntity(each)
}

func (o *UnitIndex) CheckUnit(unitId string) bool {
	return o.EntityIndex.Check(unitId)
}

func (o *UnitIndex) GetUnit(unitId string) (Unit basis.IUnitEntity, ok bool) {
	Unit, ok = o.EntityIndex.Get(unitId).(basis.IUnitEntity)
	return
}

func (o *UnitIndex) AddUnit(unit basis.IUnitEntity) (rsCode int32, err error) {
	o.lock.Lock()
	defer o.lock.Unlock()
	return o.addUnit(unit)
}

func (o *UnitIndex) AddUnits(units []basis.IUnitEntity, mustAll bool) (rsCode int32, err error) {
	if len(units) == 0 {
		return
	}
	o.lock.Lock()
	defer o.lock.Unlock()
	index := 0
	for index = range units {
		rsCode, err = o.addUnit(units[index])
		if rsCode != server.CodeSuc {
			if mustAll {
				goto undo
			}
			return
		}
	}
	return
undo:
	for index >= 0 {
		_, _, _ = o.EntityIndex.Remove(units[index].UID())
		index--
	}
	return
}

func (o *UnitIndex) RemoveUnit(unitId string) (Unit basis.IUnitEntity, rsCode int32, err error) {
	o.lock.Lock()
	defer o.lock.Unlock()
	return o.removeUnit(unitId)
}

func (o *UnitIndex) RemoveUnits(match func(entity basis.IUnitEntity) bool) (units []basis.IUnitEntity) {
	o.lock.Lock()
	defer o.lock.Unlock()
	var idArr []string
	o.ForEachEntity(func(e basis.IEntity) (interrupt bool) {
		if match(e.(basis.IUnitEntity)) {
			idArr = append(idArr, e.UID())
		}
		return false
	})
	if len(idArr) == 0 {
		return
	}
	for index := range idArr {
		e, _, _ := o.removeUnit(idArr[index])
		if nil != e {
			units = append(units, e.(basis.IUnitEntity))
		}
	}
	return
}

func (o *UnitIndex) UpdateUnit(unit basis.IUnitEntity) (rsCode int32, err error) {
	o.lock.Lock()
	defer o.lock.Unlock()
	return o.updateUnit(unit)
}

func (o *UnitIndex) addUnit(unit basis.IUnitEntity) (rsCode int32, err error) {
	num, err1 := o.EntityIndex.Add(unit)
	if nil == err1 {
		return
	}
	if num == 1 || num == 2 {
		return basis.CodeMMOIndexType, err1
	}
	return basis.CodeMMOUnitExist, err1
}

func (o *UnitIndex) removeUnit(unitId string) (Unit basis.IUnitEntity, rsCode int32, err error) {
	c, _, err1 := o.EntityIndex.Remove(unitId)
	if nil != c {
		return c.(basis.IUnitEntity), 0, nil
	}
	return nil, basis.CodeMMOUnitNotExist, err1
}

func (o *UnitIndex) updateUnit(Unit basis.IUnitEntity) (rsCode int32, err error) {
	_, err1 := o.EntityIndex.Update(Unit)
	if nil != err1 {
		return basis.CodeMMOIndexType, err1
	}
	return
}
