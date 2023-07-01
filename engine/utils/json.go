// Package utils
// Create on 2023/7/1
// @author xuzhuoxi
package utils

import (
	"encoding/json"
	"github.com/xuzhuoxi/infra-go/filex"
	"io/ioutil"
	"os"
)

func MarshalToJson(in interface{}, jsonPath string) error {
	out, err := json.Marshal(in)
	if nil != err {
		return err
	}
	return filex.WriteFile(jsonPath, out, os.ModePerm)
}

func UnmarshalFromJson(jsonPath string, out interface{}) error {
	bs, err1 := ioutil.ReadFile(jsonPath)
	if nil != err1 {
		return err1
	}
	return json.Unmarshal(bs, out)
}
