package initialize

import (
	"encoding/json"
	"fmt"

	"github.com/nacos-group/nacos-sdk-go/clients"
	"github.com/nacos-group/nacos-sdk-go/common/constant"
	"github.com/nacos-group/nacos-sdk-go/vo"
	"github.com/spf13/viper"
	"go.uber.org/zap"

	"userop_srv/global"
)

func GetEnvInfo(env string) string {
	viper.AutomaticEnv()
	return viper.GetString(env)
}

func InitConfig() {
	fmt.Println("初始化viper")
	//通过这个方法很好的隔离线上和本地的区别
	data := GetEnvInfo("Debug")

	var configFileName string
	configFileNamePrefix := "config"
	if data == "true" {
		configFileName = fmt.Sprintf("./%s-debug.yaml", configFileNamePrefix)
	} else {
		configFileName = fmt.Sprintf("./%s-pro.yaml", configFileNamePrefix)
	}

	v := viper.New()
	//文件的路径如何设置
	v.SetConfigFile(configFileName)
	if err := v.ReadInConfig(); err != nil {
		panic(err)
	}

	if err := v.Unmarshal(&global.NacosConfig); err != nil {
		panic(err.Error())
	}

	zap.S().Infof("配置信息 ：%v", global.NacosConfig)

	LogRollingConfig := constant.ClientLogRollingConfig{
		MaxAge:  3,
		MaxSize: 100,
	}
	// 创建clientConfig
	clientConfig := constant.ClientConfig{
		NamespaceId:         global.NacosConfig.Namespace, // 如果需要支持多namespace，我们可以场景多个client,它们有不同的NamespaceId。当namespace是public时，此处填空字符串。
		TimeoutMs:           5000,
		NotLoadCacheAtStart: true,
		LogDir:              "tmp/nacos/log",
		CacheDir:            "tmp/nacos/cache",
		LogRollingConfig:    &LogRollingConfig,
		LogLevel:            "debug",
	}
	// 至少一个ServerConfig
	serverConfigs := []constant.ServerConfig{
		{
			IpAddr:      global.NacosConfig.Host,
			ContextPath: "/nacos",
			Port:        global.NacosConfig.Port,
			Scheme:      "http",
		},
	}

	// 创建动态配置客户端
	configClient, _ := clients.CreateConfigClient(map[string]interface{}{
		"serverConfigs": serverConfigs,
		"clientConfig":  clientConfig,
	})
	content, _ := configClient.GetConfig(vo.ConfigParam{
		DataId: global.NacosConfig.DataId,
		Group:  global.NacosConfig.Group,
	})

	//想要将一个json字符串转换成struct，需要去设置这个struct的tag
	json.Unmarshal([]byte(content), &global.ServerConfig)

	//监听动态变动
	//go func() {
	//	v.WatchConfig()
	//	v.OnConfigChange(func(e fsnotify.Event) {
	//
	//		zap.S().Infof("配置信息变动 ：%v",e.Name)
	//		_ = v.ReadInConfig() // 读取配置数据
	//		_ = v.Unmarshal(global.ServerConfig)
	//
	//		zap.S().Infof("配置信息变动 ：%v",global.ServerConfig)
	//	})
	//}()

}
