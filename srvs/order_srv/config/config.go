package config

type ServerConfig struct {
	Name             string         `mapstructure:"name" json:"name"`
	Host             string         `mapstructure:"host" json:"host"`
	Tags             []string       `mapstructure:"tags" json:"tags"`
	Port             int            `mapstructure:"port" json:"port"`
	MysqlInfo        MysqlConfig    `mapstructure:"mysql" json:"mysql"`
	GoodsSrvInfo     GoodsSrvConfig `mapstructure:"goods_srv" json:"goods_srv"`
	InventorySrvInfo GoodsSrvConfig `mapstructure:"inventory_srv" json:"inventory_srv"`
	ConsulInfo       ConsulConfig   `mapstructure:"consul" json:"consul"`
	JaegerInfo       JaegerConfig   `mapstructure:"jaeger" json:"jaeger"`
	RocketMq         ConsulConfig   `mapstructure:"rocketmq" json:"rocketmq"`
}

type MysqlConfig struct {
	Host     string `mapstructure:"host" json:"host"`
	Port     int    `mapstructure:"port" json:"port"`
	UserName string `mapstructure:"username" json:"username"`
	BaseName string `mapstructure:"basename" json:"basename"`
	PassWord string `mapstructure:"password" json:"password"`
}

type ConsulConfig struct {
	Host string `mapstructure:"host" json:"host"`
	Port int    `mapstructure:"port" json:"port"`
}

type GoodsSrvConfig struct {
	Name string `mapstructure:"name" json:"name"`
}

type InventorySrvConfig struct {
	Name string `mapstructure:"name" json:"name"`
}
type JaegerConfig struct {
	Name string `mapstructure:"name" json:"name"`
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
