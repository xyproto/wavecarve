package main

import (
	"fmt"
	"image/png"
	"os"

	"github.com/xyproto/wavecarve"
)

func main() {
	fmt.Print("Reading example.wav...")

	streamer, format, err := wavecarve.ReadWavFile("example.wav")
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
	defer streamer.Close()

	fmt.Println("ok")
	fmt.Print("Creating spectrogram...")

	// A larger FFT size will give better frequency resolution
	const fftSize = 512

	spectrogram, err := wavecarve.CreateSpectrogramFromAudio(streamer, format, fftSize)
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
	fmt.Print("Carving the spectrogram...")

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
	fmt.Print("Converting spectogram to audio...")

	// Convert the carved image back to audio
	audioStreamer := wavecarve.CreateAudioFromSpectrogram(carvedImage, fftSize)

	// Define the audio format
	//format = beep.Format{
	//SampleRate:  beep.SampleRate(44100),
	//NumChannels: 2, // assuming stereo audio, adjust as necessary
	//Precision:   2, // 16-bit precision
	//}

	fmt.Println("ok")
	fmt.Print("Writing output.wav...")

	// Write the audio to a .wav file
	err = wavecarve.WriteWavFile("output.wav", audioStreamer, format)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}

	fmt.Println("ok")

}
