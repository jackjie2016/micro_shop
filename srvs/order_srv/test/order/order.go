package main

import (
	"context"
	"fmt"
	"go.uber.org/zap"
	"google.golang.org/grpc"

	"order_srv/proto"
	"order_srv/utils"
)

var OrderClient proto.OrderClient
var conn *grpc.ClientConn

func TestCreateOrder() {
	rsp, err := OrderClient.CreateOrder(context.Background(), &proto.OrderRequest{
		UserId:  85,
		Name:    "周正华",
		Mobile:  "15958615799",
		Address: "杭州市",
		Post:    "杭州市",
	})
	if err != nil {
		zap.S().Errorw("新建订单失败")

	}
	if err != nil {
		panic(err)
	}
	fmt.Println(rsp.Id)
}

func Init() {
	ip, _ := utils.ExternalIP()

	IP := ip.String()

	var err error
	conn, err = grpc.Dial(fmt.Sprintf("%s:63002", IP), grpc.WithInsecure())
	if err != nil {
		panic(err)
	}
	OrderClient = proto.NewOrderClient(conn)
}

func main() {
	Init()

	TestCreateOrder()
	conn.Close()
}
