package handler

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"

	"goods_srv/global"
	"goods_srv/model"
	"goods_srv/proto"
)

func (g *GoodsServer) BrandList(ctx context.Context, req *proto.BrandFilterRequest) (*proto.BrandListResponse, error) {
	var brands []model.Brands
	result := global.DB.Find(&brands)
	//查询没有错误
	if result.Error != nil {
		return nil, result.Error
	}
	brandListResponse := &proto.BrandListResponse{}
	brandListResponse.Total = int32(result.RowsAffected)
	global.DB.Scopes(Paginate(int(req.Pages), int(req.PagePerNums))).Find(&brands)

	var total int64
	global.DB.Model(&model.Brands{}).Count(&total)
	brandListResponse.Total = int32(total)

	for _, brand := range brands {
		brandListResponse.Data = append(brandListResponse.Data, &proto.BrandInfoResponse{
			Id:   brand.ID,
			Name: brand.Name,
			Logo: brand.Logo,
		})
	}
	return brandListResponse, nil
}
func (g *GoodsServer) CreateBrand(ctx context.Context, req *proto.BrandRequest) (*proto.BrandInfoResponse, error) {

	type Result struct {
		Name string
		Age  int
	}
	//
	var result Result
	//global.DB.Table("users").Select("name", "age").Where("name = ?", "Antonio").Scan(&result)
	//
	//// Raw SQL
	global.DB.Raw("SELECT name, age FROM users WHERE name = ?", "Antonio").Scan(&result)
	//新建品牌
	if result := global.DB.Where("name=?", req.Name).First(&model.Brands{}); result.RowsAffected == 1 {
		return nil, status.Errorf(codes.InvalidArgument, "品牌已存在")
	}
	brands := &model.Brands{
		Name: req.Name,
		Logo: req.Logo,
	}

	global.DB.Save(brands)
	return &proto.BrandInfoResponse{Id: brands.ID, Name: brands.Name, Logo: brands.Logo}, nil
}
func (g *GoodsServer) DeleteBrand(ctx context.Context, req *proto.BrandRequest) (*emptypb.Empty, error) {
	if result := global.DB.Delete(&model.Brands{}, req.Id); result.RowsAffected == 0 {
		return nil, status.Errorf(codes.NotFound, "品牌不存在")
	}
	return &emptypb.Empty{}, nil

}
func (g *GoodsServer) UpdateBrand(ctx context.Context, req *proto.BrandRequest) (*emptypb.Empty, error) {
	brands := model.Brands{}
	if result := global.DB.Where("name=?", req.Id).First(&brands); result.RowsAffected == 0 {
		return nil, status.Errorf(codes.InvalidArgument, "品牌不存在")
	}
	if req.Name != "" {
		brands.Name = req.Name
	}
	if req.Logo != "" {
		brands.Logo = req.Logo
	}
	global.DB.Save(&brands)

	return &emptypb.Empty{}, nil

}
