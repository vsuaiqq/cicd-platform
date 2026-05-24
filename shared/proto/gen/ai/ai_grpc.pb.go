package ai

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

const _ = grpc.SupportPackageIsVersion9

const (
	AIService_AnalyzeFailure_FullMethodName   = "/ai.AIService/AnalyzeFailure"
	AIService_GeneratePipeline_FullMethodName = "/ai.AIService/GeneratePipeline"
	AIService_PipelineCopilot_FullMethodName  = "/ai.AIService/PipelineCopilot"
)

type AIServiceClient interface {
	AnalyzeFailure(ctx context.Context, in *AnalyzeFailureRequest, opts ...grpc.CallOption) (*AnalyzeFailureResponse, error)
	GeneratePipeline(ctx context.Context, in *GeneratePipelineRequest, opts ...grpc.CallOption) (*GeneratePipelineResponse, error)
	PipelineCopilot(ctx context.Context, in *PipelineCopilotRequest, opts ...grpc.CallOption) (*PipelineCopilotResponse, error)
}

type aIServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewAIServiceClient(cc grpc.ClientConnInterface) AIServiceClient {
	return &aIServiceClient{cc}
}

func (c *aIServiceClient) AnalyzeFailure(ctx context.Context, in *AnalyzeFailureRequest, opts ...grpc.CallOption) (*AnalyzeFailureResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(AnalyzeFailureResponse)
	err := c.cc.Invoke(ctx, AIService_AnalyzeFailure_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *aIServiceClient) GeneratePipeline(ctx context.Context, in *GeneratePipelineRequest, opts ...grpc.CallOption) (*GeneratePipelineResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(GeneratePipelineResponse)
	err := c.cc.Invoke(ctx, AIService_GeneratePipeline_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *aIServiceClient) PipelineCopilot(ctx context.Context, in *PipelineCopilotRequest, opts ...grpc.CallOption) (*PipelineCopilotResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(PipelineCopilotResponse)
	err := c.cc.Invoke(ctx, AIService_PipelineCopilot_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

type AIServiceServer interface {
	AnalyzeFailure(context.Context, *AnalyzeFailureRequest) (*AnalyzeFailureResponse, error)
	GeneratePipeline(context.Context, *GeneratePipelineRequest) (*GeneratePipelineResponse, error)
	PipelineCopilot(context.Context, *PipelineCopilotRequest) (*PipelineCopilotResponse, error)
	mustEmbedUnimplementedAIServiceServer()
}

type UnimplementedAIServiceServer struct{}

func (UnimplementedAIServiceServer) AnalyzeFailure(context.Context, *AnalyzeFailureRequest) (*AnalyzeFailureResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method AnalyzeFailure not implemented")
}
func (UnimplementedAIServiceServer) GeneratePipeline(context.Context, *GeneratePipelineRequest) (*GeneratePipelineResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GeneratePipeline not implemented")
}
func (UnimplementedAIServiceServer) PipelineCopilot(context.Context, *PipelineCopilotRequest) (*PipelineCopilotResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method PipelineCopilot not implemented")
}
func (UnimplementedAIServiceServer) mustEmbedUnimplementedAIServiceServer() {}
func (UnimplementedAIServiceServer) testEmbeddedByValue()                   {}

type UnsafeAIServiceServer interface {
	mustEmbedUnimplementedAIServiceServer()
}

func RegisterAIServiceServer(s grpc.ServiceRegistrar, srv AIServiceServer) {

	if t, ok := srv.(interface{ testEmbeddedByValue() }); ok {
		t.testEmbeddedByValue()
	}
	s.RegisterService(&AIService_ServiceDesc, srv)
}

func _AIService_AnalyzeFailure_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(AnalyzeFailureRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AIServiceServer).AnalyzeFailure(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: AIService_AnalyzeFailure_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AIServiceServer).AnalyzeFailure(ctx, req.(*AnalyzeFailureRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _AIService_GeneratePipeline_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GeneratePipelineRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AIServiceServer).GeneratePipeline(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: AIService_GeneratePipeline_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AIServiceServer).GeneratePipeline(ctx, req.(*GeneratePipelineRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _AIService_PipelineCopilot_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(PipelineCopilotRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AIServiceServer).PipelineCopilot(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: AIService_PipelineCopilot_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AIServiceServer).PipelineCopilot(ctx, req.(*PipelineCopilotRequest))
	}
	return interceptor(ctx, in, info, handler)
}

var AIService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "ai.AIService",
	HandlerType: (*AIServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "AnalyzeFailure",
			Handler:    _AIService_AnalyzeFailure_Handler,
		},
		{
			MethodName: "GeneratePipeline",
			Handler:    _AIService_GeneratePipeline_Handler,
		},
		{
			MethodName: "PipelineCopilot",
			Handler:    _AIService_PipelineCopilot_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "ai.proto",
}
