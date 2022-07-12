package main

import (
	"context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"net"

	"google.golang.org/grpc"

	"OldPackageTest/grpc_validate_test/proto"
)


type Server struct{}

func (s *Server) SayHello(ctx context.Context, request *proto.Person) (*proto.Person,
	error){
	return &proto.Person{
		Id: 32,
	}, nil
}

type Validator interface {
	Validate(all bool) error
}

func main(){
	//p:=new(proto.Person)
	//err:=p.Validate(true)
	//if err!=nil{
	//	panic(err.Error())
	//}
	var interceptor grpc.UnaryServerInterceptor
	interceptor = func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
		// 继续处理请求
		if r, ok := req.(Validator); ok {

			if err := r.Validate(true); err != nil {
				return nil, status.Error(codes.InvalidArgument, err.Error())
			}
		}

		return handler(ctx, req)
	}
	var opts []grpc.ServerOption
	opts = append(opts, grpc.UnaryInterceptor(interceptor))

	g := grpc.NewServer(opts...)
	proto.RegisterGreeterServer(g, &Server{})
	lis, err := net.Listen("tcp", "0.0.0.0:50051")
	if err != nil{
		panic("failed to listen:"+err.Error())
	}
	err = g.Serve(lis)
	if err != nil{
		panic("failed to start grpc:"+err.Error())
	}
}

