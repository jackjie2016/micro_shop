package initialize

import (
	"fmt"
	opentracing "github.com/opentracing/opentracing-go"
	"goods_web/otgrpc"

	_ "github.com/mbobakov/grpc-consul-resolver" // 服务发现自动内嵌
	"go.uber.org/zap"
	"google.golang.org/grpc"

	"goods_web/global"
	"goods_web/proto"
)

func InitSrvConn() {
	consulInfo := global.ServerConfig.ConsulInfo
	zap.S().Infof("consul://%s:%d/%s?wait=14s", consulInfo.Host, consulInfo.Port, global.ServerConfig.GoodsSrvInfo.Name)
	userConn, err := grpc.Dial(
		fmt.Sprintf("consul://%s:%d/%s?wait=14s", consulInfo.Host, consulInfo.Port, global.ServerConfig.GoodsSrvInfo.Name),
		grpc.WithInsecure(),
		grpc.WithDefaultServiceConfig(`{"loadBalancingPolicy": "round_robin"}`),
		grpc.WithUnaryInterceptor(otgrpc.OpenTracingClientInterceptor(opentracing.GlobalTracer())),
	)
	if err != nil {
		zap.S().Fatal("[InitSrvConn] 连接 【用户服务失败】")
	}

	goodsSrvClient := proto.NewGoodsClient(userConn)
	global.GoodsSrvClient = goodsSrvClient
}
