package main

import (
	"fmt"
	"github.con/reward-rabieth/Troll-Trekler/types"
	"log/slog"
)

const basePrice = 3.15

type Storer interface {
	Insert(data types.Distance) error
	Get(int) (float64, error)
}

type Aggregator interface {
	AggregateDistance(distance types.Distance) error
	CalculateInvoice(int) (*types.Invoice, error)
}

type InvoiceAggregator struct {
	store Storer
}

func NewInvoiceAggregator(store Storer) Aggregator {
	return &InvoiceAggregator{
		store: store,
	}
}

func (i *InvoiceAggregator) AggregateDistance(distance types.Distance) error {
	fmt.Println("processing and inserting distance in the storage", distance)
	slog.Info("aggregating distance", "obuID", distance.OBUID, "unix", distance.Unix)
	return i.store.Insert(distance)

}

func (i *InvoiceAggregator) CalculateInvoice(obuID int) (*types.Invoice, error) {
	dist, err := i.store.Get(obuID)
	if err != nil {
		return nil, err
	}
	inv := &types.Invoice{
		OBUID:         obuID,
		TotalDistance: dist,
		TotalAmount:   basePrice * dist,
	}
	return inv, nil
}
