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
	ReadBufferSize:   1024,
	WriteBufferSize:  1024,
}

var stop = make(chan struct{})

// Define a struct to represent the data format

type AudioMessage struct {
	Metadata  map[string]string `json:"metadata"`
	AudioData []byte            `json:"audioData"`
}

// ...

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

			_, err = outFile.Write(p)
			if err != nil {
				fmt.Println("Error writing audio data to file:", err)
				return
			}
			// pcm, decodeErr := decode_opus_data.DecodeOpusData(audioData)
			// if decodeErr != nil {
			// 	fmt.Println("Error decoding audio data:", decodeErr)
			// }
			// fmt.Println("Decoded PCM data:", pcm)
			// TODO: Decode the audio data into PCM data and write it to the buffer
			// You'll need to use an appropriate decoding library for this
			// For example, you might use the opus package's DecodeFloat32 function
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
