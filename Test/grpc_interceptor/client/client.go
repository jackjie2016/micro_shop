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
	Interceptor:=func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error{
        start:=time.Now();

        err:=invoker(ctx, method, req, reply, cc, opts ...)
        fmt.Printf("执行时间：%s\n",time.Since(start))
		return err;
	}
	var opt []grpc.DialOption

	opt=append(opt,grpc.WithInsecure())
	opt=append(opt,grpc.WithUnaryInterceptor(Interceptor))
	conn,err:=grpc.Dial("127.0.0.1:50053",opt...)
	//opt:=grpc.WithUnaryInterceptor(Interceptor)
	//conn,err:=grpc.Dial("127.0.0.1:50053",grpc.WithInsecure(),opt)
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