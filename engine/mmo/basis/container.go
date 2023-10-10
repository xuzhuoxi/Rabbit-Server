// Package basis
// Created by xuzhuoxi
// on 2019-03-14.
// @author xuzhuoxi
package basis

import (
	"github.com/xuzhuoxi/infra-go/encodingx"
)

type IEntityContainer interface {
	// NumChildren 子实体数量
	NumChildren() int
	// Full 判断容量是否已满
	Full() bool

	// Contains 是否包含实体
	Contains(child IEntity) (isContains bool)
	// ContainsById 是否包含实体
	ContainsById(childId string) (isContains bool)
	// GetChildById 取子实体
	GetChildById(childId string) (entity IEntity, ok bool)

	// UpdateChild 更新子实体
	// errNum=1: child=nil
	// errNum=2: full
	UpdateChild(child IEntity) (old IEntity, errNum int, err error)
	// AddChild 添加子实体
	// errNum=1: child=nil
	// errNum=2: exist
	// errNum=3: full
	AddChild(child IEntity) (errNum int, err error)
	// RemoveChild 移除子实体
	// errNum=1: child=nil
	// errNum=2: not exist
	RemoveChild(child IEntity) (errNum int, err error)
	// RemoveChildById 通过Id移除子实体
	RemoveChildById(childId string) (child IEntity, ok bool)

	// UndoUpdate 取消子实体的更新
	UndoUpdate(old IEntity, new IEntity)
	// UndoAdd 取消子实体的添加
	UndoAdd(added IEntity)
	// UndoRemove 取消子实体的移除
	UndoRemove(removed IEntity)

	// ForEachChild 遍历子实体
	ForEachChild(each func(child IEntity) (interruptCurrent bool, interruptRecurse bool))
	// ForEachChildByType 遍历指定类型的子实体
	ForEachChildByType(entityType EntityType, each func(child IEntity), recurse bool)
}

type UnitParams struct {
	Owner string
	Vars  encodingx.IKeyValue
}

type IUnitContainer interface {
	// Units 全部 Unit
	Units() []IUnitEntity
	// CreateUnit 创建 Unit
	CreateUnit(params UnitParams) (unit IUnitEntity, rsCode int32, err error)
	// CreateUnits  创建 Unit
	CreateUnits(params []UnitParams, mustAll bool) (units []IUnitEntity, rsCode int32, err error)
	// DestroyUnit 删除 Unit
	DestroyUnit(unitId string) (unit IUnitEntity, rsCode int32, err error)
	// DestroyUnitsByOwner 删除指定持有都的全部 Unit
	DestroyUnitsByOwner(owner string) (units []IUnitEntity)
	// ForEachUnit 遍历 Unit
	ForEachUnit(each func(child IUnitEntity) (interrupt bool))
}
