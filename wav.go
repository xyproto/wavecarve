package wavecarve

import (
	"encoding/binary"
	"math"
	"os"
)

// WAVHeader represents the header of a WAV file.
type WAVHeader struct {
	ChunkID       [4]byte
	ChunkSize     uint32
	Format        [4]byte
	Subchunk1ID   [4]byte
	Subchunk1Size uint32
	AudioFormat   uint16
	NumChannels   uint16
	SampleRate    uint32
	ByteRate      uint32
	BlockAlign    uint16
	BitsPerSample uint16
	Subchunk2ID   [4]byte
	Subchunk2Size uint32
}

const (
	// Assume 44100 Hz sample rate
	SampleRate = 44100

	// Assume 16-bit depth
	BitsPerSample = 16

	// Assume mono audio
	NumChannels = 1

	// Assume PCM audio
	AudioFormat = 1

	// Compute other header values
	ByteRate   = SampleRate * NumChannels * BitsPerSample / 8
	BlockAlign = NumChannels * BitsPerSample / 8

	// Size of the FFT
	FFTSize = 1024
)

// Convert a slice of bytes to a slice of int16s
func bytesToInt16s(bytes []byte) []int16 {
	int16s := make([]int16, len(bytes)/2)
	for i := range int16s {
		int16s[i] = int16(binary.LittleEndian.Uint16(bytes[i*2 : (i+1)*2]))
	}
	return int16s
}

// Convert a slice of int16s to a slice of bytes
func int16sToBytes(int16s []int16) []byte {
	bytes := make([]byte, len(int16s)*2)
	for i, int16 := range int16s {
		binary.LittleEndian.PutUint16(bytes[i*2:(i+1)*2], uint16(int16))
	}
	return bytes
}

// Convert a slice of int16s to a slice of float64s
func int16sToFloat64s(int16s []int16) []float64 {
	float64s := make([]float64, len(int16s))
	for i, int16 := range int16s {
		float64s[i] = float64(int16) / math.MaxInt16
	}
	return float64s
}

// Convert a slice of float64s to a slice of int16s
func float64sToInt16s(float64s []float64) []int16 {
	int16s := make([]int16, len(float64s))
	for i, float64 := range float64s {
		int16s[i] = int16(float64 * math.MaxInt16)
	}
	return int16s
}

// Read a .wav file
func ReadWavFile(filePath string) ([]int16, WAVHeader, error) {
	// Open the .wav file
	file, err := os.Open(filePath)
	if err != nil {
		return nil, WAVHeader{}, err
	}
	defer file.Close()

	// Read the header
	var header WAVHeader
	if err := binary.Read(file, binary.LittleEndian, &header); err != nil {
		return nil, WAVHeader{}, err
	}

	// Read the audio data
	data := make([]byte, header.Subchunk2Size)
	if _, err := file.Read(data); err != nil {
		return nil, WAVHeader{}, err
	}

	// Convert the audio data to int16s
	int16s := bytesToInt16s(data)

	return int16s, header, nil
}

// Write a .wav file
func WriteWavFile(filePath string, int16s []int16, header WAVHeader) error {
	// Open the .wav file
	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	// Convert the int16s to bytes
	bytes := int16sToBytes(int16s)

	// Update the ChunkSize and Subchunk2Size in the header
	header.ChunkSize = 36 + uint32(len(bytes)) // 4 (ChunkID) + (8 + Subchunk1Size) + (8 + Subchunk2Size)
	header.Subchunk2Size = uint32(len(bytes))

	// Write the header
	if err := binary.Write(file, binary.LittleEndian, &header); err != nil {
		return err
	}

	// Write the audio data
	if _, err := file.Write(bytes); err != nil {
		return err
	}

	return nil
}
