package initialize

import (
	"fmt"

	"github.com/hashicorp/consul/api"
	_ "github.com/mbobakov/grpc-consul-resolver" // It's important
	"go.uber.org/zap"
	"google.golang.org/grpc"

	"user_web/global"
	"user_web/proto"
)

func InitSrvConn() {
	consulInfo := global.ServerConfig.ConsulInfo
	userConn, err := grpc.Dial(
		fmt.Sprintf("consul://%s:%d/%s?wait=14s", consulInfo.Host, consulInfo.Port, global.ServerConfig.UserSrvInfo.Name),
		grpc.WithInsecure(),
		grpc.WithDefaultServiceConfig(`{"loadBalancingPolicy": "round_robin"}`),
	)
	if err != nil {
		zap.S().Fatal("[InitSrvConn] 连接 【用户服务失败】")
	}
	zap.S().Infof("[InitSrvConn] 连接 【用户服务成功】")
	userSrvClient := proto.NewUserClient(userConn)
	global.UserSrvClient = userSrvClient
}
func InitSrvConn1() {

	consulInfo := global.ServerConfig.ConsulInfo
	conn, err := grpc.Dial(
		fmt.Sprintf("consul://%s:%d/%s?wait=14s", consulInfo.Host, consulInfo.Port, global.ServerConfig.UserSrvInfo.Name),
		grpc.WithInsecure(),
		grpc.WithDefaultServiceConfig(`{"loadBalancingPolicy": "round_robin"}`),
	)
	if err != nil {
		zap.S().Fatal("[InitSrvConn] 连接 【用户服务失败】")
	}

	UserClient := proto.NewUserClient(conn)
	global.UserSrvClient = UserClient
}
func InitSrvConn2() {
	var err error

	//服务发现
	cfg := api.DefaultConfig()
	consulInfo := global.ServerConfig.ConsulInfo
	cfg.Address = fmt.Sprintf("%s:%d", consulInfo.Host, consulInfo.Port)

	client, err := api.NewClient(cfg)
	if err != nil {
		panic(err)
	}

	data, err := client.Agent().ServicesWithFilter(fmt.Sprintf(`Service == "%s"`, global.ServerConfig.UserSrvInfo.Name))
	if err != nil {
		panic(err)
	}
	userSrvHost := ""
	userSrvPort := 0
	for key, value := range data {
		userSrvHost = value.Address
		userSrvPort = value.Port
		zap.S().Infof("consul：%s", value)
		zap.S().Infof("consul 键值：%s", key)
		break
	}
	if userSrvHost == "" {
		zap.S().Fatal("用户服务[不可达]")
		return
	}

	conn, err := grpc.Dial(fmt.Sprintf("%s:%d", userSrvHost, userSrvPort), grpc.WithInsecure())
	zap.S().Infof("链接[用户服务] :%s:%d",
		userSrvHost, userSrvPort,
	)
	if err != nil {
		//zap.S().Fatal("[InitSrvConn] 连接 【用户服务失败】")
		zap.S().Errorw("用户服务[链接失败]",
			"msg", err.Error(),
		)
	}
	userSrvClient := proto.NewUserClient(conn)
	global.UserSrvClient = userSrvClient
}
