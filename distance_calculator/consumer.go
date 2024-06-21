package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.con/reward-rabieth/Troll-Trekler/aggregator/client"
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

type KafkaConsumer struct {
	Consumer   *kafka.Consumer
	isRunning  bool
	calService CalculatorServicer
	aggClient  client.Client
}

func NewKafkaCosumer(topic string, svc CalculatorServicer, aggClient client.Client) (*KafkaConsumer, error) {
	c, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers": "localhost",
		"group.id":          "myGroup",
		"auto.offset.reset": "earliest",
	})
	if err != nil {
		return nil, err
	}
	err = c.SubscribeTopics([]string{topic}, nil)
	if err != nil {
		return nil, err
	}
	return &KafkaConsumer{
		Consumer:   c,
		calService: svc,
		aggClient:  aggClient,
	}, nil
}

func (c *KafkaConsumer) Start() {
	slog.Info("kafka transport started")
	c.isRunning = true
	c.ReadMessageLoop()
}

func (c *KafkaConsumer) ReadMessageLoop() {
	for c.isRunning {
		msg, err := c.Consumer.ReadMessage(-1)
		if err != nil {
			slog.Info("kafka consumer error", slog.Any("error", err))
			continue
		}
		var data types.OBUData
		if err := json.Unmarshal(msg.Value, &data); err != nil {
			slog.Info(" JSON serialization error: %s", err)

		}
		distance, err := c.calService.CalculateDistance(data)
		if err != nil {
			slog.Info("Calculation error", err)
			continue
		}
		req := &types.AggregateRequest{
			Value: distance,
			Unix:  time.Now().UnixNano(),
			ObuID: int32(data.OBUID),
		}
		if err := c.aggClient.Aggregate(context.Background(), req); err != nil {
			slog.Info(fmt.Sprintf("aggregate error %s", err))
			continue
		}
	}

}
