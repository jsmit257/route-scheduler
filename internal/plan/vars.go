package plan

import (
	"fmt"
	"os"
	"strings"

	log "github.com/sirupsen/logrus"

	"github.com/jsmit257/route-scheduler/internal/config"
)

func init() {
	log.SetOutput(os.Stderr)
	log.SetFormatter(&log.JSONFormatter{})
	if logLevel, ok := map[string]log.Level{
		"trace": log.TraceLevel,
		"debug": log.DebugLevel,
		"info":  log.InfoLevel,
		"warn":  log.WarnLevel,
		"error": log.ErrorLevel,
		"fatal": log.FatalLevel,
		"panic": log.PanicLevel,
	}[strings.ToLower(cfg.LogLevel)]; ok {
		log.SetLevel(logLevel)
	} else {
		log.SetLevel(log.DebugLevel)
		log.
			WithField("app", "route-scheduler").
			Infof("no valid value found for LOG_LEVEL(%s)", cfg.LogLevel)
	}

	logger = log.WithField("app", "route-scheduler")
}

var (
	Origin   = Point{X: cfg.OriginX, Y: cfg.OriginY}
	MaxDepth = cfg.HoursPerShift * 60

	cfg    = config.NewConfig()
	logger *log.Entry

	NoMoreTime = fmt.Errorf("all vehicles' schedules are full")
)

const (
	MaxDeliveries = 200
)

func GetLogger() *log.Entry {
	return logger
}

func GetConfig() *config.Config {
	return cfg
}
