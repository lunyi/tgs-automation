package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"tgs-automation/internal/log"
	"tgs-automation/internal/util"
	"tgs-automation/pkg/postgresql"
	"tgs-automation/pkg/telegram"
	"time"

	"github.com/tealeg/xlsx"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	go Run(ctx)

	<-ctx.Done() // Wait for signal
	log.LogInfo("Shutting down...")
}

type BrandTelegramChannel struct {
	Code      string
	ChatID    int64
	IsAllData bool
}

func getBrandTelegramChannels(config util.TgsConfig) []BrandTelegramChannel {
	return []BrandTelegramChannel{
		{"MOVN2", config.MomoTelegram.Movn2ChatId, false},
		{"MOPH", config.MomoTelegram.MophChatId, true},
	}
}

func Run(ctx context.Context) {
	config := util.GetConfig()
	startDate := time.Now().AddDate(0, 0, -7).Format("20060102+8")
	endDate := time.Now().AddDate(0, 0, 0).Format("20060102+8")
	log.LogInfo(fmt.Sprintf("startDate: %s, endDate: %s", startDate, endDate))

	brands := getBrandTelegramChannels(config)
	services := CreateBrandStatSvc(config.Postgresql)

	for _, brand := range brands {
		file := xlsx.NewFile()
		params := processReport(file, brand.Code, startDate, endDate, services, brand.IsAllData)
		err := file.Save(params.Filename)

		log.LogInfo(fmt.Sprintf("Sending file %s to telegram", params.Filename))
		telegram.SendFile(config.MomoTelegram.Token, fmt.Sprintf("%d", brand.ChatID), params.Filename)

		if err != nil {
			log.LogFatal(fmt.Sprintf("Save failed:: %s", err))
		}
	}
	deleteFiles()
	services.PromotionSvc.Close()
	services.PlayersAdjustSvc.Close()
}

func processReport(file *xlsx.File, brand string, startDate string, endDate string, services BrandStatSvc, IsAllData bool) BrandStatParams {
	params := CreateBrandStatParams(file, brand, startDate, endDate)
	if IsAllData {
		exportPlayerCount(&services.BonusPlayerCountSvc, params, "領取紅利人數")
		exportPromotionDistributes(services.PromotionSvc, params)
	}
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

func deleteFiles() {
	patterns := []string{
		"./*.xlsx",
		"./*.zip",
	}

	for _, pattern := range patterns {
		// Use filepath.Glob to find all files that match the pattern
		matches, err := filepath.Glob(pattern)
		if err != nil {
			log.LogFatal(err.Error())
		}

		// Loop through the matching files and delete them
		for _, match := range matches {
			err := os.Remove(match)
			if err != nil {
				log.LogInfo(fmt.Sprintf("Failed to delete %s: %s", match, err))
			} else {
				log.LogInfo(fmt.Sprintf("Deleted %s", match))
			}
		}
	}
}
