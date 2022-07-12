package main

import (
	"context"
	"github.com/go-redis/redis/v8"
	"time"
)
var ctx = context.Background()
func main() {
	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
	})

	err := rdb.Set(ctx, "15958615799", "2222", time.Duration(86400)*time.Second).Err()
	if err != nil {
		panic(err)
	}

	//val, err := rdb.Get(ctx, "15958615799").Result()
	//if err != nil {
	//	panic(err)
	//}
	//fmt.Println("key", val)

	//val2, err := rdb.Get(ctx, "key2").Result()
	//if err == redis.Nil {
	//	fmt.Println("key2 does not exist")
	//} else if err != nil {
	//	panic(err)
	//} else {
	//	fmt.Println("key2", val2)
	//}
	// Output: key value
	// key2 does not exist
}
