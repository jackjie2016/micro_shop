package handler

import (
	"context"
	"crypto/sha512"
	"fmt"
	"github.com/anaskhan96/go-password-encoder"
	"github.com/golang/protobuf/ptypes/empty"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"strings"

	"gorm.io/gorm"

	"user_srv/global"
	"user_srv/model"
	"user_srv/proto"

	"time"
)

type UserServer struct {
	proto.UnimplementedUserServer
}

func ModelToResponse(user model.User) proto.UserInfoResponse {
	UserInfoRsp := proto.UserInfoResponse{
		Id:       uint32(user.ID),
		Password: user.Password,
		Mobile:   user.Mobile,
		NickName: user.NickName,
		Gender:   user.Gender,
		Role:     int32(user.Role),
	}
	if user.Birthday != nil {
		UserInfoRsp.BirthDay = uint64(user.Birthday.Unix())
	}

	return UserInfoRsp
}

func Paginate(page, pageSize int) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {

		if page == 0 {
			page = 1
		}

		switch {
		case pageSize > 100:
			pageSize = 100
		case pageSize <= 0:
			pageSize = 10
		}

		offset := (page - 1) * pageSize
		return db.Offset(offset).Limit(pageSize)
	}
}
func en() {

}
func (s *UserServer) GetUserList(ctx context.Context, req *proto.PageInfo) (*proto.UserListResponse, error) {
	var users []model.User
	result := global.DB.Find(&users)
	//查询没有错误
	if result.Error != nil {
		return nil, result.Error
	}
	rsp := &proto.UserListResponse{}
	rsp.Total = int32(result.RowsAffected)
	global.DB.Scopes(Paginate(int(req.Pn), int(req.PSize))).Find(&users)

	for _, user := range users {
		userInfoRsp := ModelToResponse(user)
		rsp.Data = append(rsp.Data, &userInfoRsp)
	}
	return rsp, nil
}
func (s *UserServer) GetUserByMobile(ctx context.Context, req *proto.MobileRequest) (*proto.UserInfoResponse, error) {
	var user model.User
	result := global.DB.Where("mobile", req.Mobile).First(&user)
	if result.RowsAffected == 0 {
		return nil, status.Errorf(codes.NotFound, "用户不存在")
	}
	if result.Error != nil {
		return nil, result.Error
	}
	UserInfoRsp := ModelToResponse(user)
	return &UserInfoRsp, nil
}
func (s *UserServer) GetUserById(ctx context.Context, req *proto.IdRequest) (*proto.UserInfoResponse, error) {
	var user model.User
	result := global.DB.First(&user, req.Id)
	if result.RowsAffected == 0 {
		return nil, status.Errorf(codes.NotFound, "用户不存在")
	}
	if result.Error != nil {
		return nil, result.Error
	}
	UserInfoRsp := ModelToResponse(user)
	return &UserInfoRsp, nil
}
func (s *UserServer) CreateUser(ctx context.Context, req *proto.CreateUserInfo) (*proto.UserInfoResponse, error) {
	var user model.User
	result := global.DB.Where(&model.User{Mobile: req.Mobile}).First(&user)
	if result.RowsAffected != 0 {
		return nil, status.Errorf(codes.AlreadyExists, "用户已存在")
	}
	user.Mobile = req.Mobile
	user.NickName = req.NickName

	options := &password.Options{6, 100, 30, sha512.New}
	salt, encodedPwd := password.Encode(req.Password, options)
	fmt.Printf("密码是：%s", req.Password)
	user.Password = fmt.Sprintf("$zifeng-sha512$%s$%s", salt, encodedPwd)

	result = global.DB.Create(&user) // 通过数据的指针来创建

	if result.Error != nil {
		return nil, status.Errorf(codes.Internal, "创建用户失败")
	}
	Rsq := ModelToResponse(user)
	return &Rsq, nil

}
func (s *UserServer) UpdateUser(ctx context.Context, req *proto.UpdateUserInfo) (*empty.Empty, error) {
	var user model.User
	result := global.DB.First(&user, req.Id)
	if result.RowsAffected == 0 {
		return nil, status.Errorf(codes.AlreadyExists, "用户不存在")
	}
	Birthday := time.Unix(int64(req.BirthDay), 0)
	user.NickName = req.NickName
	user.Birthday = &Birthday
	user.Gender = req.Gender
	result = global.DB.Save(&user)
	if result.Error != nil {
		return nil, status.Errorf(codes.Internal, result.Error.Error())
	}
	return &empty.Empty{}, nil
}
func (s *UserServer) CheckPassWord(ctx context.Context, req *proto.CheckInfo) (*proto.CheckResponse, error) {
	options := &password.Options{6, 100, 30, sha512.New}
	passwordinfo := strings.Split(req.EncryptedPassword, "$")
	check := password.Verify(req.Password, passwordinfo[2], passwordinfo[3], options)
	return &proto.CheckResponse{Success: check}, nil
}
