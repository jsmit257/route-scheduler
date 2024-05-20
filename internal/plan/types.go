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

	Edges []Edge

	Pickup struct { // one discrete segment of travel
		Travel, Work Edge
	}

	Driver []*Pickup

	Shift []Driver
)
