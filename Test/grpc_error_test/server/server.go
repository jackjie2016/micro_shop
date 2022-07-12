package main

import (
	"context"
	"net"
	"time"

	"OldPackageTest/grpc_error_test/proto"
	"google.golang.org/grpc"
)
type Server struct{}

func (s *Server) SayHello(ctx context.Context, request *proto.HelloRequest) (*proto.HelloReply,
	error) {
	//return nil,status.Errorf(codes.NotFound,"参数未找到 %s",request.Name)
	time.Sleep(time.Second*5)
	return &proto.HelloReply{
		Message: "hello " + request.Name,
	}, nil
}

func main() {
	g := grpc.NewServer()
	proto.RegisterGreeterServer(g, &Server{})
	lis, err := net.Listen("tcp", "0.0.0.0:8088")
	if err != nil {
		panic("failed to listen:" + err.Error())
	}
	err = g.Serve(lis)
	if err != nil {
		panic("failed to start grpc:" + err.Error())
	}
}
