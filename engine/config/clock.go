// Package config
// Create on 2023/8/23
// @author xuzhuoxi
package config

type CfgClock struct {
	// IANA time zone, https://en.wikipedia.org/wiki/List_of_tz_database_time_zones
	// 例如：Asia/Shanghai
	GameLocation string `yaml:"game_loc"`
	GameZero     string `yaml:"game_zero"`  // 格式为 GameZeroLayout
	DailyZero    string `yaml:"daily_zero"` // 格式为 DailyZeroLayout
}
