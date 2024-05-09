package plan

type (
	Point struct {
		X, Y float64
	}

	Edge struct {
		ID           int
		Source, Dest Point
		l            float64 // expensive math, better to WORM
	}

	Pickup struct { // one discrete segment of travel
		Unbillable, Work Edge
	}

	Driver []*Pickup

	Edges []Edge

	Shift []*Driver
)
