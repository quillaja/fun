package main

import (
	"math"

	px "github.com/faiface/pixel"
)

// G is the universal gravitational constant
const G = 300.0
const maxGForce = 1000

// Torque calculates the components of "force" applied at "loc" which becomes
// a torque on "a" as well as the non-torque translation-producing remainder.
// return value is a [2]pixel.Vec where [0] is the torque on the center of
// "a" (in Nm or a similar unit) and [1] is the translational remainder (N, etc).
func Torque(a Bodyer, force px.Vec, loc px.Vec) (forces [2]px.Vec) {
	bA := a.GetBody()
	// find vector from a.Pos to loc, then find component (projection) of
	// "force" normal to said vector and the component parallel.
	aToF := bA.Pos.To(loc)
	norm := aToF.Normal()
	forces[0] = force.Project(norm)
	forces[1] = force.Project(aToF)
	return
}

// Gravitation calculates the force exerted on "a" by "b". (Hint: the force
// exerted on b by a is the inverse.)
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
// calculated using the formula for "Stoke's drag".
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
	Torque   px.Vec  // Nm

	RepulsorDistance float64 // m
}

// ResetForce does exactly what you think it does.
func (b *Body) ResetForce() {
	b.Force.X, b.Force.Y = 0, 0
}

// ResetTorque does exactly what you think it does.
func (b *Body) ResetTorque() {
	b.Torque.X, b.Torque.Y = 0, 0
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
func (b *Body) ApplyTorque(torques ...px.Vec) {
	for _, t := range torques {
		b.Torque.Y += t.Y
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
	acc := b.Torque.Y / b.MomentOfInertia
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
