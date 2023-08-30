// Package basis
// Created by xuzhuoxi
// on 2019-03-14.
// @author xuzhuoxi
package basis

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
	UpdateChild(child IEntity) (old IEntity, err error)
	// AddChild 添加子实体
	AddChild(child IEntity) error
	// RemoveChild 移除子实体
	RemoveChild(child IEntity) error
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
