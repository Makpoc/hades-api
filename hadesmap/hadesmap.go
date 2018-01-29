package hadesmap

import (
	"fmt"
	"image"
	"image/draw"
	"image/png"
	"os"

	"github.com/nfnt/resize"
)

// GenerateBaseImage generates the base image, composed of the real in game map with overlayed coordinates.
func GenerateBaseImage(screenFilePath, mapFilePath string) (draw.Image, error) {
	screenshotImage, err := loadImage(screenFilePath)
	if err != nil {
		return nil, err
	}
	mapImage, err := loadImage(mapFilePath)
	if err != nil {
		return nil, err
	}

	mapImageResized := resize.Resize(uint(screenshotImage.Bounds().Dx()), uint(screenshotImage.Bounds().Dy()), mapImage, resize.Lanczos3)

	baseImage := image.NewRGBA(image.Rect(0, 0, screenshotImage.Bounds().Dx(), screenshotImage.Bounds().Dy()))
	draw.Draw(baseImage, baseImage.Bounds(), screenshotImage, image.Point{0, 0}, draw.Over)
	draw.Draw(baseImage, baseImage.Bounds(), mapImageResized, image.Point{0, 0}, draw.Over)

	return baseImage, nil
}

// DrawCoords draws an arrow, pointing to the given coordinate.
func DrawCoords(baseImage draw.Image, coords string) (image.Image, error) {
	if !isValidCoord(coords) {
		return nil, fmt.Errorf("invalid coordinate: %s", coords)
	}
	arrowImage, err := loadImage("res/arrow_" + getImagePathForCoords(coords) + ".png")
	if err != nil {
		return nil, err
	}

	arrowRect := getTargetPoint(coords, baseImage.Bounds(), arrowImage.Bounds())

	draw.Draw(baseImage, arrowRect, arrowImage, image.Point{0, 0}, draw.Over)
	return baseImage, nil
}

func isValidCoord(coord string) bool {
	directions := []string{
		"a1", "a2", "a3", "a4",
		"b1", "b2", "b3", "b4", "b5",
		"c1", "c2", "c3", "c4", "c5", "c6",
		"d1", "d2", "d3", "d4", "d5", "d6", "d7",
		"e2", "e3", "e4", "e5", "e6", "e7",
		"f3", "f4", "f5", "f6", "f7",
		"g4", "g5", "g6", "g7",
	}

	for _, c := range directions {
		if coord == c {
			return true
		}
	}
	return false
}

func getImagePathForCoords(coords string) string {
	if coords == "" {
		return ""
	}

	const nw, ne, sw, se = "nw", "ne", "sw", "se"
	directions := map[string]string{
		"a1": nw, "a2": nw, "a3": nw, "a4": nw,
		"b1": nw, "b2": nw, "b3": nw, "b4": nw, "b5": nw,
		"c1": nw, "c2": nw, "c3": nw, "c4": nw, "c5": nw, "c6": nw,
		"d1": nw, "d2": nw, "d3": nw, "d4": nw, "d5": nw, "d6": nw, "d7": nw,
		"e2": nw, "e3": nw, "e4": nw, "e5": nw, "e6": nw, "e7": nw,
		"f3": nw, "f4": nw, "f5": nw, "f6": nw, "f7": nw,
		"g4": nw, "g5": nw, "g6": nw, "g7": nw,
	}

	return directions[coords]
}

func getTargetPoint(coords string, base, arrow image.Rectangle) image.Rectangle {
	coordPoint := image.Point{
		int(float64(base.Dx()) * newCellPoint(coords).x),
		int(float64(base.Dy()) * newCellPoint(coords).y),
	}

	arrowRect := image.Rectangle{coordPoint, coordPoint.Add(arrow.Max)}

	return arrowRect
}

func loadImage(filePath string) (image.Image, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	screenshotImage, err := png.Decode(file)
	if err != nil {
		return nil, err
	}
	return screenshotImage, nil
}
