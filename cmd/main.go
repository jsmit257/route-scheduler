package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/jsmit257/route-scheduler/internal/config"
	"github.com/jsmit257/route-scheduler/internal/plan"
)

func main() {
	var err error

	cfg := config.NewConfig()

	logger := initLogger(cfg).WithField("function", "main")

	logger.WithField("cfg", cfg).Debug("initialized config and logger")

	edges := getEntries(logger)

	logger.WithField("edges", len(edges)).Debug("done reading")

	// mapReduce(edges, logger)

	var sorted plan.Shift
	sorted, err = plan.NewShift(cfg.FleetSize).Sort(plan.Origin, &edges)
	if err != nil {
		logger.
			WithError(err).
			WithField("edges-remaining", len(edges)).
			Fatal("sorting failed")
	}

	report(sorted.Graph(), logger) // stdout for the client

	shiftCost := 0.0
	for _, d := range sorted {
		msg := "driver cost in range"
		if d.Minutes() > plan.MaxDepth {
			msg = "driver cost exceeded range"
		}

		c := d.TotalCost()
		shiftCost += c
		logger.
			WithFields(log.Fields{
				"efficiency": fmt.Sprintf("%.2f%%", d.Efficiency()*100),
				"total_cost": time.Duration(c) * time.Nanosecond, // actually, it's minutes
				"pickups":    d.Graph(),
			}).
			Debug(msg)
	}

	logger.
		WithFields(log.Fields{
			"shift_cost": shiftCost,
			"total_cost": 500*float64(cfg.FleetSize) + shiftCost,
		}).
		Info("done!")
}

func report(graph [][]int, l *log.Entry) {
	l = l.WithField("function", "report")
	for _, path := range graph {
		if len(path) == 0 {
			continue
		}
		if text, err := json.Marshal(path); err != nil {
			l.WithError(err).Error("json failed")
		} else {
			fmt.Printf("%s\n", text)
		}
	}
}

func getEntries(l *log.Entry) plan.Edges {
	var err error

	l = l.WithField("function", "getReader")

	var f *os.File
	if len(os.Args) < 2 {
		f = os.Stdin
		l.
			WithError(fmt.Errorf("using stdin for input")).
			Warnf("usage: %s 'path-to-manifest'", os.Args[0])
	} else if f, err = os.Open(os.Args[1]); err != nil {
		l.
			WithField("filename", os.Args[1]).
			WithError(err).
			Fatal("failed to open file")
	}
	defer f.Close()

	return processLines(bufio.NewReader(f), l)
}

func processLines(r *bufio.Reader, l *log.Entry) plan.Edges {
	var err error
	// var line []byte
	var e = plan.Edge{}

	l = l.WithField("function", processLines)

	result := make([]plan.Edge, 0, plan.MaxDeliveries)

	line, _, err := r.ReadLine()
	for err == nil {

		if scanned, err := fmt.Sscanf(string(line), "%d (%f,%f) (%f,%f)",
			&e.ID,
			&e.Source.X,
			&e.Source.Y,
			&e.Dest.X,
			&e.Dest.Y,
		); err != nil {
			l.
				WithError(err).
				WithField("line", string(line)).
				Error("scan failed")
		} else if scanned == 0 {
			l.WithField("input", string(line)).Error("couldn't scan tokens")
		} else {
			result = append(result, *plan.NewEdge(e.Source, e.Dest).SetID(e.ID))
		}

		line, _, err = r.ReadLine()
	}

	if err != io.EOF {
		l.WithError(err).Fatalf("unexpected end of input: '%s'", os.Args[1])
	}

	return result
}

func initLogger(cfg *config.Config) *log.Entry {
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
	}
	log.SetOutput(os.Stderr)
	log.SetFormatter(&log.JSONFormatter{})

	return log.WithField("app", "route-scheduler")
}
