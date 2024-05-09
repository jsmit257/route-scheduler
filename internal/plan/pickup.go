package plan

func NewPickup(o Point, e Edge) *Pickup {
	return &Pickup{
		Origin: *NewEdge(o, e.Source),
		Work:   e,
	}
}

func (p *Pickup) TotalCost() float64 {
	return p.Origin.Cost() + p.Work.Cost()
}

func (p *Pickup) Efficiency() float64 {
	return p.Work.Cost() / p.TotalCost()
}

func (p *Pickup) CompareTo(o *Pickup) int {
	if l, r := p.TotalCost(), o.TotalCost(); l == r {
		return 0
	} else if l < r {
		return -1
	} else {
		return 1
	}
}
