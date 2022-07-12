package global

import (
	"gorm.io/gorm"

	"order_srv/config"
	"order_srv/proto"
)

var (
	DB           *gorm.DB
	ServerConfig *config.ServerConfig = &config.ServerConfig{}
	NocasConfig  config.NacosConfig

	GoodsSrvClient     proto.GoodsClient
	InventorySrvClient proto.InventoryClient
)
