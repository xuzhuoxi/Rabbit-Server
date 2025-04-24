// Package login
// Create on 2023/6/17
// @author xuzhuoxi
package login

import (
	"fmt"
	"github.com/xuzhuoxi/Rabbit-Server/demo/client/net"
	"github.com/xuzhuoxi/infra-go/bytex"
	"github.com/xuzhuoxi/infra-go/cryptox"
)

func TestLoginExtension(uc *net.UserClient, cipher cryptox.IEncryptCipher) {
	dataBlock := bytex.NewBuffDataBlock(bytex.NewDefaultDataBlockHandler())
	dataBlock.WriteString("Login")   //ExtensionName
	dataBlock.WriteString("LI")      //ProtoId
	dataBlock.WriteString(uc.UserId) //Uid
	dataBlock.WriteString(uc.UserId) //Data(Password)
	bs := dataBlock.ReadBytes()
	//fmt.Println("[TestLoginExtension] data1:", bs)
	var err error
	if nil != cipher {
		bs, err = cipher.Encrypt(bs)
		if nil != err {
			fmt.Println("[TestLoginExtension] error:", err)
		}
	}
	//fmt.Println("[TestLoginExtension] data2:", bs)
	uc.SockClient.SendPackTo(bs)
}

func TestReLoginExtension(uc *net.UserClient, cipher cryptox.IEncryptCipher) {
	buffToBlock := bytex.NewBuffDataBlock(bytex.NewDefaultDataBlockHandler())
	buffToBlock.WriteString("Login")   //ExtensionName
	buffToBlock.WriteString("RLI")     //ProtoId
	buffToBlock.WriteString(uc.UserId) //Uid
	buffToBlock.WriteString(uc.UserId) //Data(Password)
	bs := buffToBlock.ReadBytes()
	fmt.Println("[TestReLoginExtension] data1:", bs)
	var err error
	if nil != cipher {
		bs, err = cipher.Encrypt(bs)
		if nil != err {
			fmt.Println("[TestReLoginExtension] error:", err)
		}
	}
	fmt.Println("[TestReLoginExtension] data2:", bs)
	uc.SockClient.SendPackTo(bs)
}
