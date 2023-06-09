package main

import (
	"fmt"
	"github.com/xyproto/wavecarve"
	"os"
)

func main() {
	fmt.Print("Reading input.wav...")

	audioInts, header, err := wavecarve.ReadWavFile("input.wav")
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}

	fmt.Println("ok")
	fmt.Printf("First 10 audio samples: %v\n", audioInts[:10]) // Print first 10 samples
	fmt.Printf("Header: %+v\n", header)                        // Print header data

	fmt.Print("Creating spectrogram...")

	// A larger FFT size will give better frequency resolution
	const fftSize = 512

	spectrogram, err := wavecarve.CreateSpectrogramFromAudio(audioInts)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}

	fmt.Println("ok")
	fmt.Print("Creating audio from spectrogram...")

	// Convert the image back to audio data
	audioInts, err = wavecarve.CreateAudioFromSpectrogram(spectrogram)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}

	fmt.Println("ok")
	fmt.Printf("First 10 audio samples after conversion: %v\n", audioInts[:10]) // Print first 10 samples
	fmt.Printf("Header: %+v\n", header)                                         // Print header data

	fmt.Print("Writing output.wav...")

	// Write the audio data to the output file
	if err = wavecarve.WriteWavFile("output.wav", audioInts, header); err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}

	fmt.Println("ok")
}
