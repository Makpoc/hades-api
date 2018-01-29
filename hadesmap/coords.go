package hadesmap

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

const cellSize = 1.0 / 7.0 // 7 cells in a map both horizontally and vertically

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
	0: 2,
	1: 1,
	2: 1,
	3: 0,
	4: 0,
	5: -1,
	6: -1,
}

func getCellPoint(r, c int) cellPoint {
	row := float64(r)
	col := float64(c)

	var verticalCellCenterOffset float64

	if r%2 != 0 {
		verticalCellCenterOffset = cellSize / 2
	}

	return cellPoint{row*cellSize + cellSize/2, col*cellSize + verticalOffset[r]*cellSize + verticalCellCenterOffset}
}
