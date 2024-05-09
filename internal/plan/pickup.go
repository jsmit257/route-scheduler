package plan

func NewPickup(o Point, e Edge) *Pickup {
	return &Pickup{
		Unbillable: *NewEdge(o, e.Source),
		Work:       e,
	}
}

func (p *Pickup) TotalCost() float64 {
	return p.Unbillable.Cost() + p.Work.Cost()
}

func (p *Pickup) Efficiency() float64 {
	return p.Work.Cost() / p.TotalCost()
}

func (p *Pickup) MostEfficient(o *Pickup) int {
	if l, r := p.Efficiency(), o.Efficiency(); l == r {
		return 0
	} else if l > r {
		return -1
	} else {
		return 1
	}
}
