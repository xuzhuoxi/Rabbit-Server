// Package config
// Create on 2023/8/23
// @author xuzhuoxi
package config

import (
	"fmt"
	"github.com/xuzhuoxi/infra-go/timex"
	"time"
)

// CfgConn
// 网络连接配置
type CfgConn struct {
	NetName string `yaml:"net-name"` // 网络连接名称
	Network string `yaml:"network"`  // 网络连接协议
	NetAddr string `yaml:"net-addr"` // 网络连接地址
}

func (o CfgConn) String() string {
	return fmt.Sprintf("{NetName=%s, Network=%s, NetAddr=%s}", o.NetName, o.Network, o.NetAddr)
}

// CfgRabbitServer
// 服务器配置，对应server.yaml文件
type CfgRabbitServer struct {
	Servers []CfgRabbitServerItem `yaml:"servers"`
}

// CfgRabbitServerItem
// 服务器单项配置
type CfgRabbitServerItem struct {
	Id         string    `yaml:"id"`                // 服务器ID
	PlatformId string    `yaml:"p-id"`              // 服务器平台ID，用于分组
	TypeName   string    `yaml:"name"`              // 服务器类型名称
	LogRef     string    `yaml:"log-ref,omitempty"` // 日志引用配置，可在log.yaml文件查看
	Home       CfgHome   `yaml:"home"`              // 与Rabbit-Home的注册通信相关配置
	Client     CfgClient `yaml:"client"`            // 客户端用户相关配置
}

func (o CfgRabbitServerItem) String() string {
	return fmt.Sprintf("{Id=%s, PId=%s, TypeName=%s, LogRef=%s, Home=%v, Client=%v}",
		o.Id, o.PlatformId, o.TypeName, o.LogRef, o.Home, o.Client)
}

// CfgHome
// 与Rabbit-Home相关的配置
type CfgHome struct {
	CfgConn `yaml:",inline"`                   // 连接配置, 匿名字段必须使用inline否则没有值
	Enable  bool   `yaml:"enable,omitempty"`   // 是否启用与 Rabbit-Home 的通信
	Post    bool   `yaml:"post,omitempty"`     // 是否使用Post模式的Http通信
	Rate    string `yaml:"rate"`               // 与Home更新频率
	Encrypt bool   `yaml:"encrypt,omitempty"`  // 是否启用加密验证
	KeyPath string `yaml:"key-path,omitempty"` // 私钥文件路径
}

func (o CfgHome) String() string {
	return fmt.Sprintf("{Enable=%v, Net=%v, Post=%v, Rate=%s, Encrypt=%v, KeyPath=%s}", o.Enable, o.Network, o.Post, o.Rate, o.Encrypt, o.KeyPath)
}

func (o CfgHome) RateDuration() time.Duration {
	return timex.ParseDuration(o.Rate)
}

// CfgClient
// 与客户端相关的配置
type CfgClient struct {
	CfgConn   `yaml:",inline"`                          // 连接配置, 匿名字段必须使用inline否则没有值
	Encrypt   bool         `yaml:"encrypt,omitempty"`   // 是否使用自定义配置，值为true时，blocks与allows才生效
	Extension CfgExtension `yaml:"extension,omitempty"` // 逻辑扩展配置
}

func (o CfgClient) String() string {
	return fmt.Sprintf("{Net=%v, Encrypt=%v, Extension=%v}", o.CfgConn, o.Encrypt, o.Extension)
}

// CfgExtension
// 逻辑扩展配置
type CfgExtension struct {
	Custom bool     `yaml:"custom,omitempty"` // 是否启用定制，优先级：blocks > allows
	Blocks []string `yaml:"blocks,omitempty"` // 禁止的扩展列表, "all"代表全部
	Allows []string `yaml:"allows,omitempty"` // 允许的扩展列表, "all"代表全部
}

func (o CfgExtension) String() string {
	return fmt.Sprintf("{Custom=%v, Blocks=%v, Allows=%v}", o.Custom, o.Blocks, o.Allows)
}
