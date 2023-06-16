module github.com/xuzhuoxi/Rabbit-Server

go 1.16

require (
	github.com/json-iterator/go v1.1.12
	github.com/xuzhuoxi/Rabbit-Home v0.0.0-00010101000000-000000000000
	github.com/xuzhuoxi/infra-go v1.0.3
	gopkg.in/yaml.v2 v2.4.0
)

replace (
	//github.com/xuzhuoxi/Rabbit-Home => ../Rabbit-Home
	//github.com/xuzhuoxi/infra-go v1.0.3 => ../infra-go
)
