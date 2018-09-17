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

// Decode decodes an image from io.Reader
func Decode(r io.Reader) (image.Image, error) {
	im, _, err := image.Decode(r)
	if err != nil {
		return nil, err
	}
	return im, nil
}

// Resize resizes an image
func Resize(img image.Image) image.Image {
	return imaging.Resize(img, 256, 0, imaging.Linear)
}

// Encode encodes an image into io.Reader
func Encode(im image.Image, ext string) (io.Reader, error) {
	reader, writer := io.Pipe()

	format, err := imaging.FormatFromExtension(ext)
	if err != nil {
		return nil, err
	}
	go func() {
		err := imaging.Encode(writer, im, format)
		defer writer.Close()
		if err != nil {
			panic(err)
		}
	}()

	return reader, nil
}
