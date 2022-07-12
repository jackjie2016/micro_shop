package main

import (
	redigo "github.com/gomodule/redigo/redis"
	"reflect"
	"unsafe"

	"time"
)
type MapStruct struct {
	Head uint16
	Tail uint16
}
type RedisKey struct {
	LockName string
	QueueName string
	StatisticsName string
}
var isHead bool
func main() {
	var addr = "127.0.0.1:6379"
	var password = ""

	pool := PoolInitRedis(addr, password)

	test(pool)

}
func test(pool *redigo.Pool)  {
	c1 := pool.Get()

	defer c1.Close()
	//Send + Flush + Receive = Do
	//rec1,err := c1.Do("Get","gwyy")
	//if err!=nil{
	//	fmt.Println(err.Error())
	//}
	//fmt.Println(rec1)
	//
	//
	//_,err = c1.Do("Set","xhc","xhc hahaha")
	//if err != nil {
	//	panic(err)
	//}

	//hash
	//_,err = c1.Do("hmset","runoobkey","name","zifeng","name2","zifeng2")
	//if err != nil {
	//	panic(err)
	//}
	//
	//_, err = c1.Do("hset", "myhash", "bike1", "mobike")
	//if err != nil {
	//	fmt.Println("haset failed", err.Error())
	//}
	//5秒过期
	//_,err = c1.Do("Set","xhc1","xhc hahaha","EX",5)
	//if err != nil {
	//	panic(err)
	//}

	//检查key是否存在
	//is_key_exit,err := redigo.Bool(c1.Do("EXISTS","gwyy"))
	//if err != nil {
	//	panic(err)
	//}
	//fmt.Println(is_key_exit)

	//删除一个key
	//_, err = c1.Do("DEL", "Num_1_header")
	//if err != nil {
	//	fmt.Println("redis delelte failed:", err)
	//}
	//当前rfid值

	//nameAndData,err := redigo.String(c1.Do("rpop","rfidList"))
	//if err!=nil{
	//	//if err == redigo.ErrNil {
	//	//	err = nil
	//	//	return
	//	//}
	//	//return
	//	panic(err)
	//}
	//fmt.Println(len(nameAndData))
	//if len(nameAndData) > 1 {
	//	fmt.Println(nameAndData)
	//}




	//设置过期时间
	//n, _ := c1.Do("EXPIRE", "xhc", 24*3600)
	//fmt.Println(n)
	//list 操作
	//_,_ = c1.Do("lpush","rfidList","L0001")
	//_,_ = c1.Do("lpush","rfidList","L0002")
	//_,_ = c1.Do("lpush","rfidList","L0003")
	//_,_ = c1.Do("lpush","rfidList","L0004")
	//_,_ = c1.Do("lpush","rfidList","L0005")
	//_,_ = c1.Do("lpush","rfidList","L0006")
	//_,_ = c1.Do("lpush","rfidList","L0007")
	//_,_ = c1.Do("lpush","rfidList","L0008")
	//_,_ = c1.Do("lpush","rfidList","L0009")
	//_,_ = c1.Do("lpush","rfidList","L000A")
	//_,_ = c1.Do("lpush","rfidList","L000B")
	//_,_ = c1.Do("lpush","rfidList","L000C")


	for k := 0; k <= 100; k++ {
		_,_ = c1.Do("lpush","HeadRfidList","H0003")
		_,_ = c1.Do("lpush","HeadRfidList","H0004")
		_,_ = c1.Do("lpush","HeadRfidList","H0005")
		_,_ = c1.Do("lpush","HeadRfidList","H0006")
	}

	for k := 0; k <= 100; k++ {
		_,_ = c1.Do("lpush","HeadRfidList","H0004")
		_,_ = c1.Do("lpush","HeadRfidList","H0005")
		_,_ = c1.Do("lpush","HeadRfidList","H0006")
		_,_ = c1.Do("lpush","HeadRfidList","H0007")
	}

	for k := 0; k <= 100; k++ {
		_,_ = c1.Do("lpush","HeadRfidList","H0005")
		_,_ = c1.Do("lpush","HeadRfidList","H0006")
		_,_ = c1.Do("lpush","HeadRfidList","H0007")
		_,_ = c1.Do("lpush","HeadRfidList","H0008")
	}

	//​ value 将一组命令结果转换为 []interface{}。如果err不等于nil，那么Values返回nil，err
	//values,_ := redigo.Values(c1.Do("lrange","rfidList","0","-1"))
	//for _,v := range values {
	//	fmt.Println(string(v.([]byte)))
	//}
}
func String2Bytes(s string) []byte {
	sh := (*reflect.StringHeader)(unsafe.Pointer(&s))
	bh := reflect.SliceHeader{
		Data: sh.Data,
		Len:  sh.Len,
		Cap:  sh.Len,
	}
	return *(*[]byte)(unsafe.Pointer(&bh))
}

func Bytes2String(b []byte) string {
	return *(*string)(unsafe.Pointer(&b))
}

// redis pool
func PoolInitRedis(server string, password string) *redigo.Pool {
	return &redigo.Pool{
		MaxIdle:     2,//空闲数
		IdleTimeout: 240 * time.Second,
		MaxActive:   3,//最大数
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