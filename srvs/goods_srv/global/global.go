package global

import (
	"github.com/olivere/elastic/v7"
	"gorm.io/gorm"

	"goods_srv/config"
)

var (
	DB           *gorm.DB
	ServerConfig *config.ServerConfig = &config.ServerConfig{}
	NocasConfig  config.NacosConfig

	Esclient *elastic.Client
)

//func init() {
//	fmt.Println("什么时候调用的")
////	dsn := "root:admin123@tcp(127.0.0.1:3306)/mxshop_goods_srv?charset=utf8mb4&parseTime=True&loc=Local"
////	//dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local", ServerConfig.MysqlInfo.UserName, ServerConfig.MysqlInfo.PassWord, ServerConfig.MysqlInfo.Host, ServerConfig.MysqlInfo.Port, ServerConfig.MysqlInfo.BaseName)
////fmt.Println(fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local", ServerConfig.MysqlInfo.UserName, ServerConfig.MysqlInfo.PassWord, ServerConfig.MysqlInfo.Host, ServerConfig.MysqlInfo.Port, ServerConfig.MysqlInfo.BaseName))
////	newLogger := logger.New(
////		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
////		logger.Config{
////			SlowThreshold: time.Second, // 慢 SQL 阈值
////			LogLevel:      logger.Info, // Log level
////			Colorful:      true,        // 禁用彩色打印
////		},
////	)
////
////	// 全局模式
////	//NamingStrategy和Tablename不能同时配置，
////	var err error
////	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
////		NamingStrategy: schema.NamingStrategy{
////			//TablePrefix: "mxshop_",
////			SingularTable: true,
////		},
////		Logger: newLogger,
////	})
////	if err != nil {
////		panic(err)
////	}
//
//}
