# Rabbit-Server

## 概述

Rabbit-Server 是一个游戏服务器框架。开源、免费、扩展性良好，
配合 Rabbit-Home 能够很好地进行动态集群扩展。

## 功能特性

- 支持多种网络连接方式，包括http,tcp,udp,quic,ws
- 支持自定义通信协议
- 支持逻辑扩展
- 支持MMO世界
- 支持集群部署，动态集群节点更新，自动发现，无限制横向扩展。
- 支持运行时命令行操作控制

## 安装与运行

通过 `go mod` 或 下载整个仓库到 `gopath` 中便在项目中使用 Rabbit-Server.

### 前提条件

- go版本要求： 1.16 或更高

### 安装步骤

支持使用go.mod或gopath管理仓库

- 通过gopath加载仓库

```
go get -u github.com/xuzhuoxi/Rabbit-Server
```

- 通过go.mod加载仓库

1. 在项目中的go.mod中加入
```
require (
	github.com/xuzhuoxi/infra-go
	github.com/xuzhuoxi/Rabbit-Home
	github.com/xuzhuoxi/Rabbit-Server
)
```

2. 安装依赖项
```shell
go mod tidy
```

### 添加配置

Rabbit-Server 中的配置文件使用yaml格式，由一个根配置文件与若干个二级配置文件组成。

这里以 (res)[/res] 目录中的配置文件举例说明。

目标结构如下：

```
./
rabbit.yaml
conf/log.yaml
conf/clock.yaml
conf/db.yaml
conf/mmo.yaml
conf/server.yaml
conf/verify.yaml
```

- rabbit.yaml
  + **根配置文件**，文件名可自定义，通过在项目运行时添加参数`-conf=rabbit.yaml`进行关联。
  + 支持 log、clock、db、server、verify、mmo 共6个属性项目配置，属性值为对应二级配置文件的路径，建议使用相对路径。

- log.yaml
  + 用于配置日志行为。
  + 支持 日志级别(level)、 日志类别(type)、日志路径(path)、日志文件大小限制(size) 等配置。
  + 可以配置多个日志行为， 可能指定默认的日志行为。

- clock.yaml
  + 用于配置服务器时钟行为。
  + 支持 时钟时区(game_loc)、 开服时刻(game_zero)、 每日零点时刻(daily_zero) 3个属性配置。

- mmo.yaml
  + 用于配置 MMO 相关的 房间、 区域、世界等相关项
  + 房间: entities/rooms 属性下，为数组项目
    + id: 房间唯一id, 在整个mmo中的全部实体中唯一。
    + rid: 引用id，用于关联其它配置，如果地图配置等。可重复。
    + name: 名称，用于显示。
    + cap: 容量，用于限制玩家数量。
    + tags: 标签，数组项。 可以为房间设置多个标签属性。
  + 区域: relations/zones 属性下，为数组项目
    + id: 区域唯一id, 在整个mmo中的全部实体中唯一。
    + name: 区域名称，用于显示。
    + list: 关联房间，数组项。 可配置房间id。
  + 世界: relations/worlds 属性下，为数组项目
    + id: 世界唯一id, 在整个mmo中的全部实体中唯一。
    + name: 世界名称，用于显示。
    + list: 关联实体，数组项。 可配置 房间id 或 区域id。
  + default: 玩家登录后的默认房间配置，数组项，支持多个，随机选择。
  + log_ref: 世界相关日志配置， 值为log.yaml中的日志配置名称。

- server.yaml
  + 服务器逻辑配置行为。
  + 支持 


- verify.yaml


### How to config

[Here](/res/conf) is config example.

[Here](/demo/server) is a example project base for Rabbit-Server.

### How to add snail to your game server.

```go
  mgr.DefaultManager.GetServerManager().StartServers()
```


## Related Library

- infra-go [https://github.com/xuzhuoxi/infra-go](https://github.com/xuzhuoxi/infra-go)<br>
  基础库支持库。

- goxc [https://github.com/laher/goxc](https://github.com/laher/goxc)<br>
  打包依赖库，主要用于交叉编译

- json-iterator [https://github.com/json-iterator/go](https://github.com/json-iterator/go)<br>
  带对应结构体的Json解释库

## Contact
xuzhuoxi<br>
<xuzhuoxi@gmail.com> or <mailxuzhuoxi@163.com>

## License
IconGen source code is available under the MIT [License](/LICENSE).