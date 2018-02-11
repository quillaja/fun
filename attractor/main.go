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
)

func run() {

	var numBalls int
	flag.IntVar(&numBalls, "n", 10, "Number of balls.")
	flag.Parse()

	grand.Seed(time.Now().UnixNano())
	cfg := pixelgl.WindowConfig{
		Title:  "Particles",
		Bounds: pixel.R(0, 0, width, height),
		VSync:  true,
	}
	win, err := pixelgl.NewWindow(cfg)
	if err != nil {
		panic(err)
	}

	// 10 particles randomly positioned on screen
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

	// mass is kinda "high"
	attractor := NewParticleParams(width/3, height/3, 1000, 6, colornames.Red) //lower left
	// attractor.RepulsorDistance = 20
	attractor2 := NewParticleParams((width*2)/3, (height*2)/3, 1000, 6, colornames.Red) // upper right
	// attractor2.RepulsorDistance = 50

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

		// draw attractor
		attractor.Draw(false) // won't move anyway
		attractor.GetVisual().Draw(win)
		attractor2.Draw(false)
		attractor2.GetVisual().Draw(win)

		win.Update()

		// update state
		for _, ball := range balls {
			ball.ResetForce()
			if win.Pressed(pixelgl.KeySpace) {
				ball.ApplyStatic(gravity)
			}
			if win.Pressed(pixelgl.KeyR) {
				ball.ApplyResistance(resistance)
			}
			ball.ApplyGravitation(attractor.Body, attractor2.Body) // gravity between ball and attractor
			ball.UpdatePosition(dt.Seconds())
		}

		select {
		case <-timer.C:
			win.SetTitle(fmt.Sprintf("%d fps", frames))
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
