package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true // You may need to adjust this depending on your production setup
	},
	HandshakeTimeout: 5 * time.Second,
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
		fmt.Printf("Received: %b\n", p)
		message := fmt.Sprintf("Message received at %s", time.Now().Format("2006-01-02 15:04:05"))
		messageBytes := []byte(message)
		// Echo the message back to the client
		if err := conn.WriteMessage(websocket.TextMessage, messageBytes); err != nil {
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
