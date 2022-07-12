package main

import "github.com/gin-gonic/gin"

func pong(c *gin.Context)  {
	c.JSON(200, gin.H{
		"message": "pong",
	})
}
func main() {
	router := gin.Default()//使用默认中间件 logger和 recovery
	//r := gin.New()
	router.GET("/ping", pong)
	goodsGrouter:=router.Group("/goods")
	{
		goodsGrouter.GET("/list",Goodslist)
		goodsGrouter.GET("/detail",GoodsDteail)
		goodsGrouter.GET("/add",CteateGoods)
	}

	router.Run(":8082") // listen and serve on 0.0.0.0:8080
}

func CteateGoods(context *gin.Context) {
	
}

func GoodsDteail(context *gin.Context) {
	context.JSON(200, gin.H{
		"message": "detail",
	})
}

func Goodslist2(context *gin.Context) {
	
}

func Goodslist(context *gin.Context) {
	
}
