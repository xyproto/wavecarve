package main

import (
	"fmt"
	"image/png"
	"os"

	"github.com/xyproto/wavecarve"
)

func main() {
	fmt.Print("Reading example.wav...")

	audioInts, _, err := wavecarve.ReadWavFile("example.wav")
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}

	fmt.Println("ok")
	fmt.Print("Creating spectrogram...")

	// A larger FFT size will give better frequency resolution
	const fftSize = 512

	spectrogram, err := wavecarve.CreateSpectrogramFromAudio(audioInts)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}

	fmt.Println("ok")
	fmt.Print("Writing spectogram.png...")

	// Create the output file
	spectrogramImageFile, err := os.Create("spectrogram.png")
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
	defer spectrogramImageFile.Close()

	// Encode the image to the output file
	err = png.Encode(spectrogramImageFile, spectrogram)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}

	fmt.Println("ok")
}
