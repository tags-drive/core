package resizing

import (
	"image"
	"io"

	"github.com/disintegration/imaging"
)

const (
	width  = 256
	height = 0 // for github.com/disintegration/imaging (If one of width or height is 0, the image aspect ratio is preserved)
)

// Resize resizes an image
func Resize(originalImage io.Reader, newFilename string) error {
	im, _, err := image.Decode(originalImage)
	if err != nil {
		return err
	}
	resizedImage := imaging.Resize(im, 256, 0, imaging.Linear)
	return imaging.Save(resizedImage, newFilename)
}
