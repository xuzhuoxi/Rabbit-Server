// Package config
// Create on 2023/8/23
// @author xuzhuoxi
package config

import "github.com/xuzhuoxi/infra-go/logx"

type CfgRabbitLogItem struct {
	Name string      `yaml:"name"` // 日志配置名称
	Conf logx.CfgLog `yaml:"conf"` // 日志配置内容
}

// CfgRabbitLog
// 日志配置文件log.yaml对应结构
type CfgRabbitLog struct {
	Default string             `yaml:"default"` // 默认日志配置
	Logs    []CfgRabbitLogItem `yaml:"logs"`    // 其它日志配置
}
