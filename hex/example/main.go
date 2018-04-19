package main

import (
	"fmt"
	"fun/hex"
	"math"
	"math/rand"
	"time"

	"golang.org/x/image/colornames"

	"github.com/lucasb-eyer/go-colorful"

	"github.com/quillaja/goutil/num"
	"github.com/quillaja/goutil/pxu"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
	"github.com/faiface/pixel/pixelgl"
)

func run() {
	rand.Seed(time.Now().UnixNano())

	cfg := pixelgl.WindowConfig{
		Title:  "HexGrid example",
		Bounds: pixel.R(0, 0, 800, 800),
		VSync:  true,
	}
	win, err := pixelgl.NewWindow(cfg)
	if err != nil {
		panic(err)
	}

	const hexRadius = 40
	grid := hex.NewHexGrid(hexRadius, hex.PointyTop)
	imd := imdraw.New(nil)

	// create some initial data
	pts, min, max := genHexPoints(5)
	for _, p := range pts {
		grid.Set(p.C(), p.R(), 0.0)
	}
	fmt.Printf("min and max grid coordinates: %d, %d\n", min, max)

	// func to draw the hex tiles with data
	render := func() {
		imd.Reset()
		for k, v := range grid.Data {
			imd.Color = colorful.Hsv(v.(float64), 1, 1)
			for _, vert := range grid.Vertices(k.CR()) {
				imd.Push(pixel.V(vert.X(), vert.Y()))
			}
			imd.Polygon(0)
		}
	}
	render() // do once

	cam := pxu.NewMouseCamera(win.Bounds().Center())

	for !win.Closed() {

		if win.JustPressed(pixelgl.MouseButtonRight) {
			click := cam.Unproject(win.MousePosition())
			x, y := grid.ToGrid(click.XY())
			c, r := hex.AxialRoundInt(x, y)
			loc := hex.Loc{c, r}
			fmt.Printf("raw grid: %0.2f, %0.2f\tloc: %v\n", x, y, loc)
			v, ok := grid.Data[loc].(float64)
			if ok {
				grid.Data[loc] = math.Mod(v+5, 360)
			} else {
				grid.Data[loc] = 0.0
			}
			render()
		}

		cam.Update(win)
		win.SetMatrix(cam.GetMatrix())

		win.Clear(colornames.Gray)
		imd.Draw(win)
		win.Update()
	}
}

func main() {
	pixelgl.Run(run)
}

func genHexPoints(edgeLength int) (points []hex.Loc, min, max int) {
	points = make([]hex.Loc, 0)

	min = -edgeLength * edgeLength
	max = ((1 + edgeLength) / 2) * (edgeLength / 2)
	edgeLength--

	for r, cStart := -edgeLength, 0; r <= edgeLength; r++ {
		cHeight := (edgeLength + 1) + (edgeLength - int(math.Abs(float64(r))))
		for c := cStart; cHeight > 0; cHeight, c = cHeight-1, c+1 {
			points = append(points, hex.Loc{c, r})
		}
		cStart = num.ClampInt(cStart-1, -edgeLength, 0)
	}

	return
}
