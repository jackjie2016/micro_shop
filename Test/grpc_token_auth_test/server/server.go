package main

import (
	"context"
	"fmt"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"net"

	"google.golang.org/grpc"

	"OldPackageTest/grpc_token_auth_test/proto"
)
type Server struct{}

func (s *Server) SayHello(ctx context.Context, request *proto.HelloRequest) (*proto.HelloReply,
	error) {



	return &proto.HelloReply{
		Message: "hello " + request.Name,
	}, nil
}





func main() {
	Interceptor:= func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error){
		fmt.Println("接收一个请求")

		md,ok:=metadata.FromIncomingContext(ctx)
		if !ok{
			return resp,status.Error(codes.Unauthenticated,"无token认证信息")
		}
		var(
			appid string
			appkey string
		)

		if val1,ok:=md["appid"];ok{
			appid=val1[0]
		}
		if val2,ok:=md["appkey"];ok{
			appkey=val2[0]
		}

		if appid!="zifeng" || appkey!="123567"{
			return resp,status.Error(codes.Unauthenticated,"token认证失败")
		}else {
			fmt.Println("认证通过")
		}

		res,err:= handler(ctx , req );
		fmt.Println("请求结束")
		return res,err;
	}






	opt:=grpc.UnaryInterceptor(Interceptor)
	g := grpc.NewServer(opt)
	proto.RegisterGreeterServer(g, &Server{})
	lis, err := net.Listen("tcp", "0.0.0.0:50053")
	if err != nil {
		panic("failed to listen:" + err.Error())
	}
	err = g.Serve(lis)
	if err != nil {
		panic("failed to start grpc:" + err.Error())
	}
}
