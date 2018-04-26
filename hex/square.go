package hex

import (
	"math"

	"github.com/go-gl/mathgl/mgl64"
)

// Grid is an interface for 2D polygon grids, such as a hexagonal or square grid.
type Grid interface {
	ToWorld(c, r float64) (float64, float64)
	ToGrid(x, y float64) (float64, float64)
	Vertices(c, r int) []mgl64.Vec2
	Get(c, r int) (interface{}, bool)
	Set(c, r int, data interface{})
	Map() map[Loc]interface{}
	Tile(c, r float64) (int, int) // converts fractional grid coords to the integer location of the grid unit
}

// SquareGrid represents a grid of squares.
type SquareGrid struct {
	SideLength   float64
	Circumradius float64
	Inradius     float64
	Orientation  float64
	Data         map[Loc]interface{}
	toWorldMat   mgl64.Mat2
}

// NewSquareGrid creates the data structure to represent a square grid. The
// squares of the grid can be rotated counterclockwise by angleRadians.
func NewSquareGrid(sideLength, angleRadians float64) *SquareGrid {
	grid := &SquareGrid{
		SideLength:   sideLength,
		Circumradius: math.Sqrt(2) / 2 * sideLength,
		Inradius:     sideLength / 2,
		Orientation:  angleRadians,
		Data:         make(map[Loc]interface{}),
	}

	grid.toWorldMat = mgl64.Rotate2D(angleRadians).Mul(sideLength)

	return grid
}

// ToWorld converts grid coordinates to world (screen) coordinates.
func (grid *SquareGrid) ToWorld(c, r float64) (float64, float64) {
	world := grid.toWorldMat.Mul2x1(mgl64.Vec2{c, r})
	return world.X(), world.Y()
}

// ToGrid converts world (screen) coordinates to grid coordinates.
func (grid *SquareGrid) ToGrid(x, y float64) (float64, float64) {
	g := grid.toWorldMat.Inv().Mul2x1(mgl64.Vec2{x, y})
	return g.X(), g.Y()
}

// Verticies gets the 4 vertices of the square at (c,r) in world coordinates,
// starting in the "top right" and going counter-clockwise.
func (grid *SquareGrid) Vertices(c, r int) (verts []mgl64.Vec2) {
	verts = make([]mgl64.Vec2, 4, 4)
	offset := grid.Orientation + math.Pi/4 // rotation + 45 deg

	x, y := grid.ToWorld(float64(c), float64(r))
	for i := 0.0; i < 4; i++ {
		theta := i*math.Pi/2 + offset
		sin, cos := math.Sincos(theta)
		verts[int(i)][0] = grid.Circumradius*cos + x
		verts[int(i)][1] = grid.Circumradius*sin + y
	}

	return
}

// Get returns the data at the grid coordinate (c,r) and a boolean indicating
// whether or not the data existed at that location.
func (grid *SquareGrid) Get(c, r int) (data interface{}, ok bool) {
	data, ok = grid.Data[Loc{c, r}]
	return
}

// Set sets the data at the grid coordinates (c,r). If data is nil, the value
// at (c,r) is deleted.
func (grid *SquareGrid) Set(c, r int, data interface{}) {
	k := Loc{c, r}
	grid.Data[k] = data
	if data == nil {
		delete(grid.Data, k)
	}
}

// Map gets access to the grid's data, for use in "range" etc.
func (grid *SquareGrid) Map() map[Loc]interface{} {
	return grid.Data
}

// Tile returns the grid coords (column and row) of the square containing
// the given fractional grid coordinates.
func (grid *SquareGrid) Tile(c, r float64) (int, int) {
	return int(math.Round(c)), int(math.Round(r))
}
