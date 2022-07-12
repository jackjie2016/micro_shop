package main

import (
	"context"
	"fmt"
	"goods_srv/proto"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"
)

var goodsClient proto.GoodsClient
var conn *grpc.ClientConn

func init() {
	var err error
	conn, err = grpc.Dial("192.168.31.147:52702", grpc.WithInsecure())
	if err != nil {
		panic(err)
	}
	goodsClient = proto.NewGoodsClient(conn)
}
func TestGetBrandsList() {

	r, err := goodsClient.BrandList(context.Background(), &proto.BrandFilterRequest{Pages: 1, PagePerNums: 10})

	if err != nil {
		panic(err)
	}
	fmt.Println(r)
	for _, brand := range r.Data {
		fmt.Println(brand.Name)

		if err != nil {
			fmt.Println(err.Error())
		}

	}

}
func GetAllCategorysList() {
	r, err := goodsClient.GetAllCategorysList(context.Background(), &emptypb.Empty{})

	if err != nil {
		panic(err)
	}
	fmt.Println(r)

}
func GetSubCategory() {
	r, err := goodsClient.GetSubCategory(context.Background(), &proto.CategoryListRequest{
		Id:    130366,
		Level: 1,
	})

	if err != nil {
		panic(err)
	}
	fmt.Println(r)

}
func main() {
	GetAllCategorysList()
	//TestGetBrandsList()
	//TestGetUserByMobile()
	//TestCreateUser()
	//GetSubCategory()
	conn.Close()
}
