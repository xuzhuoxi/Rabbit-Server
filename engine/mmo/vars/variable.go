// Package vars
// Created by xuzhuoxi
// on 2019-03-03.
// @author xuzhuoxi
package vars

import (
	"bytes"
	"fmt"
	"github.com/xuzhuoxi/Rabbit-Server/engine/mmo/basis"
	"github.com/xuzhuoxi/Rabbit-Server/engine/mmo/events"
	"github.com/xuzhuoxi/infra-go/binaryx"
	"github.com/xuzhuoxi/infra-go/encodingx"
	"github.com/xuzhuoxi/infra-go/eventx"
	"strings"
	"sync"
)

var (
	DefaultVarSetPool encodingx.IPoolKeyValue
)

func init() {
	DefaultVarSetPool = encodingx.NewPoolKeyValue()
	DefaultVarSetPool.Register(NewVarSet)
}

// ---------------------------------------------

func NewIVariableSupport(currentTarget basis.IEntity) basis.IVariableSupport {
	return NewVariableSupport(currentTarget)
}

func NewVariableSupport(currentTarget basis.IEntity) *VariableSupport {
	return &VariableSupport{currentTarget: currentTarget, vars: NewVarSet()}
}

func NewVarSet() encodingx.IKeyValue {
	return encodingx.NewCodingMap()
}

//---------------------------------------------

type VariableSupport struct {
	currentTarget basis.IEntity
	eventx.EventDispatcher
	vars encodingx.IKeyValue
	lock sync.RWMutex
}

func (o *VariableSupport) GetVar(key string) (interface{}, bool) {
	o.lock.RLock()
	defer o.lock.RUnlock()
	return o.vars.Get(key)
}

func (o *VariableSupport) CheckVar(key string) bool {
	o.lock.RLock()
	defer o.lock.RUnlock()
	return o.vars.Check(key)
}

func (o *VariableSupport) Vars() encodingx.IKeyValue {
	o.lock.RLock()
	defer o.lock.RUnlock()
	return o.vars
}

func (o *VariableSupport) SetVar(kv string, value interface{}) {
	if len(kv) == 0 {
		return
	}
	o.lock.Lock()
	defer o.lock.Unlock()
	var ok bool
	if value == nil {
		_, ok = o.vars.Delete(kv)
	} else {
		_, ok = o.vars.Set(kv, value)
	}
	if ok {
		o.DispatchEvent(events.EventEntityVarChanged, o.currentTarget,
			&events.VarEventData{Entity: o.currentTarget, Key: kv, Value: value})
	}
}

func (o *VariableSupport) SetVars(kv encodingx.IKeyValue) {
	if nil == kv {
		return
	}
	o.lock.Lock()
	defer o.lock.Unlock()
	diff, _ := o.vars.Merge(kv)
	if nil != diff {
		o.DispatchEvent(events.EventEntityVarsChanged, o.currentTarget,
			&events.VarsEventData{Entity: o.currentTarget, Vars: diff})
	}
}

func (o *VariableSupport) SetArrayVars(keys []string, vals []interface{}) {
	if len(keys) == 0 || len(keys) != len(vals) {
		return
	}
	o.lock.Lock()
	defer o.lock.Unlock()
	diff, _ := o.vars.MergeArray(keys, vals)
	if nil != diff {
		o.DispatchEvent(events.EventEntityVarsChanged, o.currentTarget,
			&events.VarsEventData{Entity: o.currentTarget, Vars: diff})
	}
}

type VarSet struct {
	intSet    map[string]int32
	stringSet map[string]string
}

func (o *VarSet) String() string {
	if o.Len() == 0 {
		return "{}"
	}
	builder := &strings.Builder{}
	builder.WriteString("{")
	index := 0
	ln := o.Len()
	for key, val := range o.intSet {
		builder.WriteString(key + ":" + fmt.Sprint(val))
		index++
		if index < ln {
			builder.WriteString(",")
		}
	}
	for key, val := range o.stringSet {
		builder.WriteString(key + ":" + val)
		index++
		if index < ln {
			builder.WriteString(",")
		}
	}
	builder.WriteString("}")
	return builder.String()
}

func (o *VarSet) EncodeToBytes() []byte {
	order := encodingx.DefaultOrder
	buff := bytes.NewBuffer(nil)
	err := binaryx.WriteLen(buff, order, o.Len())
	if nil != err {
		return nil
	}
	for key, intVal := range o.intSet {
		_ = binaryx.WriteString(buff, order, key)         //Key
		_ = binaryx.Write(buff, order, binaryx.KindInt32) //Kind
		_ = binaryx.Write(buff, order, intVal)            //Kind
	}
	for key, strVal := range o.stringSet {
		_ = binaryx.WriteString(buff, order, key)          //Key
		_ = binaryx.Write(buff, order, binaryx.KindString) //Kind
		_ = binaryx.Write(buff, order, strVal)             //Kind
	}
	return buff.Bytes()
}

func (o *VarSet) DecodeFromBytes(bs []byte) bool {
	order := encodingx.DefaultOrder
	buff := bytes.NewBuffer(bs)
	ln, err := binaryx.ReadLen(buff, order)
	if nil != err {
		return false
	}
	for ln >= 0 && buff.Len() > 0 {
		key, _ := binaryx.ReadString(buff, order)
		var kind binaryx.ValueKind
		_ = binaryx.Read(buff, order, &kind) //Kind
		if kind == binaryx.KindInt32 {
			var val int32
			err = binaryx.Read(buff, order, &val)
			if nil != err {
				return false
			}
			o.intSet[key] = val
			continue
		}
		if kind == binaryx.KindString {
			var val string
			err = binaryx.Read(buff, order, &val)
			if nil != err {
				return false
			}
			o.stringSet[key] = val
			continue
		}
		ln--
	}
	return true
}

func (o *VarSet) Len() int {
	return len(o.intSet) + len(o.stringSet)
}

func (o *VarSet) CheckVar(key string) bool {
	if len(key) == 0 {
		return false
	}
	if _, ok1 := o.intSet[key]; ok1 {
		return true
	}
	if _, ok2 := o.stringSet[key]; ok2 {
		return true
	}
	return false
}

func (o *VarSet) GetStringVar(key string) (value string, ok bool) {
	if len(key) == 0 {
		return
	}
	value, ok = o.stringSet[key]
	return
}

func (o *VarSet) GetIntVar(key string) (value int32, ok bool) {
	if len(key) == 0 {
		return
	}
	value, ok = o.intSet[key]
	return
}

func (o *VarSet) SetStringVar(key string, value string) (old string, ok bool) {
	if len(key) == 0 {
		return
	}
	return o.setStringVar(key, value)
}

func (o *VarSet) setStringVar(key string, value string) (old string, ok bool) {
	old, ok = o.stringSet[key]
	o.stringSet[key] = value
	return
}

func (o *VarSet) SetIntVar(key string, value int32) (old int32, ok bool) {
	if len(key) == 0 {
		return
	}
	return o.setIntVar(key, value)
}

func (o *VarSet) setIntVar(key string, value int32) (old int32, ok bool) {
	old, ok = o.intSet[key]
	o.intSet[key] = value
	return
}

func (o *VarSet) Delete(key string) (value interface{}, ok bool) {
	if len(key) == 0 {
		return
	}
	if old1, ok1 := o.intSet[key]; ok1 {
		delete(o.intSet, key)
		return old1, true
	}
	if old2, ok2 := o.stringSet[key]; ok2 {
		delete(o.stringSet, key)
		return old2, true
	}
	return
}

func (o *VarSet) Merge(vs basis.IVarSet) (update basis.IVarSet) {
	if vs == nil || vs.Len() == 0 {
		return
	}
	var rm []string
	vs.ForEach(func(key string, intValue int32) {
		_, ok := o.setIntVar(key, intValue)
		if !ok {
			rm = append(rm, key)
		}
	}, func(key string, strValue string) {
		_, ok := o.setStringVar(key, strValue)
		if !ok {
			rm = append(rm, key)
		}
	})
	if len(rm) > 0 { //有重复
		for _, key := range rm {
			vs.Delete(key)
		}
	}
	return vs
}

func (o *VarSet) ForEach(intHandler func(key string, value int32), stringHandler func(key string, value string)) {
	for key, value := range o.intSet {
		intHandler(key, value)
	}
	for key, value := range o.stringSet {
		stringHandler(key, value)
	}
}
