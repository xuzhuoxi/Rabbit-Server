// Package runtime
// Create on 2023/6/14
// @author xuzhuoxi
package runtime

import (
	"github.com/xuzhuoxi/Rabbit-Server/engine/server"
	"github.com/xuzhuoxi/infra-go/netx"
	"time"
)

var (
	// RabbitUserConnMapper uid与连接id的交叉映射
	RabbitUserConnMapper = netx.NewIUserConnMapperWithName("Rabbit")
)

const (
	// DefaultStatsInterval 统计时间区间
	DefaultStatsInterval = int64(5 * time.Minute)
)

func init() {
	server.RegisterRabbitServerDefault(NewIRabbitServer)
}
