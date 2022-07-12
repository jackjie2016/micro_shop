package main

import (
	opentracing "github.com/opentracing/opentracing-go"
	"github.com/uber/jaeger-client-go"
	jaegercfg "github.com/uber/jaeger-client-go/config"
	"time"
)
func main()  {
	cfg := jaegercfg.Configuration{
		Sampler: &jaegercfg.SamplerConfig{
			Type:  jaeger.SamplerTypeConst,
			Param: 1,
		},
		Reporter: &jaegercfg.ReporterConfig{
			LogSpans: true,//打印日志
			LocalAgentHostPort:"127.0.0.1:6831",
		},
		ServiceName:"mxshop",
	}
	Tracer, closer, err:= cfg.NewTracer(jaegercfg.Logger(jaeger.StdLogger))
	if err!=nil{
		panic(err)
	}

	opentracing.SetGlobalTracer(Tracer) //把Tracer 设置为全局的
	defer closer.Close()



	parentSpan:=Tracer.StartSpan("main")
	span1:=Tracer.StartSpan("funA",opentracing.ChildOf(parentSpan.Context()))
	time.Sleep(time.Second)
	span1.Finish()

	globalspan:=opentracing.StartSpan("global funA",opentracing.ChildOf(parentSpan.Context()))
	time.Sleep(time.Second)
	globalspan.Finish()

	span2:=Tracer.StartSpan("funB",opentracing.ChildOf(span1.Context()))
	time.Sleep(time.Second)
	span2.Finish()

	parentSpan.Finish()
}