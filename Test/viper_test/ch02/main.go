package main

import (
	"fmt"
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
	"time"
)

type MysqlConfig struct {
	Host string `mapstructure:"host"` // mapstructure 是viper的struct类型
	Port int    `mapstructure:"port"`
}
type ServerConfig struct {
	ServerName  string      `mapstructure:"name"` // mapstructure 是viper的struct类型
	MysqlConfig MysqlConfig `mapstructure:"mysql"`
}

func GetEnvInfo(env string) string {
	viper.AutomaticEnv()
	return viper.GetString(env)
}

func main() {
	//通过这个方法很好的隔离线上和本地的区别
	data := GetEnvInfo("Debug")
	fmt.Println(data)
	var configFileName string
	configFileNamePrefix := "config"
	if data == "true" {
		configFileName = fmt.Sprintf("viper_test/ch02/%s-debug.yaml", configFileNamePrefix)
	} else {
		configFileName = fmt.Sprintf("viper_test/ch02/%s-pro.yaml", configFileNamePrefix)
	}
	fmt.Println(configFileName)
	v := viper.New()
	//文件的路径如何设置
	v.SetConfigFile(configFileName)
	if err := v.ReadInConfig(); err != nil {
		panic(err)
	}
	ServerConfig := ServerConfig{}
	if err := v.Unmarshal(&ServerConfig); err != nil {
		panic(err.Error())
	}
	fmt.Println(ServerConfig)

	//监听动态变动
	go func() {
		v.WatchConfig()
		v.OnConfigChange(func(e fsnotify.Event) {
			fmt.Println("Config file changed:", e.Name)
			_ = v.ReadInConfig() // 读取配置数据
			_ = v.Unmarshal(&ServerConfig)
			fmt.Println(ServerConfig)
		})
	}()

	time.Sleep(time.Second * 3000)
}
