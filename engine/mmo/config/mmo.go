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
	Entities     *CfgMMOEntities `json:"entities" yaml:"entities"`
	Relations    *CfgRelations   `json:"relations" yaml:"relations"`
	DefaultWorld string          `json:"default" yaml:"default"`
	Log          *logx.CfgLog    `json:"log,omitempty" yaml:"log,omitempty"`
	LogRef       string          `json:"log_ref,omitempty" yaml:"log_ref,omitempty"`
}

func (o *MMOConfig) CheckEntity(entityId string) bool {
	if len(entityId) == 0 {
		return false
	}
	return o.Entities.ExistWorld(entityId) ||
		o.Entities.ExistZone(entityId) ||
		o.Entities.ExistRoom(entityId)
}

func (o *MMOConfig) HandleData() {
	if err := o.Entities.CheckEntities(); err != nil {
		panic(err)
	}
	if err := o.checkRelations(); nil != err {
		panic(err)
	}
}

func (o *MMOConfig) checkRelations() error {
	if nil == o.Relations || len(o.Relations.Relations) == 0 {
		return nil
	}
	for _, wr := range o.Relations.Relations {
		if _, ok := o.Entities.FindWorld(wr.WorldId); !ok {
			return errors.New(fmt.Sprintf("Relation Entity World[%s] Undefined!", wr.WorldId))
		}
		for _, zr := range wr.Zones {
			if _, ok := o.Entities.FindZone(zr.ZoneId); !ok {
				return errors.New(fmt.Sprintf("Relation Entity Zone[%s] Undefined!", zr.ZoneId))
			}
			for _, rId := range zr.Rooms {
				if _, ok := o.Entities.FindRoom(rId); !ok {
					return errors.New(fmt.Sprintf("Relation Entity Room[%s] Undefined!", rId))
				}
			}
		}
	}
	return nil
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
	cfg.HandleData()
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
	cfg.HandleData()
	return cfg, nil
}
