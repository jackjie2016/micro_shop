package main

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	elastic "github.com/olivere/elastic/v7"
	"log"
	"os"
	"strconv"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"

	"io"

	_ "github.com/anaskhan96/go-password-encoder"

	"goods_srv/model"
)

func genMd5(code, salt string) string {
	Md5 := md5.New()
	_, _ = io.WriteString(Md5, code+salt)
	return hex.EncodeToString(Md5.Sum(nil))
}

var DB *gorm.DB
var err error

func Mysql2Es() {
	client, err := elastic.NewClient(
		elastic.SetURL("http://47.97.215.11:9200"),
		elastic.SetSniff(false),
		elastic.SetInfoLog(log.New(os.Stdout, "mxshop", log.LstdFlags)),
	)

	if err != nil {
		panic(err.Error())
	}

	var goods []model.Goods
	DB.Find(&goods)
	for _, g := range goods {
		Esmodel := model.EsGoods{
			ID:          g.ID,
			CategoryID:  g.CategoryID,
			BrandsID:    g.BrandsID,
			OnSale:      g.OnSale,
			ShipFree:    g.ShipFree,
			IsNew:       g.IsNew,
			IsHot:       g.IsHot,
			Name:        g.Name,
			ClickNum:    g.ClickNum,
			SoldNum:     g.SoldNum,
			FavNum:      g.FavNum,
			MarketPrice: g.MarketPrice,
			GoodsBrief:  g.GoodsBrief,
			ShopPrice:   g.ShopPrice,
		}
		_, err := client.Index().Index(model.EsGoods{}.GetindexName()).BodyJson(Esmodel).Id(strconv.Itoa(int(g.ID))).Do(context.Background())
		if err != nil {
			panic(err.Error())
		}
	}
}
func main() {
	dsn := "root:root123@tcp(127.0.0.1:3306)/mxshop_srv?charset=utf8mb4&parseTime=True&loc=Local"

	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		logger.Config{
			SlowThreshold: time.Second, // 慢 SQL 阈值
			LogLevel:      logger.Info, // Log level
			Colorful:      true,        // 禁用彩色打印
		},
	)

	//全局模式
	//NamingStrategy和Tablename不能同时配置，
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
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
	//_ = DB.AutoMigrate(&model.Category{},
	//	&model.Brands{}, &model.GoodsCategoryBrand{}, &model.Banner{}, &model.Goods{})
	//
	Mysql2Es()

}
