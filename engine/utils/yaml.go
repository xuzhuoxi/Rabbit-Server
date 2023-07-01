// Package utils
// Create on 2023/7/1
// @author xuzhuoxi
package utils

import (
	"github.com/xuzhuoxi/infra-go/filex"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
)

func MarshalToYaml(in interface{}, yamlPath string) error {
	out, err := yaml.Marshal(in)
	if nil != err {
		return err
	}
	return filex.WriteFile(yamlPath, out, os.ModePerm)
}

func UnmarshalFromYaml(yamlPath string, out interface{}) error {
	bs, err1 := ioutil.ReadFile(yamlPath)
	if nil != err1 {
		return err1
	}
	return yaml.Unmarshal(bs, out)
}
