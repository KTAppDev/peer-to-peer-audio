package main

import (
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
	HandshakeTimeout: 5 * time.Second,
}

var stop = make(chan struct{})

func handleConnections(c *gin.Context) {
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer conn.Close()

	// Open the Opus file for writing
	outFile, err := os.Create("output.opus")
	if err != nil {
		fmt.Println("Error creating Opus file:", err)
		return
	}
	defer outFile.Close()

	for {
		select {
		case <-stop:
			return // Gracefully stop the server
		default:
			_, p, err := conn.ReadMessage()
			if err != nil {
				fmt.Println(err)
				return
			}
			fmt.Println("Received audio data:", len(p))

			_, err = outFile.Write(p)
			if err != nil {
				fmt.Println("Error writing Opus data to file:", err)
				return
			}
		}
	}
}

func main() {
	router := gin.Default()

	router.GET("/ws", func(c *gin.Context) {
		handleConnections(c)
	})

	fmt.Println("WebSocket server running on :8080")
	go func() {
		if err := router.Run(":8080"); err != nil {
			fmt.Println(err)
		}
	}()

	// Gracefully stop the server on interrupt signal
	stopSignal := make(chan os.Signal, 1)
	signal.Notify(stopSignal, os.Interrupt)
	<-stopSignal
	close(stop)
	fmt.Println("Server stopped gracefully.")
}
