package plan

import (
	"time"
)

func (d *Driver) Graph() []int {
	result := []int{}
	for _, s := range d.Segments {
		result = append(result, s.Work.ID)
	}
	return result
}

func (d *Driver) Efficiency() float64 {
	var result float64
	for _, s := range d.Segments {
		result += s.Efficiency()
	}
	return result
}

// func (d *Driver) Depth() time.Duration {
// 	result := NewEdge(Origin, d.Segments[0].Source).Cost()

// 	for _, s := range d.Segments {
// 		result += s.Cost()
// 	}

// 	return time.Duration(result) // DIXME: needs to be adjusted by something
// }

func (d *Driver) Location(when time.Duration) float64 {
	return 0
}

func (d *Driver) Last() bool {
	if d.Current == -1 {
		return true
	}
	return (d.Current + 1) < len(d.Segments)
}

func (d *Driver) Done() bool {
	return d.Current == -1
}
