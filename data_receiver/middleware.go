package main

import (
	"github.con/reward-rabieth/Troll-Trekler/types"
	"log/slog"
	"time"
)

//var (
//	w      = os.Stderr
//	logger = slog.New(tint.NewHandler(w, &tint.Options{
//		NoColor: !isatty.IsTerminal(w.Fd()),
//	}))
//)

type LogMiddleware struct {
	next DataProducer
}

func NewLogMiddleware(next DataProducer) *LogMiddleware {
	return &LogMiddleware{
		next: next,
	}
}

func (l *LogMiddleware) ProduceData(data types.OBUData) error {
	defer func(start time.Time) {
		//logger = slog.New(slog.NewJSONHandler(os.Stdout, nil))
		slog.Info("producing data to kafka",
			"ObuData", data.OBUID,
			"Latitude", data.Lat,
			"Longitude", data.Long,
			"took", time.Since(start),
		)

	}(time.Now())
	return l.next.ProduceData(data)
}
