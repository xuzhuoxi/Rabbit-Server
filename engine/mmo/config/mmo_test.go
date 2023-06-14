package config

import (
	"fmt"
	"github.com/xuzhuoxi/infra-go/osxu"
	"testing"
)

var path = osxu.GetRunningDir() + "/conf/mmo.yaml"

func TestParseMMOConfig(t *testing.T) {
	cfg, err := ParseByYamlPath(path)
	if nil != err {
		t.Fatal(err)
	}
	fmt.Println(cfg)
}
