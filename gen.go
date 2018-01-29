package main

import (
	"fmt"
	"image"
	"image/png"
	"os"
	"path/filepath"

	"github.com/makpoc/hadesmap/hadesmap"
)

const outputPath = "output/"

func getImagePath(coord string) string {
	return filepath.Join(outputPath, fmt.Sprintf("map_%s.png", coord))
}

func generateImage(coord string) error {

	resPath, err := filepath.Abs("./res")
	if err != nil {
		return err
	}

	if _, err := os.Stat(getImagePath(coord)); err == nil {
		return nil
	}

	baseImage, err := hadesmap.GenerateBaseImage(resPath+"/screenshot.png", resPath+"/map.png")
	if err != nil {
		return err
	}

	var result image.Image
	if coord != "" {
		result, err = hadesmap.DrawCoords(baseImage, coord)
		if err != nil {
			return err
		}
	} else {
		result = baseImage
	}

	dstImage, err := os.Create(getImagePath(coord))
	if err != nil {
		return err
	}
	defer dstImage.Close()

	err = png.Encode(dstImage, result)
	if err != nil {
		return err
	}

	return nil
}
