package client

import (
	"context"
	"time"

	"google.golang.org/grpc"

	pb "github.com/vsuaiqq/cicd/shared/proto/gen/auth"
)

type AuthClient struct {
	conn    *grpc.ClientConn
	client  pb.AuthServiceClient
	timeout time.Duration
}

func NewAuthClient(address string, timeout time.Duration) (*AuthClient, error) {
	conn, err := newGRPCConn(address)
	if err != nil {
		return nil, err
	}
	return &AuthClient{
		conn:    conn,
		client:  pb.NewAuthServiceClient(conn),
		timeout: timeout,
	}, nil
}

func (c *AuthClient) Close() error {
	return c.conn.Close()
}

func (c *AuthClient) Register(ctx context.Context, email, username, password string) (*pb.RegisterResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, c.timeout)
	defer cancel()
	return c.client.Register(ctx, &pb.RegisterRequest{
		Email:    email,
		Username: username,
		Password: password,
	})
}

func (c *AuthClient) Login(ctx context.Context, email, password string) (*pb.LoginResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, c.timeout)
	defer cancel()
	return c.client.Login(ctx, &pb.LoginRequest{
		Email:    email,
		Password: password,
	})
}

func (c *AuthClient) ValidateToken(ctx context.Context, token string) (*pb.ValidateTokenResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, c.timeout)
	defer cancel()
	return c.client.ValidateToken(ctx, &pb.ValidateTokenRequest{Token: token})
}

func (c *AuthClient) RefreshToken(ctx context.Context, refreshToken string) (*pb.RefreshTokenResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, c.timeout)
	defer cancel()
	return c.client.RefreshToken(ctx, &pb.RefreshTokenRequest{RefreshToken: refreshToken})
}

func (c *AuthClient) Health(ctx context.Context) (*pb.HealthResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, c.timeout)
	defer cancel()
	return c.client.Health(ctx, &pb.HealthRequest{})
}

func (c *AuthClient) GetUserByEmail(ctx context.Context, email string) (*pb.GetUserByEmailResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, c.timeout)
	defer cancel()
	return c.client.GetUserByEmail(ctx, &pb.GetUserByEmailRequest{Email: email})
}
