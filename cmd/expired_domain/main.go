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

	googleSheetSvc, err := googlesheet.New(config.GoogleSheet)
	if err != nil {
		fmt.Println(err.Error())
	}

	sheetName := time.Now().Format("01/2006") // Format: MM/YYYY

	googlesheet.CreateExpiredDomainExcel(googleSheetSvc.(*googlesheet.GoogleSheetService), sheetName, domainsForExcel)
}
