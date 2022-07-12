package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"os"
	"path/filepath"
)

func main()  {
	r:=gin.Default()

	dir,_:=filepath.Abs(filepath.Dir(os.Args[0]))
	fmt.Println(dir)
	r.Static("/static","./static")
	r.LoadHTMLGlob("temp/**/*")//加载当前目录下所有的文件

	//r.LoadHTMLFiles("temp/index.tmpl")//可以加载多文件

	r.GET("/index", func(c *gin.Context) {
		c.HTML(http.StatusOK,"default/index.tmpl",gin.H{
			"title":"首页",
		})

	})
	r.GET("/goods/list", func(c *gin.Context) {
		c.HTML(http.StatusOK,"goods/index.tmpl",gin.H{
			"title":"商品",

		})
	})
	//var title ="用户页面"
	//r.GET("/user/list", func(c *gin.Context) {
	//	c.HTML(http.StatusOK,"user/index.tmpl",gin.H{
	//		"title":"用户页面",
	//	})
	//})

	r.GET("/user/list", func(c *gin.Context) {
		c.HTML(http.StatusOK,"user/index.tmpl",gin.H{
			"title":"商品",
		})
	})

	r.Run(":8082")
}
