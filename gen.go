package main

import (
	"fmt"
	"image"
	"image/draw"
	"image/jpeg"
	"strings"

	"os"
	"path/filepath"

	"github.com/makpoc/hadesmap/hadesmap"
)

const outputPath = "output/"

func getImagePath(args []string) string {
	if len(args) == 0 {
		return filepath.Join(outputPath, "base.jpeg")
	}
	return filepath.Join(outputPath, fmt.Sprintf("map_%s.jpeg", strings.Join(args, "_")))
}

func generateImage(args []string) error {
	resPath, err := filepath.Abs("./res")
	if err != nil {
		return err
	}

	// if the image with coords already exists
	if _, err := os.Stat(getImagePath(args)); err == nil {
		// use it
		return nil
	}

	var baseImage image.Image
	var dBaseImage draw.Image
	baseImagePath := getImagePath(nil)
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

	var layers = []string{
		resPath + "/screenshot.jpeg",
		resPath + "/map.png",
		resPath + "/labels.png",
	}
	if dBaseImage == nil {
		// if we cannot convert to draw.Image - generate it again.
		dBaseImage, err = hadesmap.GenerateBaseImage(layers)
		if err != nil {
			return err
		}

		err = saveImage(baseImagePath, dBaseImage)
		if err != nil {
			return err
		}
	}

	var result draw.Image
	color := hadesmap.DefaultColor
	if len(args) > 0 {
		for _, arg := range args {
			if isColor(arg) {
				color = hadesmap.Color(arg)
				continue
			}
			result, err = hadesmap.HighlightCoord(dBaseImage, arg, color)
			if err != nil {
				return err
			}
		}
	} else {
		result = dBaseImage
	}

	err = saveImage(getImagePath(args), result)
	if err != nil {
		fmt.Println("Warning: Failed to save image cache to disk", err)
	}

	return nil
}

// isColor checks if the given string corresponds to a known color
func isColor(arg string) bool {
	allColors := []hadesmap.Color{hadesmap.Yellow, hadesmap.Green, hadesmap.Pink, hadesmap.Orange}
	return contains(allColors, hadesmap.Color(arg))
}

// contains checks if a set of colors contains given value
func contains(set []hadesmap.Color, val hadesmap.Color) bool {
	for _, c := range set {
		if val == c {
			return true
		}
	}
	return false
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
