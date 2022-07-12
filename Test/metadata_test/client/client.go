package main

import (

	"context"
	"fmt"
	"google.golang.org/grpc/metadata"
	"time"

	"google.golang.org/grpc"
	timestamppb "google.golang.org/protobuf/types/known/timestamppb"


	"OldPackageTest/metadata_test/proto"
)

func main()  {
	conn,err:=grpc.Dial("127.0.0.1:50053",grpc.WithInsecure())
	if err !=nil {
		panic(err)
	}
	defer conn.Close()

	c:=proto.NewGreeterClient(conn)
	//md:=metadata.New(map[string]string{
	//	"name":[]string{"zifeng","zifeng"},
	//
	//	"password":"123567",
	//})
	md:=metadata.Pairs(
		"name","zifeng",
		)
	ctx:=metadata.NewOutgoingContext(context.Background(),md)

	r,err:=c.SayHello(ctx,&proto.HelloRequest{
		Name:"bobby",
		Url: "baidu.com",
		G:proto.Gender_MALE,

		AddTime:timestamppb.New(time.Now()),

	})

	if err !=nil {
		panic(err)
	}
	fmt.Println(r.Message)

}