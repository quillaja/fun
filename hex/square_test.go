package hex

import (
	"math"
	"testing"
)

func TestNewSquareGrid(t *testing.T) {
	t.SkipNow()
	grid := NewSquareGrid(2, math.Pi/2)
	t.Log(grid.toWorldMat)
}
