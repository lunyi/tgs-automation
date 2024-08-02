package main

// go install github.com/google/wire/cmd/wire@latest
import (
	"net/http"
	"tgs-automation/internal/httpclient"
	"tgs-automation/internal/util"
	"tgs-automation/pkg/googlesheet"
	"tgs-automation/pkg/namecheap"
	"tgs-automation/pkg/postgresql"
	"time"
)

func NewExpiredDomainsService(config util.TgsConfig) *ExpiredDomainsService {
	client := providerHttpClient()
	return &ExpiredDomainsService{
		httpClient:           client,
		googlesheetInterface: providerGoogleSheetInterface(config.GoogleSheet),
		googlesheetSvc:       providerGoogleSheetSvc(config.GoogleSheet),
		postgresqlInterface:  providerPostgresql(config.Postgresql),
		namecheapInterface:   providerNameCheap(config, client),
	}
}

func providerNameCheap(cfg util.TgsConfig, client *httpclient.StandardHTTPClient) namecheap.NamecheapApi {
	return namecheap.New(cfg)
}

func providerPostgresql(cfg util.PostgresqlConfig) postgresql.GetAgentServiceInterface {
	return postgresql.New(cfg)
}

func providerGoogleSheetSvc(cfg util.GoogleSheetConfig) *googlesheet.GoogleSheetService {
	itf, svc, err := googlesheet.New(cfg)
	if err != nil {
		return svc
	}

	if itf != nil {
		return svc
	}
	return nil
}

func providerGoogleSheetInterface(cfg util.GoogleSheetConfig) googlesheet.GoogleSheetServiceInterface {
	itf, svc, err := googlesheet.New(cfg)
	if err != nil {
		return itf
	}
	if svc != nil {
		return svc
	}
	return nil
}

func providerHttpClient() *httpclient.StandardHTTPClient {
	return &httpclient.StandardHTTPClient{
		Client: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}
