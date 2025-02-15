# Rabbit-Server

English | [简体中文](./README.md)

## Overview

Rabbit-Server is a comprehensive game server framework that supports multiple network connection methods, custom communication protocols, logic extensions, MMO worlds, and cluster deployment. It can be dynamically scaled with Rabbit-Home for seamless horizontal expansion.

## Features

- Supports multiple network connection methods including http, tcp, udp, quic, ws.
- Supports custom communication protocols.
- Supports logic extensions.
- Supports MMO world.
- Supports cluster deployment, dynamic node updates, automatic discovery, and unlimited horizontal scaling.
- Supports runtime command-line operations.

## Quick Start

You can use Rabbit-Server in your project by using `go mod` or downloading the entire repository to `gopath`.

### Prerequisites

- Go version requirement: 1.16 or higher

### Installation Steps

Supports using go.mod or gopath to manage the repository

- Through gopath loading the repository

    ```
    go get -u github.com/xuzhuoxi/Rabbit-Server
    ```
- Through go.mod loading the repository

1. Add in the go.mod of your project

   ```
   require (
       github.com/xuzhuoxi/infra-go
       github.com/xuzhuoxi/Rabbit-Home
       github.com/xuzhuoxi/Rabbit-Server
   )     
   ```  

2. Install dependencies

   ```shell
   go mod tidy
   ```

### Adding Configuration

Configuration files in Rabbit-Server use YAML format, consisting of a root configuration file and several secondary configuration files.

Here is an example based on the configuration files in the (res)[/res] directory.

Target structure:

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
    + **Root configuration file**, the filename can be customized, associated by adding parameter `-conf=rabbit.yaml` when running the project.
    + Supports configuring 6 attribute projects: log, clock, db, server, verify, mmo. The attribute values are paths to corresponding secondary configuration files, it's recommended to use relative paths.

- log.yaml
    + Used for configuring logging behavior.
    + Supports configurations such as log level(level), log type(type), log path(path), log file size limit(size).
    + Can configure multiple logging behaviors, may specify default logging behavior.

- clock.yaml
    + Used for configuring server clock behavior.
    + Supports 3 attribute configurations: clock timezone(game_loc), server start time(game_zero), daily reset time(daily_zero).

- mmo.yaml
    + Used for configuring MMO-related rooms, zones, worlds.
    + Rooms: under entities/rooms attribute, array items
        - id: unique room id, unique among all entities in the MMO.
        - rid: reference id, used to associate other configurations like map configurations. Can be repeated.
        - name: display name.
        - cap: capacity, limits player numbers.
        - tags: tag attributes, array items. Multiple tags can be set for a room.
    + Zones: under relations/zones attribute, array items
        - id: unique zone id, unique among all entities in the MMO.
        - name: display name.
        - list: associated rooms, array items. Configurable with room ids.
    + Worlds: under relations/worlds attribute, array items
        - id: unique world id, unique among all entities in the MMO.
        - name: display name.
        - list: associated entities, array items. Configurable with room ids or zone ids.
    + default: default room configuration after player login, array items, supports multiple, randomly selected.
    + log_ref: related log configuration for the world, value is the log configuration name in log.yaml.

- server.yaml
    + Used for server logic configuration.
    + Supports Rabbit-Home cluster configuration, including connection protocol, connection url, heartbeat interval, etc.
    + Supports configuration based on Extension.
    + log_ref: related log configuration for server logic, value is the log configuration name in log.yaml.

- verify.yaml
    + Used for configuring server verification behavior.
    + Supports response verification behavior based on Extension and ProtoId, including settings like maximum requests per second, minimum request interval.

### Using in Projects

1. Create a startup manager `IStartupManager` using `NewStartupManager()` from `github.com/xuzhuoxi/infra-go/serialx`.
   ```go
   startup := serialx.NewStartupManager()
   ``` 

2. Create a startup module, which is actually a struct implementing `IStartupModule`.
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
       mgr.GetInitManager().LoadRabbitConfig("rabbit.yaml")            // Load configuration file
       mgr.GetInitManager().InitLoggerManager()                        // Initialize logger manager
       mgr.GetServerManager().StartServers()                           // Start servers
       o.DispatchEvent(serialx.EventOnStartupModuleStarted, o, nil)
   }
    
   func (o *initServer) StopModule() {
	   o.DispatchEvent(serialx.EventOnStartupModuleStopped, o, nil)
   }
    
   func (o *initServer) SaveModule() {
       o.DispatchEvent(serialx.EventOnStartupModuleSaved, o, nil)
   }
   ```
   **Note:**
   - In actual projects, there should be at least one startup module, generally more than one.
   - After completing `StartModule`, `StopModule`, and `SaveModule` functions, events **must** be dispatched; otherwise, it will remain in a waiting state.

3. Add the module struct to `startup` and start the manager.
   ```go
   startup.AppendModule(&rabbitServer{})
   startup.StartManager()
   ```

   **Note:**
   - In actual projects, there should be at least one startup module, generally only one.
   - After completing `StartModule`, `StopModule`, and `SaveModule` functions, events **must** be dispatched; otherwise, it will remain in a waiting state.
   - Generally, functionalities like loading configuration, initializing logging, database management, cache management, business logic management, MMO world management, etc., are divided into different modules and added in order.
   - Modules are started strictly in the order they are added, and the execution ends after the event is dispatched.

## Key Features and Code Logic

### Directory Explanation

1. **engine/**: Contains core game server logic.
   - **config.go**: Configuration file parsing and management.
   - **config/**: Support for configuration file parsing and management.
   - **db/**: Database management, currently mainly supports MySQL.
   - **mmo/**: Multiplayer online game (MMO) related logic.
     + **basis/**: Basic entities and managers.
     + **config/**: MMO configuration management.
     + **events/**: Event handling.
     + **index/**: Entity index management.
     + **manager/**: Entity management.
     + **proto/**: Protocol definitions.
     + **vars/**: Variable definitions.
   - **server/**: Core server logic.
     + **core/**: Core server management.
     + **extension/**: Extension management.
     + **packet/**: Packet processing.
     + **status/**: Status management.
   - **verify/**: Verification logic.
   - **clock/**: Clock management.
   - **mgr/**: Managers.
   - **utils/**: Utility functions.

2. **res/**: Resource files.
   - **conf/**: Configuration files.

3. **demo/**: Example code.
   - **client/**: Client examples.
     + **net/**: Network client.
     + **proto/**: Protocol examples.
     + **main.go**: Entry point for client example program.
   - **server/**: Server examples.
     + **cmd/**: Command line tools.
     + **extension/**: Extension examples.
     + **main.go**: Entry point for server example program.

### Key Feature Explanations

#### 1. Configuration Management

##### `engine/config.go`
- **LoadRabbitRootConfig**: Loads the root configuration file.
- **CfgRabbitRoot**: Root configuration structure.

##### `engine/config/clock.go`
- **CfgClock**: Clock configuration structure.

##### `engine/config/log.go`
- **CfgLog**: Log configuration structure.

##### `engine/config/server.go`
- **CfgServer**: Server configuration structure.

##### `engine/config/verify.go`
- **CfgVerify**: Verification configuration structure.

#### 2. Database Management

##### `engine/db/mysql/manager.go`
- **DataSourceManager**: Data source manager responsible for initializing, opening, updating, and closing data sources.
- **IDataSourceManager**: Data source manager interface.

#### 3. MMO Logic

##### `engine/mmo/basis/`
- **IEntity**: Entity interface.
- **IPlayerEntity**: Player entity interface.
- **IRoomEntity**: Room entity interface.
- **ITeamEntity**: Team entity interface.
- **IChannelEntity**: Channel entity interface.
- **IManagerBase**: Base manager interface.
- **IEventDispatcher**: Event dispatcher interface.

##### `engine/mmo/events/`
- **Event Handling**: Contains various event handling logic, such as player events, room events, team events, etc.

##### `engine/mmo/index/`
- **Index Management**: Contains index management for various entities, such as player indexes, room indexes, team indexes, etc.

##### `engine/mmo/manager/`
- **EntityManager**: Entity manager responsible for creating and managing various entities.
- **PlayerManager**: Player manager responsible for player entry, exit, transfer, etc.

##### `engine/mmo/proto/`
- **Protocol Definitions**: Contains definitions for various protocols.

##### `engine/mmo/vars/`
- **Variable Definitions**: Contains definitions for various variables, such as player variables, room variables, unit variables, etc.

#### 4. Core Server Logic

##### `engine/server/core/`
- **RabbitServer**: Core server class responsible for starting, stopping, and restarting the server.
- **CustomRabbitManager**: Custom manager responsible for handling game packets.

##### `engine/server/extension/`
- **Extension Management**: Contains management for various extensions, such as request extensions, response extensions, notification extensions, etc.

##### `engine/server/packet/`
- **Packet Processing**: Contains logic for creating, parsing, and validating packets.

##### `engine/server/status/`
- **Status Management**: Contains logic for managing server status.

#### 5. Verification Logic

##### `engine/verify/verify.go`
- **Verification Logic**: Contains various verification logic, such as login verification, data validation, etc.

#### 6. Clock Management

##### `engine/clock/clock.go`
- **Clock Management**: Contains logic for initializing, setting, and getting the clock.

#### 7. Managers

##### `engine/mgr/`
- **Managers**: Contains various managers, such as connection managers, server managers, status managers, etc.

#### 8. Utility Functions

##### `engine/utils/`
- **Utility Functions**: Contains various utility functions, such as JSON processing, path processing, YAML processing, etc.

### Key Code Civilizations

#### 1. Player Manager (`engine/mmo/manager/player.go`)

##### `forwardTransfer`
- **Function**: Transfers a player from one room to another.
- **Logic**:
  1. Check if the target room exists.
  2. Check if the player is already in the target room.
  3. Remove the player from the current room.
  4. Add the player to the target room.
  5. Confirm the player's next room.
  6. Trigger relevant events.

##### `backwardTransfer`
- **Function**: Transfers a player from the current room to the previous room.
- **Logic**:
  1. Get the player's previous room.
  2. Check if the previous room exists.
  3. Remove the player from the current room.
  4. Add the player to the previous room.
  5. Confirm the player's previous room.
  6. Trigger relevant events.

#### 2. Entity Manager (`engine/mmo/manager/entity.go`)

##### `CreateRoom`
- **Function**: Creates a new room.
- **Logic**:
  1. Check if the room already exists.
  2. Create a new room entity.
  3. Initialize the room entity.
  4. Set room variables and tags.
  5. Add the room to the room index.
  6. Add entity event listeners.
  7. Trigger room initialization events.

##### `CreatePlayer`
- **Function**: Creates a new player.
- **Logic**:
  1. Check if the player already exists.
  2. Create a new player entity.
  3. Initialize the player entity.
  4. Set player variables.
  5. Add the player to the player index.
  6. Add entity event listeners.
  7. Trigger player initialization events.

#### 3. Packet Processing (`engine/server/packet/packet.go`)

##### `AppendJson`
- **Function**: Appends JSON data to the packet.
- **Logic**:
  1. Check if the data is empty.
  2. Convert each data item to a JSON string.
  3. Append the JSON string to the packet.

##### `AppendObject`
- **Function**: Appends object data to the packet.
- **Logic**:
  1. Check if the data is empty.
  2. Check if the parameter processor is empty.
  3. Encode each data item using the parameter processor.
  4. Append the encoded data to the packet.

#### 4. Clock Management (`engine/clock/clock.go`)

##### `Init`
- **Function**: Initializes the clock manager.
- **Logic**:
  1. Check if the configuration is empty.
  2. Initialize the configuration.
  3. Set the server start time and timestamp.

##### `NowGameClock`
- **Function**: Gets the current game clock.
- **Logic**:
  1. Get the game runtime duration.
  2. Calculate the current game clock.

#### 5. Data Source Management (`engine/db/mysql/manager.go`)

##### `Init`
- **Function**: Initializes the data source manager.
- **Logic**:
  1. Fix the configuration file path.
  2. Parse the configuration file.
  3. Create data source instances.
  4. Set table metadata.
  5. Trigger manager initialization events.

##### `OpenAll`
- **Function**: Opens all data source connections.
- **Logic**:
  1. Initialize indexes.
  2. Open data source connections.

#### 6. Extension Management (`engine/server/core/manager.go`)

##### `onRabbitGamePack`
- **Function**: Handles game packets.
- **Logic**:
  1. Parse message data.
  2. Validate extension.
  3. Get recycling parameters.
  4. Pre-processing.
  5. Request processing.
  6. Post-processing.

## Related Libraries
- infra-go [https://github.com/xuzhuoxi/infra-go](https://github.com/xuzhuoxi/infra-go)<br>
  Basic library support.
- goxc [https://github.com/laher/goxc](https://github.com/laher/goxc)<br>
  Packaging dependency library, mainly used for cross-compilation.
- json-iterator [https://github.com/json-iterator/go](https://github.com/json-iterator/go)<br>
  JSON interpretation library with corresponding structures.

## Contact
xuzhuoxi<br>
<xuzhuoxi@gmail.com> or <mailxuzhuoxi@163.com>

## License
Rabbit-Server source code is available under the MIT [License](/LICENSE).
        