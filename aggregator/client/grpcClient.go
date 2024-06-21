package client

import (
	"context"
	"github.con/reward-rabieth/Troll-Trekler/types"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type GRPCClient struct {
	Endpoint string
	client   types.AggregatorClient
}

func NEWGRPCClient(endPoint string) (*GRPCClient, error) {
	conn, err := grpc.Dial(endPoint, grpc.WithTransportCredentials(insecure.NewCredentials()))

	if err != nil {
		return nil, err
	}
	c := types.NewAggregatorClient(conn)
	return &GRPCClient{
		Endpoint: endPoint,
		client:   c,
	}, nil
}

func (c *GRPCClient) Aggregate(ctx context.Context, req *types.AggregateRequest) error {
	_, err := c.client.Aggregate(ctx, req)
	if err != nil {
		return err
	}
	return nil
}
