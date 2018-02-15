package main

import (
	"math"

	px "github.com/faiface/pixel"
)

// G is the universal gravitational constant
const G = 300.0
const maxGForce = 1000

// Torque calculates the torque of "force" applied at "loc" at the center of
// "a" (which is assumed to be its Body.Pos). One can apply the torque returned
// from this method using Body.ApplyTorque().
//
// If "force" is not perpendicular to the vector from a.Pos to loc, only a
// portion of the total force is actually converted to a torque.
// One can use TorqueTrans() to find the
// (non-torque) translational component of "force".
func Torque(a Bodyer, force px.Vec, loc px.Vec) float64 {
	// find vector from a.Pos to loc, then the torque is a 1D scalar equal to
	// the radial vector cross force (r x F).
	radial := a.GetBody().Pos.To(loc)
	return radial.Cross(force) // alternative: f.Project(r.Normal()).Scaled(r.Len())
}

// TorqueTrans find the (non-torque) translational component of "force" applied
// at "loc" on "a". If "force" is perpendicular to the vector from a.Pos to loc,
// the translational force produced is (0,0). One can apply the force returned
// from this function using the Body.ApplyForce() method.
func TorqueTrans(a Bodyer, force px.Vec, loc px.Vec) px.Vec {
	radial := a.GetBody().Pos.To(loc)
	return force.Project(radial)
}

// Gravitation calculates the force exerted on "a" by "b". (Hint: the force
// exerted on b by a is the inverse.) One can apply this force by using the
// Body.ApplyForce() method.
func Gravitation(a Bodyer, b Bodyer) px.Vec {
	// F = G * (m1*m2)/|d2-d1|^2 * rhat
	// var fnet px.Vec
	bA := a.GetBody()
	// for _, ob := range others {
	bB := b.GetBody()
	atob := bA.Pos.To(bB.Pos) // vector from a to b
	d := atob.Len()           // d should always be positive
	dsq := d * d
	if dsq <= 1 { // prevent divide by zero
		dsq = 1
	}
	if d <= bA.RepulsorDistance || d <= bB.RepulsorDistance {
		dsq = -dsq
	}
	f := atob.Unit().Scaled((G * bA.Mass * bB.Mass) / dsq)
	if f.Len() > maxGForce {
		f = f.Unit().Scaled(maxGForce)
	}
	// fnet.X += f.X
	// fnet.Y += f.Y
	// }
	return f
}

// Resistance creates a force(s) of the given magnitude(s) in the direction
// opposite the Body's velocity. This is called viscous resistance and is
// calculated using the formula for "Stoke's drag". One can apply this force
// using the Body.ApplyForce() method.
func Resistance(a Bodyer, magitudes ...float64) px.Vec {
	bA := a.GetBody()
	var fnet px.Vec
	for _, m := range magitudes {
		f := bA.Vel.Unit().Scaled(-m)
		fnet.X += f.X
		fnet.Y += f.Y
	}
	return fnet
}

// Bodyer is an interface that all structs containg a "Body" struct should
// implement.
type Bodyer interface {
	GetBody() *Body
}

// Body is the "physics" part of an object.
type Body struct {
	Mass            float64 // kg
	MomentOfInertia float64 // kg*m^2

	Pos   px.Vec // m
	Vel   px.Vec // m/2
	Force px.Vec // N

	Rotation float64 // radians
	RotVel   float64 // rad/sec
	Torque   float64 // Nm

	RepulsorDistance float64 // m
}

// ResetForce does exactly what you think it does.
func (b *Body) ResetForce() {
	b.Force.X, b.Force.Y = 0, 0
}

// ResetTorque does exactly what you think it does.
func (b *Body) ResetTorque() {
	b.Torque = 0
}

// ApplyGravitation calculates the gravitational force between this body and
// the others, accumulating the total force applied in Body.Force.
// func (b *Body) ApplyGravitation(others ...*Body) {
// 	// accumulate forces acting on this body
// 	// ... could use ApplyStatic() here, but I'd have to pack/unpack all the
// 	// force vectors. This is simple enough to just duplicate the calculations.
// 	for _, o := range others {
// 		f := Gravitation(b, o)
// 		b.Force.X += f.X
// 		b.Force.Y += f.Y
// 	}
// }

// ApplyForce applies the static force(s) to Body.Force.
func (b *Body) ApplyForce(forces ...px.Vec) {
	for _, f := range forces {
		b.Force.X += f.X
		b.Force.Y += f.Y
	}
}

// ApplyTorque applies the torque(s) to the Body.Torque.
func (b *Body) ApplyTorque(torques ...float64) {
	for _, t := range torques {
		b.Torque += t
	}
}

// ApplyResistance applies force(s) to the body in the direction *opposite*
// its velocity and with the given magnitudes.
// TODO: check for "real" way of calculating this.
// func (b *Body) ApplyResistance(magnitudes ...float64) {
// 	// As in ApplyGravitation(), just redo calcs.
// 	for _, m := range magnitudes {
// 		f := b.Vel.Unit().Scaled(-m)
// 		b.Force.X += f.X
// 		b.Force.Y += f.Y
// 	}
// }

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

// UpdateRotation tells the object to recalculate it's rotation based upon the
// accumulated total torque (Body.Torque) applied to it in the duration of "dt"
// seconds. Assumes all torques have already been applied.
func (b *Body) UpdateRotation(dt float64) {
	// calculate the new angular acceleration
	// Ia = T, I = (m*r^2)/2, a = (2T)/(m*r^s)
	acc := b.Torque / b.MomentOfInertia
	// calculate new rotational velocity
	// dw = a*dt
	b.RotVel += acc * dt
	// calculate new rotation
	// dr = r + dw
	// if rotation > 2PI, scale back to range [0, 2PI]
	b.Rotation = NormalizeAngle(b.Rotation + b.RotVel*dt)
}

// NormalizeAngle takes an angle in radians and scales it to [-2PI,2PI].
func NormalizeAngle(theta float64) float64 {
	if -2*math.Pi <= theta && theta <= 2*math.Pi {
		return theta
	}
	f := theta / (2 * math.Pi)
	if f < 0 {
		return 2 * math.Pi * (f - math.Ceil(f))
	}
	return 2 * math.Pi * (f - math.Floor(f))
}

// MomentOfInertia calculates a moment of inertia for a disk/circle/wheel,
// which is I = (m*r^2)/2 .
func MomentOfInertia(mass, radius float64) float64 {
	return 0.5 * mass * radius * radius
}
