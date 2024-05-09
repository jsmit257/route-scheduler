package plan

import "math"

func NewEdge(src, dst Point) *Edge {
	return &Edge{
		Source: src,
		Dest:   dst,
		l:      math.Sqrt(math.Pow(src.X-dst.X, 2) + math.Pow(src.Y-dst.Y, 2)),
	}
}

func (e *Edge) SetID(id int) *Edge {
	e.ID = id
	return e
}

func (e *Edge) Cost() float64 {
	return e.l
}
