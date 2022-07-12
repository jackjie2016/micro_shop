package global

import (
	ut "github.com/go-playground/universal-translator"
	"orders_web/config"
	"orders_web/proto"
)

var (
	Trans        ut.Translator
	ServerConfig *config.ServerConfig = &config.ServerConfig{}

	NacosConfig config.NacosConfig

	OrderSrvClient     proto.OrderClient
	GoodsSrvClient     proto.GoodsClient
	InventorySrvClient proto.InventoryClient
)
