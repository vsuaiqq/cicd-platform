package client

import (
	"context"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"

	pb "github.com/vsuaiqq/cicd/shared/proto/gen/ai"
)

type AIClient struct {
	conn    *grpc.ClientConn
	client  pb.AIServiceClient
	timeout time.Duration
}

func NewAIClient(address string, timeout time.Duration) (*AIClient, error) {
	if timeout == 0 {
		timeout = 120 * time.Second
	}
	conn, err := newGRPCConn(address)
	if err != nil {
		return nil, err
	}
	return &AIClient{
		conn:    conn,
		client:  pb.NewAIServiceClient(conn),
		timeout: timeout,
	}, nil
}

func (c *AIClient) Close() error {
	return c.conn.Close()
}

type AnalyzeStep struct {
	Name      string
	Status    string
	ExitCode  int
	LogOutput string
}

type AnalyzeRequest struct {
	JobName      string
	JobStatus    string
	PipelineYAML string
	Steps        []AnalyzeStep
	Lang         string
}

type AnalyzeResponse struct {
	Summary       string   `json:"summary"`
	RootCause     string   `json:"root_cause"`
	Fix           string   `json:"fix"`
	RelevantLines []string `json:"relevant_lines"`
}

type GeneratePipelineResponse struct {
	YAML string `json:"yaml"`
}

func (c *AIClient) GeneratePipeline(ctx context.Context, description string) (*GeneratePipelineResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, c.timeout)
	defer cancel()

	resp, err := c.client.GeneratePipeline(ctx, &pb.GeneratePipelineRequest{
		Description: description,
	})
	if err != nil {
		return nil, err
	}
	return &GeneratePipelineResponse{YAML: resp.Yaml}, nil
}

func (c *AIClient) Analyze(ctx context.Context, req *AnalyzeRequest) (*AnalyzeResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, c.timeout)
	defer cancel()

	lang := req.Lang
	if lang == "" {
		lang = "en"
	}
	md := metadata.Pairs("lang", lang)
	ctx = metadata.NewOutgoingContext(ctx, md)

	protoSteps := make([]*pb.StepContext, len(req.Steps))
	for i, s := range req.Steps {
		protoSteps[i] = &pb.StepContext{
			Name:      s.Name,
			Status:    s.Status,
			ExitCode:  int32(s.ExitCode),
			LogOutput: s.LogOutput,
		}
	}

	resp, err := c.client.AnalyzeFailure(ctx, &pb.AnalyzeFailureRequest{
		JobName:      req.JobName,
		JobStatus:    req.JobStatus,
		PipelineYaml: req.PipelineYAML,
		Steps:        protoSteps,
	})
	if err != nil {
		return nil, err
	}

	return &AnalyzeResponse{
		Summary:       resp.Summary,
		RootCause:     resp.RootCause,
		Fix:           resp.Fix,
		RelevantLines: resp.RelevantLines,
	}, nil
}
