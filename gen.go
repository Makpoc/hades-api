package main

import (
	"image"
	"image/png"
	"os"
	"path/filepath"

	"github.com/makpoc/hadesmap/hadesmap"
)

func generateImage(coord string) error {

	resPath, err := filepath.Abs("./res")
	if err != nil {
		return err
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

	dstImage, err := os.Create("output/dest.png")
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
