package main

import (
	"flag"
	"fmt"
	grand "math/rand"
	"os"
	"time"

	"golang.org/x/image/font/gofont/gomono"

	"github.com/faiface/pixel/text"
	"github.com/golang/freetype/truetype"

	"github.com/quillaja/goutil/rand"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"golang.org/x/image/colornames"
)

const (
	width      = 1200
	height     = 800
	title      = "Particles"
	trailAlpha = 0.2
)

var (
	atlas *text.Atlas
)

func run() {

	var numBalls int
	var numAttractors int
	flag.IntVar(&numBalls, "n", 10, "Number of balls.")
	flag.IntVar(&numAttractors, "a", 1, "Number of attractors.")
	flag.Usage = func() {
		msg :=
			`Shows a simple particle physics simulation.
Hold SPACE to view the objects' velocity (thin line) and force (thick line) vectors. 
Hold T to view the paths of all objects. 
Hold R to apply a resistance and G to apply gravity in the downward direction.
Press Ctrl+Q to exit. 
(requires OpenGL 3.3+)

Optional parameters:`
		fmt.Fprintln(os.Stderr, msg)
		flag.PrintDefaults()
	}
	flag.Parse()

	grand.Seed(time.Now().UnixNano()) // gotta seed that RNG

	cfg := pixelgl.WindowConfig{
		Title:  title,
		Bounds: pixel.R(0, 0, width, height),
		VSync:  true,
	}
	win, err := pixelgl.NewWindow(cfg)
	if err != nil {
		panic(err)
	}

	// open and make font, then pixel/text stuff
	font, err := truetype.Parse(gomono.TTF)
	if err != nil {
		panic(err)
	}
	fface := truetype.NewFace(font, &truetype.Options{
		Size:              10,
		GlyphCacheEntries: 1})
	atlas = text.NewAtlas(fface, text.ASCII) // required for pixel/text package
	txt := text.New(pixel.ZV, atlas)
	txt.Color = colornames.Black

	// particles randomly positioned on screen
	// with random starting velocities
	balls := []*Particle{}
	for i := 0; i < numBalls; i++ {
		ball := NewParticleDefault()
		ball.Pos.X, ball.Pos.Y = rand.Float64NM(0, width), rand.Float64NM(0, height)
		ball.Vel.X = rand.Float64NM(-20, 20)
		ball.Vel.Y = rand.Float64NM(-20, 20)
		ball.Radius = 2
		ball.Mass = 1
		ball.Color = colornames.Royalblue
		balls = append(balls, ball)
	}

	// attractors randomly positioned on screen
	attractors := []*Particle{}
	for i := 0; i < numAttractors; i++ {
		m := rand.Float64NM(200, 1000)
		a := NewParticleParams(
			rand.Float64NM(0, width),
			rand.Float64NM(0, height),
			m,
			m/100,
			colornames.Orangered)
		// a.RepulsorDistance = 10
		attractors = append(attractors, a)
	}

	// gravity acts down
	gravity := pixel.V(0, -10)
	resistance := 5.0

	// data of particle trails
	trails := new(AlphaTrail)
	trailsMatrix := pixel.IM.Moved(pixel.V(width/2, height/2))

	// other useful things in the main loop
	frames := 0
	timer := time.NewTicker(time.Second)
	start := time.Now()

	for !win.Closed() &&
		!(win.Pressed(pixelgl.KeyLeftControl) && win.Pressed(pixelgl.KeyQ)) {
		dt := time.Since(start)
		start = time.Now()

		// drawing ///////////////////////////////////////////////////

		// clear window
		win.Clear(colornames.White)

		// draw trails as background
		// OR draw the objects (and optionally their HUD)
		if win.Pressed(pixelgl.KeyT) {
			trails.GetSprite().DrawColorMask(win, trailsMatrix, colornames.Black)
		} else {
			showVecs := win.Pressed(pixelgl.KeySpace)

			// draw balls
			for _, ball := range balls {
				ball.Draw(showVecs)
				ball.GetVisual().Draw(win)
				if showVecs {
					drawString(win, txt,
						fmt.Sprintf("F(%0.1f, %0.1f)\nV(%0.1f, %0.1f)",
							ball.Force.X, ball.Force.Y, ball.Vel.X, ball.Vel.Y),
						ball.Pos.Add(pixel.V(ball.Radius+2, 0)))
				}
			}

			// draw attractors
			for _, attractor := range attractors {
				attractor.Draw(false) // won't move (no forces applied, etc)
				attractor.GetVisual().Draw(win)
				if showVecs {
					drawString(win, txt, fmt.Sprintf("M(%0.0f)", attractor.Mass),
						attractor.Pos.Add(pixel.V(attractor.Radius+2, 0)))
				}
			}
		}

		// finish window drawing
		win.Update()

		// update state //////////////////////////////////////////////////

		for _, ball := range balls {
			ball.ResetForce()
			if win.Pressed(pixelgl.KeyG) {
				ball.ApplyForce(gravity)
			}
			if win.Pressed(pixelgl.KeyR) {
				ball.ApplyForce(Resistance(&ball.Body, resistance))
			}
			for _, a := range attractors {
				ball.ApplyForce(Gravitation(&ball.Body, &a.Body)) // gravity between ball and attractor
			}
			ball.UpdatePosition(dt.Seconds())

			trails.AddAlpha(ball.Pos, trailAlpha)
		}

		// update framerate in window title
		select {
		case <-timer.C:
			win.SetTitle(fmt.Sprintf("%s - %d fps", title, frames))
			frames = 0
		default:
			frames++
		}
	}

	timer.Stop()
}

func main() {
	pixelgl.Run(run)
}

// draws a string to a location.
func drawString(w *pixelgl.Window, t *text.Text, text string, pos pixel.Vec) {
	fmt.Fprint(t, text)
	t.Draw(w, pixel.IM.Moved(pos))
	t.Clear()
}

// AlphaTrail is a simple data structure to allow the object paths
// to be tracked in a quick and (relatively) memory efficient manner.
// Each object increments the value at AlphaTrail[x][y] by some fractional
// amount. The value at any particular (x,y) is an alpha value used
// when drawing the image to screen. The values for the data structure
// should not exceed 1.0.
type AlphaTrail [width][height]float64

// Bounds provides the bounding rectangle, for implementing pixel.Picture
func (a *AlphaTrail) Bounds() pixel.Rect {
	return pixel.R(0, 0, width, height)
}

// Color provides a color value at a particular location. For implementing
// pixel.PictureColor.
func (a *AlphaTrail) Color(at pixel.Vec) pixel.RGBA {
	x, y := int(at.X), int(at.Y)
	return pixel.RGBA{R: a[x][y], G: a[x][y], B: a[x][y], A: a[x][y]}
}

// GetSprite is a convenient way to get a pixel.Sprite of this data, for drawing.
func (a *AlphaTrail) GetSprite() *pixel.Sprite {
	return pixel.NewSprite(a, a.Bounds())
}

// AddAlpha adds the given alpha to the pixel at coords "pos" if the current
// value is less than 1.0. The function checks that "pos" is within the data
// struct's bounds before doing anything.
func (a *AlphaTrail) AddAlpha(pos pixel.Vec, alpha float64) {
	if a.Bounds().Contains(pos) {
		x, y := int(pos.X), int(pos.Y)
		if a[x][y] < 1 {
			a[x][y] += alpha
		}
	}
}
