package main

import (
	"cdnetwork/internal/httpclient"
	"cdnetwork/internal/log"
	"cdnetwork/internal/util"
	"cdnetwork/pkg/googlesheet"
	"cdnetwork/pkg/namecheap"
	"cdnetwork/pkg/postgresql"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type ExpiredDomainsService struct {
	googlesheetInterface googlesheet.GoogleSheetServiceInterface
	googlesheetSvc       *googlesheet.GoogleSheetService
	postgresqlInterface  postgresql.GetAgentServiceInterface
	namecheapInterface   namecheap.NamecheapAPI
	httpClient           *httpclient.StandardHTTPClient
}

func (app *ExpiredDomainsService) Create(sheetName string) error {
	domains, err := app.namecheapInterface.GetExpiredDomains()

	if err != nil {
		log.LogFatal(err.Error())
		return err
	}

	filterDomains, err := app.postgresqlInterface.GetAgentDomains(domains)

	if err != nil {
		log.LogFatal(err.Error())
		return err
	}

	err = googlesheet.CreateExpiredDomainExcel(app.googlesheetInterface, app.googlesheetSvc, sheetName, filterDomains)

	if err != nil {
		return err
	}
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
