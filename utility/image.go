package utility

import (
	"image"
	"os"

	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"

	"github.com/disintegration/imaging"
	_ "golang.org/x/image/webp"
)

func CreateThumbnail(inputPath string, outputPath string, maxWidth, maxHeight int) error {
	inputFile, err := os.Open(inputPath)
	if err != nil {
		return err
	}
	defer inputFile.Close()

	inputImage, _, err := image.Decode(inputFile)
	if err != nil {
		return err
	}

	thumbImage := imaging.Thumbnail(inputImage, maxWidth, maxHeight, imaging.Lanczos)

	outputFile, err := os.Create(outputPath)
	if err != nil {
		return err
	}
	defer outputFile.Close()

	err = imaging.Encode(outputFile, thumbImage, imaging.JPEG, imaging.JPEGQuality(90))
	if err != nil {
		return err
	}

	return nil
}
