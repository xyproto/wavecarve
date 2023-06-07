package wavecarve

import (
	"image"
	"image/png"
	"os"
	"testing"

	"github.com/xyproto/wavecarve"
)

func TestReadWavFile(t *testing.T) {
	streamer, format, err := wavecarve.ReadWavFile("example.wav")
	if err != nil {
		t.Fatalf("Failed to read WAV file: %v", err)
	}
	defer streamer.Close()

	// Test the streamer and format if necessary
	// ...
}

func TestCreateSpectrogramFromAudio(t *testing.T) {
	streamer, format, err := wavecarve.ReadWavFile("example.wav")
	if err != nil {
		t.Fatalf("Failed to read WAV file: %v", err)
	}
	defer streamer.Close()

	const fftSize = 512

	spectrogram, err := wavecarve.CreateSpectrogramFromAudio(streamer, format, fftSize)
	if err != nil {
		t.Fatalf("Failed to create spectrogram: %v", err)
	}

	// Test the spectrogram if necessary
	// ...
}

func TestCarveSeams(t *testing.T) {
	streamer, format, err := wavecarve.ReadWavFile("example.wav")
	if err != nil {
		t.Fatalf("Failed to read WAV file: %v", err)
	}
	defer streamer.Close()

	const fftSize = 512

	spectrogram, err := wavecarve.CreateSpectrogramFromAudio(streamer, format, fftSize)
	if err != nil {
		t.Fatalf("Failed to create spectrogram: %v", err)
	}

	carvedImage, err := wavecarve.CarveSeams(spectrogram, 50.0)
	if err != nil {
		t.Fatalf("Failed to carve seams: %v", err)
	}

	// Test the carved image if necessary
	// ...
}

func TestCreateAudioFromSpectrogram(t *testing.T) {
	streamer, format, err := wavecarve.ReadWavFile("example.wav")
	if err != nil {
		t.Fatalf("Failed to read WAV file: %v", err)
	}
	defer streamer.Close()

	const fftSize = 512

	spectrogram, err := wavecarve.CreateSpectrogramFromAudio(streamer, format, fftSize)
	if err != nil {
		t.Fatalf("Failed to create spectrogram: %v", err)
	}

	carvedImage, err := wavecarve.CarveSeams(spectrogram, 50.0)
	if err != nil {
		t.Fatalf("Failed to carve seams: %v", err)
	}

	audioStreamer := wavecarve.CreateAudioFromSpectrogram(carvedImage, fftSize)

	// Test the audio streamer if necessary
	// ...
}

func TestWriteWavFile(t *testing.T) {
	streamer, format, err := wavecarve.ReadWavFile("example.wav")
	if err != nil {
		t.Fatalf("Failed to read WAV file: %v", err)
	}
	defer streamer.Close()

	const fftSize = 512

	spectrogram, err := wavecarve.CreateSpectrogramFromAudio(streamer, format, fftSize)
	if err != nil {
		t.Fatalf("Failed to create spectrogram: %v", err)
	}

	carvedImage, err := wavecarve.CarveSeams(spectrogram, 50.0)
	if err != nil {
		t.Fatalf("Failed to carve seams: %v", err)
	}

	audioStreamer := wavecarve.CreateAudioFromSpectrogram(carvedImage, fftSize)

	format = beep.Format{
		SampleRate:  beep.SampleRate(44100),
		NumChannels: 2,
		Precision:   2,
	}

	err = wavecarve.WriteWavFile("output.wav", audioStreamer, format)
	if err != nil {
		t.Fatalf("Failed to write WAV file: %v", err)
	}

	// Test the written WAV file if necessary
	// ...
}

func saveImage(filePath string, img image.Image) error {
	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	err = png.Encode(file, img)
	if err != nil {
		return err
	}

	return nil
}
