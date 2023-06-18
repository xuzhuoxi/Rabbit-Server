// Package login
// Create on 2023/6/17
// @author xuzhuoxi
package login

import (
	"github.com/xuzhuoxi/Rabbit-Server/demo/client/net"
	"github.com/xuzhuoxi/infra-go/bytex"
)

func TestLoginExtension(uc *net.UserClient) {
	dataBlock := bytex.NewBuffDataBlock(bytex.NewDefaultDataBlockHandler())
	dataBlock.WriteString("Login")   //ExtensionName
	dataBlock.WriteString("LI")      //ProtoId
	dataBlock.WriteString(uc.UserId) //Uid
	dataBlock.WriteString(uc.UserId) //Data(Password)
	uc.SockClient.SendPackTo(dataBlock.ReadBytes())
}

func TestReLoginExtension(uc *net.UserClient) {
	buffToBlock := bytex.NewBuffDataBlock(bytex.NewDefaultDataBlockHandler())
	buffToBlock.WriteString("Login")   //ExtensionName
	buffToBlock.WriteString("RLI")     //ProtoId
	buffToBlock.WriteString(uc.UserId) //Uid
	buffToBlock.WriteString(uc.UserId) //Data(Password)
	uc.SockClient.SendPackTo(buffToBlock.ReadBytes())
}
