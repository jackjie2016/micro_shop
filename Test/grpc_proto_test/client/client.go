package main

import (

	"context"
	"fmt"
	"time"

	"google.golang.org/grpc"
	timestamppb "google.golang.org/protobuf/types/known/timestamppb"

	"OldPackageTest/grpc_proto_test/proto"
)

func main()  {
	conn,err:=grpc.Dial("127.0.0.1:50053",grpc.WithInsecure())
	if err !=nil {
		panic(err)
	}
	defer conn.Close()

	c:=proto.NewGreeterClient(conn)
	r,err:=c.SayHello(context.Background(),&proto.HelloRequest{
		Name:"bobby",
		Url: "baidu.com",
		G:proto.Gender_MALE,
		Mp: map[string]string{
			"name":"zifeng",
			"company":"慕课网",
		},
		AddTime:timestamppb.New(time.Now()),

	})

	if err !=nil {
		panic(err)
	}
	fmt.Println(r.Message)

}