package main

import (
	"fmt"
	"math"
	"sort"

	"github.com/go-gl/mathgl/mgl64"
)

// Node is a node in the kdtree
type Node struct {
	Axis   int
	Range  []mgl64.Vec2
	Data   mgl64.Vec2
	Left   *Node
	Right  *Node
	Parent *Node
}

// IsLeaf says if the node is a leaf (has no children) or not
func (n *Node) IsLeaf() bool {
	return n.Left == nil && n.Right == nil
}

// BuildTree makes the kd tree and returns its root node. The order of "items"
// will not be preserved.
func BuildTree(items []mgl64.Vec2) (root *Node) {
	return buildTree(items, 0, nil, []mgl64.Vec2{{min, max}, {min, max}})
}

// does actual tree build
func buildTree(items []mgl64.Vec2, depth int, parent *Node, rng []mgl64.Vec2) (node *Node) {
	if len(items) == 0 {
		return nil
	}

	// ascending sort items by axis
	axis := depth % 2 // 0=x, 1=y
	sort.Slice(items, func(i, j int) bool {
		return items[i][axis] < items[j][axis]
	})

	// create node
	median := len(items) / 2
	n := &Node{
		Axis:   axis,
		Range:  rng,
		Data:   items[median],
		Parent: parent}

	// create the "ranges" for the left and right children.
	// the ranges are used for plotting.
	var l, r []mgl64.Vec2
	if axis == 0 {
		// split horizontal range
		l = []mgl64.Vec2{{rng[0][0], n.Data[axis]}, rng[1]}
		r = []mgl64.Vec2{{n.Data[axis], rng[0][1]}, rng[1]}
	} else {
		// split vertical range
		l = []mgl64.Vec2{rng[0], {rng[1][0], n.Data[axis]}}
		r = []mgl64.Vec2{rng[0], {n.Data[axis], rng[1][1]}}

	}

	n.Left = buildTree(items[:median], depth+1, n, l)
	n.Right = buildTree(items[median+1:], depth+1, n, r)

	return n
}

// used in nearest neighbor searches for best candidate(s)
type neigh struct {
	node *Node
	dist float64
}

// NearestNeighbor finds the nearest neighbor to searchPt. Returns
// nil if none found (shouldn't do that).
func NearestNeighbor(root *Node, searchPt mgl64.Vec2) *Node {
	return nnSearch(root, searchPt, &neigh{nil, math.Inf(0)}).node
}

// does actual search algorithm
func nnSearch(root *Node, searchPt mgl64.Vec2, curBest *neigh) (newBest *neigh) {

	// if the current node is nil, then just return the current bests
	if root == nil {
		return curBest
	}
	fmt.Println("examining", root.Data, "current best", curBest)

	// check if current node is better than current best
	// if current best == nil/inf, set current node to best
	if dist := distSq(root.Data, searchPt); curBest.node == nil || dist < curBest.dist {
		curBest.node = root
		curBest.dist = dist
	}

	// check if points could possibly exist on the other side of root's splitting
	// axis by checking if the distace from the searchPt to axis is less than
	// the distance to the current best.
	// search-to-plane = abs(root.Data[axis] - search[axis])
	// if search-to-plane <= curbest_dist then go down both branches.
	// else choose the correct branch.
	if math.Pow(root.Data[root.Axis]-searchPt[root.Axis], 2) < curBest.dist {
		newBest = nnSearch(root.Left, searchPt, curBest)
		newBest = nnSearch(root.Right, searchPt, curBest)
	} else {
		// go down left or right branch based on axial comparison to current node.
		if searchPt[root.Axis] < root.Data[root.Axis] {
			newBest = nnSearch(root.Left, searchPt, curBest)
		} else {
			newBest = nnSearch(root.Right, searchPt, curBest)
		}
	}

	return
}

// NearestKNeighbors returns the nearest [0,k] neighbors to the search point.
// If fewer than k are found, the returned slice of nodes will be as long as
// the number found.
func NearestKNeighbors(root *Node, k int, searchPt mgl64.Vec2) (nodes []*Node) {
	bests := make([]*neigh, k)
	knnSearch(root, searchPt, bests) // will alter bests
	for _, b := range bests {
		nodes = append(nodes, b.node)
	}
	return
}

// does nn search for k nodes
// curBests is a best-worst ORDERED list of k elements (some of which may be nil)
func knnSearch(root *Node, searchPt mgl64.Vec2, curBests []*neigh) {

	if root == nil {
		return
	}

	dist := distSq(root.Data, searchPt)
	for i := 0; i < len(curBests); i++ {
		// check each. if found a best.dist > root.dist, insert
		// to keep order, and remove the worst best from the end.
		// if nil is encountered, insert.
		if curBests[i] == nil {
			curBests[i] = &neigh{root, dist}
		} else if dist < curBests[i].dist {
			insertAndTrim(&neigh{root, dist}, i, curBests)
		}
	}

	// go down branches

	return
}

func insertAndTrim(item *neigh, at int, s []*neigh) {
	// insert
	s = append(s, nil)
	copy(s[at+1:], s[at:])
	s[at] = item

	// remove end
	s[len(s)-1] = nil
	s = s[:len(s)-1]
}

// finds distance squared between a and b
func distSq(a, b mgl64.Vec2) float64 {
	delta := b.Sub(a)
	return delta[0]*delta[0] + delta[1]*delta[1]
}

// PreOrderTraversal traverses the tree in a depth-first manner, performing
// "action" on the node before visiting children.
func PreOrderTraversal(root *Node, action func(node *Node)) {
	action(root)

	if root.Left != nil {
		PreOrderTraversal(root.Left, action)
	}

	if root.Right != nil {
		PreOrderTraversal(root.Right, action)
	}
}

// InOrderTraversal traverses the tree in a depth-first manner, visiting the
// left child, then performing "action" on the node, then visiting the right child.
func InOrderTraversal(root *Node, action func(node *Node)) {
	if root.Left != nil {
		InOrderTraversal(root.Left, action)
	}

	action(root)

	if root.Right != nil {
		InOrderTraversal(root.Right, action)
	}

}

// PostOrderTraversal traverses the tree in a depth-first manner, first visitng
// the children, then performing "action" on the node.
func PostOrderTraversal(root *Node, action func(node *Node)) {
	if root.Left != nil {
		PostOrderTraversal(root.Left, action)
	}

	if root.Right != nil {
		PostOrderTraversal(root.Right, action)
	}

	action(root)
}
