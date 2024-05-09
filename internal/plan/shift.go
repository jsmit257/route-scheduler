package plan

import (
	"slices"
)

func NewShift(size int) Shift {
	if size < 1 {
		panic("what's the point of a shift with no drivers?")
	}

	result := make([]Driver, size)
	for _, d := range result {
		d.Segments = make([]*Pickup, 100)
	}

	return result
}

func (s Shift) Sort(o Point, edges []Edge) Shift {
	segments := []*Pickup{}

	for _, e := range edges {
		segments = append(segments, NewPickup(o, e))
	}

	slices.SortFunc(segments, func(a, b *Pickup) int { return a.CompareTo(b) })

	for i, d := range s {
		if i == len(segments) {
			break
		}
		s[i].Segments = append(d.Segments, segments[i])
	}

	return s.Balance(edges[:len(s)])
}

func (s Shift) Balance([]Edge) Shift {
	return s
}

func (s Shift) Graph() [][]int {
	result := [][]int{}
	for _, driver := range s {
		result = append(result, driver.Graph())
	}
	return result
}
