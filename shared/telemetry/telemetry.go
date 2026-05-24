package telemetry

import (
	"context"
	"net/http"

	"github.com/vsuaiqq/cicd/shared/logger"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.24.0"
	"go.opentelemetry.io/otel/trace"
)

type Config struct {
	ServiceName string

	OTLPEndpoint string

	SampleRatio float64
}

func TracerProvider(ctx context.Context, cfg Config) (*sdktrace.TracerProvider, func(), error) {
	if cfg.ServiceName == "" {
		cfg.ServiceName = "unknown"
	}

	res, err := resource.New(ctx,
		resource.WithSchemaURL(semconv.SchemaURL),
		resource.WithAttributes(semconv.ServiceName(cfg.ServiceName)),
	)
	if err != nil {
		return nil, nil, err
	}

	opts := []sdktrace.TracerProviderOption{
		sdktrace.WithResource(res),
	}
	if cfg.OTLPEndpoint != "" {
		exporter, err := newOTLPExporter(ctx, cfg.OTLPEndpoint)
		if err != nil {
			return nil, nil, err
		}
		opts = append(opts, sdktrace.WithBatcher(exporter))
	}
	if cfg.SampleRatio >= 0 && cfg.SampleRatio <= 1 {
		opts = append(opts, sdktrace.WithSampler(sdktrace.TraceIDRatioBased(cfg.SampleRatio)))
	}

	tp := sdktrace.NewTracerProvider(opts...)
	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(
		propagation.TraceContext{},
		propagation.Baggage{},
	))

	shutdown := func() {
		_ = tp.Shutdown(context.Background())
	}
	return tp, shutdown, nil
}

func Middleware(serviceName string) func(http.Handler) http.Handler {
	tracer := otel.Tracer(serviceName)
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := otel.GetTextMapPropagator().Extract(r.Context(), propagation.HeaderCarrier(r.Header))
			ctx, span := tracer.Start(ctx, r.Method+" "+r.URL.Path,
				trace.WithAttributes(
					attribute.String("http.method", r.Method),
					attribute.String("http.url", r.URL.String()),
				),
			)
			defer span.End()

			traceID := span.SpanContext().TraceID().String()
			spanID := span.SpanContext().SpanID().String()
			ctx = context.WithValue(ctx, logger.TraceIDKey, traceID)
			ctx = context.WithValue(ctx, logger.SpanIDKey, spanID)

			w.Header().Set("X-Trace-ID", traceID)
			w.Header().Set("X-Span-ID", spanID)

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func SpanFromContext(ctx context.Context, tracerName, spanName string, attrs ...attribute.KeyValue) (context.Context, func()) {
	tracer := otel.Tracer(tracerName)
	opts := []trace.SpanStartOption{}
	if len(attrs) > 0 {
		opts = append(opts, trace.WithAttributes(attrs...))
	}
	ctx, span := tracer.Start(ctx, spanName, opts...)
	traceID := span.SpanContext().TraceID().String()
	spanID := span.SpanContext().SpanID().String()
	ctx = context.WithValue(ctx, logger.TraceIDKey, traceID)
	ctx = context.WithValue(ctx, logger.SpanIDKey, spanID)
	return ctx, func() { span.End() }
}

func RecordError(ctx context.Context, err error) {
	span := trace.SpanFromContext(ctx)
	span.RecordError(err)
	span.SetStatus(codes.Error, err.Error())
}
