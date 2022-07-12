package server_proxy

import (
	"OldPackageTest/new_helloworld/hanlder"
	"net/rpc"
)

type HelloServicer interface {
	Hello(request string, reply *string) error
}
func RegisterHelloService(srv HelloServicer) error  {
	return  rpc.RegisterName(hanlder.HelloServerName,srv)
}