package plan

import (
	"math"
	"slices"
)

func NewDriver(p ...*Pickup) *Driver {
	result := Driver(p)
	return &result
}

func (d *Driver) Graph() []int {
	result := []int{}
	for _, s := range *d {
		result = append(result, s.Work.ID)
	}
	return result
}

func (d *Driver) TotalCost() float64 {
	last := d.Last()
	if last == nil {
		return 0
	}

	result := math.Sqrt(math.Pow(last.Work.Dest.X-Origin.X, 2) + math.Pow(last.Work.Dest.Y-Origin.Y, 2))
	for _, p := range *d {
		result += p.TotalCost()
	}
	return result
}

func (d *Driver) Minutes() float64 {
	return d.TotalCost()
}

func (d *Driver) Efficiency() float64 {
	var result float64
	for _, p := range *d {
		result += p.Efficiency()
	}
	return result / float64(len(*d)+1)
}

func (d *Driver) Vacancy(p *Pickup) bool {
	d = d.Push(p)
	err := d.Minutes() > MaxDepth
	if err {
		d.Pop()
	}
	return !err
}

func (d *Driver) Pop() *Driver {
	if l := len(*d); l != 0 {
		*d = (*d)[:l-1]
	}
	return d
}

func (d *Driver) Push(p *Pickup) *Driver {
	*d = append(*d, p)
	return d
}

func (d *Driver) FindClosest(edges Edges, head int) (pickups []*Pickup) {
	for _, e := range edges {
		pickups = append(pickups, NewPickup(d.End(), e))
	}
	slices.SortFunc(pickups, func(a, b *Pickup) int {
		if l, r := a.Efficiency(), b.Efficiency(); l == r {
			return 0
		} else if l > r {
			return -1
		} else {
			return 1
		}
	})
	if lenPickups := len(pickups); lenPickups < head {
		head = lenPickups
	}
	return pickups[:head]
}

func (d *Driver) Last() *Pickup {
	if l := len(*d); l == 0 {
		return nil
	} else {
		return (*d)[l-1]
	}
}

func (d *Driver) End() Point {
	if l := d.Last(); l == nil {
		return Origin
	} else {
		return l.Work.Dest
	}
}
