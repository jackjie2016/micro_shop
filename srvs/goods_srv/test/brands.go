package main

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"mxshop_srvs/goods_srv/proto"
)
var goodsClient proto.GoodsClient
var conn *grpc.ClientConn
func init()  {
	var err error
	conn,err=grpc.Dial("192.168.31.134:50051",grpc.WithInsecure())
	if err !=nil {
		panic(err)
	}
	goodsClient=proto.NewGoodsClient(conn)
}
func TestGetBrandsList()  {

	r,err:=goodsClient.BrandList(context.Background(),&proto.BrandFilterRequest{Pages:1,PagePerNums: 10})

	if err !=nil {
		panic(err)
	}
	fmt.Println(r)
	for _,brand :=range r.Data{
		fmt.Println(brand.Name)


		if err!=nil{
			fmt.Println(err.Error())
		}

	}

}
func TestCreateBrands(){
	r,err:=goodsClient.CreateBrand(context.Background(),&proto.BrandRequest{Name:"cesss",Logo: ""})

	if err !=nil {
		panic(err)
	}
	fmt.Println(r)

}
func main()  {
	TestCreateBrands()
	//TestGetBrandsList()
	//TestGetUserByMobile()
	//TestCreateUser()
	conn.Close()
}
