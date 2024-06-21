package main

import (
	"github.con/reward-rabieth/Troll-Trekler/aggregator/client"
	"log"
)

const (
	kafkaTopic         = "obudata"
	aggregatorEndPoint = "http://127.0.0.1:5000"
)

//Transport(HTTP,GRPC,Kafka) =>Attach business logic to this transport

func main() {
	var (
		svc CalculatorServicer
		err error
	)
	svc = NewCalculatorService()
	svc = NewLogMiddleware(svc)
	httpClient := client.NewHTTPClient(aggregatorEndPoint)
	//grpcClient, err := client.NEWGRPCClient(aggregatorEndPoint)
	//if err != nil {
	//	log.Fatal(err)
	//}

	kafkaConsumer, err := NewKafkaCosumer(kafkaTopic, svc, httpClient)
	if err != nil {
		log.Fatal(err)
	}
	kafkaConsumer.Start()
}
