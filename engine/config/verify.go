// Package config
// Create on 2023/8/23
// @author xuzhuoxi
package config

import (
	"fmt"
	"github.com/xuzhuoxi/infra-go/timex"
	"sort"
	"strings"
	"time"
)

type CfgVerify struct {
	MaxPerSec  int    `yaml:"max_per_sec"`
	MinFreq    string `yaml:"min_freq"`
	MinFreqVal time.Duration
}

func (o CfgVerify) String() string {
	return fmt.Sprint("{", o.MaxPerSec, o.MinFreq, o.MinFreqVal, "}")
}

type CfgProtoVerify struct {
	Name string `yaml:"name"`
	PId  string `yaml:"pid"`
	CfgVerify
}

type CfgVerifyRoot struct {
	Default CfgVerify        `yaml:"default"`
	Customs []CfgProtoVerify `yaml:"custom"`
}

func (o *CfgVerifyRoot) HandleData() {
	sort.Sort(o)
	o.Default.MinFreqVal = timex.ParseDuration(o.Default.MinFreq)
	for index := range o.Customs {
		o.Customs[index].MinFreqVal = timex.ParseDuration(o.Customs[index].MinFreq)
	}
}

func (o *CfgVerifyRoot) Len() int {
	return len(o.Customs)
}

func (o *CfgVerifyRoot) Less(i, j int) bool {
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

func (o *CfgVerifyRoot) Swap(i, j int) {
	o.Customs[i], o.Customs[j] = o.Customs[j], o.Customs[i]
}

func (o *CfgVerifyRoot) FindVerify(name string, pid string) (v CfgVerify) {
	for index := len(o.Customs) - 1; index >= 0; index -= 1 {
		if name == o.Customs[index].Name && pid == o.Customs[index].PId {
			return o.Customs[index].CfgVerify
		}
	}
	return o.Default
}
