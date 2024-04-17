package main

import (
	"fmt"
	"tgs-automation/internal/log"
	"tgs-automation/internal/util"
	"tgs-automation/pkg/postgresql"
	"tgs-automation/pkg/signalhandler"
	"time"

	"github.com/tealeg/xlsx"
)

func main() {
	go signalhandler.StartListening()
	initializeReports()
}

func initializeReports() {
	config := util.GetConfig()
	startDate := time.Now().AddDate(0, 0, -8).Format("20060102+8")
	endDate := time.Now().AddDate(0, 0, -1).Format("20060102+8")
	log.LogInfo(fmt.Sprintf("startDate: %s, endDate: %s", startDate, endDate))

	brands := []string{"MOPH", "MOVN2"}
	services := CreateBrandStatSvc(config.Postgresql)

	for _, brand := range brands {
		file := xlsx.NewFile()
		params := processReport(file, brand, startDate, endDate, services)
		err := file.Save(params.Filename)
		if err != nil {
			log.LogFatal(fmt.Sprintf("Save failed:: %s", err))
		}
	}
	services.PromotionSvc.Close()
	services.PlayersAdjustSvc.Close()
}

func processReport(file *xlsx.File, brand string, startDate string, endDate string, services BrandStatSvc) BrandStatParams {
	params := CreateBrandStatParams(file, brand, startDate, endDate)
	exportPlayerCount(&services.BonusPlayerCountSvc, params, "領取紅利人數")
	exportPromotionDistributes(services.PromotionSvc, params)
	exportPlayerCount(&services.WithdrawPlayerCountSvc, params, "提款人數")
	exportPlayerAdjustFile(services.PlayersAdjustSvc, params)
	return params
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

type BrandStatParams struct {
	File      *xlsx.File
	Filename  string
	Brand     string
	StartDate string
	EndDate   string
}

func CreateBrandStatParams(file *xlsx.File, brand string, startDate string, endDate string) BrandStatParams {
	filenameStart := time.Now().AddDate(0, 0, -8).Format("060102")
	filenameEnd := time.Now().AddDate(0, 0, -2).Format("0102")
	filename := fmt.Sprintf("%s-%s_%s.xlsx", filenameStart, filenameEnd, brand)
	return BrandStatParams{
		File:      file,
		Filename:  filename,
		Brand:     brand,
		StartDate: startDate,
		EndDate:   endDate,
	}
}

type BrandStatSvc struct {
	PromotionSvc           postgresql.GetPromotionInterface
	PlayersAdjustSvc       postgresql.GetPlayersAdjustAmountInterface
	BonusPlayerCountSvc    postgresql.BonusPlayerCountService
	WithdrawPlayerCountSvc postgresql.WithdrawPlayerCountService
}

func CreateBrandStatSvc(config util.PostgresqlConfig) BrandStatSvc {

	return BrandStatSvc{
		PromotionSvc:           postgresql.NewPromotionInterface(config),
		PlayersAdjustSvc:       postgresql.NewGetPlayersAdjustAmountInterface(config),
		BonusPlayerCountSvc:    postgresql.BonusPlayerCountService{Config: config},
		WithdrawPlayerCountSvc: postgresql.WithdrawPlayerCountService{Config: config},
	}
}
