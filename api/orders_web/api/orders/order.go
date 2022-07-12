package orders

import (
	"context"
	alipay "github.com/smartwalle/alipay/v3"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"orders_web/api"
	"orders_web/forms"
	"orders_web/global"
	"orders_web/models"
	"orders_web/proto"
)

func List(ctx *gin.Context) {

	//订单的列表
	userId, _ := ctx.Get("userId")
	claims, _ := ctx.Get("claims")

	request := proto.OrderFilterRequest{}

	//如果是管理员用户则返回所有的订单
	model := claims.(*models.CustomClaims)
	if model.AuthorityID == 1 {
		request.UserId = int32(userId.(uint))
	}

	pages := ctx.DefaultQuery("p", "0")
	pagesInt, _ := strconv.Atoi(pages)
	request.Pages = int32(pagesInt)

	perNums := ctx.DefaultQuery("pnum", "0")
	perNumsInt, _ := strconv.Atoi(perNums)
	request.PagePerNums = int32(perNumsInt)

	request.Pages = int32(pagesInt)
	request.PagePerNums = int32(perNumsInt)

	rsp, err := global.OrderSrvClient.OrderList(context.Background(), &request)
	if err != nil {
		zap.S().Errorw("获取订单列表失败")
		api.HandleGrpcErrorToHttp(err, ctx)
		return
	}

	reMap := gin.H{
		"total": rsp.Total,
	}
	orderList := make([]interface{}, 0)

	for _, item := range rsp.Data {
		tmpMap := map[string]interface{}{}

		tmpMap["id"] = item.Id
		tmpMap["status"] = item.Status
		tmpMap["pay_type"] = item.PayType
		tmpMap["user"] = item.UserId
		tmpMap["post"] = item.Post
		tmpMap["total"] = item.Total
		tmpMap["address"] = item.Address
		tmpMap["name"] = item.Name
		tmpMap["mobile"] = item.Mobile
		tmpMap["order_sn"] = item.OrderSn
		tmpMap["id"] = item.Id
		tmpMap["add_time"] = item.AddTime

		orderList = append(orderList, tmpMap)
	}
	reMap["data"] = orderList
	ctx.JSON(http.StatusOK, reMap)
}

func New(ctx *gin.Context) {

	orderForm := forms.CreateOrderForm{}

	if err := ctx.ShouldBindJSON(&orderForm); err != nil {
		api.HandleValitorError(ctx, err)
		return
	}

	userId, _ := ctx.Get("userId")
	rsp, err := global.OrderSrvClient.CreateOrder(context.WithValue(context.Background(), "ginContext", ctx), &proto.OrderRequest{
		UserId:  int32(userId.(uint)),
		Name:    orderForm.Name,
		Mobile:  orderForm.Mobile,
		Address: orderForm.Address,
		Post:    orderForm.Post,
	})
	if err != nil {
		zap.S().Errorf("失败原因：%v", err.Error())
		api.HandleGrpcErrorToHttp(err, ctx)
		return
	}

	var appId = global.ServerConfig.AlipayInfo.AppId
	var privateKey = global.ServerConfig.AlipayInfo.PrivateKey // 必须，上一步中使用 RSA签名验签工具 生成的私钥
	//var aliPublicKey="MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAkBui320sHC572qEBT8oCmaswX32x3YO2Di07XlkG71B+tdSJ7R6GzjVacrmGD2H9METl5b1Et0hbwrcBwnt/ER5Pyk1eprIy1UoTzcsJFH7sGfPNEET0jl0f3axwvFFNsIf9vi4mGjuR+YYZPcjYKpQ7ZhGT2I9mkHfKbkViQshe5LtlW9Qi4ag+vUl8vV0vkD5/OO1qP0c8zbwUmPpXO572TbuMVz4ZR+QSXsOUK/XTQbto3I2aBj12QnIwStOlH1WHxqnuxE3JQSYhS1p0jAb8dXTB2b2gdBwEDv9jahSNv1i1xRFPrnvH4tjkrn1frnJ2pEIICD1+dYg0624ztwIDAQAB"

	client, err := alipay.New(appId, privateKey, false)
	if err != nil {
		zap.S().Errorw("实例化支付宝失败失败")
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg": "内部错误",
		})
		return
	}
	//client.LoadAliPayPublicKey(aliPublicKey)
	//
	err = client.LoadAppPublicCertFromFile("alipay/appCertPublicKey.crt") // 加载应用公钥证书
	// 将 key 的验证调整到初始化阶段
	if err != nil {
		zap.S().Errorw("加载应用公钥证书失败")
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg": "内部错误",
		})
		return
	}
	err = client.LoadAliPayRootCertFromFile("alipay/alipayRootCert.crt") // 加载支付宝根证书
	if err != nil {
		zap.S().Errorw("加载支付宝根证书 失败")
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg": "内部错误",
		})
		return
	}
	err = client.LoadAliPayPublicCertFromFile("alipay/alipayCertPublicKey_RSA2.crt") // 加载支付宝公钥证书
	if err != nil {
		zap.S().Errorw("加载支付宝公钥证书 失败")
		return
	}
	// 将 key 的验证调整到初始化阶段
	if err != nil {
		zap.S().Errorw(err.Error())
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg": "内部错误",
		})
		return
	}

	var p = alipay.TradePagePay{}
	p.NotifyURL = global.ServerConfig.AlipayInfo.NotifyURL
	p.ReturnURL = global.ServerConfig.AlipayInfo.ReturnURL
	p.Subject = "慕学生鲜订单-" + rsp.OrderSn
	p.OutTradeNo = rsp.OrderSn
	p.TotalAmount = strconv.FormatFloat(float64(rsp.Total), 'f', 2, 64)
	p.ProductCode = "FAST_INSTANT_TRADE_PAY"

	url, err := client.TradePagePay(p)
	if err != nil {
		zap.S().Errorw("生成支付url失败")
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"id":         rsp.Id,
		"alipay_url": url.String(),
	})
}
func Detail(ctx *gin.Context) {
	id := ctx.Param("id")
	userId, _ := ctx.Get("userId")
	i, err := strconv.Atoi(id)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"msg": "url格式出错",
		})
		return
	}

	//如果是管理员用户则返回所有的订单
	request := proto.OrderRequest{
		Id: int32(i),
	}
	claims, _ := ctx.Get("claims")
	model := claims.(*models.CustomClaims)
	if model.AuthorityID == 1 {
		request.UserId = int32(userId.(uint))
	}

	rsp, err := global.OrderSrvClient.OrderDetail(context.Background(), &request)
	if err != nil {
		zap.S().Errorw("获取订单详情失败")
		api.HandleGrpcErrorToHttp(err, ctx)
		return
	}
	reMap := gin.H{}
	reMap["id"] = rsp.OrderInfo.Id
	reMap["status"] = rsp.OrderInfo.Status
	reMap["user"] = rsp.OrderInfo.UserId
	reMap["post"] = rsp.OrderInfo.Post
	reMap["total"] = rsp.OrderInfo.Total
	reMap["address"] = rsp.OrderInfo.Address
	reMap["name"] = rsp.OrderInfo.Name
	reMap["mobile"] = rsp.OrderInfo.Mobile
	reMap["pay_type"] = rsp.OrderInfo.PayType
	reMap["order_sn"] = rsp.OrderInfo.OrderSn

	goodsList := make([]interface{}, 0)
	for _, item := range rsp.Goods {
		tmpMap := gin.H{
			"id":    item.GoodsId,
			"name":  item.GoodsName,
			"image": item.GoodsImage,
			"price": item.GoodsPrice,
			"nums":  item.Nums,
		}

		goodsList = append(goodsList, tmpMap)
	}
	reMap["goods"] = goodsList
	//支付
	if rsp.OrderInfo.Status == "PAYING" {
		var appId = global.ServerConfig.AlipayInfo.AppId
		var privateKey = global.ServerConfig.AlipayInfo.PrivateKey // 必须，上一步中使用 RSA签名验签工具 生成的私钥
		//var aliPublicKey="MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAkBui320sHC572qEBT8oCmaswX32x3YO2Di07XlkG71B+tdSJ7R6GzjVacrmGD2H9METl5b1Et0hbwrcBwnt/ER5Pyk1eprIy1UoTzcsJFH7sGfPNEET0jl0f3axwvFFNsIf9vi4mGjuR+YYZPcjYKpQ7ZhGT2I9mkHfKbkViQshe5LtlW9Qi4ag+vUl8vV0vkD5/OO1qP0c8zbwUmPpXO572TbuMVz4ZR+QSXsOUK/XTQbto3I2aBj12QnIwStOlH1WHxqnuxE3JQSYhS1p0jAb8dXTB2b2gdBwEDv9jahSNv1i1xRFPrnvH4tjkrn1frnJ2pEIICD1+dYg0624ztwIDAQAB"

		client, err := alipay.New(appId, privateKey, false)
		if err != nil {
			zap.S().Errorw("实例化支付宝失败失败")
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"msg": "内部错误",
			})
			return
		}
		//client.LoadAliPayPublicKey(aliPublicKey)
		//
		err = client.LoadAppPublicCertFromFile("alipay/appCertPublicKey.crt") // 加载应用公钥证书
		// 将 key 的验证调整到初始化阶段
		if err != nil {
			zap.S().Errorw("加载应用公钥证书失败")
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"msg": "内部错误",
			})
			return
		}
		err = client.LoadAliPayRootCertFromFile("alipay/alipayRootCert.crt") // 加载支付宝根证书
		if err != nil {
			zap.S().Errorw("加载支付宝根证书 失败")
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"msg": "内部错误",
			})
			return
		}
		err = client.LoadAliPayPublicCertFromFile("alipay/alipayCertPublicKey_RSA2.crt") // 加载支付宝公钥证书
		if err != nil {
			zap.S().Errorw("加载支付宝公钥证书 失败")
			return
		}
		// 将 key 的验证调整到初始化阶段
		if err != nil {
			zap.S().Errorw(err.Error())
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"msg": "内部错误",
			})
			return
		}

		var p = alipay.TradePagePay{}
		p.NotifyURL = global.ServerConfig.AlipayInfo.NotifyURL
		p.ReturnURL = global.ServerConfig.AlipayInfo.ReturnURL
		p.Subject = "慕学生鲜订单-" + rsp.OrderInfo.OrderSn
		p.OutTradeNo = rsp.OrderInfo.OrderSn
		p.TotalAmount = strconv.FormatFloat(float64(rsp.OrderInfo.Total), 'f', 2, 64)
		p.ProductCode = "FAST_INSTANT_TRADE_PAY"

		url, err := client.TradePagePay(p)
		if err != nil {
			zap.S().Errorw("生成支付url失败")
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"msg": err.Error(),
			})
			return
		}
		reMap["alipay_url"] = url.String()
	}

	ctx.JSON(http.StatusOK, reMap)
}

func Update(ctx *gin.Context) {

}
func Delete(ctx *gin.Context) {

}
