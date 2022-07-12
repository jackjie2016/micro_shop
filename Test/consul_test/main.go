package main

import (
	"fmt"

	"github.com/hashicorp/consul/api"
)

func Register(address string, port int, name string, tags []string, id string) error {
	cfg := api.DefaultConfig()
	cfg.Address = "192.168.1.101:8500"

	client, err := api.NewClient(cfg)
	if err != nil {
		panic(err)
	}
	//生成对应的检查对象
	check := &api.AgentServiceCheck{
		HTTP:                           "http://192.168.1.101:8023/health",
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

func AllServices() {
	cfg := api.DefaultConfig()
	cfg.Address = ":8500"

	client, err := api.NewClient(cfg)
	if err != nil {
		panic(err)
	}

	data, err := client.Agent().Services()
	if err != nil {
		panic(err)
	}
	for key, _ := range data {
		fmt.Println(key)
	}
}
func FilterSerivice() {
	cfg := api.DefaultConfig()
	cfg.Address = "192.168.1.101:8500"

	client, err := api.NewClient(cfg)
	if err != nil {
		panic(err)
	}

	data, err := client.Agent().ServicesWithFilter(`Service == "mxshop_web"`)
	if err != nil {
		panic(err)
	}
	for key, _ := range data {
		fmt.Println(key)
	}
}

func main() {
	err:= Register("192.168.1.101", 8023, "user-web2", []string{"mxshop2", "bobby2"}, "user-web2")
	if err!=nil{
		panic(err)
	}

	//AllServices()
	//FilterSerivice()
	fmt.Println(fmt.Sprintf(`Service == "%s"`, "user-srv"))
}
