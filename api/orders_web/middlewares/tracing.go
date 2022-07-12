package middlewares

import (
	"github.com/gin-gonic/gin"
	opentracing "github.com/opentracing/opentracing-go"
	"github.com/uber/jaeger-client-go"
	jaegercfg "github.com/uber/jaeger-client-go/config"

)
func Tracing() gin.HandlerFunc {
	return func(c *gin.Context) {
		cfg := jaegercfg.Configuration{
			Sampler: &jaegercfg.SamplerConfig{
				Type:  jaeger.SamplerTypeConst,
				Param: 1,
			},
			Reporter: &jaegercfg.ReporterConfig{
				LogSpans:           true, //打印日志
				LocalAgentHostPort: "192.168.1.105:6831",
			},
			ServiceName: "mxshop",
		}
		Tracer, closer, err := cfg.NewTracer(jaegercfg.Logger(jaeger.StdLogger))
		if err != nil {
			panic(err)
		}

		opentracing.SetGlobalTracer(Tracer) //把Tracer 设置为全局的
		defer closer.Close()
		//c.Set("claims", )
		//c.Set("userId", claims.ID)
		c.Next()
	}
}