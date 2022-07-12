package main

import (
	"fmt"
	validator "github.com/go-playground/validator/v10"
	uuid "github.com/satori/go.uuid"
	"github.com/spf13/viper"
	"os"
	"os/signal"
	"syscall"

	"github.com/gin-gonic/gin/binding"
	ut "github.com/go-playground/universal-translator"
	"go.uber.org/zap"

	"user_web/global"
	"user_web/initialize"
	"user_web/utils"
	"user_web/utils/register/consul"
	myvalidator "user_web/validator"
)

func main() {

	//1、初始化log
	initialize.InitLogger()

	//2、初始化配置文件
	initialize.InitConfig()

	//3、初始化routers
	Router := initialize.Routers()

	//4、初始化翻译
	if err := initialize.InitTrans("zh"); err != nil {
		panic(err.Error())
	}

	//5. 初始化srv的连接
	initialize.InitSrvConn()

	//6、注册验证器
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		_ = v.RegisterValidation("mobile", myvalidator.ValidateMobile)
		_ = v.RegisterTranslation("mobile", global.Trans, func(ut ut.Translator) error {
			return ut.Add("mobile", "{0} 非法的手机号码!", true) // see universal-translator for details
		}, func(ut ut.Translator, fe validator.FieldError) string {
			t, _ := ut.T("mobile", fe.Field())
			return t
		})
	}

	ip, _ := utils.ExternalIP()

	viper.AutomaticEnv()
	debug := viper.GetBool("Debug")
	var port int
	var err error
	if !debug {
		port, err = utils.GetFreePort()
		if err == nil {
			port = port
		}
	} else {
		port, err = utils.GetFreePort()
	}
	zap.S().Infof("启动服务，端口：%d", port)
	zap.S().Infof("Nacos获取的配置信息 ：%v", global.ServerConfig)

	registry := consul.NewRegistry(global.ServerConfig.ConsulInfo.Host, global.ServerConfig.ConsulInfo.Port)
	serviceID := fmt.Sprintf("%s", uuid.NewV4())
	registry.Register(ip.String(), port, global.ServerConfig.Name, global.ServerConfig.Tags, serviceID)

	go func() {
		if err := Router.Run(fmt.Sprintf(":%d", port)); err != nil {
			panic(err.Error())
			zap.S().Panic("启动失败：", err.Error())
		}
	}()

	//接收终止信号
	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	if err := registry.DeRegister(serviceID); err != nil {
		zap.S().Info("注销失败")
	}
	zap.S().Info("注销成功")
}
