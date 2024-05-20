package plan

func NewPickup(o Point, e Edge) *Pickup {
	return &Pickup{
		Travel: *NewEdge(o, e.Source),
		Work:   e,
	}
}

func (p *Pickup) TotalCost() float64 {
	return p.Travel.Cost() + p.Work.Cost()
}

func (p *Pickup) Efficiency() float64 {
	return p.Work.Cost() / p.TotalCost()
}
