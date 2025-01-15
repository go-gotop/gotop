package tracing

import (
	"context"
	"errors"

	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/exporters/zipkin"

	traceSdk "go.opentelemetry.io/otel/sdk/trace"
)

// NewExporter 创建一个导出器，支持：zipkin、otlp-http、otlp-grpc
func NewExporter(exporterName, endpoint string, insecure bool) (traceSdk.SpanExporter, error) {
	ctx := context.Background()
	switch exporterName {
	case "zipkin":
		return NewZipkinExporter(ctx, endpoint)
	case "otlp-http":
		return NewOtlpHttpExporter(ctx, endpoint, insecure)
	case "otlp-grpc":
		return NewOtlpGrpcExporter(ctx, endpoint, insecure)
	}
	return nil, errors.New("invalid exporter name")
}

// NewZipkinExporter 创建一个zipkin导出器
func NewZipkinExporter(_ context.Context, endpoint string) (traceSdk.SpanExporter, error) {
	return zipkin.New(endpoint)
}

// NewOtlpHttpExporter 创建OTLP/HTTP导出器
func NewOtlpHttpExporter(ctx context.Context, endpoint string, insecure bool, options ...otlptracehttp.Option) (traceSdk.SpanExporter, error) {
	var opts []otlptracehttp.Option
	opts = append(opts, otlptracehttp.WithEndpoint(endpoint))
	if insecure {
		opts = append(opts, otlptracehttp.WithInsecure())
	}
	opts = append(opts, options...)
	return otlptrace.New(
		ctx,
		otlptracehttp.NewClient(opts...),
	)
}

// NewOtlpGrpcExporter 创建OTLP/gRPC导出器
func NewOtlpGrpcExporter(ctx context.Context, endpoint string, insecure bool, options ...otlptracegrpc.Option) (traceSdk.SpanExporter, error) {
	var opts []otlptracegrpc.Option
	opts = append(opts, otlptracegrpc.WithEndpoint(endpoint))
	if insecure {
		opts = append(opts, otlptracegrpc.WithInsecure())
	}
	opts = append(opts, options...)
	return otlptrace.New(
		ctx,
		otlptracegrpc.NewClient(opts...),
	)
}
