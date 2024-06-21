package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/joho/godotenv"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.con/reward-rabieth/Troll-Trekler/aggregator/client"
	"github.con/reward-rabieth/Troll-Trekler/types"
	"google.golang.org/grpc"
	"log"
	"net"
	"net/http"
	"os"
	"time"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal(err)
	}

	var (
		store          = makeStore()
		svc            = NewInvoiceAggregator(store)
		grpcListenAddr = os.Getenv("AGG_GRPC_ENDPOINT")
		httpListenAddr = os.Getenv("AGG_HTTP_ENDPOINT")
	)
	svc = NewMetricsMiddleware(svc)
	svc = NewLogMiddleware(svc)

	go func() {
		log.Fatal(makeGRPCTransport(grpcListenAddr, svc))
	}()

	c, err := client.NEWGRPCClient(grpcListenAddr)
	if err != nil {
		log.Fatal(err)
	}
	err = c.Aggregate(context.Background(), &types.AggregateRequest{
		ObuID: 1,
		Value: 23.4,
		Unix:  time.Now().UnixNano(),
	})
	if err != nil {
		log.Fatal(err)
	}
	log.Fatal(makeHTTPTransport(httpListenAddr, svc))

}

func makeHTTPTransport(listenAddr string, svc Aggregator) error {
	//aggMetricHandler := newHTTPMetricsHandler("aggregate")
	invMetricHandler := newHTTPMetricsHandler("invoice")
	invoiceHandler := makeHTTPHandlerFunc(invMetricHandler.instrument(handleGetInvoice(svc)))
	http.HandleFunc("/invoice", invoiceHandler)
	//http.HandleFunc("/aggregate", aggMetricHandler.instrument(handleAggregate(svc)))
	//http.HandleFunc("/invoice", invMetricHandler.instrument(handleGetInvoice(svc)))
	http.Handle("/metrics", promhttp.Handler())
	fmt.Println("HTTP transport running on port", listenAddr)
	return http.ListenAndServe(listenAddr, nil)
}

func makeGRPCTransport(listenAddr string, svc Aggregator) error {
	fmt.Println("GRPC transport running on port", listenAddr)
	//Make a TCP listener
	ln, err := net.Listen("tcp", listenAddr)
	if err != nil {
		return err
	}
	defer ln.Close()
	//Make a new GRPC native server
	var opts []grpc.ServerOption
	server := grpc.NewServer(opts...)

	//Register GRPC server implementation to GRPC package
	types.RegisterAggregatorServer(server, NewGRPCServer(svc))
	return server.Serve(ln)
}

func WriteJson(w http.ResponseWriter, status int, V any) error {
	w.WriteHeader(status)
	w.Header().Set("Content-Type", "application/json")
	return json.NewEncoder(w).Encode(V)
}

func makeStore() Storer {
	storeType := os.Getenv("AGG_STORE_TYPE")
	switch storeType {
	case "memory":
		return NewMemoryStore()
	default:
		log.Fatalf("invalid store type given %s", storeType)
		return nil
	}

}
