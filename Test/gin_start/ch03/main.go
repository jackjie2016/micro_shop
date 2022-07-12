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
		goodsGrouter.GET("/v1/:id/:action",GoodsDteail)
		goodsGrouter.GET("/v2/:id/*action",GoodsDteail2)
		goodsGrouter.GET("/list",Goodslist)
		goodsGrouter.POST("/",CteateGoods)
	}

	router.Run(":8082") // listen and serve on 0.0.0.0:8080
}

func GoodsDteail2(c *gin.Context) {
	var id = c.Param("id")
	var action = c.Param("action")
	c.JSON(200, gin.H{
		"message": "detail",
		"id":id,
		"action":action,
	})
}

func CteateGoods(c *gin.Context) {

	c.JSON(200, gin.H{
		"message": "add",
	})
}

func GoodsDteail(c *gin.Context) {
	var id = c.Param("id")
	var action = c.Param("action")
	c.JSON(200, gin.H{
		"message": "detail",
		"id":id,
		"action":action,
	})
}



func Goodslist(context *gin.Context) {
	context.JSON(200, gin.H{
		"message": "list",
	})
}
