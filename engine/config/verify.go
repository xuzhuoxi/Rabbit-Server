// Package config
// Create on 2023/8/23
// @author xuzhuoxi
package config

import (
	"github.com/xuzhuoxi/infra-go/timex"
	"sort"
	"strings"
	"time"
)

// ICfgProtoVerify
// 协议响应时间验证单元
type ICfgProtoVerify interface {
	// PerSecLimitOn
	// 是否启用每秒响应次数限制
	PerSecLimitOn() bool
	// GetMaxPerSec
	// 第秒最大响应次数
	GetMaxPerSec() int
	// FreqLimitOn
	// 是否启用最小响应间隔限制
	FreqLimitOn() bool
	// GetMinFreq
	// 获取最小响应间隔时间
	GetMinFreq() time.Duration
}

// CfgProtoVerify
// 协议响应时间验证单元
type CfgProtoVerify struct {
	Name       string        `yaml:"name,omitempty"` // Extension中的Name
	PId        string        `yaml:"pid,omitempty"`  // Extension中的ProtoId
	MaxPerSec  int           `yaml:"max-per-sec"`    // 每秒最大响应次数
	MinFreq    string        `yaml:"min-freq"`       // 同一客户端最小响应间隔时间(字符串表示)
	MinFreqVal time.Duration // 同一客户端最小响应间隔时间, 在执行HandleData后值更新
}

func (o CfgProtoVerify) PerSecLimitOn() bool {
	return o.MaxPerSec > 0
}

func (o CfgProtoVerify) GetMaxPerSec() int {
	return o.MaxPerSec
}

func (o CfgProtoVerify) FreqLimitOn() bool {
	return o.MinFreqVal > 0
}

func (o CfgProtoVerify) GetMinFreq() time.Duration {
	return o.MinFreqVal
}

// CfgVerifyRoot
// 协议响应时间验证配置, 对应verify.yaml文件
type CfgVerifyRoot struct {
	Default CfgProtoVerify   `yaml:"default"` // 默认配置
	Customs []CfgProtoVerify `yaml:"custom"`  // 自定义配置
}

// FindVerify
// 查找限制配置，如果找不到就返回默认配置
func (o *CfgVerifyRoot) FindVerify(name string, pid string) (v ICfgProtoVerify) {
	for index := len(o.Customs) - 1; index >= 0; index -= 1 {
		if name == o.Customs[index].Name && pid == o.Customs[index].PId {
			return o.Customs[index]
		}
	}
	return o.Default
}

// HandleData
// 处理
func (o *CfgVerifyRoot) HandleData() {
	o.Default.MinFreqVal = timex.ParseDuration(o.Default.MinFreq)
	// 排序
	sort.Slice(o.Customs, o.less)
	for index := range o.Customs {
		o.Customs[index].MinFreqVal = timex.ParseDuration(o.Customs[index].MinFreq)
	}
}

// 排序比较函数
// 判断顺序：
// 1. PId都不为空时
// 1.1 Name相同时，按字符比较PId.
// 1.2 Name不相同时，按字符比较Name.
//
// 2. PId有一个为空时
// 2.1 PId相同时，按字符比较Name.
// 2.2 PId不相同时，PId为空时排在前.
func (o *CfgVerifyRoot) less(i, j int) bool {
	if o.Customs[i].PId == "" || o.Customs[j].PId == "" {
		if o.Customs[i].PId == o.Customs[j].PId {
			return o.compareStr(o.Customs[i].Name, o.Customs[j].Name)
		}
		if o.Customs[i].PId == "" {
			return true
		}
		return false
	} else {
		if o.Customs[i].Name == o.Customs[j].Name {
			o.compareStr(o.Customs[i].PId, o.Customs[j].PId)
		}
		return o.compareStr(o.Customs[i].Name, o.Customs[j].Name)
	}
}

func (o *CfgVerifyRoot) compareStr(a, b string) bool {
	return strings.Compare(a, b) < 0
}
