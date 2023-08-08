package main

import (
	"context"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/app/server"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
	"github.com/hertz-contrib/obs-opentelemetry/provider"
	"github.com/hertz-contrib/obs-opentelemetry/tracing"
	"github.com/hertz-contrib/requestid"
	"go.opentelemetry.io/otel/trace"
)

func main() {
	provider.NewOpenTelemetryProvider(
		provider.WithServiceName("trace"),
		provider.WithExportEndpoint(":4317"),
		provider.WithInsecure(),
	)
	tracer, cfg := tracing.NewServerTracer()
	h := server.Default(tracer)
	h.Use(requestid.New(requestid.WithGenerator(func(ctx context.Context, c *app.RequestContext) string {
		traceID := trace.SpanFromContext(ctx).SpanContext().TraceID().String()
		return traceID
	})), tracing.ServerMiddleware(cfg))
	h.GET("/get", func(ctx context.Context, c *app.RequestContext) {
		c.String(consts.StatusOK, "success")
	})
	h.POST("/post", func(ctx context.Context, c *app.RequestContext) {
		c.String(consts.StatusOK, "success")
	})
	h.Spin()
}
