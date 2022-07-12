package config

type ServerConfig struct {
	Name        string        `mapstructure:"name" json:"name"`
	Port        int           `mapstructure:"port" json:"port"`
	Tags        []string      `mapstructure:"tags" json:"tags"`
	UserSrvInfo UserSrvConfig `mapstructure:"user_srv" json:"user_srv"`
	JWTInfo     JWTConfig     `mapstructure:"jwt" json:"jwt"`
	AliSmsInfo  AliSmsConfig  `mapstructure:"sms" json:"sms"`
	RedisInfo   RedisConfig   `mapstructure:"redis" json:"redis"`
	ConsulInfo  ConsulConfig  `mapstructure:"consul" json:"consul"`
}
type UserSrvConfig struct {
	Name string `mapstructure:"name" json:"name"`
}

type JWTConfig struct {
	SigningKey string `mapstructure:"key" json:"key"`
}

type AliSmsConfig struct {
	ApiKey       string `mapstructure:"key"`
	ApiSecret    string `mapstructure:"secret"`
	SignName     string `mapstructure:"sign"`
	TemplateCode string `mapstructure:"template_code"`
}

type RedisConfig struct {
	Host   string `mapstructure:"host" json:"host"`
	Port   int    `mapstructure:"port" json:"port"`
	Expire int    `mapstructure:"expire" json:"expire"`
}
type ConsulConfig struct {
	Host string `mapstructure:"host" json:"host"`
	Port int    `mapstructure:"port" json:"port"`
}

type NacosConfig struct {
	Host      string `mapstructure:"host"`
	Port      uint64 `mapstructure:"port"`
	Namespace string `mapstructure:"namespace"`
	User      string `mapstructure:"user"`
	Password  string `mapstructure:"password"`
	DataId    string `mapstructure:"dataid"`
	Group     string `mapstructure:"group"`
}
