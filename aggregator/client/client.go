package client

import (
	"context"
	"github.con/reward-rabieth/Troll-Trekler/types"
)

type Client interface {
	Aggregate(ctx context.Context, request *types.AggregateRequest) error
	GetInvoice(context.Context, int) (*types.Invoice, error)
}
