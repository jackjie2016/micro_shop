package main

import (
	"context"
	"fmt"
	"google.golang.org/protobuf/types/known/emptypb"
	"net"

	"google.golang.org/grpc"

	"OldPackageTest/grpc_proto_test/proto"
)
type Server struct{}

func (s *Server) SayHello(ctx context.Context, request *proto.HelloRequest) (*proto.HelloReply,
	error) {
	fmt.Println(fmt.Sprintf("性别：%d",request.G))
	fmt.Println(fmt.Sprintf("addtime：%d",request.AddTime.Seconds))
	for e,i:=range request.Mp{
		fmt.Println(i)
		fmt.Println(e)
	}
	return &proto.HelloReply{
		Message: "hello " + request.Name+" URl: " + request.Url,
	}, nil
}
func (s *Server) Ping(ctx context.Context, request *emptypb.Empty) (*proto.Pong, error)  {
	return nil,nil
}

func main() {

	g := grpc.NewServer()
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
