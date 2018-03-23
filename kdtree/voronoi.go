package main

import (
	"bytes"
	"image"
	"image/color"
	"image/png"
	"io/ioutil"
	"math"
	"sync"

	"github.com/go-gl/mathgl/mgl64"
	"github.com/lucasb-eyer/go-colorful"
	"github.com/quillaja/goutil/rand"
)

// Voronoi creates a Voronoi plot of the given points.
func Voronoi(points []mgl64.Vec2, width, height int, d DistMetric, filename string) {

	// 0. set distance metric
	Dist = d

	// 1. use given point set to build tree and assign a random
	// color to each point
	root := BuildTree(points)
	PreOrderTraversal(root, func(node *Node) {
		node.Data = colorful.Hsv(rand.Float64NM(0, 360), 1, 1)
	})

	// 2. build image pixel by pixel. color pixel based on the
	// nearest neighbor in the tree. Run NN search and image draw in
	// a gofunc to improve performance.
	wg := sync.WaitGroup{}
	img := image.NewRGBA(image.Rect(0, 0, width, height))
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			wg.Add(1)
			go func(x, y int) {
				nn := NearestNeighbor(root, mgl64.Vec2{float64(x), float64(y)})
				img.Set(x, y, nn.Data.(color.Color))
				wg.Done()
			}(x, y)
		}
	}
	wg.Wait()

	// 2.5 draw point set
	for _, p := range points {
		img.Set(int(math.Round(p.X())), int(math.Round(p.Y())), color.Black)
	}

	// 3. write image to disk
	data := []byte{}
	buf := bytes.NewBuffer(data)
	err := png.Encode(buf, img)
	if err != nil {
		panic(err)
	}
	err = ioutil.WriteFile(filename, buf.Bytes(), 0644)
	if err != nil {
		panic(err)
	}
}
