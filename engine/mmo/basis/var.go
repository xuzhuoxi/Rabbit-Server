// Package basis
// Create on 2023/9/7
// @author xuzhuoxi
package basis

import (
	"github.com/xuzhuoxi/infra-go/encodingx"
	"github.com/xuzhuoxi/infra-go/eventx"
)

// IVariableSupport 变量列表
type IVariableSupport interface {
	eventx.IEventDispatcher
	// GetVar 取变量数据
	GetVar(key string) (interface{}, bool)
	// CheckVar 检查变量是否存在
	CheckVar(key string) bool
	// Vars 取变量数据集合
	Vars() encodingx.IKeyValue

	// SetVar 设置变量
	SetVar(kv string, value interface{}, notify bool) (ok bool)
	// SetVars 批量设置变量
	SetVars(kv encodingx.IKeyValue, notify bool) (diff []string)
	// SetArrayVars 批量设置变量
	SetArrayVars(keys []string, vals []interface{}, notify bool) (diff []string)
}

//type IVarSet interface {
//	encodingx.ICodingData
//	// Len 键值对数量
//	Len() int
//	// CheckVar 检查变量是否存在
//	CheckVar(key string) bool
//	// GetStringVar 取值
//	GetStringVar(key string) (value string, ok bool)
//	// GetIntVar 取值
//	GetIntVar(key string) (value int32, ok bool)
//	// SetStringVar 设置字符变量
//	SetStringVar(key string, value string) (old string, ok bool)
//	// SetIntVar 设置整数变量
//	SetIntVar(key string, value int32) (old int32, ok bool)
//	// Delete 删除键值
//	Delete(key string) (value interface{}, ok bool)
//	// Merge 合并
//	Merge(vs IVarSet) (update IVarSet)
//	// ForEach 遍历
//	ForEach(intHandler func(key string, value int32), stringHandler func(key string, value string))
//}
