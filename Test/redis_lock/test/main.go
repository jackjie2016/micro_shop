package main

import (
	"fmt"
	goredislib "github.com/go-redis/redis/v8"
	"github.com/go-redsync/redsync/v4"
	"github.com/go-redsync/redsync/v4/redis/goredis/v8"
	"sync"
	"time"
)

func main() {

	client := goredislib.NewClient(&goredislib.Options{
		Addr: "localhost:6379",
	})
	pool := goredis.NewPool(client)

	rs := redsync.New(pool)

	// Obtain a new mutex by using the same name for all instances wanting the
	// same lock.
	gNum :=10
	mutexname := "my-global-mutex"

	var wg sync.WaitGroup
	wg.Add(gNum)
	for i :=0 ;i<gNum;i++ {
		go func() {
			defer wg.Done()

			mutex := rs.NewMutex(mutexname)


			if err := mutex.Lock(); err != nil {
				fmt.Printf(err.Error())
			}
			fmt.Println("获取锁成功")
			time.Sleep(time.Second*6)

			fmt.Println("开始释放锁")

			if ok, err := mutex.Unlock(); !ok || err != nil {
				fmt.Printf("unlock failed")
			}
			fmt.Println("释放锁")


		}()
	}
	wg.Wait()
}