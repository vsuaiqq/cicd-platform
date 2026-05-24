package analytics

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

const _ = grpc.SupportPackageIsVersion9

const (
	AnalyticsService_GetDashboard_FullMethodName      = "/analytics.AnalyticsService/GetDashboard"
	AnalyticsService_GetCopilotContext_FullMethodName = "/analytics.AnalyticsService/GetCopilotContext"
)

type AnalyticsServiceClient interface {
	GetDashboard(ctx context.Context, in *DashboardRequest, opts ...grpc.CallOption) (*DashboardResponse, error)

	GetCopilotContext(ctx context.Context, in *CopilotContextRequest, opts ...grpc.CallOption) (*CopilotContextResponse, error)
}

type analyticsServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewAnalyticsServiceClient(cc grpc.ClientConnInterface) AnalyticsServiceClient {
	return &analyticsServiceClient{cc}
}

func (c *analyticsServiceClient) GetDashboard(ctx context.Context, in *DashboardRequest, opts ...grpc.CallOption) (*DashboardResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(DashboardResponse)
	err := c.cc.Invoke(ctx, AnalyticsService_GetDashboard_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *analyticsServiceClient) GetCopilotContext(ctx context.Context, in *CopilotContextRequest, opts ...grpc.CallOption) (*CopilotContextResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(CopilotContextResponse)
	err := c.cc.Invoke(ctx, AnalyticsService_GetCopilotContext_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

type AnalyticsServiceServer interface {
	GetDashboard(context.Context, *DashboardRequest) (*DashboardResponse, error)

	GetCopilotContext(context.Context, *CopilotContextRequest) (*CopilotContextResponse, error)
	mustEmbedUnimplementedAnalyticsServiceServer()
}

type UnimplementedAnalyticsServiceServer struct{}

func (UnimplementedAnalyticsServiceServer) GetDashboard(context.Context, *DashboardRequest) (*DashboardResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetDashboard not implemented")
}
func (UnimplementedAnalyticsServiceServer) GetCopilotContext(context.Context, *CopilotContextRequest) (*CopilotContextResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetCopilotContext not implemented")
}
func (UnimplementedAnalyticsServiceServer) mustEmbedUnimplementedAnalyticsServiceServer() {}
func (UnimplementedAnalyticsServiceServer) testEmbeddedByValue()                          {}

type UnsafeAnalyticsServiceServer interface {
	mustEmbedUnimplementedAnalyticsServiceServer()
}

func RegisterAnalyticsServiceServer(s grpc.ServiceRegistrar, srv AnalyticsServiceServer) {

	if t, ok := srv.(interface{ testEmbeddedByValue() }); ok {
		t.testEmbeddedByValue()
	}
	s.RegisterService(&AnalyticsService_ServiceDesc, srv)
}

func _AnalyticsService_GetDashboard_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DashboardRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AnalyticsServiceServer).GetDashboard(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: AnalyticsService_GetDashboard_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AnalyticsServiceServer).GetDashboard(ctx, req.(*DashboardRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _AnalyticsService_GetCopilotContext_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CopilotContextRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AnalyticsServiceServer).GetCopilotContext(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: AnalyticsService_GetCopilotContext_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AnalyticsServiceServer).GetCopilotContext(ctx, req.(*CopilotContextRequest))
	}
	return interceptor(ctx, in, info, handler)
}

var AnalyticsService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "analytics.AnalyticsService",
	HandlerType: (*AnalyticsServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "GetDashboard",
			Handler:    _AnalyticsService_GetDashboard_Handler,
		},
		{
			MethodName: "GetCopilotContext",
			Handler:    _AnalyticsService_GetCopilotContext_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "analytics.proto",
}
