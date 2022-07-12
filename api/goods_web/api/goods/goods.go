package goods

import (
	"context"
	"net/http"
	"strconv"

	sentinel "github.com/alibaba/sentinel-golang/api"
	"github.com/alibaba/sentinel-golang/core/base"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"goods_web/api"
	"goods_web/forms"
	"goods_web/global"
	"goods_web/proto"
)

func List(ctx *gin.Context) {
	//request := &proto.GoodsFilterRequest{}
	PriceMin := ctx.DefaultQuery("pmin", "0")
	PriceMinInt, _ := strconv.Atoi(PriceMin)
	PriceMax := ctx.DefaultQuery("pmax", "0")
	PriceMaxInt, _ := strconv.Atoi(PriceMax)

	IsHot := ctx.DefaultQuery("ih", "")

	IH := false
	if IsHot != "" {
		IH = true
	}

	IsNew := ctx.DefaultQuery("in", "")
	IN := false
	if IsNew != "" {
		IN = true
	}

	IsTab := ctx.DefaultQuery("it", "")
	IT := false
	if IsTab != "" {
		IT = true
	}

	q := ctx.DefaultQuery("q", "")

	b := ctx.DefaultQuery("b", "")
	brandId, _ := strconv.Atoi(b)

	CategoryID := ctx.DefaultQuery("c", "")

	Cid, _ := strconv.Atoi(CategoryID)

	pn := ctx.DefaultQuery("pn", "")
	pnInt, _ := strconv.Atoi(pn)
	pSize := ctx.DefaultQuery("pnum", "10")
	pSizeInt, _ := strconv.Atoi(pSize)

	ee, bb := sentinel.Entry("goods-list", sentinel.WithTrafficType(base.Inbound))
	if bb != nil {
		zap.S().Error("限流了")
		ctx.JSON(http.StatusTooManyRequests, gin.H{
			"msg": "请求过于频繁，稍后重试",
		})
		return
	}

	r, err := global.GoodsSrvClient.GoodsList(
		context.WithValue(context.Background(), "ginContext", ctx), //为了jaeger调用链路的改造
		&proto.GoodsFilterRequest{
			PriceMin:    int32(PriceMinInt),
			PriceMax:    int32(PriceMaxInt),
			IsHot:       IH,
			IsNew:       IN,
			IsTab:       IT,
			TopCategory: int32(Cid),
			Pages:       int32(pnInt),
			PagePerNums: int32(pSizeInt),
			KeyWords:    q,
			Brand:       int32(brandId),
		},
	)
	if err != nil {
		zap.S().Errorw("【GetGoodslist】查询【商品列表】失败")
		api.HandleGrpcErrorToHttp(err, ctx)
		return
	}
	ee.Exit()

	goodsList := make([]interface{}, 0)
	for _, value := range r.Data {
		goodsList = append(goodsList, map[string]interface{}{
			"id":          value.Id,
			"name":        value.Name,
			"goods_brief": value.GoodsBrief,
			"desc":        value.GoodsDesc,
			"ship_free":   value.ShipFree,
			"images":      value.Images,
			"desc_images": value.DescImages,
			"front_image": value.GoodsFrontImage,
			"shop_price":  value.ShopPrice,
			"category": map[string]interface{}{
				"id":   value.Category.Id,
				"name": value.Category.Name,
			},
			"brand": map[string]interface{}{
				"id":   value.Brand.Id,
				"name": value.Brand.Name,
				"logo": value.Brand.Logo,
			},
			"is_hot":  value.IsHot,
			"is_new":  value.IsNew,
			"on_sale": value.OnSale,
		})

	}
	//rsMap:=map[string]interface{}{
	//	"total":r.Total,
	//	"data":goodsList,
	//}
	//ctx.JSON(http.StatusOK,rsMap)

	ctx.JSON(http.StatusOK, gin.H{
		"total": r.Total,
		"data":  goodsList,
	})

}

func New(ctx *gin.Context) {
	GoodsForm := forms.GoodsForm{}
	if err := ctx.ShouldBind(&GoodsForm); err != nil {
		api.HandleValitorError(ctx, err)
		return
	}

	//链接grpc
	//request := &proto.GoodsFilterRequest{}
	rsp, err := global.GoodsSrvClient.CreateGoods(
		context.WithValue(context.Background(), "ginContext", ctx),
		&proto.CreateGoodsInfo{
			Name:        GoodsForm.Name,
			GoodsSn:     GoodsForm.GoodsSn,
			MarketPrice: GoodsForm.MarketPrice,
			ShopPrice:   GoodsForm.ShopPrice,
			GoodsBrief:  GoodsForm.GoodsBrief,
			//GoodsDesc:GoodsForm.GoodsDesc,
			ShipFree:        *GoodsForm.ShipFree,
			Images:          GoodsForm.Images,
			DescImages:      GoodsForm.DescImages,
			GoodsFrontImage: GoodsForm.FrontImage,
			CategoryId:      GoodsForm.CategoryId,
			BrandId:         GoodsForm.Brand,
		})
	if err != nil {
		api.HandleGrpcErrorToHttp(err, ctx)
		return
	}
	//如何设置库存
	//TODO 商品的库存 - 分布式事务
	ctx.JSON(http.StatusOK, rsp)
}

func Detail(ctx *gin.Context) {
	id := ctx.Param("id")
	i, err := strconv.ParseInt(id, 10, 32)
	if err != nil {
		ctx.Status(http.StatusNotFound)
		return
	}

	r, err := global.GoodsSrvClient.GetGoodsDetail(context.WithValue(context.Background(), "ginContext", ctx), &proto.GoodInfoRequest{
		Id: int32(i),
	})
	if err != nil {
		api.HandleGrpcErrorToHttp(err, ctx)
		return
	}

	rsp := map[string]interface{}{
		"id":          r.Id,
		"name":        r.Name,
		"goods_brief": r.GoodsBrief,
		"desc":        r.GoodsDesc,
		"ship_free":   r.ShipFree,
		"images":      r.Images,
		"desc_images": r.DescImages,
		"front_image": r.GoodsFrontImage,
		"shop_price":  r.ShopPrice,
		"ctegory": map[string]interface{}{
			"id":   r.Category.Id,
			"name": r.Category.Name,
		},
		"brand": map[string]interface{}{
			"id":   r.Brand.Id,
			"name": r.Brand.Name,
			"logo": r.Brand.Logo,
		},
		"is_hot":  r.IsHot,
		"is_new":  r.IsNew,
		"on_sale": r.OnSale,
	}
	ctx.JSON(http.StatusOK, rsp)
}

func Delete(ctx *gin.Context) {
	id := ctx.Param("id")
	i, err := strconv.ParseInt(id, 10, 32)
	if err != nil {
		ctx.Status(http.StatusNotFound)
		return
	}

	_, err = global.GoodsSrvClient.DeleteGoods(context.WithValue(context.Background(), "ginContext", ctx), &proto.DeleteGoodsInfo{
		Id: int32(i),
	})
	if err != nil {
		api.HandleGrpcErrorToHttp(err, ctx)
		return
	}
	ctx.Status(http.StatusOK)
	return
}
func Update(ctx *gin.Context) {
	goodsForm := forms.GoodsForm{}
	if err := ctx.ShouldBindJSON(&goodsForm); err != nil {
		api.HandleValitorError(ctx, err)
		return
	}

	id := ctx.Param("id")
	i, err := strconv.ParseInt(id, 10, 32)
	if _, err = global.GoodsSrvClient.UpdateGoods(context.WithValue(context.Background(), "ginContext", ctx), &proto.CreateGoodsInfo{
		Id:              int32(i),
		Name:            goodsForm.Name,
		GoodsSn:         goodsForm.GoodsSn,
		Stocks:          goodsForm.Stocks,
		MarketPrice:     goodsForm.MarketPrice,
		ShopPrice:       goodsForm.ShopPrice,
		GoodsBrief:      goodsForm.GoodsBrief,
		ShipFree:        *goodsForm.ShipFree,
		Images:          goodsForm.Images,
		DescImages:      goodsForm.DescImages,
		GoodsFrontImage: goodsForm.FrontImage,
		CategoryId:      goodsForm.CategoryId,
		BrandId:         goodsForm.Brand,
	}); err != nil {
		api.HandleGrpcErrorToHttp(err, ctx)
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"msg": "更新成功",
	})

}
func Stocks(ctx *gin.Context) {
	id := ctx.Param("id")
	_, err := strconv.ParseInt(id, 10, 32)
	if err != nil {
		ctx.Status(http.StatusNotFound)
		return
	}

	//TODO 商品的库存
	return
}

func UpdateStatus(ctx *gin.Context) {
	goodsStatusForm := forms.GoodsStatusForm{}
	if err := ctx.ShouldBindJSON(&goodsStatusForm); err != nil {
		api.HandleValitorError(ctx, err)
		return
	}

	id := ctx.Param("id")
	i, err := strconv.ParseInt(id, 10, 32)
	if _, err = global.GoodsSrvClient.UpdateGoods(context.WithValue(context.Background(), "ginContext", ctx), &proto.CreateGoodsInfo{
		Id:     int32(i),
		IsHot:  *goodsStatusForm.IsHot,
		IsNew:  *goodsStatusForm.IsNew,
		OnSale: *goodsStatusForm.OnSale,
	}); err != nil {
		api.HandleGrpcErrorToHttp(err, ctx)
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"msg": "修改成功",
	})
}
