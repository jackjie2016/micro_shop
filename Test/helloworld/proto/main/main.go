package main

import (
	helloworld "OldPackageTest/helloworld/proto"
	"fmt"
	"github.com/golang/protobuf/proto"
)

type Hello struct{
	 Name string `json:"name"`
	 Age int `json:"age"`
	 Courses []string `json:"courses"`
}
func main()  {
	req :=helloworld.HelloRequest{
		Name:"bobby",
		Age:12,
		Courses:[]string{"go","php","java"},
	}
	//reqstruct:=Hello{
	//	Name:"bobby",
	//	Age:12,
	//	Courses:[]string{"go","php","java"},
	//}
	//reqjson,_:=json.Marshal(&reqstruct)
	//fmt.Println(reqjson)
	//fmt.Println(len(reqjson))
	//var Hello = new(Hello)
	//_=json.Unmarshal(reqjson,&Hello)
	//fmt.Println(*Hello)

	rsp,_:=proto.Marshal(&req)
	//fmt.Println(rsp)
	//fmt.Println(len(rsp))

	newReq:=helloworld.HelloRequest{}
	_=proto.Unmarshal(rsp,&newReq)
	fmt.Println(newReq.Age,newReq.Courses)
}
