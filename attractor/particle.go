package main

import (
	"image/color"

	"github.com/faiface/pixel"

	"golang.org/x/image/colornames"

	"github.com/faiface/pixel/imdraw"
)

// Drawable provides methods to allow objects to be drawn.
type Drawable interface {
	Draw(bool)
	GetVisual() *imdraw.IMDraw
}

// Particle is the overall object in the world, consisting of a "body"
// for physics-like movement and a graphical representation. One must call
// the "New" methods to initialize private members.
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
		Body:   Body{Mass: 1, MomentOfInertia: 0.5},
		Radius: 1,
		Color:  colornames.White,
		visual: imdraw.New(nil),
		dirty:  true}
}

// NewParticleParams makes a new Particle with the parameters YOU can decide!
func NewParticleParams(x, y, mass, radius float64, color color.Color) *Particle {
	p := NewParticleDefault()
	p.Pos.X = x
	p.Pos.Y = y
	p.Mass = mass
	p.Radius = radius
	p.MomentOfInertia = MomentOfInertia(mass, radius)
	p.Color = color
	return p
}

// GetBody implements the "Bodyer" interface.
func (p *Particle) GetBody() *Body {
	return &p.Body
}

// Draw tells the particle to redraw itself.
func (p *Particle) Draw(showVectors bool) {
	if p.dirty {
		p.visual.Reset()
		p.visual.Clear()

		// set matrix to do rotation
		p.visual.SetMatrix(pixel.IM.Rotated(p.Pos, p.Rotation))

		// draw particle
		p.visual.Color = p.Color
		p.visual.Push(p.Pos)
		p.visual.Circle(p.Radius, 0)
		p.visual.Color = colornames.Black
		p.visual.Push(p.Pos, p.Pos.Add(pixel.V(p.Radius, 0)))
		p.visual.Line(1)

		if showVectors {
			// draw velocity vector
			p.visual.Color = colornames.Black
			p.visual.Push(p.Pos, p.Pos.Add(p.Vel))
			p.visual.Line(1)

			// draw force vector
			p.visual.Color = colornames.Black
			p.visual.Push(p.Pos, p.Pos.Add(p.Force))
			p.visual.Line(2)

			// draw angular velocity
			p.visual.Color = colornames.Black
			rtip := p.Pos.Add(pixel.V(p.Radius, 0))
			p.visual.Push(rtip, rtip.Add(pixel.V(0, p.RotVel)))
			p.visual.Line(1)

			// what i wanted to do, but it wasn't working... pixel said
			// "panic: (*pixel.Batch).MakePicture: Picture is not the Batch's Picture"
			// txt := text.New(pixel.ZV, atlas)
			// fmt.Fprintf(txt, "F:(%0.1f, %0.1f)", p.Force.X, p.Force.Y)
			// txt.Draw(p.visual, pixel.IM.Moved(p.Pos.Add(pixel.V(p.Radius+2, 0))))
		}

		// p.dirty = false
	}

}

// GetVisual returns the particle's visual representation.
func (p *Particle) GetVisual() *imdraw.IMDraw {
	return p.visual
}

// CollidePoint checks if the given pixel.Vec is within the particle's radius.
func (p *Particle) CollidePoint(pt pixel.Vec) bool {
	return p.Pos.To(pt).Len() <= p.Radius
}
