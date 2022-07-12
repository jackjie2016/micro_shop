package main

import (
	"context"
	"fmt"
	"log"

_ "github.com/mbobakov/grpc-consul-resolver" // It's important

"google.golang.org/grpc"
"OldPackageTest/grpclb_test/proto"
)

func main() {
	conn, err := grpc.Dial(
		"consul://192.168.31.134:8500/user-srv?wait=14s",
		grpc.WithInsecure(),
		grpc.WithDefaultServiceConfig(`{"loadBalancingPolicy": "round_robin"}`),
	)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	UserClient:=proto.NewUserClient(conn)
	rsp,err:=UserClient.GetUserList(context.Background(),&proto.PageInfo{
		Pn:1,
		PSize: 1,
	})
	for i,v:=range  rsp.Data{
		fmt.Println(i,v)
	}


}