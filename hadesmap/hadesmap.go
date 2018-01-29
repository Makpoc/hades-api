package hadesmap

import (
	"fmt"
	"image"
	"image/draw"
	"image/png"
	"os"

	"github.com/nfnt/resize"
)

const cellSize = 1.0 / 7.0 // 7 cells in a map both horizontally and vertically

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
	hexSelectorImg, err := loadImage("res/hex.png")
	if err != nil {
		return nil, err
	}
	hexImageResized := resize.Resize(0, uint(cellSize*float64(baseImage.Bounds().Dy())), hexSelectorImg, resize.Lanczos3)

	hexRect := getTargetPoint(coords, baseImage.Bounds(), hexImageResized.Bounds())

	fmt.Println(hexRect)
	draw.Draw(baseImage, hexRect, hexImageResized, image.Point{0, 0}, draw.Over)
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

func getTargetPoint(coords string, base, hex image.Rectangle) image.Rectangle {
	coordPoint := image.Point{
		int(float64(base.Dx()) * newCellPoint(coords).x),
		int(float64(base.Dy()) * newCellPoint(coords).y),
	}

	hexRect := image.Rectangle{coordPoint, coordPoint.Add(hex.Max)}

	return hexRect
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
