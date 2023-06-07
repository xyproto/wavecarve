package wavecarve

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"math"
	"math/cmplx"
	"os"

	"github.com/esimov/caire"
	"github.com/faiface/beep"
	"github.com/faiface/beep/wav"
	"github.com/mjibson/go-dsp/fft"
)

// ReadWavFile can read a .wav file
func ReadWavFile(filePath string) (beep.StreamSeekCloser, beep.Format, error) {
	// Open the .wav file
	file, err := os.Open(filePath)
	if err != nil {
		return nil, beep.Format{}, err
	}

	// Create a new WAV decoder
	streamer, format, err := wav.Decode(file)
	if err != nil {
		return nil, beep.Format{}, err
	}

	return streamer, format, nil
}

// Function to create a spectrogram from an *audio.IntBuffer
func CreateSpectrogramFromAudio(streamer beep.StreamSeekCloser, format beep.Format, fftSize int) (*image.RGBA, error) {
	// Create buffer
	buffer := make([][2]float64, fftSize)

	// Create a new image with the width of the number of samples and height of fftSize
	img := image.NewRGBA(image.Rect(0, 0, streamer.Len()/fftSize, fftSize))

	// Iterate over the audio data
	for i := 0; i < streamer.Len()/fftSize; i++ {
		// Get the audio frame
		_, ok := streamer.Stream(buffer)

		if !ok {
			break
		}

		// Convert the frame to a float slice for the FFT function
		floatFrame := make([]float64, fftSize)
		for j, sample := range buffer {
			floatFrame[j] = sample[0] // Take first channel for simplicity
		}

		// Compute the FFT of the frame
		fftFrame := fft.FFTReal(floatFrame)

		// Set the pixels in the image
		for j, val := range fftFrame {
			// Compute the magnitude of the FFT value (log scale)
			mag := 20 * math.Log10(cmplx.Abs(val))

			// Normalize the magnitude to the range of 0-255
			mag = (mag + 140) * 255 / 140

			// Cap the values at 0 and 255
			if mag < 0 {
				mag = 0
			} else if mag > 255 {
				mag = 255
			}

			// Set the pixel in the image
			img.Set(i, j, color.RGBA{uint8(mag), uint8(mag), uint8(mag), 255})
		}
	}

	return img, nil
}

// CarveSeams removes seams from the image to reduce its width by the given percentage.
func CarveSeams(img *image.RGBA, newWidthInPercentage float64) (*image.RGBA, error) {
	// Calculate the new width
	newWidth := int(float64(img.Bounds().Dx()) * newWidthInPercentage / 100.0)

	// Convert img to *image.NRGBA
	nrgba := image.NewNRGBA(img.Bounds())
	draw.Draw(nrgba, nrgba.Bounds(), img, img.Bounds().Min, draw.Src)

	// Create the caire Processor with the desired options
	p := &caire.Processor{
		NewWidth:   newWidth,
		Percentage: false,
	}

	// Perform the seam carving
	resizedImage, err := p.Resize(nrgba)
	if err != nil {
		return nil, fmt.Errorf("could not resize image: %w", err)
	}

	// Convert back to *image.RGBA
	resizedRGBA := image.NewRGBA(resizedImage.Bounds())
	draw.Draw(resizedRGBA, resizedRGBA.Bounds(), resizedImage, resizedImage.Bounds().Min, draw.Src)

	return resizedRGBA, nil
}

// Function to convert a spectrogram back to audio samples
func CreateAudioFromSpectrogram(img *image.RGBA, fftSize int) ([]int16, error) {
	// Compute the number of FFT frames
	numFrames := img.Bounds().Dx()

	// Create a buffer for audio samples
	samples := make([]int16, numFrames*fftSize)

	// Iterate over the pixels in the image
	for i := 0; i < numFrames; i++ {
		// Initialize the FFT frame
		fftFrame := make([]complex128, fftSize)

		for j := 0; j < fftSize; j++ {
			// Get the pixel color
			pixel := img.RGBAAt(i, j)

			// Convert the pixel color to a magnitude
			mag := float64(pixel.R)*140/255 - 140

			// Convert the magnitude to an FFT value (assume phase is zero)
			fftFrame[j] = cmplx.Rect(math.Pow(10, mag/20), 0)
		}

		// Compute the inverse FFT of the frame
		floatFrame := fft.IFFT(fftFrame)

		// Convert float samples to int16
		for j, val := range floatFrame {
			sample := int16(real(val) * math.MaxInt16)
			samples[i*fftSize+j] = sample
		}
	}

	return samples, nil
}

// Function to write audio data to a .wav file
func WriteWavFile(filePath string, samples []int16, sampleRate int) error {
	// Create the .wav file
	file, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}
	defer file.Close()

	// Create a beep.Streamer from the samples
	streamer := beep.StreamerFunc(func(samples [][2]float64) (n int, ok bool) {
		for i := range samples {
			sample := float64(samples[i][0])
			sampleInt := int16(sample)
			samples[i][0] = float64(sampleInt)
			samples[i][1] = float64(sampleInt)
		}
		return len(samples), true
	})

	// Define the audio format
	format := beep.Format{
		SampleRate:  beep.SampleRate(sampleRate),
		NumChannels: 1,
		Precision:   2,
	}

	// Write the audio data to the .wav file
	err = wav.Encode(file, beep.Seq(streamer), format)
	if err != nil {
		return fmt.Errorf("failed to encode wav file: %w", err)
	}

	return nil
}
