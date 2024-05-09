package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strconv"

	log "github.com/sirupsen/logrus"

	"github.com/jsmit257/route-scheduler/internal/plan"
)

func main() {
	var err error
	var fleetSize int

	log.SetLevel(log.DebugLevel) // TODO: grab this from the config
	log.SetOutput(os.Stderr)
	log.SetFormatter(&log.JSONFormatter{})
	l := log.WithFields(log.Fields{
		"app":      "route-scheduler",
		"function": "main",
	})

	if trucks := os.Getenv("FLEET_SIZE"); trucks == "" {
		fleetSize = 12
	} else if fleetSize, err = strconv.Atoi(trucks); err != nil {
		l.WithError(err).Fatalf("FLEET_SIZE=%s cannot be converted to an int", trucks)
	}

	l.WithField("drivers", fleetSize).Debug("fleet")

	edges := getEntries(l)

	l.WithField("edges", edges).Debug("finished reading")

	for _, path := range plan.NewShift(fleetSize).Sort(plan.Origin, edges).Graph() {
		if text, err := json.Marshal(path); err != nil {
			l.WithError(err).Error("json failed")
		} else {
			fmt.Printf("%s\n", text)
		}
	}
}

func getEntries(l *log.Entry) []plan.Edge {
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

func processLines(r *bufio.Reader, l *log.Entry) []plan.Edge {
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
