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
	return getCellPoint(int(letterCoord), numericCoord)
}

// verticalOffset is the offset of all cells, part of the given second coordinate
var verticalOffset = map[int]float64{
	0: 1,
	1: 1,
	2: 0,
	3: 0,
	4: 0,
	5: 1,
	6: 1,
}

// getCellPoint gets the point the row/column argument select
func getCellPoint(r, c int) cellPoint {
	row := float64(r)
	col := float64(c)

	var verticalCellCenterOffset float64

	if r%2 == 0 {
		verticalCellCenterOffset = cellSizeHight / 2
	}

	return cellPoint{row * cellSizeWight, col*cellSizeHight + verticalCellCenterOffset + verticalOffset[r]*cellSizeHight}

	//	return cellPoint{row*cellSize - shiftFixCoefficient*row, col*cellSize + verticalOffset[r]*cellSize + verticalCellCenterOffset}
}
