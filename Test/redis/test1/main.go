package main

import (
	"context"
	"fmt"
	"time"

	goredislib "github.com/go-redis/redis/v8"

	"github.com/go-redsync/redsync/v4/redis/goredis/v8"

)
func main()  {
	client := goredislib.NewClient(&goredislib.Options{
		Addr: "localhost:6379",
	})
	pool := goredis.NewPool(client)
	ctx := context.Background()
	c1,err:=pool.Get(ctx)
	if err!=nil{
		panic(err)
	}
	defer  c1.Close()
	c1.Set("test","222")
	t,err:=c1.PTTL("test")
	if err!=nil{
		panic(err)
	}
	c1.SetNX("test2","222",10000*time.Millisecond)

	fmt.Println(t)
}
