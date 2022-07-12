package initialize

import (
	"fmt"
	"log"
	"os"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"

	"order_srv/global"
)

func init() {
	fmt.Println("会不会多次执行init")
}
func InitDb() {
	//dsn := "root:admin123@tcp(127.0.0.1:3306)/mxshop_order_srv?charset=utf8mb4&parseTime=True&loc=Local"
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local", global.ServerConfig.MysqlInfo.UserName, global.ServerConfig.MysqlInfo.PassWord, global.ServerConfig.MysqlInfo.Host, global.ServerConfig.MysqlInfo.Port, global.ServerConfig.MysqlInfo.BaseName)

	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		logger.Config{
			SlowThreshold: time.Second, // 慢 SQL 阈值
			LogLevel:      logger.Info, // Log level
			Colorful:      true,        // 禁用彩色打印
		},
	)

	// 全局模式
	//NamingStrategy和Tablename不能同时配置，
	var err error
	global.DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			//TablePrefix: "mxshop_",
			SingularTable: true,
		},
		Logger: newLogger,
	})
	if err != nil {
		panic(err)
	}

}
