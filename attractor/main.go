package main

import (
	"time"

	"github.com/quillaja/goutil/rand"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"golang.org/x/image/colornames"
)

func run() {
	cfg := pixelgl.WindowConfig{
		Title:  "Pixel Rocks!",
		Bounds: pixel.R(0, 0, 1024, 768),
		VSync:  true,
	}
	win, err := pixelgl.NewWindow(cfg)
	if err != nil {
		panic(err)
	}

	// 10 particles randomly positioned on screen
	// with random starting velocities
	balls := []*Particle{}
	for i := 0; i < 10; i++ {
		ball := NewParticleDefault()
		ball.Pos.X, ball.Pos.Y = rand.Float64NM(10, 1000), rand.Float64NM(10, 700) // top center of screen
		ball.Vel.X = rand.Float64NM(-20, 20)
		ball.Vel.Y = rand.Float64NM(-20, 20)
		ball.Radius = 2
		ball.Mass = 1
		ball.Color = colornames.Blue
		balls = append(balls, ball)
	}

	// mass is kinda "high"
	attractor := NewParticleParams(512, 384, 10000, 5, colornames.Red)

	// gravity acts down
	gravity := pixel.V(0, -10)

	start := time.Now()

	for !win.Closed() {
		dt := time.Since(start)
		start = time.Now()

		showVecs := win.Pressed(pixelgl.KeyV)

		// update state
		for _, ball := range balls {
			ball.ResetForce() // have to do when using gravitation
			if win.Pressed(pixelgl.KeySpace) {
				// have to reapply gravity each time
				ball.ApplyStatic(gravity)
			}
			ball.ApplyGravitation(attractor.Body) // gravity between ball and attractor
			ball.UpdatePosition(dt.Seconds())
			// fmt.Printf("ball: P: %s\tF: %s\n", ball.Pos, ball.Force)
		}

		// clear window
		win.Clear(colornames.Whitesmoke)

		// draw balls
		for _, ball := range balls {
			ball.Draw(showVecs)
			ball.GetVisual().Draw(win)
		}

		// draw attractor
		attractor.Draw(false) // won't move anyway
		attractor.GetVisual().Draw(win)

		win.Update()
	}
}

func main() {
	pixelgl.Run(run)
}
