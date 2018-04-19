package hex

import (
	"math"

	"github.com/go-gl/mathgl/mgl64"
)

// HexagonOrientation describes if a hexagon is flat topped or pointy topped.
type HexagonOrientation int

// constants for the 2 types of hexagon orientations
const (
	FlatTop   HexagonOrientation = iota
	PointyTop HexagonOrientation = iota
)

// Loc is a [2]int representing the (column, row) of a hexagonal grid tile's
// center.
type Loc [2]int

// R gets the row coordinate of Loc.
func (l Loc) R() int { return l[0] }

// C gets the column coordinate of Loc.
func (l Loc) C() int { return l[1] }

// CR gets the column and row coordinates of Loc.
func (l Loc) CR() (int, int) { return l[0], l[1] }

// HexGrid represents a grid of regular hexagons of either the "flat topped" or
// "pointy topped" variety.
//
// Arbitrary user data can be associated with a
// particular hexagon by using the 'Data' map. The grid is indexed by "columns"
// and "rows" using the "axial" style coordinates described by
// https://www.redblobgames.com/grids/hexagons/#coordinates-axial. However,
// this grid follows the normal Y-orientation (+y = up) instead of the inverted
// one in the link
type HexGrid struct {
	Circumradius float64
	Inradius     float64
	Orientation  HexagonOrientation
	Data         map[Loc]interface{}
	toWorldMat   mgl64.Mat2
}

// NewHexGrid creates the data structure to represent a hexagonal grid in
// either flat-topped or pointy-topped regular hexagons.
func NewHexGrid(circumradius float64, orientation HexagonOrientation) *HexGrid {
	grid := &HexGrid{
		Circumradius: circumradius,
		Inradius:     circumradius * 0.86602540378, // = sqrt(3)/2 = cos(Pi/6)
		Orientation:  orientation,
		Data:         make(map[Loc]interface{}),
	}

	switch {
	case orientation == FlatTop:
		i := mgl64.Vec2{grid.Circumradius * 1.5, grid.Inradius}
		j := mgl64.Vec2{0, 2 * grid.Inradius}
		grid.toWorldMat = mgl64.Mat2FromCols(i, j)
	case orientation == PointyTop:
		i := mgl64.Vec2{2 * grid.Inradius, 0}
		j := mgl64.Vec2{grid.Inradius, grid.Circumradius * 1.5}
		grid.toWorldMat = mgl64.Mat2FromCols(i, j)
	default:
		panic("incorrect orientation")
	}

	return grid
}

// ToWorld converts axial grid coordinates to world/carteasian coordinates.
func (grid *HexGrid) ToWorld(c, r float64) (float64, float64) {
	world := grid.toWorldMat.Mul2x1(mgl64.Vec2{c, r})
	return world.X(), world.Y()
}

// ToGrid converts world coordinates to axial grid coordinates.
func (grid *HexGrid) ToGrid(x, y float64) (float64, float64) {
	g := grid.toWorldMat.Inv().Mul2x1(mgl64.Vec2{x, y})
	return g.X(), g.Y()
}

// Vertices gets the 6 vertices (corners) of the hexagon at (c,r) in world
// coordinates, starting on the right and going in a counter-clockwise direction.
func (grid *HexGrid) Vertices(c, r int) (verts []mgl64.Vec2) {
	verts = make([]mgl64.Vec2, 6, 6)
	var offset float64
	if grid.Orientation == FlatTop {
		offset = 0
	} else { // PointyTop
		offset = math.Pi / 6 // 30 deg
	}

	x, y := grid.ToWorld(float64(c), float64(r))
	for i := 0.0; i < 6; i++ {
		theta := i*math.Pi/3 + offset
		sin, cos := math.Sincos(theta)
		verts[int(i)][0] = grid.Circumradius*cos + x
		verts[int(i)][1] = grid.Circumradius*sin + y
	}

	return
}

// Get returns the data at axial coordinates (c,r) and a boolean indicating
// whether or not data existed at that location. Really it's just a convenience
// method for accessing the Data member.
func (grid *HexGrid) Get(c, r int) (data interface{}, ok bool) {
	data, ok = grid.Data[Loc{c, r}]
	return
}

// Set sets the data at axial coordinates (c,r). Really it's just a convenience
// method for accessing the Data member. If data is nil, the map value at (c,r)
// is deleted.
func (grid *HexGrid) Set(c, r int, data interface{}) {
	k := Loc{c, r}
	grid.Data[k] = data
	if data == nil {
		delete(grid.Data, k)
	}
}

// Axial converts cube coordinates to axial coordinates.
func Axial(x, y, z float64) (float64, float64) {
	return x, y
}

// Cube converts axial coordinates to cube coordinates.
func Cube(c, r float64) (float64, float64, float64) {
	return c, r, -c - r
}

// CubeRound rounds fractional cube coordinats to the center of the
// hexagon they're in.
func CubeRound(x, y, z float64) (float64, float64, float64) {
	rx := math.Round(x)
	ry := math.Round(y)
	rz := math.Round(z)

	xDiff := math.Abs(rx - x)
	yDiff := math.Abs(ry - y)
	zDiff := math.Abs(rz - z)

	if xDiff > yDiff && xDiff > zDiff {
		rx = -ry - rz
	} else if yDiff > zDiff {
		ry = -rx - rz
	} else {
		rz = -rx - ry
	}

	return rx, ry, rz
}

// CubeRoundInt does the same as CubeRound() but converts to int for you.
func CubeRoundInt(x, y, z float64) (int, int, int) {
	a, b, c := CubeRound(x, y, z)
	return int(a), int(b), int(c)
}

// AxialRound rounds fractional axial coordinates to the center of the
// hexagon they're in.
//
// First converts coords to cube, uses CubeRound(), then converts back.
func AxialRound(c, r float64) (float64, float64) {
	return Axial(CubeRound(Cube(c, r)))
}

// AxialRoundInt does the same as AxialRound() but converts to int for you.
func AxialRoundInt(c, r float64) (int, int) {
	a, b := AxialRound(c, r)
	return int(a), int(b)
}
