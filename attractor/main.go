package main

import (
	"flag"
	"fmt"
	grand "math/rand"
	"time"

	"github.com/quillaja/goutil/rand"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"golang.org/x/image/colornames"
)

const (
	width  = 1200
	height = 800
	title  = "Particles"
)

func run() {

	var numBalls int
	var numAttractors int
	flag.IntVar(&numBalls, "n", 10, "Number of balls.")
	flag.IntVar(&numAttractors, "a", 1, "Number of attractors.")
	flag.Parse()

	grand.Seed(time.Now().UnixNano())
	cfg := pixelgl.WindowConfig{
		Title:  title,
		Bounds: pixel.R(0, 0, width, height),
		VSync:  true,
	}
	win, err := pixelgl.NewWindow(cfg)
	if err != nil {
		panic(err)
	}

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
		ball.Color = colornames.Blue
		balls = append(balls, ball)
	}

	// attractors randomly positioned on screen
	attractors := []*Particle{}
	for i := 0; i < numAttractors; i++ {
		// mass is kinda "high"?
		a := NewParticleParams(rand.Float64NM(0, width), rand.Float64NM(0, height), 1000, 6, colornames.Red)
		// a.RepulsorDistance = 20
		attractors = append(attractors, a)
	}

	// gravity acts down
	gravity := pixel.V(0, -10)
	resistance := 5.0

	frames := 0
	timer := time.NewTicker(time.Second)
	start := time.Now()
	for !win.Closed() {
		dt := time.Since(start)
		start = time.Now()

		showVecs := win.Pressed(pixelgl.KeyV)

		// clear window
		win.Clear(colornames.White)

		// draw balls
		for _, ball := range balls {
			ball.Draw(showVecs)
			ball.GetVisual().Draw(win)
		}

		// draw attractors
		for _, attractor := range attractors {
			attractor.Draw(false) // won't move (no forces applied, etc)
			attractor.GetVisual().Draw(win)
		}

		win.Update()

		// update state
		for _, ball := range balls {
			ball.ResetForce()
			if win.Pressed(pixelgl.KeySpace) {
				ball.ApplyForce(gravity)
			}
			if win.Pressed(pixelgl.KeyR) {
				ball.ApplyForce(Resistance(&ball.Body, resistance))
			}
			for _, a := range attractors {
				ball.ApplyForce(Gravitation(&ball.Body, &a.Body)) // gravity between ball and attractor
			}
			ball.UpdatePosition(dt.Seconds())
		}

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
