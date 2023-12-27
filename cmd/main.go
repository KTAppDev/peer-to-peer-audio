package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/ktappdev/peer_to_peer_audio/getbrowser"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
	HandshakeTimeout: 5 * time.Second,
	ReadBufferSize:   1024,
	WriteBufferSize:  1024,
}

var stop = make(chan struct{})

var clients = make(map[*websocket.Conn]bool) // connected clients

func handleConnections(c *gin.Context, browser string) {
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer conn.Close()

	// Set Content-Type based on browser
	contentType := ""
	switch browser {
	case "Chrome or Firefox":
		contentType = "audio/opus"
	case "Safari":
		contentType = "audio/aac"
	default:
		contentType = "application/octet-stream" // Fallback content type for unknown browsers
	}

	// Set Content-Type header
	c.Header("Content-Type", contentType)

	for client := range clients {
		fmt.Println(client.RemoteAddr())
	}

	var fileName string
	fmt.Println("Browser: ", browser)

	switch browser {
	case "Chrome or Firefox":
		fileName = "audio.opus"
	case "Safari":
		fileName = "audio.aac"
	default:
		// Handle other cases or provide a default value
		fileName = "audio.unknown"
	}

	fmt.Print(fileName)
	// // Open the Opus file for writing
	outFile, err := os.Create(fileName)
	if err != nil {
		fmt.Println("Error creating Opus file:", err)
		return
	}
	defer outFile.Close()

	// Add the new WebSocket connection to the map of clients
	clients[conn] = true

	for {
		select {
		case <-stop:
			return // Gracefully stop the server
		default:
			messageType, p, err := conn.ReadMessage()
			if err != nil {
				fmt.Println(err)
				return
			}

			fmt.Println(len(p), "bytes", messageType)

			// Broadcast to all clients
			for client := range clients {
				if err := client.WriteMessage(websocket.BinaryMessage, p); err != nil {
					log.Printf("error: %v", err)
					client.Close()
					delete(clients, client)
				}
			}

			_, err = outFile.Write(p)
			if err != nil {
				fmt.Println("Error writing audio data to file:", err)
				return
			}
		}
	}
}

func main() {
	router := gin.Default()

	router.GET("/ws", func(c *gin.Context) {
		userAgent := c.GetHeader("User-Agent")
		fmt.Println("User Agent: ", userAgent)
		browser := getbrowser.GetBrowser(userAgent)

		handleConnections(c, browser)
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
