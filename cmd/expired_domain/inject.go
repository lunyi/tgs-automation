// go install github.com/google/wire/cmd/wire@latest
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
	providerHttpClient,
	providerNameCheap,
	providerPostgresql,
	providerGoogleSheetSvc,
	providerGoogleSheetInterface,
)

func providerNameCheap(cfg util.NamecheapConfig, client *httpclient.StandardHTTPClient) namecheap.NamecheapAPI {
	return namecheap.New(cfg, client)
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
