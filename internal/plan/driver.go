package plan

import (
	"math"
	"slices"
	"time"
)

func NewDriver(len int) *Driver {
	result := make(Driver, 0, len)
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

func (d *Driver) Minutes() time.Duration {
	return time.Duration(d.TotalCost()) * time.Minute
}

func (d *Driver) Efficiency() float64 {
	var result float64
	for _, p := range *d {
		result += p.Efficiency()
	}
	return result / float64(len(*d)+1)
}

func (d *Driver) Vacancy(p *Pickup) bool {
	d.Push(p)
	err := d.Minutes() > MaxDepth
	// fmt.Fprintf(os.Stderr, "result %5v, test: %d, max: %d\n", !err, d.Minutes(), MaxDepth)
	if err {
		d.Pop()
	}
	return !err
}

func (d *Driver) Pop() *Driver {
	*d = (*d)[:len(*d)-1]
	return d
}

func (d *Driver) Push(p *Pickup) *Driver {
	*d = append(*d, p)
	return d
}

func (d *Driver) FindClosest(edges Edges, head int) []*Pickup {
	pickups := []*Pickup{}
	for _, e := range edges {
		pickups = append(pickups, NewPickup(d.End(), e))
	}
	slices.SortFunc(pickups, func(a, b *Pickup) int {
		return a.MostEfficient(b)
	})
	if lenPickups := len(pickups); lenPickups < head {
		head = lenPickups
	}
	return pickups[:head]
}

func (d *Driver) ReportWork() (result Edges) {
	for _, s := range *d {
		result = append(result, s.Work)
	}
	return result
}

func (d *Driver) Last() *Pickup {
	if l := len(*d) - 1; l < 0 {
		return nil
	} else {
		return (*d)[l]
	}
}

func (d *Driver) End() Point {
	if l := d.Last(); l == nil {
		return Origin
	} else {
		return l.Work.Dest
	}
}
