package main

import (
	"github.con/reward-rabieth/Troll-Trekler/types"
	"log/slog"
	"time"
)

type LogMiddleware struct {
	next CalculatorServicer
}

func NewLogMiddleware(next CalculatorServicer) CalculatorServicer {
	return &LogMiddleware{
		next: next,
	}
}

func (m *LogMiddleware) CalculateDistance(data types.OBUData) (dist float64, err error) {

	defer func(start time.Time) {
		slog.Info("calculate distance",
			"took", time.Since(start),
			"err", err,
			"dist", dist)
	}(time.Now())
	dist, err = m.next.CalculateDistance(data)
	return
}
