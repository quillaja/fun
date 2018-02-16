package main

import (
	"fmt"
	"math"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
)

type clamp struct {
	Low, High float64
}

type Camera struct {
	Position  pixel.Vec
	PanSpeed  float64
	Zoom      float64
	ZoomSpeed float64

	XExtents clamp
	YExtents clamp
	ZExtents clamp

	UpButton, DownButton    pixelgl.Button
	LeftButton, RightButton pixelgl.Button
	ZoomLevel               func() float64

	prevWinBounds pixel.Rect
}

func NewCamera() *Camera {
	return &Camera{
		pixel.ZV,
		200,
		1,
		1.1,
		clamp{-5000, 5000},
		clamp{-5000, 5000},
		clamp{-50, 50},
		pixelgl.KeyUp, pixelgl.KeyDown, pixelgl.KeyLeft, pixelgl.KeyRight,
		nil,
		pixel.Rect{}}
}

func (cam *Camera) Update(win *pixelgl.Window, timeElapsed float64) {
	if cam.prevWinBounds != win.Bounds() {
		cam.prevWinBounds = win.Bounds()
		fmt.Println("set bounds", cam.prevWinBounds)
	}
	// update user controlled things
	// using switch because pressing both left/right or up/down at the same
	// time would just cancel out anyway
	if win.Pressed(cam.LeftButton) {
		cam.Position.X -= cam.PanSpeed * timeElapsed
	}
	if win.Pressed(cam.RightButton) {
		cam.Position.X += cam.PanSpeed * timeElapsed
	}
	if win.Pressed(cam.DownButton) {
		cam.Position.Y -= cam.PanSpeed * timeElapsed
	}
	if win.Pressed(cam.UpButton) {
		cam.Position.Y += cam.PanSpeed * timeElapsed
	}

	var zlvl float64
	if cam.ZoomLevel != nil {
		zlvl = cam.ZoomLevel()
	} else {
		zlvl = win.MouseScroll().Y
	}
	cam.Zoom *= math.Pow(cam.ZoomSpeed, zlvl)

	// clamp to extents
	cam.Position.X = pixel.Clamp(cam.Position.X, cam.XExtents.Low, cam.XExtents.High)
	cam.Position.Y = pixel.Clamp(cam.Position.Y, cam.YExtents.Low, cam.YExtents.High)
	cam.Zoom = pixel.Clamp(cam.Zoom, cam.ZExtents.Low, cam.ZExtents.High)

}

func (cam *Camera) GetMatrix() pixel.Matrix {
	return pixel.IM.Scaled(cam.Position, cam.Zoom).
		Moved(cam.prevWinBounds.Center().Sub(cam.Position))
}

func (cam *Camera) Unproject(point pixel.Vec) pixel.Vec {
	return cam.GetMatrix().Unproject(point)
}

func (cam *Camera) Reset() {
	cam.ResetPan()
	cam.ResetZoom()
}

func (cam *Camera) ResetPan() {
	cam.ResetXPan()
	cam.ResetYPan()
}

func (cam *Camera) ResetXPan() { cam.Position.X = 0 }

func (cam *Camera) ResetYPan() { cam.Position.Y = 0 }

func (cam *Camera) ResetZoom() { cam.Zoom = 1 }
