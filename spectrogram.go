package wavecarve

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"math"
	"math/cmplx"
	"os"

	"github.com/mjibson/go-dsp/fft"
)

// Function to create a spectrogram from an []int16
func CreateSpectrogramFromAudio(int16s []int16) (*image.RGBA, error) {
	// Convert the int16s to float64s
	float64s := int16sToFloat64s(int16s)

	// Create a new image with the width of the number of samples and height of FFTSize
	img := image.NewRGBA(image.Rect(0, 0, len(float64s)/FFTSize, FFTSize))

	// Iterate over the audio data
	for i := 0; i < len(float64s)/FFTSize; i++ {
		// Get the audio frame
		frame := float64s[i*FFTSize : (i+1)*FFTSize]

		// Compute the FFT of the frame
		fftFrame := fft.FFTReal(frame)

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

// Function to read a spectrogram PNG file
func ReadSpectrogramPNG(filePath string) (*image.RGBA, error) {
	// Open the PNG file
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	// Decode the PNG image
	img, err := png.Decode(file)
	if err != nil {
		return nil, fmt.Errorf("failed to decode PNG image: %w", err)
	}

	// Convert the image to RGBA format if necessary
	rgba, ok := img.(*image.RGBA)
	if !ok {
		bounds := img.Bounds()
		rgba = image.NewRGBA(bounds)
		draw.Draw(rgba, bounds, img, bounds.Min, draw.Src)
	}

	return rgba, nil
}
