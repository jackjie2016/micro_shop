package config

type ServerConfig struct {
	Name       string       `mapstructure:"name" json:"name"`
	Host       string       `mapstructure:"host" json:"host"`
	Tags       []string     `mapstructure:"tags" json:"tags"`
	Port       int          `mapstructure:"port" json:"port"`
	MysqlInfo  MysqlConfig  `mapstructure:"mysql" json:"mysql"`
	ConsulInfo ConsulConfig `mapstructure:"consul" json:"consul"`
	RocketMq   ConsulConfig `mapstructure:"rocketmq" json:"rocketmq"`
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

type NacosConfig struct {
	Host      string `mapstructure:"host"`
	Port      uint64 `mapstructure:"port"`
	Namespace string `mapstructure:"namespace"`
	User      string `mapstructure:"user"`
	Password  string `mapstructure:"password"`
	DataId    string `mapstructure:"dataid"`
	Group     string `mapstructure:"group"`
}
