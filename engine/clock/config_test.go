// Create on 2023/8/17
// @author xuzhuoxi
package clock

import (
	"fmt"
	"github.com/xuzhuoxi/Rabbit-Server/engine/config"
	"testing"
)

var (
	cfg1 = &config.CfgClock{
		GameLocation: "Asia/Shanghai",
		GameZero:     "2022-08-01 10:00:00",
		DailyZero:    "10h00m00s",
	}
	cfg2 = &config.CfgClock{
		GameLocation: "Etc/UTC",
		GameZero:     "2022-08-01 02:00:00",
		DailyZero:    "10h00m00s",
	}
)

func TestRabbitClockConfig_FromCfg(t *testing.T) {
	config := RabbitClockConfig{}
	err1 := config.FromCfg(cfg1)
	if nil != err1 {
		fmt.Println(err1)
		return
	}
	fmt.Println(config.GameLocation.String())
	fmt.Println(config.GameZeroTime, config.GameZeroTime.UnixNano())
	fmt.Println(config.DailyZeroTime, int64(config.DailyZeroTime))
	err2 := config.FromCfg(cfg2)
	if nil != err2 {
		fmt.Println(err2)
		return
	}
	fmt.Println(config.GameLocation.String())
	fmt.Println(config.GameZeroTime, config.GameZeroTime.UnixNano())
	fmt.Println(config.DailyZeroTime, int64(config.DailyZeroTime))
}
