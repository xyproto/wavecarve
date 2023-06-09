package main

import (
	"fmt"
	"image/png"
	"os"

	"github.com/xyproto/wavecarve"
)

func main() {
	fmt.Print("Reading input.wav...")

	audioInts, header, err := wavecarve.ReadWavFile("input.wav")
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
	fmt.Print("Seam carving the spectrogram...")

	carvedImage, err := wavecarve.CarveSeams(spectrogram, 50.0) // Reduce the width of the spectrogram by 50%
	if err != nil {
		fmt.Fprintf(os.Stderr, "Could not carve seams: %s\n", err)
		os.Exit(1)
	}

	fmt.Println("ok")
	fmt.Print("Writing carved.png...")

	// Create the output file
	carvedImageFile, err := os.Create("carved.png")
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
	defer carvedImageFile.Close()

	// Encode the image to the output file
	err = png.Encode(carvedImageFile, carvedImage)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}

	fmt.Println("ok")
	fmt.Print("Creating audio from carved spectrogram...")

	// Convert the image back to audio data
	audioInts, err = wavecarve.CreateAudioFromSpectrogram(carvedImage)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}

	fmt.Println("ok")
	fmt.Print("Writing output.wav...")

	// Write the audio data to the output file
	if err = wavecarve.WriteWavFile("output.wav", audioInts, header); err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}

	fmt.Println("ok")
}
