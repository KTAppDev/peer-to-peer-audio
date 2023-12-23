package main

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true // You may need to adjust this depending on your production setup
	},
	HandshakeTimeout: 10,
}

func handleConnections(c *gin.Context) {
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer conn.Close()

	for {
		_, p, err := conn.ReadMessage()
		if err != nil {
			fmt.Println(err)
			return
		}

		// Your audio data processing logic here

		// For simplicity, we just print the received message
		fmt.Printf("Received: %s\n", p)

		// Echo the message back to the client
		if err := conn.WriteMessage(websocket.TextMessage, p); err != nil {
			fmt.Println(err)
			return
		}
	}
}

func main() {
	router := gin.Default()
	router.GET("/ws", func(c *gin.Context) {
		handleConnections(c)
	})

	fmt.Println("WebSocket server running on :8080")
	err := router.Run(":8080")
	if err != nil {
		fmt.Println(err)
	}
}
