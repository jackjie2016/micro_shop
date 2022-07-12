package api

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"go.uber.org/zap"
	"math/rand"
	"net/http"
	"strings"
	"time"
	"user_web/forms"
	"user_web/global"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/dysmsapi"
	"github.com/gin-gonic/gin"
)

func GenerateSmsCode(witdh int) string {
	//生成width长度的短信验证码

	numeric := [10]byte{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}
	r := len(numeric)
	rand.Seed(time.Now().UnixNano())

	var sb strings.Builder
	for i := 0; i < witdh; i++ {
		fmt.Fprintf(&sb, "%d", numeric[rand.Intn(r)])
	}
	return sb.String()
}
func SendSms(ctx *gin.Context) {
	smsSendForm := forms.SmsSendForm{}
	if err := ctx.ShouldBind(&smsSendForm); err != nil {
		HandleValitor(ctx, err)
		return
	}
	mobile := smsSendForm.Mobile
	code := GenerateSmsCode(6)
	client, err := dysmsapi.NewClientWithAccessKey("cn-beijing", global.ServerConfig.AliSmsInfo.ApiKey, global.ServerConfig.AliSmsInfo.ApiSecret)
	if err != nil {
		panic(err)
	}
	request := requests.NewCommonRequest()
	request.Method = "POST"
	request.Scheme = "https" // https | http
	request.Domain = "dysmsapi.aliyuncs.com"
	request.Version = "2017-05-25"
	request.ApiName = "SendSms"
	request.QueryParams["RegionId"] = "cn-beijing"
	request.QueryParams["PhoneNumbers"] = mobile                                      //手机号
	request.QueryParams["SignName"] = global.ServerConfig.AliSmsInfo.SignName         //阿里云验证过的项目名 自己设置
	request.QueryParams["TemplateCode"] = global.ServerConfig.AliSmsInfo.TemplateCode //阿里云的短信模板号 自己设置
	request.QueryParams["TemplateParam"] = "{\"code\":" + code + "}"                  //短信模板中的验证码内容 自己生成   之前试过直接返回，但是失败，加上code成功。
	response, err := client.ProcessCommonRequest(request)

	zap.S().Infof("%v", client.DoAction(request, response))

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"msg": "发送失败",
		})
	}
	rdb := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", global.ServerConfig.RedisInfo.Host, global.ServerConfig.RedisInfo.Port),
		Password: "", // no password set
		DB:       0,  // use default DB
	})
	err = rdb.Set(context.Background(), mobile, code, time.Duration(global.ServerConfig.RedisInfo.Expire)*time.Second).Err()
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"msg": "redis不能链接",
		})
	}
	ctx.JSON(http.StatusOK, gin.H{
		"msg": "发送成功",
	})
	//json数据解析
}
