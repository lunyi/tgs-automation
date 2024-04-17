package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"tgs-automation/internal/log"
	"tgs-automation/internal/util"
	"tgs-automation/pkg/postgresql"
	"time"

	"github.com/tealeg/xlsx"
)

func main() {
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)
	config := util.GetConfig()

	startDate := time.Now().AddDate(0, 0, -8).Format("20060102+8")
	endDate := time.Now().AddDate(0, 0, -1).Format("20060102+8")
	log.LogInfo(fmt.Sprintf("startDate: %s, endDate: %s", startDate, endDate))

	filenameStart := time.Now().AddDate(0, 0, -8).Format("060102")
	filenameEnd := time.Now().AddDate(0, 0, -2).Format("0102")

	brands := []string{"MOPH", "MOVN2"}

	bonusPlayerSvc := &postgresql.BonusPlayerCountService{Config: config.Postgresql}
	withdrawPlayerSvc := &postgresql.WithdrawPlayerCountService{Config: config.Postgresql}

	for _, brand := range brands {
		file := xlsx.NewFile()
		filename := fmt.Sprintf("%s-%s_%s.xlsx", filenameStart, filenameEnd, brand)

		exportPlayerCount(bonusPlayerSvc, file, filename, brand, startDate, endDate, "領取紅利人數")
		exportPromotionDistributes(config, file, brand, filename, startDate, endDate)
		exportPlayerCount(withdrawPlayerSvc, file, filename, brand, startDate, endDate, "提款人數")
		exportPlayerAdjustFile(config, file, filename, brand, startDate, endDate)

		err := file.Save(filename)
		if err != nil {
			log.LogFatal(fmt.Sprintf("Save failed:: %s", err))
		}
	}

	sig := <-signals
	log.LogInfo(fmt.Sprintf("Received signal: %v, initiating shutdown", sig))
	os.Exit(0)
}

func setHeaderAndFillData(file *xlsx.File, players []interface{}, excelFilename string, populate PopulatorFunc, sheetName string, dataType string) error {

	log.LogInfo(fmt.Sprintf("Creating filename %s, sheet %s", excelFilename, sheetName))

	sheet, err := file.AddSheet(sheetName)
	if err != nil {
		log.LogFatal(fmt.Sprintf("AddSheet failed: %s", err))
		return err
	}

	boldStyle := xlsx.NewStyle()
	boldStyle.Font.Bold = true
	headerRow := sheet.AddRow()

	populate(headerRow, boldStyle, sheet, players, dataType)

	log.LogInfo(fmt.Sprintf("Player adjust excel %s successfully.", sheetName))
	return nil
}
