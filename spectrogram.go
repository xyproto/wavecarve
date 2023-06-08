package wavecarve

import (
	"image"
	"image/color"
	"math"
	"math/cmplx"

	"github.com/mjibson/go-dsp/fft"
)

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
			// Compute the magnitude and phase of the FFT value (log scale for magnitude)
			mag := 20 * math.Log10(cmplx.Abs(val))
			phase := cmplx.Phase(val)

			// Normalize the magnitude to the range of 0-255
			mag = (mag + 140) * 255 / 140

			// Cap the values at 0 and 255
			if mag < 0 {
				mag = 0
			} else if mag > 255 {
				mag = 255
			}

			// Normalize the phase to the range of 0-255
			phase = (phase + math.Pi) * 255 / (2 * math.Pi)

			// Set the pixel in the image
			img.Set(i, j, color.RGBA{uint8(mag), uint8(phase), 0, 255})
		}
	}

	return img, nil
}

func CreateAudioFromSpectrogram(img *image.RGBA) ([]int16, error) {
	// Get the image bounds
	bounds := img.Bounds()

	// Create a slice to hold the audio data
	audioData := make([]int16, bounds.Dx()*FFTSize)

	// Iterate over the image pixels
	for i := 0; i < bounds.Dx(); i++ {
		// Create a slice to hold the FFT frame
		fftFrame := make([]complex128, bounds.Dy())

		// Get the FFT frame from the image
		for j := 0; j < bounds.Dy(); j++ {
			// Get the pixel color
			r, g, _, _ := img.At(i, j).RGBA()

			// Compute the magnitude and phase from the pixel color
			mag := (float64(r) * 140 / 255) - 140
			phase := (float64(g) * 2 * math.Pi / 255) - math.Pi

			// Convert the magnitude from dB to linear
			mag = math.Pow(10, mag/20)

			// Create the complex FFT value
			fftFrame[j] = cmplx.Rect(mag, phase)
		}

		// Compute the inverse FFT of the frame
		frame := fft.IFFT(fftFrame)

		// Convert the complex values to float64s and then to int16s
		for j, val := range frame {
			audioData[i*FFTSize+j] = int16(real(val) * 32767)
		}
	}

	return audioData, nil
}
