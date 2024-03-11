package main

import (
	"cdnetwork/internal/httpclient"
	"cdnetwork/internal/util"
	"cdnetwork/pkg/googlesheet"
	"cdnetwork/pkg/namecheap"
	"cdnetwork/pkg/postgresql"
	"fmt"
	"net/http"
	"time"
)

type ExpiredDomainsService struct {
	googlesheetInterface *googlesheet.GoogleSheetServiceInterface
	googlesheetSvc       *googlesheet.GoogleSheetService
	postgresqlInterface  *postgresql.GetAgentServiceInterface
	namecheapInterface   *namecheap.NamecheapAPI
}

func newExpiredDomainsService(
	googlesheetInterface *googlesheet.GoogleSheetServiceInterface,
	googlesheetSvc *googlesheet.GoogleSheetService,
	postgresqlInterface *postgresql.GetAgentServiceInterface,
	namecheapInterface *namecheap.NamecheapAPI
	) *ExpiredDomainsService {
	return &ExpiredDomainsService{
		googlesheetInterface: googlesheetInterface,
		googlesheetSvc:       googlesheetSvc,
		postgresqlInterface:  postgresqlInterface,
		namecheapInterface:   namecheapInterface,
	}
}

func main() {
	config := util.GetConfig()

	myclient := &httpclient.StandardHTTPClient{
		Client: &http.Client{
			Timeout: 10 * time.Second,
		},
	}

	namecheapClient := namecheap.New(config.Namecheap, myclient)

	domains, err := namecheapClient.GetExpiredDomains()
	if err != nil {
		fmt.Println(err.Error())
	}

	postgresqlClient := postgresql.New(config.Postgresql)
	domainsForExcel, err := postgresqlClient.GetAgentDomains(domains)

	if err != nil {
		fmt.Println(err.Error())
	}

	googleSheetSvcInterface, googleSheetSvc, err := googlesheet.New(config.GoogleSheet)
	if err != nil {
		fmt.Println(err.Error())
	}

	sheetName := time.Now().Format("01/2006") // Format: MM/YYYY

	googlesheet.CreateExpiredDomainExcel(
		googleSheetSvcInterface,
		googleSheetSvc,
		sheetName,
		domainsForExcel)
}

