package main

import (

	"context"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/status"
	"time"

	"OldPackageTest/grpc_error_test/proto"
)

func main()  {
	conn,err:=grpc.Dial("127.0.0.1:8088",grpc.WithInsecure())
	if err !=nil {
		panic(err)
	}
	defer conn.Close()

	c:=proto.NewGreeterClient(conn)
	ctx,_:=context.WithTimeout(context.Background(),time.Second*3)
	_, err = c.SayHello(ctx, &proto.HelloRequest{Name: "bobby"})

	if err !=nil {
		st, ok := status.FromError(err)
		if !ok {
			// Error was not a status error
			panic("解析error失败")
		}

		fmt.Println(st.Message())

		fmt.Println(st.Code())
	}
	//fmt.Println(r.Message)



}