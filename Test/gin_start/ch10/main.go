package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"os"
	"os/signal"
	"syscall"
)

func main()  {
	//优雅退出
	router := gin.Default()//使用默认中间件 logger和 recovery
	//r := gin.New()
	router.GET("/ping", pong)

	go func() {
		router.Run(":8082") // listen and serve on 0.0.0.0:8080
	}()

	//如果想要接收到信号 不能使用kill -9（强杀命令） 只能用 kill
	quit:=make(chan os.Signal)
	signal.Notify(quit,syscall.SIGINT,syscall.SIGTERM)
	<-quit

	fmt.Println("关闭server中...")
	fmt.Println("注销服务...")

}

func pong(c *gin.Context) {
	
}

