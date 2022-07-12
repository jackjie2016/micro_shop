package client_proxy

import (
	"net/rpc"

	"OldPackageTest/new_helloworld/hanlder"
)

type HelloServiceStub struct {
	*rpc.Client
}
//在go语言中没有类、对象 就意味着没有初始化方法
func NewHelloServiceClient(protol,address string) HelloServiceStub  {
	conn,err:=rpc.Dial(protol,address)
	if err!=nil{
		panic("connect err！")
	}
	return HelloServiceStub{conn}
}

func (c *HelloServiceStub) Hello (request string,reply *string)  error {
	err:=c.Call(hanlder.HelloServerName+".Hello",request,reply)
	if err!=nil{
		return err
	}
	return nil
}