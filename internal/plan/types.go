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
		Origin, Work Edge
	}

	Driver struct {
		Segments []*Pickup
		Current  int // for real-time add/delete
	}

	Shift []Driver
)
