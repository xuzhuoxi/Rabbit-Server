// Package clock
// Create on 2023/8/16
// @author xuzhuoxi
package clock

import (
	"github.com/xuzhuoxi/infra-go/timex"
	"time"
)

func NewIRabbitClockManager() IRabbitClockManager {
	return NewRabbitClockManager()
}

func NewRabbitClockManager() *RabbitClockManager {
	return &RabbitClockManager{}
}

// IRabbitClockManager
// 时钟管理器
type IRabbitClockManager interface {
	// GetConfig 获取时钟配置
	GetConfig() RabbitClockConfig
	// Init 初始化时钟管理器
	Init(cfg *CfgClock) error
	// SetOperaDuration 设置运行时长
	SetOperaDuration(d time.Duration)

	// NowGameClock 获取游戏当前时钟
	NowGameClock() (now time.Time, pass time.Duration)
	// NowGameDuration 获取游戏运行时长
	NowGameDuration() time.Duration
	// NowRuntimeDuration 获取服务器运行时长
	NowRuntimeDuration() time.Duration
}

type RabbitClockManager struct {
	Config            RabbitClockConfig
	RuntimeStart      time.Time     // 服务器启动时间
	RuntimeStartStamp int64         // 服务器启动时间戳
	OperationDuration time.Duration // 服务器运行时长
}

func (o *RabbitClockManager) Init(cfg *CfgClock) error {
	if nil == cfg {
		o.initConfigDefault()
	} else {
		err := o.initConfig(cfg)
		if nil != err {
			return err
		}
	}
	o.RuntimeStart = time.Now().In(o.Config.GameLocation)
	o.RuntimeStartStamp = o.RuntimeStart.UnixNano()
	return nil
}

func (o *RabbitClockManager) SetOperaDuration(d time.Duration) {
	o.OperationDuration = d
}

func (o *RabbitClockManager) GetConfig() RabbitClockConfig {
	return o.Config
}

func (o *RabbitClockManager) NowGameClock() (now time.Time, pass time.Duration) {
	pass = o.NowGameDuration()
	now = o.Config.GameZeroTime.Add(pass)
	return
}

func (o *RabbitClockManager) NowGameDuration() time.Duration {
	return o.NowRuntimeDuration() + o.OperationDuration
}

func (o *RabbitClockManager) NowRuntimeDuration() time.Duration {
	return time.Now().Sub(o.RuntimeStart)
}

func (o *RabbitClockManager) initConfigDefault() {
	o.Config = RabbitClockConfig{
		GameLocation:  time.UTC,
		GameZeroTime:  timex.Zero1970UTC,
		DailyZeroTime: time.Duration(0),
	}
}

func (o *RabbitClockManager) initConfig(cfg *CfgClock) error {
	config := &RabbitClockConfig{}
	err := config.FromCfg(cfg)
	if nil != err {
		return err
	}
	o.Config = *config
	return nil
}
