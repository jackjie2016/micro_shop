package main

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"inventory_srv/proto"
	"inventory_srv/utils"
	"sync"
)

var InventoryClient proto.InventoryClient
var conn *grpc.ClientConn

func TestInvDetail(GoodsId int32) {
	res, err := InventoryClient.InvDetail(context.Background(), &proto.GoodsInvInfo{
		GoodsId: GoodsId,
	})
	if err != nil {
		panic(err)
	}

	fmt.Println(res)
}

func TestSetInv(GoodsId, Num int32) {
	_, err := InventoryClient.SetInv(context.Background(), &proto.GoodsInvInfo{
		GoodsId: GoodsId,
		Num:     Num,
	})
	if err != nil {
		panic(err)
	}
	fmt.Println("设置成功")
}

func TestSell(wg *sync.WaitGroup) {
	defer wg.Done()
	_, err := InventoryClient.Sell(context.Background(), &proto.SellInfo{
		GoodsInfo: []*proto.GoodsInvInfo{
			{GoodsId: 1, Num: 1},
			//{GoodsId: 2, Num: 1},
		},
	})
	if err != nil {
		panic(err)
	}
	fmt.Println("设置成功")
}
func Init() {
	ip, _ := utils.ExternalIP()

	IP := ip.String()
	fmt.Printf("%s:50011", IP)
	var err error
	conn, err = grpc.Dial(fmt.Sprintf("%s:50011", IP), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		panic(err)
	}
	InventoryClient = proto.NewInventoryClient(conn)
}

func main() {
	Init()
	//TestCreateUser()
	//TestSetInv(2,2)

	//TestGetGoodsDetail()
	var wg sync.WaitGroup
	wg.Add(17)
	for i := 0; i < 17; i++ {
		go func() {
			TestSell(&wg)
		}()
	}
	wg.Wait()

	conn.Close()
}
