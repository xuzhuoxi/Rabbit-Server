// Package clock
// Create on 2023/8/17
// @author xuzhuoxi
package clock

import (
	"time"
)

const (
	GameZeroLayout  = "2006-01-02 15:04:05"
	DailyZeroLayout = "15h04m05s"
)

// CfgClock
// Yaml配置对应的文件结构
type CfgClock struct {
	// IANA time zone, https://en.wikipedia.org/wiki/List_of_tz_database_time_zones
	// 例如：Asia/Shanghai
	GameLocation string `yaml:"game_loc"`
	GameZero     string `yaml:"game_zero"`  // 格式为 GameZeroLayout
	DailyZero    string `yaml:"daily_zero"` // 格式为 DailyZeroLayout
}

// RabbitClockConfig
// 时钟配置
type RabbitClockConfig struct {
	GameLocation  *time.Location // 游戏时区
	GameZeroTime  time.Time      // 游戏开始时间
	DailyZeroTime time.Duration  // 游戏零点时间
}

func (o *RabbitClockConfig) FromCfg(cfg *CfgClock) error {
	loc, err1 := time.LoadLocation(cfg.GameLocation)
	if nil != err1 {
		return err1
	}
	gz, err2 := time.ParseInLocation(GameZeroLayout, cfg.GameZero, loc)
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
