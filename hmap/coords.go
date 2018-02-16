package hmap

import (
	"strconv"
	"strings"
)

type cellPoint struct {
	x float64
	y float64
}

type offset struct {
	x float64
	y float64
}

func newCellPoint(coords string) cellPoint {
	if coords == "" || len(coords) != 2 {
		return cellPoint{}
	}

	c := strings.ToLower(coords)

	letterCoord := []rune(c)[0] - 97
	if letterCoord < 0 || letterCoord > 6 {
		return cellPoint{}
	}
	numericCoord, err := strconv.Atoi(c[1:2])
	if err != nil {
		return cellPoint{}
	}

	// start from 0 like the letterCoord
	numericCoord--

	if numericCoord < 0 || numericCoord > 6 {
		return cellPoint{}
	}
	return getCellPoint(numericCoord, int(letterCoord))
}

// verticalOffset is the offset of all cells, part of the given second coordinate
var verticalOffset = map[int]float64{
	0: 1,
	1: 1,
	2: 0,
	3: 0,
	4: -1,
	5: -1,
	6: -2,
}

const shiftFixCoefficient = 0.007

func getCellPoint(r, c int) cellPoint {
	row := float64(r)
	col := float64(c)

	var verticalCellCenterOffset float64

	if r%2 == 0 {
		verticalCellCenterOffset = cellSize / 2
	}

	return cellPoint{row*cellSize - shiftFixCoefficient*row, col*cellSize + verticalOffset[r]*cellSize + verticalCellCenterOffset}
}
