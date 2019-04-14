package main

import (
	"fmt"
	"math"
	"math/rand"
	"time"
)

func main() {
	// a := []int{5, 1, 10, 3, 4, 2, 0} // 5 units
	// a := []int{1, 2, 3, 4, 5, 4, 3, 2, 1} // 0 units
	// a := []int{1, 0, 1} // 1
	// a := []int{1, 0, 1, 0, 1, 0, 1, 0, 1} // 4
	// a := []int{5, 0, 2, -50, 2, 1}
	// a := []int{0, 0, 4, 0, 0, 6, 0, 0, 3, 0, 5, 0, 1, 0, 0, 0} // 26

	rand.Seed(time.Now().UnixNano())
	// a := rand.Perm(16000)
	a := createArray(0, 10, 20)

	fmt.Println(a)
	fmt.Println("ans", fillVolume(a))
}

func createArray(minHeight, maxHeight, length int) []int {
	array := make([]int, length)
	for h := maxHeight; h >= minHeight; h-- {
		for n := rand.Intn(length / 2); n >= 0; n-- {
			i := rand.Intn(length)
			array[i] = h
		}
	}
	return array
}

type span struct {
	lo, hi int
}

func (s span) width() int {
	return s.hi - s.lo
}

func fillVolume(a []int) (volume int) {
	if len(a) == 0 {
		return
	}

	s := span{0, len(a) - 1}
	inner := searchInner(a, s)
	volume += calcVolume(a, inner)

	fmt.Println("c", inner)
	for left := searchLeft(a, inner); left.lo >= 0; left = searchLeft(a, left) {
		volume += calcVolume(a, left)
		fmt.Println("l", left)
	}

	for right := searchRight(a, inner); right.hi <= len(a)-1; right = searchRight(a, right) {
		volume += calcVolume(a, right)
		fmt.Println("r", right)
	}

	return
}

// searchInner finds the left and right most indicies of the maximum value WITHIN
// the 's' span of array 'a'.
func searchInner(a []int, s span) span {
	max := math.MinInt64
	var sub span

	for i := s.lo; i <= s.hi; i++ {
		// reset max 'stuff'
		if a[i] > max {
			max = a[i]
			sub.lo = i
		}

		// set/advance the right-most max index
		if a[i] == max {
			sub.hi = i
		}
	}

	return sub
}

// searchLeft finds the index of the left most maximum value LEFT of
// span 's' in 'a'.
func searchLeft(a []int, s span) (left span) {
	max := math.MinInt64
	left.lo, left.hi = s.lo-1, s.lo-1

	for i := s.lo - 1; i >= 0; i-- {
		if a[i] > max {
			max = a[i]
			left.lo = i
		}

		if a[i] == max {
			left.lo = i
		}
	}
	return
}

// searchRight finds the index of the right most maximum value RIGHT of
// span 's' in 'a'.
func searchRight(a []int, s span) (right span) {
	max := math.MinInt64
	right.lo, right.hi = s.hi+1, s.hi+1

	for i := s.hi + 1; i < len(a); i++ {
		if a[i] > max {
			max = a[i]
			right.hi = i
		}

		if a[i] == max {
			right.hi = i
		}
	}

	return
}

// calcVolume calculates the volume enclosed by span 's', assuming the value
// of one of s.lo or s.hi is the maximum in the span.
func calcVolume(a []int, s span) (vol int) {
	max := a[s.lo]
	if a[s.hi] > a[s.lo] {
		max = a[s.hi]
	}

	for i := s.lo; i <= s.hi; i++ {
		vol += max - a[i]
	}
	return
}
