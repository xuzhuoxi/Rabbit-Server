// Package server
// Create on 2023/6/13
// @author xuzhuoxi
package server

import (
	RabbitConfig "github.com/xuzhuoxi/Rabbit-Server/engine/config"
	"github.com/xuzhuoxi/Rabbit-Server/engine/mmo/config"
	"github.com/xuzhuoxi/Rabbit-Server/engine/utils"
)

const (
	GameZeroLayout  = "2006-01-02 15:04:05"
	DailyZeroLayout = "15h04m05s"
)

type CfgRabbitRoot struct {
	LogPath   string `yaml:"log,omitempty"`
	ClockPath string `yaml:"clock,omitempty"`

	DbPath     string `yaml:"db,omitempty"`
	MMOPath    string `yaml:"mmo,omitempty"`
	ServerPath string `yaml:"server,omitempty"`
	VerifyPath string `yaml:"verify,omitempty"`
}

func (o CfgRabbitRoot) LoadLogConfig() (conf *RabbitConfig.CfgRabbitLog, err error) {
	if o.LogPath == "" {
		return nil, nil
	}
	filePath := utils.FixFilePath(o.LogPath)
	conf = &RabbitConfig.CfgRabbitLog{}
	err = utils.UnmarshalFromYaml(filePath, conf)
	if nil != err {
		return nil, err
	}
	return
}

func (o CfgRabbitRoot) LoadClockConfig() (conf *RabbitConfig.CfgClock, err error) {
	if o.ClockPath == "" {
		return nil, nil
	}
	filePath := utils.FixFilePath(o.ClockPath)
	conf = &RabbitConfig.CfgClock{}
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

func (o CfgRabbitRoot) LoadServerConfig() (conf *RabbitConfig.CfgRabbitServer, err error) {
	if o.ServerPath == "" {
		return nil, nil
	}
	filePath := utils.FixFilePath(o.ServerPath)
	conf = &RabbitConfig.CfgRabbitServer{}
	err = utils.UnmarshalFromYaml(filePath, conf)
	if nil != err {
		return nil, err
	}
	return
}

func (o CfgRabbitRoot) LoadVerifyConfig() (conf *RabbitConfig.CfgVerifyRoot, err error) {
	if o.VerifyPath == "" {
		return nil, nil
	}
	filePath := utils.FixFilePath(o.VerifyPath)
	conf = &RabbitConfig.CfgVerifyRoot{}
	err = utils.UnmarshalFromYaml(filePath, conf)
	if nil != err {
		return nil, err
	}
	conf.HandleData()
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
