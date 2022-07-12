package main

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io"
	"log"
	"os"
	"time"

	_ "github.com/anaskhan96/go-password-encoder"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"

	"mxshop_srvs/inventory_srv/model"
)

func genMd5(code ,salt string) string {
	Md5:=md5.New()
	_,_=io.WriteString(Md5,code+salt)
	return  hex.EncodeToString(Md5.Sum(nil))
}

func main() {
	dsn := "root:admin123@tcp(127.0.0.1:3306)/mxshop_inventory_srv?charset=utf8mb4&parseTime=True&loc=Local"

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
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			//TablePrefix: "mxshop_",
			SingularTable: true,
		},
		Logger: newLogger,
	})
	if err != nil {
		panic(err)
	}

	//
	//_ = db.AutoMigrate(&model.Inventory{},&model.Delivery{},&model.StockSellDetail{})

     //OrderDetai:=model.StockSellDetail{
		//OrderSn: "12222",
		//Status:  1,
		//Detail:  model.GoodsDetailList{{1,3},{2,3}},
	 //}
	 //db.Create(OrderDetai)
	var SellDetail model.StockSellDetail
	if res:=db.Where(model.StockSellDetail{OrderSn: "12222"}).Find(&SellDetail);res.RowsAffected==0{

	}else {
		fmt.Println(SellDetail.Detail)
	}

}
