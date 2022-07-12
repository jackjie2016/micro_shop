package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

func TokenRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		var token string
		for k,v:=range  c.Request.Header{
			if k=="X-token"{
				token=v[0]
			}else{

			}
		}

		if token!="bobby"{
			c.JSON(http.StatusUnauthorized,gin.H{
				"msg":"未登录",
			})
			c.Abort()//只有这个才能中断
		}
		c.Next()

	}
}
func MyLogger() gin.HandlerFunc  {
	return func(c *gin.Context) {
		t:=time.Now()
		c.Set("example","123456")
		c.Next()

		end:=time.Since(t)
		fmt.Printf("耗时：%d\n",end)
		status:=c.Writer.Status()
		fmt.Println("状态",status)
	}
}
func main()  {

	r:=gin.Default()

	r.Use(MyLogger(),TokenRequired())

	r.GET("/ping",pong)
	r.POST("/auth",auth)

	r.Run(":8082")
}

func auth(c *gin.Context) {
	fmt.Println("请求开始")
	c.JSON(http.StatusOK,gin.H{
		"msg":"授权成功",
	})
}

func pong(c *gin.Context) {
	fmt.Println("请求开始")
	c.JSON(http.StatusOK,gin.H{
		"msg":"ping",
	})
	fmt.Println("请求结束")
}