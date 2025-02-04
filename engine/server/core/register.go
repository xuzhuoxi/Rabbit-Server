// Package core
// Create on 2023/6/14
// @author xuzhuoxi
package core

import (
	"github.com/xuzhuoxi/Rabbit-Server/engine/server"
	"github.com/xuzhuoxi/infra-go/netx"
	"time"
)

var (
	// RabbitAddressProxy uid与address的交叉映射
	RabbitAddressProxy = netx.NewIAddressProxyWithName("Rabbit")
)

const (
	// DefaultStatsInterval 统计时间区间
	DefaultStatsInterval = int64(5 * time.Minute)
)

func init() {
	server.RegisterRabbitServerDefault(NewIRabbitServer)
}
