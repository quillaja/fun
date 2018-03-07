package main

import (
	"flag"
	"fmt"
	stdrand "math/rand"
	"time"

	"github.com/cpmech/gosl/plt"

	"github.com/go-gl/mathgl/mgl64"
	"github.com/quillaja/goutil/rand"
)

const (
	min = 0
	max = 100
)

// construct kd tree then plot points and lines on chart
func main() {
	nPts := flag.Int("n", 10, "Number of points.")
	flag.Parse()

	//generate random points
	stdrand.Seed(time.Now().UnixNano())
	points := []mgl64.Vec2{}
	for n := 0; n < *nPts; n++ {
		points = append(points, mgl64.Vec2{
			float64(rand.IntNM(min, max)),
			float64(rand.IntNM(min, max))})
	}
	// points = []mgl64.Vec2{{2, 3}, {5, 4}, {9, 6}, {4, 7}, {8, 1}, {7, 2}} // test data

	start := time.Now()
	root := BuildTree(points) // make tree
	fmt.Println("tree build time (ms):", time.Since(start).Seconds()*1000)

	searchpt := mgl64.Vec2{50, 50}
	fmt.Println("nearest neighbor to", searchpt, "is", NearestNeighbor(root, searchpt).Data)

	// plot styles
	ptStyle := &plt.A{C: "#000000", M: "."}
	vStyle := &plt.A{C: "#FF0000"}
	hStyle := &plt.A{C: "#0000FF"}

	action := func(node *Node) {
		// print vertical line for X-median node,
		// and horizontal line for Y-median node.
		if node.Axis == 0 {
			plt.Polyline([][]float64{
				{node.Data.X(), node.Range[1][0]},
				{node.Data.X(), node.Range[1][1]}},
				vStyle)
		} else {
			plt.Polyline([][]float64{
				{node.Range[0][0], node.Data.Y()},
				{node.Range[0][1], node.Data.Y()}},
				hStyle)
		}
		// plot the point
		plt.PlotOne(node.Data.X(), node.Data.Y(), ptStyle)
	}

	PreOrderDFS(root, action) // correct way to print
	// InOrderDFS(root, action)
	// PostOrderDFS(root, action)

	// fmt.Println(points)
	plt.Show() // blocks while window is open
}
