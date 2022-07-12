package consul

import (
	"fmt"
	"github.com/hashicorp/consul/api"
)

//1、定义struct 工厂
type Registry struct {
	Host string
	Port int
}

//2、定义接口
type RegistryClient interface {
	Register(address string, port int, name string, tags []string, id string) error
	DeRegister(ServiceId string) error
}

//3、抽象工厂实现
func NewRegistry(Host string, Port int) RegistryClient {
	return &Registry{
		Host: Host,
		Port: Port,
	}
}
func (r *Registry) Register(address string, port int, name string, tags []string, id string) error {
	cfg := api.DefaultConfig()
	cfg.Address = fmt.Sprintf("%s:%d", r.Host, r.Port)

	client, err := api.NewClient(cfg)
	if err != nil {
		panic(err)
	}
	//生成对应的检查对象
	check := &api.AgentServiceCheck{
		HTTP:                           fmt.Sprintf("http://%s:%d/health", address, port), //当前服务的ip+端口http://192.168.31.134:8021/health
		Timeout:                        "5s",
		Interval:                       "5s",
		DeregisterCriticalServiceAfter: "10s",
	}

	//生成注册对象
	registration := new(api.AgentServiceRegistration)
	registration.Name = name
	registration.ID = id
	registration.Port = port
	registration.Tags = tags
	registration.Address = address
	registration.Check = check

	fmt.Println(registration)
	//client.Agent().ServiceDeregister("user-web2")
	err = client.Agent().ServiceRegister(registration)
	//client.Agent().ServiceDeregister(id)
	if err != nil {
		panic(err)
	}
	return nil
}

func (r *Registry) DeRegister(ServiceId string) error {
	cfg := api.DefaultConfig()
	cfg.Address = fmt.Sprintf("%s:%d", r.Host, r.Port)
	client, err := api.NewClient(cfg)
	if err != nil {
		panic(err)
	}

	return	client.Agent().ServiceDeregister(ServiceId)
}
