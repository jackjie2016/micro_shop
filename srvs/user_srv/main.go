package main

import (
	"flag"
	"fmt"
	"github.com/hashicorp/consul/api"
	"go.uber.org/zap"

	"net"
	"os"
	"os/signal"
	"syscall"

	uuid "github.com/satori/go.uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"

	"user_srv/global"
	"user_srv/handler"
	"user_srv/initialize"
	"user_srv/proto"
	"user_srv/utils"
)

func main() {

	IP := flag.String("ip", "0.0.0.0", "ip地址")
	Port := flag.Int("port", 0, "端口号")

	if *Port == 0 {
		ip, _ := utils.ExternalIP()

		*IP = ip.String()
		*Port, _ = utils.GetFreePort()

	}
	//1. 初始化logger
	initialize.InitLogger()

	//2、初始化配置
	initialize.InitConfig()

	//3、数据库链接
	initialize.InitDb()

	fmt.Println(global.ServerConfig.MysqlInfo)

	g := grpc.NewServer()
	proto.RegisterUserServer(g, &handler.UserServer{})
	//proto.RegisterUserServer(g, &proto.UnimplementedUserServer{})

	lis, err := net.Listen("tcp", fmt.Sprintf("%s:%d", *IP, *Port))
	if err != nil {
		panic("failed to listen:" + err.Error())
	}

	//注册健康检查
	grpc_health_v1.RegisterHealthServer(g, health.NewServer())

	//注册到consul
	cfg := api.DefaultConfig()
	cfg.Address = fmt.Sprintf("%s:%d", global.ServerConfig.ConsulInfo.Host,
		global.ServerConfig.ConsulInfo.Port)
	fmt.Println(cfg)
	client, err := api.NewClient(cfg)
	if err != nil {
		panic(err)
	}
	//生成对应的检查对象
	check := &api.AgentServiceCheck{
		GRPC:                           fmt.Sprintf("%s:%d", *IP, *Port),
		Timeout:                        "5s",
		Interval:                       "5s",
		DeregisterCriticalServiceAfter: "15s",
	}

	//生成注册对象
	registration := new(api.AgentServiceRegistration)
	registration.Name = global.ServerConfig.Name
	serviceID := fmt.Sprintf("%s", uuid.NewV4())
	registration.ID = serviceID
	registration.Port = *Port
	registration.Tags = []string{"imooc", "bobby", "user", "srv"}
	registration.Address = *IP
	registration.Check = check
	//1. 如何启动两个服务
	//2. 即使我能够通过终端启动两个服务，但是注册到consul中的时候也会被覆盖
	err = client.Agent().ServiceRegister(registration)
	if err != nil {
		panic(err)
	}

	go func() {
		err = g.Serve(lis)
		if err != nil {
			panic("failed to start grpc:" + err.Error())
		}
	}()

	//接收终止信号
	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	if err = client.Agent().ServiceDeregister(serviceID); err != nil {
		zap.S().Info("注销失败")
	}
	zap.S().Info("注销成功")
}
