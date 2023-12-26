package decode_opus_data

import (
	"io"

	"github.com/hraban/opus"
)

// DecodeOpusToPCM decodes Opus data to PCM using the hraban/opus library
func DecodeOpusData(opusData []byte) ([]int16, error) {
	decoder, err := opus.NewDecoder(48000, 1)
	if err != nil {
		return nil, err
	}

	// Create a buffer for PCM data
	pcmData := make([]int16, 4096) // Adjust the buffer size as needed

	var pcmBuffer []int16

	// Decode Opus data to PCM
	for {
		n, err := decoder.Decode(opusData, pcmData)
		if err == io.EOF {
			break
		} else if err != nil {
			return nil, err
		}

		// Append the decoded PCM data to the buffer
		pcmBuffer = append(pcmBuffer, pcmData[:n]...)
	}

	return pcmBuffer, nil
}
