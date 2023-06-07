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

// Function to convert a spectrogram back to a *beep.Streamer
func CreateAudioFromSpectrogram(img *image.RGBA, fftSize int) beep.Streamer {
	// Compute the number of FFT frames
	numFrames := img.Bounds().Dx()

	// Initialize the audio data slice
	data := make([][2]float64, numFrames*fftSize)

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

		// Store the real parts of the FFT frame in the audio data
		for j, val := range floatFrame {
			data[i*fftSize+j] = [2]float64{real(val), real(val)}
		}
	}

	// Create and return a Streamer for the audio data
	return beep.StreamerFunc(func(samples [][2]float64) (n int, ok bool) {
		n = copy(samples, data)
		return n, n > 0
	})
}

// Function to write a beep.Streamer to a .wav file
func WriteWavFile(filePath string, streamer beep.Streamer, format beep.Format) error {
	// Create the .wav file
	file, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}
	defer file.Close()

	// Write the Streamer to the file
	err = wav.Encode(file, streamer, format)
	if err != nil {
		return fmt.Errorf("failed to encode wav file: %w", err)
	}

	return nil
}
