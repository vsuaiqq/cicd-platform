package grpchandler

import (
	"context"
	"strings"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"

	"github.com/vsuaiqq/cicd/ai-service/internal/handler"
	"github.com/vsuaiqq/cicd/ai-service/internal/llm"
	"github.com/vsuaiqq/cicd/shared/logger"
	pb "github.com/vsuaiqq/cicd/shared/proto/gen/ai"
)

type AIServer struct {
	pb.UnimplementedAIServiceServer
	llm *llm.Client
}

func NewAIServer(llmClient *llm.Client) *AIServer {
	return &AIServer{llm: llmClient}
}

func (s *AIServer) GeneratePipeline(ctx context.Context, req *pb.GeneratePipelineRequest) (*pb.GeneratePipelineResponse, error) {
	if req.Description == "" {
		return nil, status.Error(codes.InvalidArgument, "description is required")
	}

	messages := handler.BuildGenerateMessages(handler.GenerateRequest{Description: req.Description})

	content, err := s.llm.Complete(ctx, messages)
	if err != nil {
		logger.L().Error().Err(err).Msg("GeneratePipeline llm error")
		return nil, status.Errorf(codes.Unavailable, "LLM request failed: %v", err)
	}

	yaml := strings.TrimSpace(content)
	if strings.HasPrefix(yaml, "```") {
		if i := strings.Index(yaml[3:], "\n"); i >= 0 {
			yaml = yaml[3+i+1:]
		}
		yaml = strings.TrimSuffix(strings.TrimSpace(yaml), "```")
		yaml = strings.TrimSpace(yaml)
	}

	return &pb.GeneratePipelineResponse{Yaml: yaml}, nil
}

func (s *AIServer) AnalyzeFailure(ctx context.Context, req *pb.AnalyzeFailureRequest) (*pb.AnalyzeFailureResponse, error) {
	if req.JobName == "" || len(req.Steps) == 0 {
		return nil, status.Error(codes.InvalidArgument, "job_name and at least one step are required")
	}

	lang := "en"
	if md, ok := metadata.FromIncomingContext(ctx); ok {
		if vals := md.Get("lang"); len(vals) > 0 {
			if vals[0] == "ru" || vals[0] == "en" {
				lang = vals[0]
			}
		}
	}

	steps := make([]handler.StepCtx, len(req.Steps))
	for i, s := range req.Steps {
		steps[i] = handler.StepCtx{
			Name:      s.Name,
			Status:    s.Status,
			ExitCode:  int(s.ExitCode),
			LogOutput: s.LogOutput,
		}
	}

	analyzeReq := handler.AnalyzeRequest{
		JobName:      req.JobName,
		JobStatus:    req.JobStatus,
		PipelineYAML: req.PipelineYaml,
		Steps:        steps,
		Lang:         lang,
	}

	messages := handler.BuildMessages(analyzeReq)

	content, err := s.llm.Complete(ctx, messages)
	if err != nil {
		logger.L().Error().Err(err).Msg("llm.Complete error")
		return nil, status.Errorf(codes.Unavailable, "LLM request failed: %v", err)
	}

	analysis, err := handler.ParseAnalysis(content)
	if err != nil {

		logger.L().Error().Err(err).Msg("parse analysis error")
		return &pb.AnalyzeFailureResponse{
			Summary:   "Analysis completed (unstructured response)",
			RootCause: content,
		}, nil
	}

	return &pb.AnalyzeFailureResponse{
		Summary:       analysis.Summary,
		RootCause:     analysis.RootCause,
		Fix:           analysis.Fix,
		RelevantLines: analysis.RelevantLines,
	}, nil
}
