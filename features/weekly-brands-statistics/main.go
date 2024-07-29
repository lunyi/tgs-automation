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

	"github.com/hashicorp/go-multierror"
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
		{"MOPH", config.MomoTelegram.MophChatId, true},
	}
}

func Run(ctx context.Context) {
	config := util.GetConfig()
	now := time.Now()

	//for month
	// firstDayOfLastMonth := time.Date(now.Year(), now.Month()-1, 1, 0, 0, 0, 0, time.UTC)
	// firstDayOfThisMonth := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, time.UTC)
	// startDate := firstDayOfLastMonth.Format("20060102") + "+8"
	// endDate := firstDayOfThisMonth.Format("20060102") + "+8"
	//for week
	startDate := now.AddDate(0, 0, -7).Format("20060102+8")
	endDate := now.AddDate(0, 0, 0).Format("20060102+8")
	log.LogInfo(fmt.Sprintf("startDate: %s, endDate: %s", startDate, endDate))

	brands := getBrandTelegramChannels(config)
	services := CreateBrandStatSvc(config.Postgresql)
	defer closeServices(services)

	for _, brand := range brands {
		err := createBrandReport(brand, startDate, endDate, services, config)
		if err != nil {
			log.LogError(fmt.Sprintf("Failed to create report: %v", err))
		}
	}

	err := deleteFiles()
	if err != nil {
		log.LogError(fmt.Sprintf("Failed to delete files: %v", err))
	}
}

func createBrandReport(brand BrandTelegramChannel, startDate string, endDate string, services BrandStatSvc, config util.TgsConfig) error {
	file := xlsx.NewFile()
	filename, err := processReport(file, brand.Code, startDate, endDate, services, brand.IsAllData)

	if err != nil {
		return fmt.Errorf("processReport failed: %s", err)
	}

	err = file.Save(filename)
	if err != nil {
		return fmt.Errorf("save failed: %s", err)
	}

	log.LogInfo(fmt.Sprintf("Sending file %s to telegram", filename))
	telegram.SendFile(config.MomoTelegram.Token, fmt.Sprintf("%d", brand.ChatID), filename)

	if err != nil {
		return fmt.Errorf("telegram sending file failed: %s", err)
	}
	return nil
}

func closeServices(services BrandStatSvc) {
	services.PromotionSvc.Close()
	services.PlayersAdjustSvc.Close()
}

func processReport(file *xlsx.File, brand string, startDate string, endDate string, services BrandStatSvc, IsAllData bool) (string, error) {
	params := CreateBrandStatParams(file, brand, startDate, endDate)
	if IsAllData {
		err := exportPlayerCount(&services.BonusPlayerCountSvc, params, "領取紅利人數")
		if err != nil {
			return "", err
		}
		err = exportPromotionDistributes(services.PromotionSvc, params)
		if err != nil {
			return "", err
		}
	}
	err := exportPlayerCount(&services.WithdrawPlayerCountSvc, params, "提款人數")
	if err != nil {
		return "", err
	}
	err = exportPlayerAdjustFile(services.PlayersAdjustSvc, params)
	if err != nil {
		return "", err
	}
	return params.Filename, nil
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
	now := time.Now()

	//for month
	// firstDayOfLastMonth := time.Date(now.Year(), now.Month()-1, 1, 0, 0, 0, 0, time.UTC)
	// lastDayOfLastMonth := firstDayOfLastMonth.AddDate(0, 1, -1)
	// filenameStart := firstDayOfLastMonth.Format("060102")
	// filenameEnd := lastDayOfLastMonth.Format("0102")

	//for week
	filenameStart := now.AddDate(0, 0, -7).Format("060102")
	filenameEnd := now.AddDate(0, 0, -1).Format("0102")
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

func deleteFiles() error {
	patterns := []string{"./*.xlsx", "./*.zip"}
	var allErrors *multierror.Error

	for _, pattern := range patterns {
		// Use filepath.Glob to find all files that match the pattern
		matches, err := filepath.Glob(pattern)
		if err != nil {
			allErrors = multierror.Append(allErrors, fmt.Errorf("error matching files with pattern %s: %w", pattern, err))
			continue
		}

		// Loop through the matching files and delete them
		for _, match := range matches {
			if err := os.Remove(match); err != nil {
				allErrors = multierror.Append(allErrors, fmt.Errorf("failed to delete %s: %w", match, err))
			}
		}
	}

	return allErrors.ErrorOrNil() // Returns nil if no errors were added
}
