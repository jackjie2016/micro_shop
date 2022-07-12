package main

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"user_srv/proto"
)

var userClient proto.UserClient
var conn *grpc.ClientConn

func init() {
	var err error
	conn, err = grpc.Dial("192.168.31.134:62946", grpc.WithInsecure())
	if err != nil {
		panic(err)
	} else {
		fmt.Println("链接成功")
	}
	userClient = proto.NewUserClient(conn)
}
func TestGetUserList() {

	r, err := userClient.GetUserList(context.Background(), &proto.PageInfo{Pn: 1, PSize: 10})

	if err != nil {
		panic(err)
	}
	for _, user := range r.Data {
		fmt.Println(user.NickName, user.Mobile, user.Password)

		rsp, err := userClient.CheckPassWord(context.Background(), &proto.CheckInfo{
			Password:          "zifeng234",
			EncryptedPassword: user.Password,
		})
		if err != nil {
			fmt.Println(err.Error())
		}
		fmt.Println(rsp.Success)
	}
	fmt.Println(r.Data)
}

func TestGetUserByMobile() {
	r, err := userClient.GetUserByMobile(context.Background(), &proto.MobileRequest{Mobile: "16888888886"})

	if err != nil {
		panic(err)
	}
	fmt.Println(r.NickName, r.Mobile, r.Id)
}
func TestLogin() {
	if rsp, err := userClient.GetUserByMobile(context.Background(), &proto.MobileRequest{
		Mobile: "15958615780",
	}); err != nil {
		panic(err.Error())
	} else {
		fmt.Println(rsp)
	}

}
func TestCreateUser() {

	for i := 0; i < 9; i++ {
		r, err := userClient.CreateUser(context.Background(), &proto.CreateUserInfo{
			Mobile:   fmt.Sprintf("1595861578%d", i),
			NickName: fmt.Sprintf("zifeng_%d", i),
			Password: "zifeng234",
		})
		if err != nil {
			panic(err)
		}
		fmt.Println(r.Id)
	}

}
func main() {
	//TestLogin()
	//TestCreateUser()
	TestGetUserList()
	//TestGetUserByMobile()
	//TestCreateUser()
	conn.Close()
}
