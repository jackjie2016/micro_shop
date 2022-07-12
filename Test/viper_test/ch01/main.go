package main

import (
	"fmt"
	"github.com/spf13/viper"
)

type ServerConfig struct {
	ServerName string `mapstructure:"name"`// mapstructure 是viper的struct类型
	Port int `mapstructure:"port"`

}
func main()  {
	v:=viper.New()
	//文件的路径如何设置
	v.SetConfigFile("viper_test/ch01/config.yaml")
	if err:=v.ReadInConfig();err!=nil{
		panic(err)
	}
	ServerConfig:=ServerConfig{}
	if err:=v.Unmarshal(&ServerConfig);err!=nil{
		panic(err.Error())
	}
	fmt.Println(ServerConfig)
}