// Package startup
// Create on 2023/7/2
// @author xuzhuoxi
package startup

import (
	"encoding/binary"
	jsoniter "github.com/json-iterator/go"
	"github.com/xuzhuoxi/Rabbit-Server/engine/server/packet"
	"github.com/xuzhuoxi/infra-go/bytex"
	"github.com/xuzhuoxi/infra-go/encodingx"
	"github.com/xuzhuoxi/infra-go/encodingx/jsonx"
	"github.com/xuzhuoxi/infra-go/eventx"
	"github.com/xuzhuoxi/infra-go/netx"
	"github.com/xuzhuoxi/infra-go/netx/tcpx"
	"github.com/xuzhuoxi/infra-go/serialx"
)

var (
	order            binary.ByteOrder = binary.LittleEndian
	dataBlockHandler                  = bytex.NewDataBlockHandler(order, bytex.DefaultDataToBlockHandler, bytex.DefaultBlockToDataHandler)
)

type initDefault struct {
	eventx.EventDispatcher
}

func (o *initDefault) Name() string {
	return "Set Default Value."
}

func (o *initDefault) StartModule() {
	bytex.DefaultOrder, bytex.DefaultDataBlockHandler = order, dataBlockHandler         // 包 bytex 下的大小端设置，封包处理
	encodingx.DefaultOrder, encodingx.DefaultDataBlockHandler = order, dataBlockHandler // 包 encodingx下的大小端设置，封包处理
	tcpx.TcpDataBlockHandler = dataBlockHandler                                         // Tcp封包处理
	jsonx.DefaultDataBlockHandler = dataBlockHandler                                    // Json封包处理
	packet.SetJsonMarshalHandler(jsoniter.Marshal)                                      // 设置Json序列化处理器
	netx.ReceiverBuffLen = 4096                                                         // 设置数据包缓存的长度
	o.DispatchEvent(serialx.EventOnStartupModuleStarted, o, nil)
}

func (o *initDefault) StopModule() {
	o.DispatchEvent(serialx.EventOnStartupModuleStopped, o, nil)
}

func (o *initDefault) SaveModule() {
	o.DispatchEvent(serialx.EventOnStartupModuleSaved, o, nil)
}
