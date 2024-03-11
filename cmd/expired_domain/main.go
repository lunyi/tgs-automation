package main

import (
	"cdnetwork/internal/log"
	"cdnetwork/internal/util"
	"cdnetwork/pkg/googlesheet"
	"cdnetwork/pkg/namecheap"
	"cdnetwork/pkg/postgresql"
	"time"
)

type ExpiredDomainsService struct {
	googlesheetInterface *googlesheet.GoogleSheetServiceInterface
	googlesheetSvc       *googlesheet.GoogleSheetService
	postgresqlInterface  *postgresql.GetAgentService
	namecheapInterface   *namecheap.NamecheapClient
}

func newExpiredDomainsService(
	googlesheetInterface *googlesheet.GoogleSheetServiceInterface,
	googlesheetSvc *googlesheet.GoogleSheetService,
	postgresqlInterface *postgresql.GetAgentService,
	namecheapInterface *namecheap.NamecheapClient,
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
	sheetName := time.Now().Format("01/2006")

	app, err := InitExpiredDomainsService(config)

	if err != nil {
		log.LogFatal(err.Error())
	}

	domains, err := app.namecheapInterface.GetExpiredDomains()

	if err != nil {
		log.LogFatal(err.Error())
	}

	filterDomains, err := app.postgresqlInterface.GetAgentDomains(domains)

	if err != nil {
		log.LogFatal(err.Error())
	}

	googlesheet.CreateExpiredDomainExcel(*app.googlesheetInterface, app.googlesheetSvc, sheetName, filterDomains)

	// myclient := &httpclient.StandardHTTPClient{
	// 	Client: &http.Client{
	// 		Timeout: 10 * time.Second,
	// 	},
	// }

	// namecheapClient := namecheap.New(config.Namecheap, myclient)

	// googlesheet.CreateExpiredDomainExcel(
	// 	googleSheetSvcInterface,
	// 	googleSheetSvc,
	// 	sheetName,
	// 	domainsForExcel)

	// domains, err := namecheapClient.GetExpiredDomains()
	// if err != nil {
	// 	fmt.Println(err.Error())
	// }

	// postgresqlClient := postgresql.New(config.Postgresql)
	// domainsForExcel, err := postgresqlClient.GetAgentDomains(domains)

	// if err != nil {
	// 	fmt.Println(err.Error())
	// }

	// googleSheetSvcInterface, googleSheetSvc, err := googlesheet.New(config.GoogleSheet)
	// if err != nil {
	// 	fmt.Println(err.Error())
	// }

	// // Format: MM/YYYY

	// googlesheet.CreateExpiredDomainExcel(
	// 	googleSheetSvcInterface,
	// 	googleSheetSvc,
	// 	sheetName,
	// 	domainsForExcel)
}