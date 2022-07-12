package main

import (
	"github.com/gin-gonic/gin"
	"golang.org/x/net/websocket"
	"log"
	"time"
)

func main() {
	log.SetFlags(log.Lshortfile | log.LstdFlags)

	var router = gin.Default()
	router.GET("/ws", ginWebsocketHandler(wsConnHandle))
	router.GET("/ping", pong)

	if err := router.Run(":9898"); err != nil {
		log.Println(err)
	}
}

func pong(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "pong",
	})
}

// websocket.Handler 转 gin HandlerFunc
func ginWebsocketHandler(wsConnHandle websocket.Handler) gin.HandlerFunc {
	return func(c *gin.Context) {
		log.Printf("new ws request: %v", c.Request.RemoteAddr)
		if c.IsWebsocket() {
			wsConnHandle.ServeHTTP(c.Writer, c.Request)
		} else {
			_, _ = c.Writer.WriteString("===not websocket request===")
		}
	}
}

// websocket连接处理
func wsConnHandle(conn *websocket.Conn) {
	for {
		var msg string
		if err := websocket.Message.Receive(conn, &msg); err != nil {
			log.Println(err)
			return
		}

		log.Printf("recv: %v", msg)

		data := []byte(time.Now().Format(time.RFC3339))
		if _, err := conn.Write(data); err != nil {
			log.Println(err)
			return
		}
	}
}