package main

import (
	"encoding/json"
	"fmt"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.con/reward-rabieth/Troll-Trekler/types"
	"log/slog"
	"net/http"
	"strconv"
	"time"
)

type HTTPFunc func(http.ResponseWriter, *http.Request) error

type APIError struct {
	Code int
	Err  error
}

func (e APIError) Error() string {
	return e.Err.Error()
}

type HTTPMetricHandler struct {
	reqCounter prometheus.Counter
	errCounter prometheus.Counter
	reqLatency prometheus.Histogram
}

func makeHTTPHandlerFunc(fn HTTPFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := fn(w, r); err != nil {
			if apiErr, ok := err.(APIError); ok {
				WriteJson(w, apiErr.Code, apiErr)
			}
		}
	}
}

func newHTTPMetricsHandler(reqName string) *HTTPMetricHandler {
	reqCounter := promauto.NewCounter(prometheus.CounterOpts{
		Namespace: fmt.Sprintf("http_%s_%s", reqName, "request_counter"),
		Name:      "aggregator",
	})

	errCounter := promauto.NewCounter(prometheus.CounterOpts{
		Namespace: fmt.Sprintf("http_%s_%s", reqName, "err_counter"),
		Name:      "aggregator",
	})

	reqLatency := promauto.NewHistogram(prometheus.HistogramOpts{
		Namespace: fmt.Sprintf("http_%s_%s", reqName, "request_latency"),
		Name:      "aggregator",
		Buckets:   []float64{0.1, 0.5},
	})
	return &HTTPMetricHandler{
		reqCounter: reqCounter,
		reqLatency: reqLatency,
		errCounter: errCounter,
	}
}

func (h *HTTPMetricHandler) instrument(next HTTPFunc) HTTPFunc {
	return func(w http.ResponseWriter, r *http.Request) error {
		var err error
		defer func(start time.Time) {
			latency := time.Since(start)
			h.reqLatency.Observe(time.Since(start).Seconds())
			slog.Info("",
				"latency", latency,
				"request", r.RequestURI,
				"err", err,
			)
		}(time.Now())
		h.reqCounter.Inc()
		if err != nil {
			h.errCounter.Inc()
		}
		err = next(w, r)
		return err

	}
}

func handleAggregate(svc Aggregator) HTTPFunc {
	return func(w http.ResponseWriter, r *http.Request) error {
		if r.Method != "POST" {
			return APIError{
				Code: http.StatusBadRequest,
				Err:  fmt.Errorf("invalid HTTP method %s", r.Method),
			}
		}
		var distance types.Distance
		if err := json.NewDecoder(r.Body).Decode(&distance); err != nil {
			return APIError{
				Code: http.StatusBadRequest,
				Err:  fmt.Errorf("failed to decode response body: %s", err),
			}
		}
		err := svc.AggregateDistance(distance)
		if err != nil {
			return APIError{
				Code: http.StatusInternalServerError,
				Err:  err,
			}
		}
		return WriteJson(w, http.StatusOK, map[string]string{"msg": "ok"})
	}
}

func handleGetInvoice(svc Aggregator) HTTPFunc {
	return func(w http.ResponseWriter, r *http.Request) error {
		if r.Method != "GET" {
			return APIError{
				Code: http.StatusBadRequest,
				Err:  fmt.Errorf("invalid HTTP method %s", r.Method),
			}

		}
		values, ok := r.URL.Query()["obu"]
		if !ok {
			return APIError{
				Code: http.StatusBadRequest,
				Err:  fmt.Errorf(" missing OBU id "),
			}

		}
		obuID, err := strconv.Atoi(values[0])
		if err != nil {
			return APIError{
				Code: http.StatusBadRequest,
				Err:  fmt.Errorf("invalid OBU id %s", values[0]),
			}
		}
		invoice, err := svc.CalculateInvoice(obuID)
		if err != nil {
			return APIError{
				Code: http.StatusInternalServerError,
				Err:  err,
			}
		}
		return WriteJson(w, http.StatusOK, invoice)
	}
}
