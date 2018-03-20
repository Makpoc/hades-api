package hmap

import (
	"fmt"
	"image"
	"image/draw"
	"image/jpeg"
	"strings"

	"os"
	"path/filepath"
)

const outputPath = "output/"
const resPath = "res/"
const staticPath = "static/"

func getImagePath(args []string) string {
	if len(args) == 0 {
		return filepath.Join(outputPath, "base.jpeg")
	}
	return filepath.Join(outputPath, fmt.Sprintf("map_%s.jpeg", strings.Join(args, "_")))
}

func generateImage(args []string) error {
	var err error
	// if the image with coords already exists
	if _, err = os.Stat(getImagePath(args)); err == nil {
		// use it
		return nil
	}

	var baseImage image.Image
	var dBaseImage draw.Image
	baseImagePath := getImagePath(nil)
	if _, err := os.Stat(baseImagePath); err == nil {
		// load the image from the disk
		baseImage, err = LoadImage(baseImagePath)
		if err != nil {
			return err
		}

		if dbi, ok := baseImage.(draw.Image); ok {
			dBaseImage = dbi
		}
	}

	var layers = getLayers()
	if dBaseImage == nil {
		// if we cannot convert to draw.Image - generate it again.
		dBaseImage, err = GenerateBaseImage(layers)
		if err != nil {
			return err
		}

		err = saveImage(baseImagePath, dBaseImage)
		if err != nil {
			return err
		}
	}

	var result draw.Image
	color := DefaultColor
	if len(args) > 0 {
		for _, arg := range args {
			if isColor(arg) {
				color = Color(arg)
				continue
			}
			result, err = HighlightCoord(dBaseImage, arg, color)
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

func getLayers() []string {
	layers := []string{
		resPath + "/screenshot.jpeg",
		staticPath + "/coords.png",
		resPath + "/labels.png",
	}

	var result []string

	for _, layer := range layers {
		if _, err := os.Stat(layer); err == nil {
			result = append(result, layer)
		}
	}

	return result
}

// isColor checks if the given string corresponds to a known color
func isColor(arg string) bool {
	allColors := []Color{Yellow, Green, Pink, Orange, Red, Warn}
	return contains(allColors, Color(arg))
}

// contains checks if a set of colors contains given value
func contains(set []Color, val Color) bool {
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
