package plan

import (
	"slices"
	// log "github.com/sirupsen/logrus"
)

func NewShift(size int) Shift {
	if size < 1 {
		panic("what's the point of a shift with no drivers?")
	}

	result := make([]*Driver, size)
	for i := range result {
		result[i] = NewDriver(100)
	}
	return result
}

func (s Shift) Graph() [][]int {
	result := [][]int{}
	for _, driver := range s {
		result = append(result, driver.Graph())
	}
	return result
}

func (s Shift) TotalCost() (result float64) {
	for _, d := range s {
		result += d.TotalCost()
	}
	return result
}

func (s Shift) Sort(edges Edges) Shift {
	var recurse func(Shift, Edges) Shift

	recurse = func(shift Shift, tail Edges) Shift {
		var l int
		if l = len(tail); l == 0 {
			return shift
		}

		var shifts []Shift

		d := shift[len(shift)-1]
		for i, t := range tail {
			shift, tail := shift[:], tail[:]
			// for j, s := range shift {
			// 	shift[j] = func(d Driver) *Driver { return &d }(*s)
			// }
			if p := NewPickup(d.End(), t); !d.Vacancy(p) {
				d = NewDriver(l).Push(p)
				shift = append(shift, d)
			}
			shifts = append(shifts, recurse(shift, append(tail[:i], tail[i+1:]...)))
		}

		slices.SortFunc(shifts, func(a, b Shift) int {
			// logger.WithFields(log.Fields{
			// 	"left":  a.TotalCost(),
			// 	"right": b.TotalCost(),
			// }).Debug("comparing")
			if a.TotalCost() < b.TotalCost() {
				return -1
			}
			return 0
		})

		return shifts[0]
	}

	return recurse(s, edges)
}
