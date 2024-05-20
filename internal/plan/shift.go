package plan

import (
	"fmt"
	"slices"
	"sync"

	log "github.com/sirupsen/logrus"
)

var (
	mincost float64
	minlock sync.Mutex
)

func NewShift(edges Edges) (Shift, error) {

	if Debug {
		count := 0
		edges[:7].tryall(func(e Edges) {
			count++
			logger.WithField("edges", e.Graph()).Trace("default function")
		})
		logger.WithField("total edges", count).Trace("done reading")
	}

	nearest := edges.nearestneighbor()

	if len(edges) > 5 {
		return nearest, nil
	}

	// FIXME: combinatorials are too time-consuming; 200! is a lot of combinations;
	//        works great for less-than 10 though
	mincost = nearest.TotalCost()
	logger.WithFields(log.Fields{
		"mincost": mincost,
		"shift":   nearest.Graph(),
		"count":   len(edges),
	}).Debug("minimum for nearest neighbors")
	return Shift{*NewDriver()}.balance(edges, 1)
}

func (s Shift) balance(tail Edges, threads int) (Shift, error) {
	if l := len(tail); l == 0 {
		minlock.Lock()
		if mycost := s.TotalCost(); mycost < mincost {
			mincost = mycost
		}
		minlock.Unlock()
		return s, nil
	}

	_ = logger.WithFields(log.Fields{"method": "balance"})

	var shifts []Shift
	shiftlock := sync.Mutex{}
	waitlock := sync.WaitGroup{}
	var maxgoroutines = make(chan struct{}, 2)

	for i, head := range tail {
		shift, head, newTail := s.Clone(), head, make(Edges, len(tail))
		copy(newTail, tail)

		newTail = append(newTail[:i], newTail[i+1:]...)

		shift[len(shift)-1] = shift[len(shift)-1][:]
		driver := &shift[len(shift)-1]

		if !driver.Vacancy(NewPickup(driver.End(), head)) {
			driver = NewDriver(NewPickup(Origin, head))
			shift = append(shift, *driver)
		}

		if temp := shift.TotalCost(); temp > mincost {
			logger.WithFields(log.Fields{
				"newest":  temp,
				"current": mincost,
				"newTail": newTail.Graph(),
			}).Trace("path is already too long")
			continue
		}

		maxgoroutines <- struct{}{}
		waitlock.Add(1)

		go func(shift Shift, tail Edges) {
			defer func() {
				waitlock.Done()
				<-maxgoroutines
			}()

			result, err := shift.balance(tail, threads*2)
			if err != nil {
				logger.WithError(err).Trace("balance returned short")
				return
			}

			shiftlock.Lock()
			defer shiftlock.Unlock()

			shifts = append(shifts, result)
		}(shift, newTail)
	}

	waitlock.Wait()

	if len(shifts) == 0 {
		return nil, fmt.Errorf("dead ends")
	}

	s = slices.MinFunc(shifts, func(a, b Shift) int {
		if a.TotalCost() < b.TotalCost() {
			return -1
		}
		return 0
	})

	// logger.WithFields(log.Fields{
	// 	"tail":           report(tail),
	// 	"shortest shift": s.Graph(),
	// 	"distance":       s.TotalCost(),
	// 	"shifts": func() (result [][][]int) {
	// 		for _, s := range shifts {
	// 			result = append(result, s.Graph())
	// 		}
	// 		return result
	// 	}(),
	// }).Info("message")

	return s, nil
}

func (s Shift) Graph() (result [][]int) {
	for _, driver := range s {
		result = append(result, driver.Graph())
	}
	return result
}

func (s Shift) TotalCost() (result float64) {
	for _, d := range s {
		result += d.TotalCost()
	}
	return result
}

func (s Shift) Clone() (result Shift) {
	for i, d := range s {
		result = append(result, d[:])
		for j, p := range d {
			result[i][j] = func(p Pickup) *Pickup { return &p }(*p)
		}
	}
	return result
}
