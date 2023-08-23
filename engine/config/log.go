// Package config
// Create on 2023/8/23
// @author xuzhuoxi
package config

import "github.com/xuzhuoxi/infra-go/logx"

type CfgRabbitLogItem struct {
	Name string      `yaml:"name"`
	Conf logx.CfgLog `yaml:"conf"`
}

type CfgRabbitLog struct {
	Default string             `yaml:"default"`
	Logs    []CfgRabbitLogItem `yaml:"logs"`
}
