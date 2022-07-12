package initialize

import (
	"fmt"

	_ "github.com/mbobakov/grpc-consul-resolver" // 服务发现自动内嵌
	"go.uber.org/zap"
	"google.golang.org/grpc"

	"orders_web/config"
	"orders_web/global"
	"orders_web/proto"
)

var consulInfo config.ConsulConfig

func OrderSrv() {
	fmt.Println("consulInfo 信息：", consulInfo)
	userConn, err := grpc.Dial(
		fmt.Sprintf("consul://%s:%d/%s?wait=14s", consulInfo.Host, consulInfo.Port, global.ServerConfig.OrderSrvInfo.Name),
		grpc.WithInsecure(),
		grpc.WithDefaultServiceConfig(`{"loadBalancingPolicy": "round_robin"}`),
	)
	if err != nil {
		zap.S().Fatal("[InitSrvConn] 连接 【用户服务失败】")
	}

	orderSrvClient := proto.NewOrderClient(userConn)
	global.OrderSrvClient = orderSrvClient
}
func GoodsSrv() {
	zap.S().Infof("consul://%s:%d/%s?wait=14s", consulInfo.Host, consulInfo.Port, global.ServerConfig.GoodsSrvInfo.Name)
	userConn, err := grpc.Dial(
		fmt.Sprintf("consul://%s:%d/%s?wait=14s", consulInfo.Host, consulInfo.Port, global.ServerConfig.GoodsSrvInfo.Name),
		grpc.WithInsecure(),
		grpc.WithDefaultServiceConfig(`{"loadBalancingPolicy": "round_robin"}`),
	)
	if err != nil {
		zap.S().Fatal("[InitSrvConn] 连接 【用户服务失败】")
	}

	goodsSrvClient := proto.NewGoodsClient(userConn)
	global.GoodsSrvClient = goodsSrvClient
}
func InventorySrv() {
	zap.S().Infof("consul://%s:%d/%s?wait=14s", consulInfo.Host, consulInfo.Port, global.ServerConfig.InventorySrvInfo.Name)
	userConn, err := grpc.Dial(
		fmt.Sprintf("consul://%s:%d/%s?wait=14s", consulInfo.Host, consulInfo.Port, global.ServerConfig.InventorySrvInfo.Name),
		grpc.WithInsecure(),
		grpc.WithDefaultServiceConfig(`{"loadBalancingPolicy": "round_robin"}`),
	)
	if err != nil {
		zap.S().Fatal("[InitSrvConn] 连接 【用户服务失败】")
	}

	InventorySrvClient := proto.NewInventoryClient(userConn)
	global.InventorySrvClient = InventorySrvClient
}
func InitSrvConn() {
	consulInfo = global.ServerConfig.ConsulInfo
	zap.S().Info(global.ServerConfig)
	zap.S().Infof("consul://%s:%d/%s?wait=14s", consulInfo.Host, consulInfo.Port, global.ServerConfig.GoodsSrvInfo.Name)
	GoodsSrv()
	OrderSrv()
	InventorySrv()
}
