package global

import (
	redigo "github.com/gomodule/redigo/redis"
)

var (
	GlobalRedi           *redigo.Pool
	RedisClient          *redigo.Conn
)
