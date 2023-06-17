// Package demo
// Created by xuzhuoxi
// on 2019-03-24.
// @author xuzhuoxi
//
package demo

import (
	"fmt"
	jsoniter "github.com/json-iterator/go"
	"github.com/xuzhuoxi/Rabbit-Server/demo/client"
	"github.com/xuzhuoxi/infra-go/bytex"
)

type testA struct {
	A string
	B int
	C bool
}

type I interface {
	Func1()
	Func2()
	Func3()
}

type B struct{}

func (b *B) Func1() {
	b.Func3()
}

func (*B) Func3() {
	fmt.Println("B3")
}

type S struct {
	B
}

//func (s *S) Func1() {
//	fmt.Println("S1")
//}

func (s *S) Func2() {
	s.Func1()
}

func (s *S) Func3() {
	fmt.Println("S3")
}

func TestDemoExtension(uc *client.UserClient) {
	bsName := []byte("ObjDemo")
	bsPid := []byte("Obj_0")
	bsUid := []byte("顶你个肺")
	data := testA{A: "A", B: 99, C: false}

	buffToBlock := bytex.NewBuffToBlock(bytex.NewDefaultDataBlockHandler())
	buffToBlock.WriteData(bsName)
	buffToBlock.WriteData(bsPid)
	buffToBlock.WriteData(bsUid)
	dataBs, _ := jsoniter.Marshal(data)
	buffToBlock.WriteData(dataBs)
	buffToBlock.WriteData(dataBs)
	uc.SockClient.SendPackTo(buffToBlock.ReadBytes())
}
