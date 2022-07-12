package global

import (
	ut "github.com/go-playground/universal-translator"
	"user_web/config"
	"user_web/proto"
)

var (
	Trans         ut.Translator
	ServerConfig  *config.ServerConfig = &config.ServerConfig{}
	UserSrvClient proto.UserClient
	NacosConfig   config.NacosConfig
)
