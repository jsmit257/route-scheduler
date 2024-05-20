package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"os"

	log "github.com/sirupsen/logrus"

	"github.com/jsmit257/route-scheduler/internal/plan"
)

func main() {
	logger := plan.GetLogger().WithField("function", "main")

	logger.Debug("initialized main")

	edges := getEntries(logger)

	logger.WithField("edges", len(edges)).Info("done reading")

	sorted, _ := plan.NewShift(edges)

	report(sorted.Graph(), logger) // stdout for the client

	shiftCost := 0.0
	for _, d := range sorted {
		msg := "driver cost in range"
		if d.Minutes() > plan.MaxDepth {
			msg = "driver cost exceeded range"
		}

		c := d.Minutes()
		shiftCost += c
		logger.WithFields(log.Fields{
			"efficiency": fmt.Sprintf("%.2f%%", d.Efficiency()/(float64(len(sorted)))*100),
			"total_cost": c,
			"pickups":    d.Graph(),
		}).Debug(msg)
	}

	logger.WithFields(log.Fields{
		"shift_cost": shiftCost,
		"total_cost": 500*float64(len(sorted)) + shiftCost,
		"sorted":     sorted.Graph(),
		"count":      len(sorted),
		"size":       len(edges),
	}).Info("done!")
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
	var e = plan.Edge{}

	l = l.WithField("function", "processLines")

	result := make(plan.Edges, 0, plan.MaxDeliveries)

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
				WithField("input", string(line)).
				Warn("scan failed")
		} else if scanned == 0 {
			l.WithField("input", string(line)).Error("couldn't scan tokens")
		} else {
			(&result).Push(*plan.NewEdge(e.Source, e.Dest).SetID(e.ID))
		}

		line, _, err = r.ReadLine()
	}

	if err != io.EOF {
		l.
			WithError(err).
			WithField("filename", os.Args[1]).
			Fatal("unexpected end of input")
	}

	return result
}
