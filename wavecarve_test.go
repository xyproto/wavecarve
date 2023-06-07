package wavecarve

import (
	"image"
	"image/png"
	"os"
	"testing"
)

func TestReadWavFile(t *testing.T) {
	streamer, format, err := ReadWavFile("example.wav")
	if err != nil {
		t.Fatalf("Failed to read WAV file: %v", err)
	}
	defer streamer.Close()
}

func TestCreateSpectrogramFromAudio(t *testing.T) {
	streamer, format, err := ReadWavFile("example.wav")
	if err != nil {
		t.Fatalf("Failed to read WAV file: %v", err)
	}
	defer streamer.Close()

	const fftSize = 512

	spectrogram, err := CreateSpectrogramFromAudio(streamer, format, fftSize)
	if err != nil {
		t.Fatalf("Failed to create spectrogram: %v", err)
	}
}

func TestCarveSeams(t *testing.T) {
	streamer, format, err := ReadWavFile("example.wav")
	if err != nil {
		t.Fatalf("Failed to read WAV file: %v", err)
	}
	defer streamer.Close()

	const fftSize = 512

	spectrogram, err := CreateSpectrogramFromAudio(streamer, format, fftSize)
	if err != nil {
		t.Fatalf("Failed to create spectrogram: %v", err)
	}

	carvedImage, err := CarveSeams(spectrogram, 50.0)
	if err != nil {
		t.Fatalf("Failed to carve seams: %v", err)
	}
}

func TestCreateAudioFromSpectrogram(t *testing.T) {
	streamer, format, err := ReadWavFile("example.wav")
	if err != nil {
		t.Fatalf("Failed to read WAV file: %v", err)
	}
	defer streamer.Close()

	const fftSize = 512

	spectrogram, err := CreateSpectrogramFromAudio(streamer, format, fftSize)
	if err != nil {
		t.Fatalf("Failed to create spectrogram: %v", err)
	}

	carvedImage, err := CarveSeams(spectrogram, 50.0)
	if err != nil {
		t.Fatalf("Failed to carve seams: %v", err)
	}

	audioStreamer := CreateAudioFromSpectrogram(carvedImage, fftSize)
}

func TestWriteWavFile(t *testing.T) {
	streamer, format, err := ReadWavFile("example.wav")
	if err != nil {
		t.Fatalf("Failed to read WAV file: %v", err)
	}
	defer streamer.Close()

	const fftSize = 512

	spectrogram, err := CreateSpectrogramFromAudio(streamer, format, fftSize)
	if err != nil {
		t.Fatalf("Failed to create spectrogram: %v", err)
	}

	carvedImage, err := CarveSeams(spectrogram, 50.0)
	if err != nil {
		t.Fatalf("Failed to carve seams: %v", err)
	}

	audioStreamer := CreateAudioFromSpectrogram(carvedImage, fftSize)

	format = beep.Format{
		SampleRate:  beep.SampleRate(44100),
		NumChannels: 2,
		Precision:   2,
	}

	err = WriteWavFile("output.wav", audioStreamer, format)
	if err != nil {
		t.Fatalf("Failed to write WAV file: %v", err)
	}
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
