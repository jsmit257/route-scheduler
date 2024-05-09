package plan

import "time"

var (
	Origin = Point{}
)

const (
	MaxDeliveries = 200
	MaxDepth      = time.Hour * 12
)
