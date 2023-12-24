package main

import (
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gordonklaus/portaudio"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
	HandshakeTimeout: 5 * time.Second,
}

var stop = make(chan struct{})

func int16ToFloat32(input []int16) []float32 {
	output := make([]float32, len(input))
	for i, v := range input {
		output[i] = float32(v) / 32767.0
	}
	return output
}

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

	// Initialize PortAudio
	portaudio.Initialize()
	defer portaudio.Terminate()

	// Create a buffer for the audio data
	buffer := make([]float32, 44100)

	// Open a stream for audio output
	stream, err := portaudio.OpenDefaultStream(0, 1, 44100, len(buffer), &buffer)
	if err != nil {
		fmt.Println("Error opening audio stream:", err)
		return
	}
	defer stream.Close()

	// Start the audio stream
	err = stream.Start()
	if err != nil {
		fmt.Println("Error starting audio stream:", err)
		return
	}

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

			// TODO: Decode the Opus data into PCM data and write it to the buffer
			// You'll need to use an Opus decoding library for this
			// For example, you might use the opus package's DecodeFloat32 function
			pcm, err := decodeOpusData(p)
			if err != nil {
				fmt.Println("Error decoding Opus data:", err)
			}

			pcm, decodeError := decodeOpusData(p)
			if decodeError != nil {
				fmt.Println("Error decoding Opus data:", err)
				return
			}
			floats := int16ToFloat32(pcm)
			copy(buffer, floats)

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
