package main

import (
	"github.com/g3n/engine/core"
	"github.com/g3n/engine/geometry"
	"github.com/g3n/engine/graphic"
	"github.com/g3n/engine/light"
	"github.com/g3n/engine/material"
	"github.com/g3n/engine/math32"
	"github.com/g3n/engine/util/application"
	"github.com/g3n/engine/window"
)

var cubeMaterial *material.Standard = material.NewStandard(math32.NewColor("red"))

type Cube struct {
	math32.Vector3
	Side float64
}

func (c *Cube) Subdivide() (cubes []*Cube) {
	cubes = make([]*Cube, 0, 20) //27 total - 7 removed
	for x := float32(-1.0); x <= 1.0; x++ {
		for y := float32(-1.0); y <= 1.0; y++ {
			for z := float32(-1.0); z <= 1.0; z++ {
				if math32.Abs(x)+math32.Abs(y)+math32.Abs(z) >= 2 {
					// make a new smaller cube
					side := float32(c.Side / 3)
					mini := &Cube{
						Side: float64(side),
						Vector3: math32.Vector3{
							X: c.X + x*side,
							Y: c.Y + y*side,
							Z: c.Z + z*side,
						}}
					cubes = append(cubes, mini)
				}
			}
		}
	}
	return
}

func (c *Cube) Draw(scene *core.Node) {
	geom := geometry.NewBox(float32(c.Side), float32(c.Side), float32(c.Side))
	box := graphic.NewMesh(geom, cubeMaterial)
	box.SetPositionVec(&c.Vector3)
	scene.Add(box)
}

func main() {
	app, err := application.Create(application.Options{
		Width:     800,
		Height:    600,
		Title:     "Sponge",
		TargetFPS: 60,
	})
	if err != nil {
		panic(err)
	}

	cubes := []*Cube{&Cube{
		Side:    10,
		Vector3: math32.Vector3{0, 0, 0},
	}}
	cubeGroup := core.NewNode()
	for _, c := range cubes {
		c.Draw(cubeGroup)
	}
	app.Scene().Add(cubeGroup)

	dirLight := light.NewDirectional(math32.NewColor("white"), 0.7)
	dirLight.SetPosition(1, 1, 1)
	app.Scene().Add(dirLight)
	dirLight2 := light.NewDirectional(math32.NewColor("white"), 0.5)
	dirLight2.SetPosition(-1, 1, 0)
	app.Scene().Add(dirLight2)
	ambLight := light.NewAmbient(math32.NewColor("white"), 0.2)
	app.Scene().Add(ambLight)

	ah := graphic.NewAxisHelper(20)
	app.Scene().Add(ah)

	app.CameraPersp().SetPositionZ(15)

	divideCount := 0
	app.Window().Subscribe(window.OnKeyDown, func(name string, ev interface{}) {
		kev, ok := ev.(*window.KeyEvent)
		if !ok {
			panic("it's not ok")
		}
		if kev.Keycode == window.KeySpace {
			if divideCount <= 2 {
				app.Log().Info("subdividing")
				cubeGroup.RemoveAll(true)
				newCubes := []*Cube{}
				for _, c := range cubes {
					for _, mini := range c.Subdivide() {
						mini.Draw(cubeGroup)
						newCubes = append(newCubes, mini)
					}
				}
				cubes = newCubes
				divideCount++
			} else {
				app.Log().Info("Can't divide more.")
			}
		}
	})

	app.Run()
}
