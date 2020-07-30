package middleware

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/go-programming-tour-book/blog-service/global"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	"github.com/uber/jaeger-client-go"
)

//func Tracing() func(c *gin.Context) {
//	return func(c *gin.Context) {
//		var ctx context.Context
//		span := opentracing.SpanFromContext(c.Request.Context())
//		if span != nil {
//			span, ctx = opentracing.StartSpanFromContextWithTracer(
//				c.Request.Context(),
//				global.Tracer,
//				c.Request.URL.Path,
//				opentracing.ChildOf(span.Context()))
//		} else {
//			span, ctx = opentracing.StartSpanFromContextWithTracer(
//				c.Request.Context(),
//				global.Tracer,
//				c.Request.URL.Path,
//			)
//		}
//		defer span.Finish()
//		c.Request = c.Request.WithContext(ctx)
//		c.Next()
//	}
//}

func Tracing() func(c *gin.Context) {
	return func(c *gin.Context) {
		var newCtx context.Context
		var span opentracing.Span
		spanCtx, err := opentracing.GlobalTracer().Extract(opentracing.HTTPHeaders, opentracing.HTTPHeadersCarrier(c.Request.Header))
		if err != nil {
			span, newCtx = opentracing.StartSpanFromContextWithTracer(c.Request.Context(), global.Tracer, c.Request.URL.Path)
		} else {
			span, newCtx = opentracing.StartSpanFromContextWithTracer(
				c.Request.Context(),
				global.Tracer,
				c.Request.URL.Path,
				opentracing.ChildOf(spanCtx),
				opentracing.Tag{Key: string(ext.Component), Value: "HTTP"},
			)
		}
		defer span.Finish()

		var traceID string
		var spanID string
		var spanContext = span.Context()
		switch spanContext.(type) {
		case jaeger.SpanContext:
			jaegerContext := spanContext.(jaeger.SpanContext)
			traceID = jaegerContext.TraceID().String()
			spanID = jaegerContext.SpanID().String()
		}
		c.Set("X-Trace-ID", traceID)
		c.Set("X-Span-ID", spanID)
		c.Request = c.Request.WithContext(newCtx)
		c.Next()
	}
}
