package main

import (
	"fmt"
	"image"
	"image/draw"
	"image/jpeg"

	"os"
	"path/filepath"

	"github.com/makpoc/hadesmap/hadesmap"
)

const outputPath = "output/"

func getImagePath(coord string) string {
	if coord == "" {
		return filepath.Join(outputPath, "base.jpeg")
	}
	return filepath.Join(outputPath, fmt.Sprintf("map_%s.jpeg", coord))
}

func generateImage(coord string) error {
	resPath, err := filepath.Abs("./res")
	if err != nil {
		return err
	}

	// if the image with coords already exists
	if _, err := os.Stat(getImagePath(coord)); err == nil {
		// use it
		return nil
	}

	var baseImage image.Image
	var dBaseImage draw.Image
	baseImagePath := getImagePath("")
	if _, err := os.Stat(baseImagePath); err == nil {
		// load the image from the disk
		baseImage, err = hadesmap.LoadImage(baseImagePath)
		if err != nil {
			return err
		}

		if dbi, ok := baseImage.(draw.Image); ok {
			dBaseImage = dbi
		}
	}

	if dBaseImage == nil {
		// if we cannot convert to draw.Image - generate it again.
		dBaseImage, err = hadesmap.GenerateBaseImage(resPath+"/screenshot.jpeg", resPath+"/map.png")
		if err != nil {
			return err
		}

		err = saveImage(baseImagePath, dBaseImage)
		if err != nil {
			return err
		}
	}

	var result draw.Image
	if coord != "" {
		result, err = hadesmap.DrawCoords(dBaseImage, coord)
		if err != nil {
			return err
		}
	} else {
		result = dBaseImage
	}

	return saveImage(getImagePath(coord), result)
}

// saveImage saves the image on the file system as a JPEG file
func saveImage(filePath string, image image.Image) error {
	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	err = jpeg.Encode(file, image, nil)
	if err != nil {
		return err
	}

	return nil
}
