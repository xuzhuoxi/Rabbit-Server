// Package server
// Create on 2023/6/13
// @author xuzhuoxi
package server

import (
	"github.com/xuzhuoxi/infra-go/filex"
	"github.com/xuzhuoxi/infra-go/logx"
	"github.com/xuzhuoxi/infra-go/osxu"
	"github.com/xuzhuoxi/infra-go/timex"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"time"
)

type CfgNet struct {
	Name    string `yaml:"name"`
	Network string `yaml:"network"`
	Addr    string `yaml:"addr"`
}

type CfgExtension struct {
	All  bool     `yaml:"all"`
	List []string `yaml:"list,omitempty"`
}

func (o CfgExtension) Extensions() []string {
	if o.All {
		return GetAllExtensions()
	}
	return o.List
}

type CfgRabbitServer struct {
	Id         string       `yaml:"id"`                  // 服务哭喊实例Id
	Name       string       `yaml:"name"`                // 服务器名称
	ToHome     CfgNet       `yaml:"to_home"`             // Home连接信息
	ToHomeRate string       `yaml:"to_home_rate"`        // Home更新频率
	FromUser   CfgNet       `yaml:"from_user"`           // 接收User请求
	FromHome   CfgNet       `yaml:"from_home,omitempty"` // 接收Home命令
	Extension  CfgExtension `yaml:"extension,omitempty"` // Extension配置
	Log        *logx.CfgLog `yaml:"log,omitempty"`       // 日志记录路径
}

func (o CfgRabbitServer) GetToHomeRate() time.Duration {
	return timex.ParseDuration(o.ToHomeRate)
}

type CfgRabbitServerConfig struct {
	Servers []CfgRabbitServer `yaml:"servers"`
	MMO     string            `yaml:"mmo,omitempty"`
	Logger  *logx.CfgLog      `yaml:"log_global,omitempty"`
}

func PauseServerConfig(filePath string) (cfg *CfgRabbitServerConfig, err error) {
	if !filex.IsFile(filePath) {
		filePath = filex.Combine(osxu.GetRunningDir(), filePath)
	}
	bs, err1 := ioutil.ReadFile(filePath)
	if nil != err1 {
		return nil, err
	}
	cfg = &CfgRabbitServerConfig{}
	err = yaml.Unmarshal(bs, cfg)
	if nil != err {
		return nil, err
	}
	return
}
