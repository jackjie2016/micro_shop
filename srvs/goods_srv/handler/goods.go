package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"goods_srv/global"
	"goods_srv/model"
	"goods_srv/proto"

	"github.com/olivere/elastic/v7"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

type GoodsServer struct {
	proto.UnimplementedGoodsServer //快速注册接口

}

func ModelToResponse(goods model.Goods) proto.GoodsInfoResponse {
	return proto.GoodsInfoResponse{
		Id:              goods.ID,
		CategoryId:      goods.CategoryID,
		Name:            goods.Name,
		GoodsSn:         goods.GoodsSn,
		ClickNum:        goods.ClickNum,
		SoldNum:         goods.SoldNum,
		FavNum:          goods.FavNum,
		MarketPrice:     goods.MarketPrice,
		ShopPrice:       goods.ShopPrice,
		GoodsBrief:      goods.GoodsBrief,
		ShipFree:        goods.ShipFree,
		GoodsFrontImage: goods.GoodsFrontImage,
		IsNew:           goods.IsNew,
		IsHot:           goods.IsHot,
		OnSale:          goods.OnSale,
		DescImages:      goods.DescImages,
		Images:          goods.Images,
		Category: &proto.CategoryBriefInfoResponse{
			Id:   goods.Category.ID,
			Name: goods.Category.Name,
		},
		Brand: &proto.BrandInfoResponse{
			Id:   goods.Brands.ID,
			Name: goods.Brands.Name,
			Logo: goods.Brands.Logo,
		},
	}
}

//商品接口
func (g *GoodsServer) GoodsList(ctx context.Context, req *proto.GoodsFilterRequest) (*proto.GoodsListResponse, error) {
	//关键词搜索、查询新品、查询热门商品、通过价格区间筛选， 通过商品分类筛选
	goodsListResponse := &proto.GoodsListResponse{}
	localDB := global.DB.Model(&model.Goods{})

	q := elastic.NewBoolQuery()
	//q = q.Must(NewTermQuery("tag", "wow"))
	//q = q.Filter(NewTermQuery("account", "1"))

	var goods []model.Goods

	if req.KeyWords != "" {
		//localDB=localDB.Where("name Like ?","%"+req.KeyWords+"%")
		//多字段查询
		q = q.Must(elastic.NewMultiMatchQuery(req.KeyWords, "name", "goods_brief"))
	}
	if req.IsHot {
		//localDB=localDB.Where(model.Goods{IsHot: true})
		q = q.Filter(elastic.NewTermQuery("is_hot", req.IsHot))
	}
	if req.IsNew {
		//localDB=localDB.Where(model.Goods{IsNew: true})
		q = q.Filter(elastic.NewTermQuery("is_new", req.IsNew))
	}

	if req.PriceMin > 0 {
		//localDB = localDB.Where("shop_price>=?", req.PriceMin)
		q = q.Filter(elastic.NewRangeQuery("shop_price").Gte(req.PriceMin))
	}

	if req.PriceMax > 0 {
		//localDB=localDB.Where("shop_price <=?",req.PriceMax)
		q = q.Filter(elastic.NewRangeQuery("shop_price").Lte(req.PriceMax))
	}

	if req.TopCategory > 0 {
		//通过category去查询商品
		var subQuery string

		var category model.Category
		if result := global.DB.First(&category, req.TopCategory); result.RowsAffected == 0 {
			return nil, status.Errorf(codes.NotFound, "商品分类不存在")
		}

		if category.Level == 1 {
			subQuery = fmt.Sprintf("select id from category where parent_category_id in (select id from category WHERE parent_category_id=%d)", req.TopCategory)
		} else if category.Level == 2 {
			subQuery = fmt.Sprintf("select id from category WHERE parent_category_id=%d", req.TopCategory)
		} else if category.Level == 3 {
			subQuery = fmt.Sprintf("select id from category WHERE id=%d", req.TopCategory)
		}
		categoryIDs := make([]interface{}, 0)

		type CateResult struct {
			ID int32
		}
		var CateResults []CateResult
		//global.DB.Where(fmt.Sprintf("category_id in (%s)", subQuery)).Find(CateResults)
		global.DB.Model(model.Category{}).Raw(subQuery).Scan(&CateResults)
		for _, c := range CateResults {
			categoryIDs = append(categoryIDs, c.ID)
		}
		//es 中根据分类id过滤
		q = q.Filter(elastic.NewTermsQuery("category_id", categoryIDs...))
	}
	if req.Pages == 0 {
		req.Pages = 1
	}

	switch {
	case req.PagePerNums > 100:
		req.PagePerNums = 100
	case req.PagePerNums <= 0:
		req.PagePerNums = 10
	}
	goodsIds := make([]int32, 0)
	esRes, err := global.Esclient.Search().Index(model.EsGoods{}.GetindexName()).Query(q).From(int(req.Pages)).Size(int(req.PagePerNums)).Do(context.Background())
	if err != nil {
		return nil, err
	}
	//res,err:=global.Esclient.Search().Index(model.EsGoods{}.GetindexName()).Query(q).From(int(req.Pages)).Size(int(req.PagePerNums)).Do(context.Background())
	//if err!=nil{
	//	return nil, err
	//}
	goodsListResponse.Total = int32(esRes.Hits.TotalHits.Value)

	if goodsListResponse.Total == 0 {
		return nil, status.Errorf(codes.ResourceExhausted, "没有查询到商品")
	}

	for _, value := range esRes.Hits.Hits {
		goods := model.EsGoods{}
		_ = json.Unmarshal(value.Source, &goods)
		goodsIds = append(goodsIds, goods.ID)
	}

	//var count int64
	//localDB.Count(&count)
	//goodsListResponse.Total = int32(count)

	result := localDB.Preload("Category").Preload("Brands").Find(&goods, goodsIds)
	if result.Error != nil {
		return nil, result.Error
	}

	for _, good := range goods {
		goodsInfoResponse := ModelToResponse(good)
		goodsListResponse.Data = append(goodsListResponse.Data, &goodsInfoResponse)
	}

	return goodsListResponse, nil

}

//现在用户提交订单有多个商品，你得批量查询商品的信息吧
func (s *GoodsServer) BatchGetGoods(ctx context.Context, req *proto.BatchGoodsIdInfo) (*proto.GoodsListResponse, error) {
	goodsListResponse := &proto.GoodsListResponse{}
	var goods []model.Goods

	//调用where并不会真正执行sql 只是用来生成sql的 当调用find， first才会去执行sql，
	result := global.DB.Where(req.Id).Find(&goods)
	for _, good := range goods {
		goodsInfoResponse := ModelToResponse(good)
		goodsListResponse.Data = append(goodsListResponse.Data, &goodsInfoResponse)
	}
	goodsListResponse.Total = int32(result.RowsAffected)
	return goodsListResponse, nil
}
func (s *GoodsServer) GetGoodsDetail(ctx context.Context, req *proto.GoodInfoRequest) (*proto.GoodsInfoResponse, error) {
	var goods model.Goods

	if result := global.DB.First(&goods, req.Id); result.RowsAffected == 0 {
		return nil, status.Errorf(codes.NotFound, "商品不存在")
	}
	goodsInfoResponse := ModelToResponse(goods)
	return &goodsInfoResponse, nil
}

func (s *GoodsServer) CreateGoods(ctx context.Context, req *proto.CreateGoodsInfo) (*proto.GoodsInfoResponse, error) {
	var category model.Category
	if result := global.DB.First(&category, req.CategoryId); result.RowsAffected == 0 {
		return nil, status.Errorf(codes.InvalidArgument, "商品分类不存在")
	}

	var brand model.Brands
	if result := global.DB.First(&brand, req.BrandId); result.RowsAffected == 0 {
		return nil, status.Errorf(codes.InvalidArgument, "品牌不存在")
	}
	//这里没有看到图片文件是如何上传， 在微服务中 普通的文件上传已经不再使用
	goods := model.Goods{
		Brands:          brand,
		BrandsID:        brand.ID,
		Category:        category,
		CategoryID:      category.ID,
		Name:            req.Name,
		GoodsSn:         req.GoodsSn,
		MarketPrice:     req.MarketPrice,
		ShopPrice:       req.ShopPrice,
		GoodsBrief:      req.GoodsBrief,
		ShipFree:        req.ShipFree,
		Images:          req.Images,
		DescImages:      req.DescImages,
		GoodsFrontImage: req.GoodsFrontImage,
		IsNew:           req.IsNew,
		IsHot:           req.IsHot,
		OnSale:          req.OnSale,
	}

	//srv之间互相调用了
	tx := global.DB
	tx.Begin()

	result := tx.Save(&goods)
	if result.Error != nil {
		tx.Rollback()
		return nil, result.Error
	}
	tx.Commit()

	return &proto.GoodsInfoResponse{
		Id: goods.ID,
	}, nil
}

func (s *GoodsServer) DeleteGoods(ctx context.Context, req *proto.DeleteGoodsInfo) (*emptypb.Empty, error) {
	tx := global.DB
	tx.Begin()

	result := tx.Delete(&model.Goods{}, req.Id)
	if result.Error != nil {
		tx.Rollback()
		return nil, result.Error
	}
	tx.Commit()

	if result := global.DB.Delete(&model.Goods{BaseModel: model.BaseModel{ID: req.Id}}, req.Id); result.RowsAffected == 0 {

		return nil, status.Errorf(codes.NotFound, "商品不存在")
	}
	return &emptypb.Empty{}, nil
}

func (s *GoodsServer) UpdateGoods(ctx context.Context, req *proto.CreateGoodsInfo) (*emptypb.Empty, error) {
	var goods model.Goods

	if result := global.DB.First(&goods, req.Id); result.RowsAffected == 0 {
		return nil, status.Errorf(codes.NotFound, "商品不存在")
	}

	var category model.Category
	if result := global.DB.First(&category, req.CategoryId); result.RowsAffected == 0 {
		return nil, status.Errorf(codes.InvalidArgument, "商品分类不存在")
	}

	var brand model.Brands
	if result := global.DB.First(&brand, req.BrandId); result.RowsAffected == 0 {
		return nil, status.Errorf(codes.InvalidArgument, "品牌不存在")
	}

	goods.Brands = brand
	goods.BrandsID = brand.ID
	goods.Category = category
	goods.CategoryID = category.ID
	goods.Name = req.Name
	goods.GoodsSn = req.GoodsSn
	goods.MarketPrice = req.MarketPrice
	goods.ShopPrice = req.ShopPrice
	goods.GoodsBrief = req.GoodsBrief
	goods.ShipFree = req.ShipFree
	goods.Images = req.Images
	goods.DescImages = req.DescImages
	goods.GoodsFrontImage = req.GoodsFrontImage
	goods.IsNew = req.IsNew
	goods.IsHot = req.IsHot
	goods.OnSale = req.OnSale

	tx := global.DB
	tx.Begin()

	result := tx.Save(&goods)
	if result.Error != nil {
		tx.Rollback()
		return nil, result.Error
	}
	tx.Commit()
	return &emptypb.Empty{}, nil
}
