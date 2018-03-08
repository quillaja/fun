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
	max = 800
)

// construct kd tree then plot points and lines on chart
func main() {
	nPts := flag.Int("n", 10, "Number of points.")
	x := flag.Float64("x", 0, "x-coord of search point.")
	y := flag.Float64("y", 0, "y-coord of search point.")
	k := flag.Int("k", 1, "number of neighbors to find")
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

	searchpt := mgl64.Vec2{*x, *y}
	start = time.Now()
	if *k == 1 {
		fmt.Println("nearest neighbor to", searchpt, "is", NearestNeighbor(root, searchpt).Data)
		fmt.Println("took (ms):", time.Since(start).Seconds()*1000)
	} else if *k > 1 {
		fmt.Println("the", *k, "nearest neighbors to", searchpt, "are:")
		result := NearestKNeighbors(root, *k, searchpt)
		fmt.Println("took (ms):", time.Since(start).Seconds()*1000)
		for _, n := range result {
			fmt.Println(n.Data)
		}

	} else {
		fmt.Println("k = ", *k)
	}

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

	PreOrderTraversal(root, action) // correct way to print tree

	plt.PlotOne(searchpt.X(), searchpt.Y(), &plt.A{C: "#00FF00", M: "x"}) // plot search pt

	plt.Show() // blocks while window is open
}
