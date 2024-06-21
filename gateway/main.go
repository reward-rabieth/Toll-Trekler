package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"github.con/reward-rabieth/Troll-Trekler/aggregator/client"
	"log"
	"log/slog"
	"net/http"
	"time"
)

type apiFunc func(w http.ResponseWriter, r *http.Request) error

func main() {

	listenAddr := flag.String("listenAddr", ":6000", "the listen Address of the HTTP server")
	aggregatorServiceAddr := flag.String("aggServiceAddr", "http://localhost:5000", "the listen Address of the aggregator server")
	flag.Parse()

	var (
		client         = client.NewHTTPClient(*aggregatorServiceAddr)
		invoiceHandler = newInvoiceHandler(client)
	)
	http.HandleFunc("/invoice", makeAPIfunc(invoiceHandler.handleGetInvoice))
	slog.Info(fmt.Sprintf("gateway HTTP server runnning on port %s", *listenAddr))
	log.Fatal(http.ListenAndServe(*listenAddr, nil))

}

type InvoiceHandler struct {
	client client.Client
}

func newInvoiceHandler(c client.Client) *InvoiceHandler {
	return &InvoiceHandler{
		client: c,
	}
}

func (h *InvoiceHandler) handleGetInvoice(w http.ResponseWriter, r *http.Request) error {
	slog.Info("hitting endpoint")
	//access aggClient
	inv, err := h.client.GetInvoice(context.Background(), 17491)
	if err != nil {
		return err
	}
	return writeJSON(w, http.StatusOK, inv)
}

func writeJSON(w http.ResponseWriter, code int, v any) error {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(code)
	return json.NewEncoder(w).Encode(v)
}

func makeAPIfunc(fn apiFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer func(start time.Time) {
			slog.Info("REQ : :", "took", time.Since(start), "uri", r.RequestURI)
		}(time.Now())
		err := fn(w, r)
		if err != nil {
			writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		}
	}
}
