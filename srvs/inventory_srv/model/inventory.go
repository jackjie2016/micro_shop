package model

import (
	"database/sql/driver"
	"encoding/json"
)

//类型， 这个字段是否能为null， 这个字段应该设置为可以为null还是设置为空， 0
//实际开发过程中 尽量设置为不为null
//https://zhuanlan.zhihu.com/p/73997266
//这些类型我们使用int32还是int
type Inventory struct{
	BaseModel
	GoodsID int32 `gorm:"type:int;index"`
	Stocks int32 `gorm:"type:int"`
	Version int32 `gorm:"type:int"`
	Freeze int32 `gorm:"type:int"` //冻结库存
}

type Delivery struct{
	Goods int32 `gorm:"type:int;index"`
	Nums int32 `gorm:"type:int"`
	OrderSn int32 `gorm:"type:varchar(200)"`
	Status int32 `gorm:"type:int"`
}

type StockSellDetail struct{
	OrderSn string `gorm:"type:varchar(200);index:idx_index_ordersn,unique"`
	Status int32 `gorm:"type:int"` //1、已扣减，2、已归还
	Detail GoodsDetailList `gorm:"type:varchar(200)"` //商品详情
}

func (StockSellDetail) TableName() string{
	return "stockselldetail"
}

type GoodsDetail struct {
	Goods int32
	Nums int32
}

type GoodsDetailList []GoodsDetail

func (g GoodsDetailList) Value() (driver.Value, error){
	return json.Marshal(g)
}

// 实现 sql.Scanner 接口，Scan 将 value 扫描至 Jsonb
func (g *GoodsDetailList) Scan(value interface{}) error {
	return json.Unmarshal(value.([]byte), &g)
}