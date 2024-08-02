package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"tgs-automation/internal/httpclient"
	"tgs-automation/internal/log"
	"tgs-automation/internal/util"
	"tgs-automation/pkg/googlesheet"
	"tgs-automation/pkg/namecheap"
	"tgs-automation/pkg/postgresql"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
	"go.opentelemetry.io/otel/sdk/resource"
	traceSDK "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.21.0"
)

type ExpiredDomainsService struct {
	googlesheetInterface googlesheet.GoogleSheetServiceInterface
	googlesheetSvc       *googlesheet.GoogleSheetService
	postgresqlInterface  postgresql.GetAgentServiceInterface
	namecheapInterface   namecheap.NamecheapApi
	httpClient           *httpclient.StandardHTTPClient
}

func initTracer() func() {
	ctx := context.Background()

	exporter, err := stdouttrace.New(stdouttrace.WithPrettyPrint())
	if err != nil {
		log.LogFatal(fmt.Sprintf("failed to initialize stdouttrace exporter: %v", err))
	}

	tp := traceSDK.NewTracerProvider(
		traceSDK.WithBatcher(exporter),
		traceSDK.WithResource(resource.NewWithAttributes(semconv.SchemaURL, semconv.ServiceNameKey.String("expired-domains-service"))),
	)

	otel.SetTracerProvider(tp)

	return func() {
		if err := tp.Shutdown(ctx); err != nil {
			log.LogFatal(fmt.Sprintf("error shutting down tracer provider: %v", err))
		}
	}
}

func (app *ExpiredDomainsService) Create(sheetName string) error {
	shutdown := initTracer()
	defer shutdown()

	tracer := otel.Tracer("expired-domains-tracer")
	ctx, span := tracer.Start(context.Background(), "namecheap-expired-domains")
	span.SetAttributes(attribute.String("api", "namecheap"))
	domains, err := app.namecheapInterface.GetExpiredDomains()
	if err != nil {
		log.LogFatal(err.Error())
		return err
	}
	span.End()

	ctx, span = tracer.Start(ctx, "postgresql-expired-domains")
	span.SetAttributes(attribute.String("db", "postgresql"))
	filterDomains, err := app.postgresqlInterface.GetAgentDomains(domains)
	if err != nil {
		log.LogFatal(err.Error())
		return err
	}
	span.End()

	_, span = tracer.Start(ctx, "googlesheet-expired-domains")
	span.SetAttributes(attribute.String("api", "googlesheet"))
	err = googlesheet.CreateExpiredDomainExcel(app.googlesheetInterface, app.googlesheetSvc, sheetName, filterDomains)

	if err != nil {
		return err
	}
	span.End()
	return nil
}

func main() {
	config := util.GetConfig()
	sheetName := time.Now().Format("01/2006")

	app, err := InitExpiredDomainsService(config)

	if err != nil {
		log.LogFatal(err.Error())
	}

	// Set up channel to receive OS signals
	signals := make(chan os.Signal, 1)
	// Notify this channel on SIGINT or SIGTERM
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		app.Create(sheetName)
		if err != nil {
			log.LogFatal(err.Error())
		}
	}()

	sig := <-signals
	log.LogInfo(fmt.Sprintf("Received signal: %v, initiating shutdown", sig))

	// Perform any cleanup before exiting
	// app.Cleanup() // Example cleanup method

	// Exit program
	os.Exit(0)
}
