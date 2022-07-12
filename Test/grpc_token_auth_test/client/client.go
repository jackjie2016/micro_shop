package main

import (
	"context"
	"fmt"
	"google.golang.org/grpc"

	"OldPackageTest/grpc_token_auth_test/proto"
)

type customCredentials struct {
	
}
//GetRequestMetadata(ctx context.Context, uri ...string) (map[string]string, error)
//// RequireTransportSecurity indicates whether the credentials requires
//// transport security.
//RequireTransportSecurity() bool
func (c customCredentials)  GetRequestMetadata(ctx context.Context, uri ...string) (map[string]string, error){
	return map[string]string{
		"appid":"zifeng",
		"appkey":"123567",
	}, nil
}
func (c customCredentials)  RequireTransportSecurity() bool{
	return false
}
func main()  {
	//Interceptor:=func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error{
    //    start:=time.Now();
	//	md:=metadata.New(map[string]string{
	//		"appid":"zifeng",
	//		"appkey":"123567",
	//	})
	//
	//	ctx=metadata.NewOutgoingContext(context.Background(),md)
	//
    //    err:=invoker(ctx, method, req, reply, cc, opts ...)
    //    fmt.Printf("执行时间：%s\n",time.Since(start))
	//	return err;
	//}


	var opt []grpc.DialOption

	opt=append(opt,grpc.WithInsecure())
	opt=append(opt,grpc.WithPerRPCCredentials(customCredentials{}))

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
	})

	if err !=nil {
		panic(err)
	}
	fmt.Println(r.Message)

}