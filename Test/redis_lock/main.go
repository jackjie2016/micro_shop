package main

import (
	"OldPackageTest/redis_lock/demo"
	"time"

)

func main() {
	for i:=0;i<10;i++ {
		go demo.MockTest("A")
		go demo.MockTest("B")
		go demo.MockTest("C")
		go demo.MockTest("D")
		go demo.MockTest("E")
	}


	// 用于测试goroutine接收到ctx.Done()信号后的打印
	time.Sleep(time.Second * 20)
}