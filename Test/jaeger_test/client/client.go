package main

import (

	"context"
	"fmt"

	"github.com/opentracing/opentracing-go"
	"github.com/uber/jaeger-client-go"
	jaegercfg "github.com/uber/jaeger-client-go/config"
	"google.golang.org/grpc"

	"OldPackageTest/jaeger_test/otgrpc"
	"OldPackageTest/jaeger_test/proto"
)

func main()  {
	cfg := jaegercfg.Configuration{
		Sampler: &jaegercfg.SamplerConfig{
			Type:  jaeger.SamplerTypeConst,
			Param: 1,
		},
		Reporter: &jaegercfg.ReporterConfig{
			LogSpans: true,//打印日志
			LocalAgentHostPort:"192.168.1.105:6831",
		},
		ServiceName:"mxshop",
	}
	Tracer, closer, err:= cfg.NewTracer(jaegercfg.Logger(jaeger.StdLogger))
	if err!=nil{
		panic(err)
	}

	opentracing.SetGlobalTracer(Tracer) //把Tracer 设置为全局的
	defer closer.Close()

	conn,err:=grpc.Dial("127.0.0.1:8088",
		grpc.WithInsecure(),
		grpc.WithUnaryInterceptor(otgrpc.OpenTracingClientInterceptor(opentracing.GlobalTracer()),
		))
	if err !=nil {
		panic(err)
	}
	defer conn.Close()

	c:=proto.NewGreeterClient(conn)
	r,err:=c.SayHello(context.Background(),&proto.HelloRequest{Name:"bobby"})

	if err !=nil {
		panic(err)
	}
	fmt.Println(r.Message)

}