package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"tgs-automation/internal/log"
	"tgs-automation/internal/util"
	"tgs-automation/pkg/postgresql"

	"github.com/tealeg/xlsx"
)

func main() {
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)

	config := util.GetConfig()
	app := postgresql.NewGetPlayersAdjustAmountInterface(config.Postgresql)
	defer app.Close()

	brands := []string{"MOPH", "MOVN2"}

	// Displaying all elements in the slice
	for _, brand := range brands {
		exportPlayerAdjustFile(app, brand)
	}

	sig := <-signals
	log.LogInfo(fmt.Sprintf("Received signal: %v, initiating shutdown", sig))
	os.Exit(0)
}

func createSheet(file *xlsx.File, players []interface{}, excelFilename string, populate PopulatorFunc, sheetName string, dataType = "" string) error {

	sheet, err := file.AddSheet(sheetName)
	if err != nil {
		log.LogFatal(fmt.Sprintf("AddSheet failed: %s", err))
		return err
	}

	boldStyle := xlsx.NewStyle()
	boldStyle.Font.Bold = true
	headerRow := sheet.AddRow()
	// Populating data
	populate(headerRow, boldStyle, sheet, players, dataType)
	// Save the file to the disk
	err = file.Save(excelFilename)
	if err != nil {
		log.LogFatal(fmt.Sprintf("Save failed:: %s", err))
		return err
	}

	log.LogInfo(fmt.Sprintf("Player adjust excel %s successfully.", sheetName))
	return nil
}
