package initialize

import (
	"context"
	"fmt"
	"goods_srv/model"
	"log"
	"os"

	"github.com/olivere/elastic/v7"

	"goods_srv/global"
)

func InitEs() {

	//global.Esclient, err = elastic.NewClient(
	//	elastic.SetURL(fmt.Sprintf("http://%s:%d",global.ServerConfig.EsInfo.Host,global.ServerConfig.EsInfo.Port)),
	//	elastic.SetSniff(false),
	//	elastic.SetInfoLog(log.New(os.Stdout, "mxshop", log.LstdFlags)),
	//)

	host := fmt.Sprintf("http://%s:%d", global.ServerConfig.EsInfo.Host, global.ServerConfig.EsInfo.Port)
	//logger := log.New(os.Stdout, "mxshop", log.LstdFlags)
	fmt.Println(host)
	var err error
	global.Esclient, err = elastic.NewClient(
		elastic.SetURL(host),
		elastic.SetSniff(false),
		elastic.SetTraceLog(log.New(os.Stdout, "mxshop", log.LstdFlags)))
	if err != nil {
		panic(err)
	}

	if err != nil {
		panic(err.Error())
	}
	//新建mapping，判断mapping是否存在
	exists, err := global.Esclient.IndexExists(model.EsGoods{}.GetindexName()).Do(context.Background())
	if err != nil {
		panic(err.Error())
	}
	//不存在新建
	if !exists {
		_, err := global.Esclient.CreateIndex(model.EsGoods{}.GetindexName()).BodyString(model.EsGoods{}.GetMapping()).Do(context.Background())
		if err != nil {
			panic(err.Error())
		}
	}

}
