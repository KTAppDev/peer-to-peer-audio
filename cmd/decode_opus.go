package main

import (
	"fmt"

	"github.com/hraban/opus"
)

const (
	sampleRate = 48000
	channels   = 1 // mono; 2 for stereo
)

func DecodeOpusData(opusData []byte) ([]int16, error) {
	fmt.Println("decoding opus data")
	// Create a new Opus decoder
	dec, err := opus.NewDecoder(sampleRate, channels)
	if err != nil {
		fmt.Printf("Error creating Opus decoder: %v", err)
		return nil, err
	}

	// Calculate the frame size
	frameSizeMs := float32(len(opusData)) / channels * 1000 / sampleRate
	frameSize := channels * int(frameSizeMs) * sampleRate / 1000

	// Create a buffer for the PCM data
	pcm := make([]int16, frameSize)

	// Decode the Opus data
	_, err = dec.Decode(opusData, pcm)
	if err != nil {
		fmt.Printf("Error decoding Opus data: %v", err)
		return nil, err
	}

	return pcm, nil
}
