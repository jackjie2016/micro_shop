package main

import (
	proto "OldPackageTest/gin_start/ch06/probuf"
	"github.com/gin-gonic/gin"
	"net/http"
)

func pong(c *gin.Context)  {
	c.JSON(200, gin.H{
		"message": "<b>Hello, world!</b>",
	})
}
func main() {
	r := gin.Default()
	r.GET("/ping", pong)
	r.GET("/moreJosn", moreJosn)
	r.GET("/protobuf", returnprotobuf)

	r.GET("/purejson", PureJSON)

	r.Run(":8082") // listen and serve on 0.0.0.0:8080
}

func PureJSON(c *gin.Context) {
	c.PureJSON(200, gin.H{
		"html": "<b>Hello, world!</b>",
	})
}

func returnprotobuf(c *gin.Context) {
	var course =[]string{"1","2"}
	user:=&proto.Teacher{
		Name: "bobby",
		Course: course,
	}
	c.ProtoBuf(http.StatusOK,user)

}

func moreJosn(c *gin.Context) {
	var msg struct {
		Name string `json:"user"`
		message string
		Number int
	}
	msg.Name="zifeng"
	msg.message="懆懆懆懆懆懆"
	msg.Number=1
	c.JSON(http.StatusOK,msg)
}
