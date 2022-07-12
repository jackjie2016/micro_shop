package main

import (
	"OldPackageTest/new_helloworld/hanlder"
	"net"
	"net/rpc"


	"OldPackageTest/new_helloworld/server_proxy"

)
//type HelloService struct {}
//func (s *HelloService) Hello(request string ,reply *string) error {
//	//返回值是通过修改reply的值
//	*reply="hello,"+request
//	return  nil
//}
func main()  {
	//1、实例化一个server
	listener,_:=net.Listen("tcp",":1235")

	//2.注册处理逻辑handle

	_ = server_proxy.RegisterHelloService(&hanlder.HelloService{})

	//3.启动服务
	for  {
		conn,_:=listener.Accept()//当一个新的连接进来的时候
		go rpc.ServeConn(conn)
	}

}