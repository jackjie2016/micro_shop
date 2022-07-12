package main

import (
	"OldPackageTest/nacos_test/config"
	"encoding/json"
	"fmt"
	"time"

	"github.com/nacos-group/nacos-sdk-go/clients"
	"github.com/nacos-group/nacos-sdk-go/common/constant"
	"github.com/nacos-group/nacos-sdk-go/vo"
)
func main()  {
	// 创建clientConfig
	clientConfig := constant.ClientConfig{
		NamespaceId:         "3c512c55-e348-40e9-8465-02eab2e0173c", // 如果需要支持多namespace，我们可以场景多个client,它们有不同的NamespaceId。当namespace是public时，此处填空字符串。
		TimeoutMs:           5000,
		NotLoadCacheAtStart: true,
		LogDir:              "/tmp/nacos/log",
		CacheDir:            "/tmp/nacos/cache",
		RotateTime:          "1h",
		MaxAge:              3,
		LogLevel:            "debug",
	}
	// 至少一个ServerConfig
	serverConfigs := []constant.ServerConfig{
		{
			IpAddr:      "127.0.0.1",
			ContextPath: "/nacos",
			Port:        8848,
			Scheme:      "http",
		},
	}

	// 创建动态配置客户端
	configClient, _ := clients.CreateConfigClient(map[string]interface{}{
		"serverConfigs": serverConfigs,
		"clientConfig":  clientConfig,
	})
	content, _ := configClient.GetConfig(vo.ConfigParam{
		DataId: "user_srv",
		Group:  "pro"})


	//fmt.Println(content) //字符串 - yaml
	serverConfig := config.ServerConfig{}
	//想要将一个json字符串转换成struct，需要去设置这个struct的tag
	json.Unmarshal([]byte(content), &serverConfig)

	fmt.Println(serverConfig.JWTInfo)

	_ = configClient.ListenConfig(vo.ConfigParam{
		DataId: "user_srv",
		Group:  "pro",
		OnChange: func(namespace, group, dataId, data string) {
			fmt.Println("配置文件变化")
			json.Unmarshal([]byte(data), &serverConfig)
			fmt.Println(serverConfig)
		},
	})
	time.Sleep(3000 * time.Second)
}
