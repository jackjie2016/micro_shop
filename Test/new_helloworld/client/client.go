package main

import (
	"OldPackageTest/new_helloworld/client_proxy"
	"fmt"
)

func main()  {
	//1.建立链接

	//client,err:=rpc.Dial("tcp",":1235")
	client:=client_proxy.NewHelloServiceClient("tcp",":1235")

	//var reply *string =new(string) //string 没有默认值
	var reply string //string 有默认值

	err:= client.Hello("bobby",&reply)

	if err !=nil{
		panic("调用失败")
	}

	fmt.Println(reply)

}

