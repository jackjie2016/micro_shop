package main

import (
	"flag"
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/apache/rocketmq-client-go/v2"
	"github.com/apache/rocketmq-client-go/v2/consumer"
	"github.com/apache/rocketmq-client-go/v2/primitive"
	"github.com/hashicorp/consul/api"
	uuid "github.com/satori/go.uuid"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"

	"inventory_srv/global"
	"inventory_srv/handler"
	"inventory_srv/initialize"
	"inventory_srv/proto"
	"inventory_srv/utils"
)

func main() {

	IP := flag.String("ip", "0.0.0.0", "ip地址")
	Port := flag.Int("port", 0, "端口号")

	if *Port == 0 {
		ip, _ := utils.ExternalIP()

		*IP = ip.String()
		*Port, _ = utils.GetFreePort()
		//*Port=50011
	}
	//1. 初始化logger
	initialize.InitLogger()

	//2、初始化配置
	initialize.InitConfig()

	//3、数据库链接
	initialize.InitDb()

	fmt.Println("MysqlInfo", global.ServerConfig.MysqlInfo)

	g := grpc.NewServer()
	proto.RegisterInventoryServer(g, &handler.InventoryServer{})
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
	fmt.Println("consul:", global.ServerConfig.ConsulInfo)
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
	registration.Tags = global.ServerConfig.Tags
	registration.Address = *IP
	registration.Check = check
	//1. 如何启动两个服务
	//2. 即使我能够通过终端启动两个服务，但是注册到consul中的时候也会被覆盖
	err = client.Agent().ServiceRegister(registration)
	if err != nil {
		panic(err)
	}

	//rocketMq 消费消息
	c, _ := rocketmq.NewPushConsumer(
		consumer.WithGroupName("mxshop-inventory"), //分布式用到，同一种实例
		consumer.WithNsResolver(primitive.NewPassthroughResolver([]string{fmt.Sprintf("%s:%d", global.ServerConfig.RocketMq.Host, global.ServerConfig.RocketMq.Port)})),
	)
	err = c.Subscribe("order_reback", consumer.MessageSelector{}, handler.AutoReback)
	if err != nil {
		zap.S().Errorf("生成 consumer 失败: %s", err.Error())
		return
	}
	// Note: start after subscribe
	err = c.Start()
	if err != nil {
		zap.S().Errorf("rocketmq 启动失败")
		return
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

	err = c.Shutdown()
	if err != nil {
		fmt.Printf("shutdown Consumer error: %s", err.Error())
	}

	if err = client.Agent().ServiceDeregister(serviceID); err != nil {
		zap.S().Info("注销失败")
	}
	zap.S().Info("注销成功")
}
