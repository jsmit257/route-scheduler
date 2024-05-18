package plan

import (
	"slices"
	"sync"
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

func (s Shift) Sort(o Point, edges *Edges) (Shift, error) {
	pickups := []*Pickup{}

	for _, e := range *edges {
		pickups = append(pickups, NewPickup(o, e))
	}

	slices.SortFunc(pickups, func(a, b *Pickup) int { return a.MostEfficient(b) })

	for i, d := range s {
		if i == len(pickups) {
			break
		}

		d.Push(pickups[i])

		edges.Remove(edges.Find(pickups[i].Work))
	}

	return s.Balance(edges)
}

func (s Shift) Balance(edges *Edges) (Shift, error) {
	if len(*edges) == 0 {
		return s, nil
	}

	found := false

	var wg sync.WaitGroup
	wg.Add(len(s))

	var edgeLock sync.Mutex
	for _, d := range s {
		d := d
		go func(edges *Edges) {
			defer wg.Done()
			closest := d.FindClosest(*edges, len(s))
			edgeLock.Lock()
			defer edgeLock.Unlock()
			for _, p := range closest {
				if i := edges.Find(p.Work); i != -1 {
					if d.Vacancy(p) {
						found = true
						edges.Remove(i)
						return
					}
				}
			}
		}(edges)
	}

	wg.Wait()

	if !found { // XXX: this shouldn't need to be here
		return s, NoMoreTime
	}

	return s.Balance(edges)
}

func (s Shift) Graph() [][]int {
	result := [][]int{}
	for _, driver := range s {
		result = append(result, driver.Graph())
	}
	return result
}
