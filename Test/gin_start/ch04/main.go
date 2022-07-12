package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type Person struct {
	ID   string `uri:"id" binding:"required,uuid"`
	Name string `uri:"name" binding:"required"`
}

func pong(c *gin.Context) {
	var person Person
	if err := c.ShouldBindUri(&person); err != nil {
		c.Status(http.StatusNotFound)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "pong",
		"id":      person.ID,
		"name":person.Name,
	})

}
func main() {
	r := gin.Default()
	r.GET("/:name/:id", pong)
	r.Run(":8082") // listen and serve on 0.0.0.0:8080
}
