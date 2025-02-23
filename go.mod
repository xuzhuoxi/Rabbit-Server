module github.com/xuzhuoxi/Rabbit-Server

go 1.16

require (
	github.com/go-sql-driver/mysql v1.7.1
	github.com/json-iterator/go v1.1.12
	github.com/xuzhuoxi/Rabbit-Home v1.1.0
	github.com/xuzhuoxi/infra-go v1.1.0
	gopkg.in/yaml.v2 v2.4.0
)

replace (
	github.com/xuzhuoxi/infra-go v1.1.0 => ../infra-go
	github.com/xuzhuoxi/Rabbit-Home v1.1.0 => ../Rabbit-Home
)
