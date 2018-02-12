package main

import (
	"fmt"
	"image/color"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/text"

	"golang.org/x/image/colornames"

	"github.com/faiface/pixel/imdraw"
)

// Drawable provides methods to allow objects to be drawn.
type Drawable interface {
	Draw(bool)
	GetVisual() *imdraw.IMDraw
}

// Particle is the overall object in the world, consisting of a "body"
// for physics-like movement and a graphical representation.
type Particle struct {
	Body

	Radius float64
	Color  color.Color
	visual *imdraw.IMDraw
	dirty  bool // does it need redrawn
}

// NewParticleDefault makes a "default" Particle.
// Has mass and diameter of 1, white color. Other fields are zero.
func NewParticleDefault() *Particle {
	return &Particle{
		Body:   Body{Mass: 1},
		Radius: 1,
		Color:  colornames.White,
		visual: imdraw.New(nil),
		dirty:  true}
}

// NewParticleParams makes a new Particle with the parameters YOU can decide!
func NewParticleParams(x, y, mass, diam float64, color color.Color) *Particle {
	p := NewParticleDefault()
	p.Pos.X = x
	p.Pos.Y = y
	p.Mass = mass
	p.Radius = diam
	p.Color = color
	return p
}

// Draw tells the particle to redraw itself.
func (p Particle) Draw(showVectors bool) {
	if p.dirty {
		p.visual.Reset()
		p.visual.Clear()

		// draw particle
		p.visual.Color = p.Color
		p.visual.Push(p.Pos)
		p.visual.Circle(p.Radius, 0)

		if showVectors {
			// draw velocity vector
			p.visual.Color = colornames.Black
			p.visual.Push(p.Pos, p.Pos.Add(p.Vel))
			p.visual.Line(1)

			// draw force vector
			p.visual.Color = colornames.Black
			p.visual.Push(p.Pos, p.Pos.Add(p.Force))
			p.visual.Line(2)

			txt := text.New(p.Pos.Add(pixel.V(p.Radius+2, 0)), atlas)
			fmt.Fprintf(txt, "F:(%0.1f, %0.1f)", p.Force.X, p.Force.Y)
			txt.Draw(p.visual, pixel.IM)
		}

		// p.dirty = false
	}

}

// GetVisual returns the particle's visual representation.
func (p Particle) GetVisual() *imdraw.IMDraw {
	return p.visual
}
