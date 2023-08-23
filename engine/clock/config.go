// Package clock
// Create on 2023/8/17
// @author xuzhuoxi
package clock

import (
	"github.com/xuzhuoxi/Rabbit-Server/engine"
	"github.com/xuzhuoxi/Rabbit-Server/engine/config"
	"time"
)

type RabbitClockConfig struct {
	GameLocation  *time.Location // 游戏时区
	GameZeroTime  time.Time      // 游戏开始时间
	DailyZeroTime time.Duration  // 游戏零点时间
}

func (o *RabbitClockConfig) FromCfg(cfg *config.CfgClock) error {
	loc, err1 := time.LoadLocation(cfg.GameLocation)
	if nil != err1 {
		return err1
	}
	gz, err2 := time.ParseInLocation(server.GameZeroLayout, cfg.GameZero, loc)
	if nil != err2 {
		return err2
	}
	daily, err3 := time.ParseDuration(cfg.DailyZero)
	if nil != err3 {
		return err3
	}

	o.GameLocation, o.GameZeroTime, o.DailyZeroTime = loc, gz, daily
	return nil
}
