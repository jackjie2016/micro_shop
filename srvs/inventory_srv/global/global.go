package global

import (
	"gorm.io/gorm"

	"inventory_srv/config"
)

var (
	DB           *gorm.DB
	ServerConfig *config.ServerConfig = &config.ServerConfig{}
	NocasConfig  config.NacosConfig
)
