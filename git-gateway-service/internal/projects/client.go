package projects

import (
	"context"
	"fmt"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/keepalive"

	pb "github.com/vsuaiqq/cicd/shared/proto/gen/projects"
)

type Client struct {
	conn    *grpc.ClientConn
	client  pb.ProjectsServiceClient
	timeout time.Duration
}

func NewClient(grpcAddress string, timeout time.Duration) (*Client, error) {
	conn, err := grpc.NewClient(
		grpcAddress,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithKeepaliveParams(keepalive.ClientParameters{
			Time:                10 * time.Second,
			Timeout:             3 * time.Second,
			PermitWithoutStream: true,
		}),
	)
	if err != nil {
		return nil, fmt.Errorf("projects grpc client: %w", err)
	}
	return &Client{
		conn:    conn,
		client:  pb.NewProjectsServiceClient(conn),
		timeout: timeout,
	}, nil
}

func (c *Client) Close() error {
	return c.conn.Close()
}

func (c *Client) GetProjectInternal(ctx context.Context, projectID string) (*pb.GetProjectInternalResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, c.timeout)
	defer cancel()
	return c.client.GetProjectInternal(ctx, &pb.GetProjectInternalRequest{ProjectId: projectID})
}
