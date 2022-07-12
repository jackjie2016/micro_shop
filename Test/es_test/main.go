package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/olivere/elastic/v7"
	"log"
	"os"
)

type Goods struct {
	Name string `json:"name"`
	GoodsBrief string `json:"goods_brief"`
	ShopPrice float64 `json:"shop_price"`
	ID int32 `json:"id"`
}
func main()  {
	client, err := elastic.NewClient(
		elastic.SetURL("http://192.168.31.134:9200"),
		elastic.SetSniff(false),
		elastic.SetInfoLog(log.New(os.Stdout, "mxshop", log.LstdFlags)),
	)

	if err!=nil{
		panic(err.Error())
	}

	//q := elastic.NewMatchQuery("address", "street")
	q := elastic.NewBoolQuery()
	q = q.Must(elastic.NewMultiMatchQuery("深海速冻", "name","goods_brief"))
	q = q.Filter(elastic.NewRangeQuery("shop_price").Gte(10))
	res,err:=client.Search().Index("es_goods").Query(q).From(1).Size(10).Do(context.Background())
	if err!=nil{
		panic(err.Error())
	}

	for _,value:= range  res.Hits.Hits{
		var goods=Goods{}
		json.Unmarshal(value.Source,&goods)
		fmt.Println(goods)
	}
}
