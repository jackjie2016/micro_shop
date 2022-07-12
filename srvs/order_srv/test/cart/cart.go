package main

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"order_srv/proto"
	"order_srv/utils"
)

var OrderClient proto.OrderClient
var conn *grpc.ClientConn

func TestCreateCartItem(userId, nums, goodsId int32) {
	rsp, err := OrderClient.CreateCartItem(context.Background(), &proto.CartItemRequest{
		UserId:  userId,
		Nums:    nums,
		GoodsId: goodsId,
	})
	if err != nil {
		panic(err)
	}
	fmt.Println(rsp.Id)
}

func getCartItemList(userId int32) {
	rsp, err := OrderClient.CartItemList(context.Background(), &proto.UserInfo{Id: userId})

	if err != nil {
		panic(err)
	}

	for _, item := range rsp.Data {
		fmt.Println(item.Id, item.GoodsId, item.Nums)
	}

}

func Init() {
	ip, _ := utils.ExternalIP()

	IP := ip.String()
	fmt.Printf("%s:50011", IP)
	var err error
	conn, err = grpc.Dial(fmt.Sprintf("%s:50011", IP), grpc.WithInsecure())
	if err != nil {
		panic(err)
	}
	OrderClient = proto.NewOrderClient(conn)
}

func main() {
	Init()
	TestCreateCartItem(1, 2, 3)
	//getCartItemList(1)

	conn.Close()
}
