package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	_ "tgs-automation/features/domain-api/docs"
	jwttoken "tgs-automation/internal/jwt_token"
	middleware "tgs-automation/internal/middleware"

	"github.com/gin-gonic/gin"
	"github.com/iris-contrib/swagger/swaggerFiles"

	ginSwagger "github.com/swaggo/gin-swagger"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.24.0"
)

func main() {

	ctx := context.Background()
	gin.SetMode(gin.ReleaseMode)
	router := gin.New()

	initTracerProvider(ctx, "domain-api", "0.1.0", "prod")
	router.Use(middleware.TraceMiddleware("domain-api"))

	server := &http.Server{
		Addr:    ":8080",
		Handler: router,
	}

	router.ForwardedByClientIP = true
	router.SetTrustedProxies([]string{"127.0.0.1", "10.139.0.0/16"})
	router.GET("/healthz", healthCheckHandler)
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	router.GET("/token", jwttoken.TokenHandler)
	router.GET("/nameserver", middleware.AuthMiddleware(), GetNameServer)
	router.PUT("/nameserver", middleware.AuthMiddleware(), UpdateNameServer)
	router.GET("/domain/price", middleware.AuthMiddleware(), GetDomainPrice)
	router.POST("/domain", middleware.AuthMiddleware(), CreateDomain)
	err := server.ListenAndServe()

	if err != nil {
		fmt.Println("Failed to start server:", err)
	}
}

func healthCheckHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "up"})
}

func initTracerProvider(ctx context.Context, service string, version string, environment string) *sdktrace.TracerProvider {
	exporter, err := otlptracehttp.New(ctx,
		otlptracehttp.WithEndpoint("localhost:4318"),
		otlptracehttp.WithInsecure(),
	)
	if err != nil {
		log.Fatalf("OTLP Trace http Creation: %s %v", service, err)
	}
	tp := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exporter),
		sdktrace.WithResource(getResource(service, version, environment)),
	)
	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{}))
	return tp
}

func getResource(
	serviceName string,
	serviceVersion string,
	deploymentEnvironment string) *resource.Resource {
	// Defines resource with service name, version, and environment.
	return resource.NewWithAttributes(
		semconv.SchemaURL,
		semconv.ServiceNameKey.String(serviceName),
		semconv.ServiceVersionKey.String(serviceVersion),
		attribute.String("environment", deploymentEnvironment),
	)
}
