// Package login
// Create on 2023/6/17
// @author xuzhuoxi
package login

import (
	"github.com/xuzhuoxi/Rabbit-Server/demo/client"
	"github.com/xuzhuoxi/infra-go/bytex"
)

func TestLoginExtension(uc *client.UserClient) {
	buffToBlock := bytex.NewBuffToBlock(bytex.NewDefaultDataBlockHandler())
	buffToBlock.WriteData([]byte("Login"))   //ExtensionName
	buffToBlock.WriteData([]byte("LI"))      //ProtoId
	buffToBlock.WriteData([]byte(uc.UserId)) //Uid
	buffToBlock.WriteData([]byte(uc.UserId)) //Data(Password)
	uc.SockClient.SendPackTo(buffToBlock.ReadBytes())
}

func TestReLoginExtension(uc *client.UserClient) {
	buffToBlock := bytex.NewBuffToBlock(bytex.NewDefaultDataBlockHandler())
	buffToBlock.WriteData([]byte("Login"))   //ExtensionName
	buffToBlock.WriteData([]byte("RLI"))     //ProtoId
	buffToBlock.WriteData([]byte(uc.UserId)) //Uid
	buffToBlock.WriteData([]byte(uc.UserId)) //Data(Password)
	uc.SockClient.SendPackTo(buffToBlock.ReadBytes())
}
