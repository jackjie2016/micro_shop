package initialize

import (
	"OldPackageTest/hardware/test2/global"
	redigo "github.com/gomodule/redigo/redis"
	"time"
)

func InitRedis(){
	var addr = "127.0.0.1:6379"
	var password = ""

	global.GlobalRedi= PoolInitRedis(addr, password)


}

func PoolInitRedis(server string, password string) *redigo.Pool {
	return &redigo.Pool{
		MaxIdle:     2,//空闲数
		IdleTimeout: 240 * time.Second,
		MaxActive:   5,//最大数
		Wait:        true,
		Dial: func() (redigo.Conn, error) {
			c, err := redigo.Dial("tcp", server)
			if err != nil {
				return nil, err
			}
			if password != "" {
				if _, err := c.Do("AUTH", password); err != nil {
					c.Close()
					return nil, err
				}
			}
			return c, err
		},
		TestOnBorrow: func(c redigo.Conn, t time.Time) error {
			_, err := c.Do("PING")
			return err
		},
	}
}