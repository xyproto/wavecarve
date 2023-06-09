package wavecarve

import (
	"image"
	"image/color"
	"math"
	"math/cmplx"

	"github.com/mjibson/go-dsp/fft"
)

// CreateSpectrogramFromAudio creates a spectrogram from an []int16 and
// encodes the length of the audio data into the image.
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
		//if i == 0 {
		//fmt.Printf("First FFT frame: %v\n", fftFrame[:10])
		//}

		// Find the maximum absolute sample in the frame
		maxSample := 0.0
		for _, val := range frame {
			absVal := math.Abs(val)
			if absVal > maxSample {
				maxSample = absVal
			}
		}
		// Normalize it to the range of 0-255
		volume := maxSample * 255

		// Set the pixels in the image
		for j, val := range fftFrame {
			// Compute the magnitude of the FFT value (log scale)
			mag := 20 * math.Log10(cmplx.Abs(val))
			// Normalize the magnitude to the range of 0-255
			mag = (mag + 140) * 255 / 140

			//fmt.Printf("Computed pixel color (r) from magnitude: %f\n", mag)

			// Cap the values at 0 and 255
			if mag < 0 {
				mag = 0
			} else if mag > 255 {
				mag = 255
			}

			// Compute the phase of the FFT value and normalize it
			phase := cmplx.Phase(val)
			// Map phase from [-pi, pi] to [0, 255]
			phase = (phase + math.Pi) * 255 / (2 * math.Pi)

			//fmt.Printf("Computed pixel color (g) from phase: %f\n", phase)

			// Set the pixel in the image, red for magnitude, green for phase, blue for volume
			img.Set(i, j, color.RGBA{uint8(mag), uint8(phase), uint8(volume), 255})
		}
	}

	// Encode length of audio data into first pixel's RGB values
	length := len(int16s)
	img.Set(0, 0, color.RGBA{uint8(length >> 16), uint8(length >> 8), uint8(length), 255})

	return img, nil
}

// CreateAudioFromSpectrogram creates audio from a spectrogram and
// extracts the length of the audio data from the image.
func CreateAudioFromSpectrogram(img *image.RGBA) ([]int16, error) {
	// Get the bounds of the image
	bounds := img.Bounds()

	// Extract length of audio data from first pixel's RGB values
	r, g, b, _ := img.At(0, 0).RGBA()
	length := int(r)<<16 | int(g)<<8 | int(b)

	// Create a slice to hold the audio data
	float64s := make([]float64, length)

	// Iterate over the pixels in the image
	for x := 0; x < bounds.Dx(); x++ {
		// Create a slice to hold the FFT frame
		iframe := make([]float64, FFTSize)

		// Extract the volume from the blue value of the first pixel in the column
		_, _, volume, _ := img.At(x, 0).RGBA()
		// Shift right by 8 bits
		volume = volume >> 8
		// Convert to float64
		fvolume := float64(volume) / 255

		for y := 0; y < FFTSize; y++ {
			// Get the pixel color
			r, g, _, _ := img.At(x, y).RGBA()

			// Shift right by 8 bits
			r = r >> 8
			g = g >> 8

			// Compute the magnitude from the red value
			mag := float64(r)*140.0/255.0 - 140.0
			// Convert magnitude from dB to linear
			mag = math.Pow(10.0, mag/20.0)

			//fmt.Printf("Computed magnitude from pixel color (r: %d): %f\n", r, mag)

			// Compute the phase from the green value
			// Map phase from [0, 255] to [-pi, pi]
			phase := float64(g)*2.0*math.Pi/255.0 - math.Pi

			//fmt.Printf("Computed phase from pixel color (g: %d): %f\n", g, phase)

			// Compute the FFT value
			fftVal := cmplx.Rect(mag, phase)

			// Add the FFT value to the frame
			iframe[y] = real(fftVal)
		}

		// Scale the inverse FFT frame with the extracted volume
		for i, val := range iframe {
			iframe[i] = val * fvolume
		}

		// Add the audio frame to the audio data
		for i, val := range iframe {
			if x*FFTSize+i < length {
				float64s[x*FFTSize+i] = val
			}
		}
	}

	// Convert the float64s to int16s
	int16s := float64sToInt16s(float64s)

	return int16s, nil
}
