package wavecarve

import (
	"fmt"
	"image"
	"image/draw"

	"github.com/esimov/caire"
)

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
