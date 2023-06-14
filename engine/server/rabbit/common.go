// Package rabbit
// Create on 2023/6/14
// @author xuzhuoxi
package rabbit

import (
	"github.com/xuzhuoxi/infra-go/bytex"
	"github.com/xuzhuoxi/infra-go/mathx"
	"github.com/xuzhuoxi/infra-go/netx"
	"time"
)

var (
	AddressProxy     = netx.NewIAddressProxy()            //uid与address的交叉映射,整个game模块共享
	DataBlockHandler = bytex.NewDefaultDataBlockHandler() //数据封包处理
)

const (
	// GameNotifyRouteInterval 通知Route间隔(默认60秒)
	GameNotifyRouteInterval = time.Second * 60
	// DefaultStatsInterval 统计时间区间
	DefaultStatsInterval = int64(5 * time.Minute)
	// LogMaxSize 日志文件最大体量
	LogMaxSize = 10 * mathx.MB
)
