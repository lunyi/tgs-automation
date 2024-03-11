package main

import (
	"cdnetwork/internal/httpclient"
	"cdnetwork/internal/util"
	"cdnetwork/pkg/googlesheet"
	"cdnetwork/pkg/namecheap"
	"cdnetwork/pkg/postgresql"
	"net/http"
	"time"

	"github.com/google/wire"
)

var expiredDomainSet = wire.NewSet(
	providerNameCheap,
	providerPosrgresql,
	providerGoogleSheet,
)

func providerNameCheap(
	cfg util.NamecheapConfig,
	client *httpclient.StandardHTTPClient) namecheap.NamecheapAPI {
	return namecheap.New(cfg, client)
}

func providerPosrgresql(
	cfg util.PostgresqlConfig) postgresql.GetAgentServiceInterface {
	return postgresql.New(cfg)
}

func providerGoogleSheet(
	cfg util.GoogleSheetConfig) (googlesheet.GoogleSheetServiceInterface, *googlesheet.GoogleSheetService, error) {
	return googlesheet.New(cfg)
}

func providerHttpClient() *httpclient.StandardHTTPClient {
	return &httpclient.StandardHTTPClient{
		Client: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}
