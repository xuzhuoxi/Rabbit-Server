// Package server
// Create on 2023/6/13
// @author xuzhuoxi
package server

import (
	"github.com/xuzhuoxi/Rabbit-Server/engine/mmo/config"
	"github.com/xuzhuoxi/Rabbit-Server/engine/utils"
	"github.com/xuzhuoxi/infra-go/logx"
	"github.com/xuzhuoxi/infra-go/timex"
	"time"
)

const (
	GameZeroLayout  = "2006-01-02 15:04:05"
	DailyZeroLayout = "15h04m05s"
)

type CfgClock struct {
	// IANA time zone, https://en.wikipedia.org/wiki/List_of_tz_database_time_zones
	// 例如：Asia/Shanghai
	GameLocation string `yaml:"game_loc"`
	GameZero     string `yaml:"game_zero"`  // 格式为 GameZeroLayout
	DailyZero    string `yaml:"daily_zero"` // 格式为 DailyZeroLayout
}

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

func (o CfgExtension) Extensions() []string {
	if o.All {
		return GetAllExtensions()
	}
	return o.List
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

type CfgRabbitLogItem struct {
	Name string      `yaml:"name"`
	Conf logx.CfgLog `yaml:"conf"`
}

type CfgRabbitLog struct {
	Default string             `yaml:"default"`
	Logs    []CfgRabbitLogItem `yaml:"logs"`
}

type CfgRabbitRoot struct {
	LogPath    string `yaml:"log,omitempty"`
	ServerPath string `yaml:"server,omitempty"`
	MMOPath    string `yaml:"mmo,omitempty"`
	DbPath     string `yaml:"db,omitempty"`
	ClockPath  string `yaml:"clock,omitempty"`
}

func (o CfgRabbitRoot) LoadLogConfig() (conf *CfgRabbitLog, err error) {
	if o.LogPath == "" {
		return nil, nil
	}
	filePath := utils.FixFilePath(o.LogPath)
	conf = &CfgRabbitLog{}
	err = utils.UnmarshalFromYaml(filePath, conf)
	if nil != err {
		return nil, err
	}
	return
}

func (o CfgRabbitRoot) LoadServerConfig() (conf *CfgRabbitServer, err error) {
	if o.ServerPath == "" {
		return nil, nil
	}
	filePath := utils.FixFilePath(o.ServerPath)
	conf = &CfgRabbitServer{}
	err = utils.UnmarshalFromYaml(filePath, conf)
	if nil != err {
		return nil, err
	}
	return
}

func (o CfgRabbitRoot) LoadMMOConfig() (conf *config.MMOConfig, err error) {
	if o.MMOPath == "" {
		return nil, nil
	}
	filePath := utils.FixFilePath(o.MMOPath)
	conf = &config.MMOConfig{}
	err = utils.UnmarshalFromYaml(filePath, conf)
	if nil != err {
		return nil, err
	}
	return
}

func (o CfgRabbitRoot) LoadClockConfig() (conf *CfgClock, err error) {
	if o.ClockPath == "" {
		return nil, nil
	}
	filePath := utils.FixFilePath(o.ClockPath)
	conf = &CfgClock{}
	err = utils.UnmarshalFromYaml(filePath, conf)
	if nil != err {
		return nil, err
	}
	return
}

func LoadRabbitRootConfig(filePath string) (cfg *CfgRabbitRoot, err error) {
	filePath = utils.FixFilePath(filePath)
	cfg = &CfgRabbitRoot{}
	err = utils.UnmarshalFromYaml(filePath, cfg)
	if nil != err {
		return nil, err
	}
	return
}
