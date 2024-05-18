package plan

import (
	"math"
	"slices"
)

func NewEdge(src, dst Point) *Edge {
	return &Edge{
		Source: src,
		Dest:   dst,
		l:      math.Sqrt(math.Pow(src.X-dst.X, 2) + math.Pow(src.Y-dst.Y, 2)),
	}
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
