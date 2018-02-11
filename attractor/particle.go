package main

import (
	"image/color"

	"golang.org/x/image/colornames"

	"github.com/faiface/pixel/imdraw"

	px "github.com/faiface/pixel"
)

// G is the universal gravitational constant
const G = 300.0
const maxGForce = 400

// Gravitation calculates the force exerted on "a" by "b".
func Gravitation(a, b Body) px.Vec {
	// F = G * (m1*m2)/|d2-d1|^2 * rhat
	atob := a.Pos.To(b.Pos) // vector from a to b
	d := atob.Len()         // d should always be positive
	dsq := d * d
	if dsq <= 1 { // prevent divide by zero
		dsq = 1
	}
	if d <= a.RepulsorDistance || d <= b.RepulsorDistance {
		dsq = -dsq
	}
	f := atob.Unit().Scaled((G * a.Mass * b.Mass) / dsq)
	if f.Len() > maxGForce {
		f = f.Unit().Scaled(maxGForce)
	}
	return f
}

// Drawable provides methods to allow objects to be drawn.
type Drawable interface {
	Draw(bool)
	GetVisual() *imdraw.IMDraw
}

// Body is the "physics" part of an object.
type Body struct {
	Pos              px.Vec
	Vel              px.Vec
	Force            px.Vec
	Mass             float64
	RepulsorDistance float64
}

// ResetForce does exactly what you think it does.
func (b *Body) ResetForce() {
	b.Force.X, b.Force.Y = 0, 0
}

// ApplyGravitation calculates the gravitational force between this body and
// the others, accumulating the total force applied in Body.Force.
func (b *Body) ApplyGravitation(others ...Body) {
	// accumulate forces acting on this body
	// ... could use ApplyStatic() here, but I'd have to pack/unpack all the
	// force vectors. This is simple enough to just duplicate the calculations.
	for _, o := range others {
		f := Gravitation(*b, o)
		b.Force.X += f.X
		b.Force.Y += f.Y
	}
}

// ApplyStatic applies the static force(s) to Body.Force.
func (b *Body) ApplyStatic(forces ...px.Vec) {
	for _, f := range forces {
		b.Force.X += f.X
		b.Force.Y += f.Y
	}
}

// ApplyResistance applies force(s) to the body in the direction *opposite*
// its velocity and with the given magnitudes.
// TODO: check for "real" way of calculating this.
func (b *Body) ApplyResistance(magnitudes ...float64) {
	// As in ApplyGravitation(), just redo calcs.
	for _, m := range magnitudes {
		f := b.Vel.Unit().Scaled(-m)
		b.Force.X += f.X
		b.Force.Y += f.Y
	}
}

// UpdatePosition tells the object to recalculate its position based on the
// accumlated total force (Body.Force) applied to it in the duration
// of "dt" seconds. Assumes all forces have already been applied.
func (b *Body) UpdatePosition(dt float64) {
	// calculate new acceleration
	// F=ma, a = F/m
	acc := b.Force.Scaled(1 / b.Mass)
	// calculate new velocity
	// dv = a*dt
	b.Vel.X += acc.X * dt
	b.Vel.Y += acc.Y * dt
	// calculate new position
	// dp = p + dv
	b.Pos.X += b.Vel.X * dt
	b.Pos.Y += b.Vel.Y * dt
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
			p.visual.Color = colornames.Green
			p.visual.Push(p.Pos, p.Pos.Add(p.Vel))
			p.visual.Line(1)

			// draw force vector
			p.visual.Color = colornames.Black
			p.visual.Push(p.Pos, p.Pos.Add(p.Force))
			p.visual.Line(2)
		}

		// p.dirty = false
	}

}

// GetVisual returns the particle's visual representation.
func (p Particle) GetVisual() *imdraw.IMDraw {
	return p.visual
}
