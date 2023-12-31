// Package config
// Created by xuzhuoxi
// on 2019-06-09.
// @author xuzhuoxi
//
package config

import (
	"errors"
	"fmt"
	"github.com/json-iterator/go"
	"github.com/xuzhuoxi/infra-go/cmdx"
	"github.com/xuzhuoxi/infra-go/filex"
	"github.com/xuzhuoxi/infra-go/logx"
	"github.com/xuzhuoxi/infra-go/osxu"
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

type MMOConfig struct {
	Entities  *CfgMMOEntities `json:"entities" yaml:"entities"`
	Relations *CfgRelations   `json:"relations" yaml:"relations"`
	Defaults  []string        `json:"default" yaml:"default"`
	Log       *logx.CfgLog    `json:"log,omitempty" yaml:"log,omitempty"`
	LogRef    string          `json:"log_ref,omitempty" yaml:"log_ref,omitempty"`
}

func (o *MMOConfig) CheckConfig() error {
	err1 := o.Entities.CheckDuplicate()
	if nil != err1 {
		return err1
	}
	return nil
}

func (o *MMOConfig) MergeRelationToTags() {
	for index, room := range o.Entities.Rooms {
		zones := o.Relations.FindMyZones(room.Id)
		if len(zones) > 0 {
			o.Entities.Rooms[index].AppendTags(zones)
		}
		worlds := o.Relations.FindMyWorlds(room.Id)
		if len(worlds) > 0 {
			o.Entities.Rooms[index].AppendTags(worlds)
		}
	}
}

//------------------------------------------

var DefaultMMOConfig *MMOConfig

func ParseMMOConfigByFlag(flagSet *cmdx.FlagSetExtend) (cfg *MMOConfig, err error) {
	if !flagSet.CheckKey("mmo") {
		return nil, errors.New("flag mmo is not exist! ")
	}
	cfgName, _ := flagSet.GetString("mmo")
	path := filex.Combine(osxu.GetRunningDir(), cfgName)
	_, _, ext := filex.SplitFileName(cfgName)
	if ext == "json" {
		return ParseByJsonPath(path)
	} else {
		return ParseByYamlPath(path)
	}
}

func ParseByJsonPath(path string) (cfg *MMOConfig, err error) {
	if !filex.IsFile(path) {
		path = filex.Combine(osxu.GetRunningDir(), path)
	}
	cfgBody, err := ioutil.ReadFile(path)
	if nil != err {
		return nil, errors.New(fmt.Sprintf("mmo does not exist: %s", err))
	}
	return ParseByJsonContent(cfgBody)
}

func ParseByJsonContent(content []byte) (cfg *MMOConfig, err error) {
	cfg = &MMOConfig{}
	err = jsoniter.Unmarshal(content, cfg)
	if nil != err {
		return nil, err
	}
	cfg.MergeRelationToTags()
	return cfg, nil
}

func ParseByYamlPath(path string) (cfg *MMOConfig, err error) {
	if !filex.IsFile(path) {
		path = filex.Combine(osxu.GetRunningDir(), path)
	}
	cfgBody, err := ioutil.ReadFile(path)
	if nil != err {
		return nil, errors.New(fmt.Sprintf("mmo does not exist: %s", err))
	}
	return ParseByYamlContent(cfgBody)
}

func ParseByYamlContent(content []byte) (cfg *MMOConfig, err error) {
	cfg = &MMOConfig{}
	err = yaml.Unmarshal(content, cfg)
	if nil != err {
		return nil, err
	}
	cfg.MergeRelationToTags()
	return cfg, nil
}
