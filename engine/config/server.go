// Package config
// Create on 2023/8/23
// @author xuzhuoxi
package config

import (
	"github.com/xuzhuoxi/infra-go/timex"
	"time"
)

type CfgNet struct {
	Name    string `yaml:"name"`
	Network string `yaml:"network"`
	Addr    string `yaml:"addr"`
	Disable bool   `yaml:"disable,omitempty"`
}

type CfgExtension struct {
	All  bool     `yaml:"all"`
	List []string `yaml:"list,omitempty"`
}

type CfgRabbitServerItem struct {
	Id         string       `yaml:"id"`                  // 服务哭喊实例Id
	Name       string       `yaml:"name"`                // 服务器名称
	ToHome     CfgNet       `yaml:"to_home"`             // Home连接信息
	ToHomeRate string       `yaml:"to_home_rate"`        // Home更新频率
	FromUser   CfgNet       `yaml:"from_user"`           // 接收User请求
	FromHome   CfgNet       `yaml:"from_home,omitempty"` // 接收Home命令
	Extension  CfgExtension `yaml:"extension,omitempty"` // Extension配置
	LogRef     string       `yaml:"log_ref,omitempty"`   // 日志记录路径
}

func (o CfgRabbitServerItem) GetToHomeRate() time.Duration {
	return timex.ParseDuration(o.ToHomeRate)
}

type CfgRabbitServer struct {
	Servers []CfgRabbitServerItem `yaml:"servers"`
}
