package plan

import (
	"fmt"
	"time"

	"github.com/jsmit257/route-scheduler/internal/config"
)

func init() {
	cfg := config.NewConfig()

	Origin = Point{X: cfg.OriginX, Y: cfg.OriginY}

	MaxDepth = time.Hour * time.Duration(cfg.HoursPerShift)
}

var (
	Origin   Point
	MaxDepth time.Duration

	NoMoreTime = fmt.Errorf("all vehicles' schedules are full")
)

const (
	MaxDeliveries = 200
)
