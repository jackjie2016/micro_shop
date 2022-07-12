package middlewares

import (
	"fmt"
	"github.com/gin-gonic/gin"
	opentracing "github.com/opentracing/opentracing-go"
	jaeger "github.com/uber/jaeger-client-go"
	jaegercfg "github.com/uber/jaeger-client-go/config"
	"goods_web/global"
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
				LocalAgentHostPort: fmt.Sprintf("%s:%d", global.ServerConfig.JaegerInfo.Host, global.ServerConfig.JaegerInfo.Port),
			},
			ServiceName: global.ServerConfig.JaegerInfo.Name,
		}
		Tracer, closer, err := cfg.NewTracer(jaegercfg.Logger(jaeger.StdLogger))
		if err != nil {
			panic(err)
		}

		opentracing.SetGlobalTracer(Tracer) //把Tracer 设置为全局的
		defer closer.Close()

		parentSpan := Tracer.StartSpan(c.Request.URL.Path)
		defer parentSpan.Finish()

		c.Set("tracer", Tracer)
		c.Set("parentSpan", parentSpan)
		c.Next()
	}
}
