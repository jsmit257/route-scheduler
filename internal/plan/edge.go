package plan

import (
	"math"
	"slices"
	"sync"
	// log "github.com/sirupsen/logrus"
)

func NewEdge(src, dst Point) *Edge {
	return &Edge{
		Source: src,
		Dest:   dst,
		l:      math.Sqrt(math.Pow(src.X-dst.X, 2) + math.Pow(src.Y-dst.Y, 2)),
	}
}

func (e *Edges) nearestneighbor() Shift {
	minshift := e.shortestpath(((*e)[0]))
	for _, edge := range (*e)[1:] {
		if shift := e.shortestpath(edge); shift.TotalCost() < minshift.TotalCost() {
			minshift = shift
		}
	}
	return minshift
}

func (e Edges) shortestpath(first Edge) Shift {
	edges := make(Edges, len(e))
	copy(edges, e)

	result := Shift{*NewDriver()}
	driver := &result[0]
	start := first.Source
	for len(edges) > 0 {
		closest := 0
		for j, next := range edges[1:] {
			if next.ID == first.ID {
				continue
			}
			if edge := *NewEdge(start, next.Source); edge.Cost() < edges[closest].Cost() {
				closest = j + 1
				start = next.Dest
			}
		}

		if !driver.Vacancy(NewPickup(driver.End(), edges[closest])) {
			result = append(result, *NewDriver(NewPickup(Origin, edges[closest])))
			driver = &result[len(result)-1]
		}

		edges = append(edges[:closest], edges[closest+1:]...)
		start = Origin
	}

	return result
}

func (es *Edges) Find(s Edge) int {
	return slices.IndexFunc(*es, func(e Edge) bool {
		return e.ID == s.ID
	})
}

func (es *Edges) Push(e Edge) *Edges {
	*es = append((*es), e)
	return es
}

func (es *Edges) Remove(ndx int) *Edges {
	if l := len(*es); l > ndx {
		*es = append((*es)[:ndx], (*es)[ndx+1:]...)
	}
	return es
}

func (e *Edge) SetID(id int) *Edge {
	e.ID = id
	return e
}

func (e *Edge) Cost() float64 {
	return e.l
}

func (e Edges) Graph() (result []int) {
	for _, t := range e {
		result = append(result, t.ID)
	}
	return result
}

func (e Edges) tryall(fn func(Edges)) {
	var permutate func(Edges, Edges, chan<- Edges)
	permutate = func(head, tail Edges, out chan<- Edges) {
		if len(tail) == 0 {
			out <- head
		}
		for i, t := range tail {
			newTail := make(Edges, len(tail))
			copy(newTail, tail)

			permutate(append(head, t), append(newTail[:i], newTail[i+1:]...), out)
		}
	}

	pipeline := make(chan Edges)

	go func(out chan<- Edges) {
		defer close(out)
		permutate(Edges{}, e, out)
	}(pipeline)

	var doitwait sync.WaitGroup
	doitwait.Add(1)
	go func(out <-chan Edges) {
		defer doitwait.Done()
		for e := range out {
			fn(e)
		}
	}(pipeline)
	doitwait.Wait()
}
