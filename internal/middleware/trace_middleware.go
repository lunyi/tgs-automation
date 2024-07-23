package middleware

import (
	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/otel"
)

func TraceMiddleware(service string) gin.HandlerFunc {
	return func(c *gin.Context) {
		tracer := otel.Tracer(service)
		ctx, span := tracer.Start(c.Request.Context(), c.Request.URL.Path)
		defer span.End()

		c.Request = c.Request.WithContext(ctx)
		c.Next()
	}
}
