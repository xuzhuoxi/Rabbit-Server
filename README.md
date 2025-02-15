# Rabbit-Server

[English](./README_EN.md) | 简体中文

## 概述

Rabbit-Server 是一个功能全面的游戏服务器框架，支持多种网络连接方式、自定义通信协议、逻辑扩展、MMO 世界、
Rabbit-Server 支持集群部署，配合 Rabbit-Home 能够很好地进行动态集群扩展。

## 功能特性

- 支持多种网络连接方式，包括http,tcp,udp,quic,ws
- 支持自定义通信协议
- 支持逻辑扩展
- 支持MMO世界
- 支持集群部署，动态集群节点更新，自动发现，无限制横向扩展。
- 支持运行时命令行操作控制

## 快速开始

通过 `go mod` 或 克隆整个仓库到 `gopath` 中便在项目中使用 Rabbit-Server.

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
res/
├── rabbit.yaml
└── conf/
    ├── clock.yaml
    ├── log.yaml
    ├── mmo.yaml
    ├── mysql.yaml
    ├── server.yaml
    └── verify.yaml
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
    - id: 房间唯一id, 在整个mmo中的全部实体中唯一。
    - rid: 引用id，用于关联其它配置，如地图配置等。可重复。
    - name: 名称，用于显示。
    - cap: 容量，用于限制玩家数量。
    - tags: 标签，数组项。 可以为房间设置多个标签属性。
  + 区域: relations/zones 属性下，为数组项目
    - id: 区域唯一id, 在整个mmo中的全部实体中唯一。
    - name: 区域名称，用于显示。
    - list: 关联房间，数组项。 可配置房间id。
  + 世界: relations/worlds 属性下，为数组项目
    - id: 世界唯一id, 在整个mmo中的全部实体中唯一。
    - name: 世界名称，用于显示。
    - list: 关联实体，数组项。 可配置 房间id 或 区域id。
  + default: 玩家登录后的默认房间配置，数组项，支持多个，随机选择。
  + log_ref: 世界相关日志配置， 值为log.yaml中的日志配置名称。

- server.yaml
  + 用于服务器逻辑配置行为。
  + 支持 与 Rabbit-Home 集群配置，包括 连接协议、连接url、心跳间隔、~~连接超时、连接重试次数~~ 等配置。
  + 支持以 Extension 为基础的逻辑扩展的配置。
  + log_ref: 服务器逻辑相关日志配置， 值为log.yaml中的日志配置名称。

- verify.yaml
  + 用于配置服务器验证行为。
  + 支持以 Extension 和 ProtoId 为基础的的响应验证行为。包括 每秒最高请求次数、 最小请求间隔 的设置。

### 在项目中使用

1. 使用 `github.com/xuzhuoxi/infra-go/serialx` 中的 `NewStartupManager()` 创建一个启动管理器 `IStartupManager`。
   ```go
   startup := serialx.NewStartupManager()
   ```
2. 创建启动模块，实际是一个实现IStartupModule的结构体。
   ```go
   import (
   "github.com/xuzhuoxi/Project2208_Server/src/core"
   "github.com/xuzhuoxi/Rabbit-Server/engine/mgr"
   "github.com/xuzhuoxi/infra-go/eventx"
   "github.com/xuzhuoxi/infra-go/serialx"
   )
   
   type rabbitServer struct {
       eventx.EventDispatcher
   }
    
   func (o *rabbitServer) Name() string {
       return "Init Servers"
   }
    
   func (o *rabbitServer) StartModule() {
       mgr := NewRabbitManager()
       mgr.GetInitManager().LoadRabbitConfig("rabbit.yaml")            // 加载配置文件
       mgr.GetInitManager().InitLoggerManager()                        // 初始化日志管理器
	   mgr.GetServerManager().StartServers()                           // 启动服务器
       o.DispatchEvent(serialx.EventOnStartupModuleStarted, o, nil)
   }
    
   func (o *initServer) StopModule() {
	   o.DispatchEvent(serialx.EventOnStartupModuleStopped, o, nil)
   }
    
   func (o *initServer) SaveModule() {
       o.DispatchEvent(serialx.EventOnStartupModuleSaved, o, nil)
   }
   ```
   **注意：**
   - 实际项目中，启动模块至少有一个，一般不只一个。
   - 在 StartModule、StopModule、SaveModule 函数完成后，事件**必须**抛出，否则会处理一直等待状态。

3. 把模块结构体添加到 startup 中去，并启动管理器。
   ```go
   startup.AppendModule(&rabbitServer{})
   startup.StartManager()
   ```
   **注意：**
   - 实际项目中，启动模块至少有一个，一般一只一个。
   - 在 StartModule、StopModule、SaveModule 函数完成后，事件**必须**抛出，否则会处理一直等待状态。
   - 一般情况下把 加载配置、日志管理、数据库管理、缓存管理、业务逻辑管理、MMO世界管理等功能划分为不同的模块，并按顺序添加。
   - 模块的启动完全按照添加顺序，在事件抛出后代表执行结束。

## 关键功能和代码逻辑

### 目录说明

1. **engine/**: 包含核心游戏服务器逻辑。
   - **config.go**: 配置文件解析和管理。
   - **config/**: 配置文件解析和管理功能支持。
   - **db/**: 数据库管理，目前主要支持 MySQL。
   - **mmo/**: 多人在线游戏（MMO）相关逻辑。
     + **basis/**: 基础实体和管理器。
     + **config/**: MMO 配置管理。
     + **events/**: 事件处理。
     + **index/**: 实体索引管理。
     + **manager/**: 实体管理。
     + **proto/**: 协议定义。
     + **vars/**: 变量定义。
   - **server/**: 服务器核心逻辑。
     + **core/**: 核心服务器管理。
     + **extension/**: 扩展管理。
     + **packet/**: 数据包处理。
     + **status/**: 状态管理。
   - **verify/**: 验证逻辑。
   - **clock/**: 时钟管理。
   - **mgr/**: 管理器。
   - **utils/**: 工具函数。

2. **res/**: 资源文件。
   - **conf/**: 配置文件。

3. **demo/**: 示例代码。
   - **client/**: 客户端示例。
     - **net/**: 网络客户端。
     - **proto/**: 协议示例。
     - **main.go**: 客户端示例程序入口。
   - **server/**: 服务器示例。
     - **cmd/**: 命令行工具。
     - **extension/**: 扩展示例。
     - **main.go**: 服务器示例程序入口。

### 关键功能说明

#### 1. 配置管理

##### `engine/config.go`
- **LoadRabbitRootConfig**: 加载根配置文件。
- **CfgRabbitRoot**: 根配置结构体。

##### `engine/config/clock.go`
- **CfgClock**: 时钟配置结构体。

##### `engine/config/log.go`
- **CfgLog**: 日志配置结构体。

##### `engine/config/server.go`
- **CfgServer**: 服务器配置结构体。

##### `engine/config/verify.go`
- **CfgVerify**: 验证配置结构体。

#### 2. 数据库管理

##### `engine/db/mysql/manager.go`
- **DataSourceManager**: 数据源管理器，负责初始化、打开、更新和关闭数据源。
- **IDataSourceManager**: 数据源管理器接口。

#### 3. MMO 逻辑

##### `engine/mmo/basis/`
- **IEntity**: 实体接口。
- **IPlayerEntity**: 玩家实体接口。
- **IRoomEntity**: 房间实体接口。
- **ITeamEntity**: 队伍实体接口。
- **IChannelEntity**: 频道实体接口。
- **IManagerBase**: 管理器基础接口。
- **IEventDispatcher**: 事件分发器接口。

##### `engine/mmo/events/`
- **事件处理**: 包含各种事件处理逻辑，如玩家事件、房间事件、队伍事件等。

##### `engine/mmo/index/`
- **索引管理**: 包含各种实体的索引管理，如玩家索引、房间索引、队伍索引等。

##### `engine/mmo/manager/`
- **EntityManager**: 实体管理器，负责创建和管理各种实体。
- **PlayerManager**: 玩家管理器，负责玩家的进入、离开、转移等操作。

##### `engine/mmo/proto/`
- **协议定义**: 包含各种协议的定义。

##### `engine/mmo/vars/`
- **变量定义**: 包含各种变量的定义，如玩家变量、房间变量、单位变量等。

#### 4. 服务器核心逻辑

##### `engine/server/core/`
- **RabbitServer**: 核心服务器类，负责启动、停止、重启服务器。
- **CustomRabbitManager**: 自定义管理器，负责处理游戏数据包。

##### `engine/server/extension/`
- **扩展管理**: 包含各种扩展的管理，如请求扩展、响应扩展、通知扩展等。

##### `engine/server/packet/`
- **数据包处理**: 包含数据包的创建、解析、验证等逻辑。

##### `engine/server/status/`
- **状态管理**: 包含服务器状态的管理。

#### 5. 验证逻辑

##### `engine/verify/verify.go`
- **验证逻辑**: 包含各种验证逻辑，如登录验证、数据验证等。

#### 6. 时钟管理

##### `engine/clock/clock.go`
- **时钟管理**: 包含时钟的初始化、设置、获取等逻辑。

#### 7. 管理器

##### `engine/mgr/`
- **管理器**: 包含各种管理器，如连接管理器、服务器管理器、状态管理器等。

#### 8. 工具函数

##### `engine/utils/`
- **工具函数**: 包含各种工具函数，如 JSON 处理、路径处理、YAML 处理等。

### 关键代码文明

#### 1. 玩家管理器 (`engine/mmo/manager/player.go`)

##### `forwardTransfer`
- **功能**: 将玩家从一个房间转移到另一个房间。
- **逻辑**:
  1. 检查目标房间是否存在。
  2. 检查玩家是否已经在目标房间中。
  3. 从当前房间移除玩家。
  4. 将玩家添加到目标房间。
  5. 确认玩家的下一个房间。
  6. 触发相应的事件。

##### `backwardTransfer`
- **功能**: 将玩家从当前房间转移到前一个房间。
- **逻辑**:
  1. 获取玩家的前一个房间。
  2. 检查前一个房间是否存在。
  3. 从当前房间移除玩家。
  4. 将玩家添加到前一个房间。
  5. 确认玩家的前一个房间。
  6. 触发相应的事件。

#### 2. 实体管理器 (`engine/mmo/manager/entity.go`)

##### `CreateRoom`
- **功能**: 创建一个新的房间。
- **逻辑**:
  1. 检查房间是否已经存在。
  2. 创建新的房间实体。
  3. 初始化房间实体。
  4. 设置房间变量和标签。
  5. 将房间添加到房间索引中。
  6. 添加实体事件监听器。
  7. 触发房间初始化事件。

##### `CreatePlayer`
- **功能**: 创建一个新的玩家。
- **逻辑**:
  1. 检查玩家是否已经存在。
  2. 创建新的玩家实体。
  3. 初始化玩家实体。
  4. 设置玩家变量。
  5. 将玩家添加到玩家索引中。
  6. 添加实体事件监听器。
  7. 触发玩家初始化事件。

#### 3. 数据包处理 (`engine/server/packet/packet.go`)

##### `AppendJson`
- **功能**: 将 JSON 数据追加到数据包中。
- **逻辑**:
  1. 检查数据是否为空。
  2. 将每个数据项转换为 JSON 字符串。
  3. 将 JSON 字符串追加到数据包中。

##### `AppendObject`
- **功能**: 将对象数据追加到数据包中。
- **逻辑**:
  1. 检查数据是否为空。
  2. 检查参数处理器是否为空。
  3. 使用参数处理器编码每个数据项。
  4. 将编码后的数据追加到数据包中。

#### 4. 时钟管理 (`engine/clock/clock.go`)

##### `Init`
- **功能**: 初始化时钟管理器。
- **逻辑**:
  1. 检查配置是否为空。
  2. 初始化配置。
  3. 设置服务器启动时间和启动时间戳。

##### `NowGameClock`
- **功能**: 获取游戏当前时钟。
- **逻辑**:
  1. 获取游戏运行时长。
  2. 计算游戏当前时钟。

#### 5. 数据源管理 (`engine/db/mysql/manager.go`)

##### `Init`
- **功能**: 初始化数据源管理器。
- **逻辑**:
  1. 修复配置文件路径。
  2. 解析配置文件。
  3. 创建数据源实例。
  4. 设置表元数据。
  5. 触发管理器初始化事件。

##### `OpenAll`
- **功能**: 打开所有数据源连接。
- **逻辑**:
  1. 初始化索引。
  2. 打开数据源连接。

#### 6. 扩展管理 (`engine/server/core/manager.go`)

##### `onRabbitGamePack`
- **功能**: 处理游戏数据包。
- **逻辑**:
  1. 解析消息数据。
  2. 验证扩展。
  3. 获取回收参数。
  4. 前置处理。
  5. 请求处理。
  6. 后置处理。

## 依赖库
- infra-go [https://github.com/xuzhuoxi/infra-go](https://github.com/xuzhuoxi/infra-go)<br>
  基础库支持库。
- goxc [https://github.com/laher/goxc](https://github.com/laher/goxc)<br>
  打包依赖库，主要用于交叉编译
- json-iterator [https://github.com/json-iterator/go](https://github.com/json-iterator/go)<br>
  带对应结构体的Json解释库

## 联系作者
xuzhuoxi<br>
<xuzhuoxi@gmail.com> or <mailxuzhuoxi@163.com> or <m_xuzhuoxi@outlook.com>

## License
Rabbit-Server source code is available under the MIT [License](/LICENSE).