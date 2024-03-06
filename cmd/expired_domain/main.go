package main

import (
	"cdnetwork/internal/util"
	"cdnetwork/pkg/googlesheet"
	"cdnetwork/pkg/namecheap"
	"cdnetwork/pkg/postgresql"
	"fmt"
)

func main() {
	config := util.GetConfig()
	namecheapClient := namecheap.New(config.Namecheap)

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

	googleSheetSvc.CreateExpiredDomainExcel(domainsForExcel)
}
