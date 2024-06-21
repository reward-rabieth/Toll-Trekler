package main

import (
	"context"
	"github.con/reward-rabieth/Troll-Trekler/types"
)

type GRPCAggregatorServer struct {
	types.UnimplementedAggregatorServer
	svc Aggregator
}

func NewGRPCServer(svc Aggregator) *GRPCAggregatorServer {
	return &GRPCAggregatorServer{
		svc: svc,
	}
}

//Transport layer
//JSON ->types.Distance
//GRPC ->types.AggregateRequest ->type.Distance
//business layer -> Business layer type (main type everyone needs to convert to)

func (s *GRPCAggregatorServer) Aggregate(ctx context.Context, req *types.AggregateRequest) (*types.None, error) {
	distance := types.Distance{
		OBUID: int(req.ObuID),
		Value: req.Value,
		Unix:  req.Unix,
	}

	return &types.None{}, s.svc.AggregateDistance(distance)
}
