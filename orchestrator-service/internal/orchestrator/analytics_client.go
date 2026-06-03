package orchestrator

import (
	"context"
	"fmt"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/keepalive"

	pb "github.com/vsuaiqq/cicd/shared/proto/gen/analytics"
)

type analyticsGateClient struct {
	client  pb.AnalyticsServiceClient
	timeout time.Duration
}

func NewAnalyticsGateClient(address string, timeout time.Duration) (*analyticsGateClient, error) {
	if timeout == 0 {
		timeout = 10 * time.Second
	}
	conn, err := grpc.NewClient(
		address,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithKeepaliveParams(keepalive.ClientParameters{
			Time:                10 * time.Second,
			Timeout:             3 * time.Second,
			PermitWithoutStream: true,
		}),
	)
	if err != nil {
		return nil, fmt.Errorf("analytics gate client: dial %s: %w", address, err)
	}
	return &analyticsGateClient{
		client:  pb.NewAnalyticsServiceClient(conn),
		timeout: timeout,
	}, nil
}

func (c *analyticsGateClient) EvaluatePerformanceGate(ctx context.Context, req *pb.EvaluatePerformanceGateRequest) (*pb.EvaluatePerformanceGateResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, c.timeout)
	defer cancel()
	resp, err := c.client.EvaluatePerformanceGate(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("EvaluatePerformanceGate: %w", err)
	}
	return resp, nil
}
