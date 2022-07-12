package pay

import (
	"context"
	alipay "github.com/smartwalle/alipay/v3"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"orders_web/global"
	"orders_web/proto"
)

func Notify(ctx *gin.Context) {
	//支付宝回调通知
	client, err := alipay.New(global.ServerConfig.AlipayInfo.AppId, global.ServerConfig.AlipayInfo.PrivateKey, false)
	if err != nil {
		zap.S().Errorw("实例化支付宝失败")
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg": err.Error(),
		})
		return
	}
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

	noti, err := client.GetTradeNotification(ctx.Request)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{})
		return
	}

	_, err = global.OrderSrvClient.UpdateOrderStatus(context.Background(), &proto.OrderStatus{
		OrderSn: noti.OutTradeNo,
		Status:  string(noti.TradeStatus),
	})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{})
		return
	}
	ctx.String(http.StatusOK, "success")
}
