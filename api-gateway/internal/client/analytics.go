package client

import (
	"context"
	"fmt"
	"time"

	"google.golang.org/grpc"

	pb "github.com/vsuaiqq/cicd/shared/proto/gen/analytics"
)

type AnalyticsClient struct {
	conn    *grpc.ClientConn
	client  pb.AnalyticsServiceClient
	timeout time.Duration
}

func NewAnalyticsClient(address string, timeout time.Duration) (*AnalyticsClient, error) {
	conn, err := newGRPCConn(address)
	if err != nil {
		return nil, fmt.Errorf("analytics client: dial %s: %w", address, err)
	}
	return &AnalyticsClient{
		conn:    conn,
		client:  pb.NewAnalyticsServiceClient(conn),
		timeout: timeout,
	}, nil
}

func (c *AnalyticsClient) Close() error {
	return c.conn.Close()
}

func (c *AnalyticsClient) GetDashboard(ctx context.Context, projectID, period string) (*pb.DashboardResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, c.timeout)
	defer cancel()

	resp, err := c.client.GetDashboard(ctx, &pb.DashboardRequest{
		ProjectId: projectID,
		Period:    period,
	})
	if err != nil {
		return nil, fmt.Errorf("analytics client: GetDashboard: %w", err)
	}
	return resp, nil
}
