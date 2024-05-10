package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"path/filepath"
	"strings"
	"syscall"
	"tgs-automation/internal/log"
	"tgs-automation/internal/util"
	"tgs-automation/pkg/postgresql"
	"tgs-automation/pkg/telegram"
	"time"

	"github.com/hashicorp/go-multierror"
	"github.com/tealeg/xlsx"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/jaeger"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.9.0"
)

func setupTracer() (*trace.TracerProvider, error) {
	exporter, err := jaeger.New(jaeger.WithCollectorEndpoint(jaeger.WithEndpoint("http://0.0.0.0:4317")))
	if err != nil {
		return nil, err
	}
	tp := trace.NewTracerProvider(
		trace.WithBatcher(exporter),
		trace.WithResource(resource.NewWithAttributes(semconv.SchemaURL, semconv.ServiceNameKey.String("YourServiceName"))),
	)
	otel.SetTracerProvider(tp)
	return tp, nil
}

// 每日首存人數和註冊玩家資料
func main() {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	tracerProvider, err := setupTracer()
	if err != nil {
		log.LogFatal(fmt.Sprintf("Error setting up tracer: %v", err))
		return
	}
	defer func() {
		if err := tracerProvider.Shutdown(ctx); err != nil {
			log.LogFatal(fmt.Sprintf("Failed to shutdown TracerProvider: %v", err))
		}
	}()

	go Run(ctx)

	<-ctx.Done() // Wait for signal
	log.LogInfo("Shutting down...")
}

func Run(ctx context.Context) {
	config := util.GetConfig()
	momoDataInterface := postgresql.NewMomoDataInterface(config.Postgresql)
	defer momoDataInterface.Close()

	now := time.Now()
	today := now.Format("2006-01-02") // Today's date in "YYYY-MM-DD" format
	yesterday := now.AddDate(0, 0, -1).Format("2006-01-02")
	prefilename := now.AddDate(0, 0, -1).Format("0102")

	brands := []struct {
		Code   string
		ChatID int64
	}{
		{"MOVN2", config.MomoTelegram.Movn2ChatId},
		{"MOPH", config.MomoTelegram.MophChatId},
	}

	for _, brand := range brands {
		playerFirstDepositFile := createExcelPlayerFirstDeposit(momoDataInterface, brand.Code, yesterday, today, prefilename)
		playerRegisteredFile := createExcelPlayerRegistered(momoDataInterface, brand.Code, yesterday, today, prefilename)
		filenames := []string{playerFirstDepositFile, playerRegisteredFile}

		err := sendFilesToTelegram(filenames, config.MomoTelegram.Token, fmt.Sprintf("%d", brand.ChatID))
		if err != nil {
			log.LogFatal(fmt.Sprintf("Failed to send files to Telegram: %v", err))
		}
		fmt.Println("-----")
	}
	err := deleteFiles()
	if err != nil {
		log.LogFatal(fmt.Sprintf("Failed to delete files: %v", err))
	}
}

func createExcelPlayerRegistered(app postgresql.GetMomoDataInterface, brand string, yesterday string, today string, prefilename string) string {
	playerRegistered, err := app.GetRegisteredPlayers(brand, yesterday, today, "+08:00")
	if err != nil {
		log.LogFatal(err.Error())
	}

	var playerRegisteredInterface []interface{}
	for _, p := range playerRegistered {
		playerRegisteredInterface = append(playerRegisteredInterface, p)
	}

	playerRegisteredFile := fmt.Sprintf("%s_%s_註冊.xlsx", prefilename, strings.ToLower(brand))

	err = createExcel(playerRegisteredInterface, playerRegisteredFile, populateSheetPlayerRegistered)
	if err != nil {
		log.LogFatal(err.Error())
	}
	return playerRegisteredFile
}

func createExcelPlayerFirstDeposit(app postgresql.GetMomoDataInterface, brand string, yesterday string, today string, prefilename string) string {
	playerFirstDeposit, err := app.GetFirstDepositedPlayers(brand, yesterday, today, "+08:00")
	if err != nil {
		log.LogFatal(err.Error())
	}

	playerFirstDepositFile := fmt.Sprintf("%s_%s_首存.xlsx", prefilename, strings.ToLower(brand))

	var playerFirstDepositInterface []interface{}
	for _, p := range playerFirstDeposit {
		playerFirstDepositInterface = append(playerFirstDepositInterface, p)
	}

	err = createExcel(playerFirstDepositInterface, playerFirstDepositFile, populateSheetFirstDeposit)
	if err != nil {
		log.LogFatal(err.Error())
	}
	return playerFirstDepositFile
}

func createExcel(players []interface{}, excelFilename string, populate PopulatorFunc) error {
	file, sheet, err := initializeExcel("PlayerInfo")
	if err != nil {
		log.LogFatal(fmt.Sprintf("Initialize Excel failed: %s", err))
		return err
	}

	if err := populateExcelData(sheet, players, populate); err != nil {
		log.LogFatal(fmt.Sprintf("Populate Excel data failed: %s", err))
		return err
	}

	if err := saveExcelFile(file, excelFilename); err != nil {
		log.LogFatal(fmt.Sprintf("Save Excel file failed: %s", err))
		return err
	}

	log.LogInfo("Player first deposit excel successfully.")
	return nil
}

func populateSheetFirstDeposit(headerRow *xlsx.Row, boldStyle *xlsx.Style, sheet *xlsx.Sheet, players []interface{}) {
	headerTitles := []string{"Agent", "Host", "PlayerName", "DailyDepositAmount", "DailyDepositCount", "FirstDepositOn"}
	for _, title := range headerTitles {
		cell := headerRow.AddCell()
		cell.Value = title
		cell.SetStyle(boldStyle)
	}

	for _, player := range players {
		row := sheet.AddRow()
		row.AddCell().Value = player.(postgresql.PlayerFirstDeposit).Agent
		row.AddCell().Value = player.(postgresql.PlayerFirstDeposit).Host
		row.AddCell().Value = player.(postgresql.PlayerFirstDeposit).PlayerName
		row.AddCell().SetFloat(player.(postgresql.PlayerFirstDeposit).DailyDepositAmount)
		row.AddCell().SetInt(player.(postgresql.PlayerFirstDeposit).DailyDepositCount)
		row.AddCell().Value = player.(postgresql.PlayerFirstDeposit).FirstDepositOn.Format(time.RFC3339)
	}
}

func populateSheetPlayerRegistered(headerRow *xlsx.Row, boldStyle *xlsx.Style, sheet *xlsx.Sheet, players []interface{}) {
	headerTitles := []string{"Agent", "Host", "PlayerName", "RealName", "RegisteredOn"}
	for _, title := range headerTitles {
		cell := headerRow.AddCell()
		cell.Value = title
		cell.SetStyle(boldStyle)
	}

	// Populating data
	for _, player := range players {
		row := sheet.AddRow()
		row.AddCell().Value = player.(postgresql.PlayerRegisterInfo).Agent
		row.AddCell().Value = player.(postgresql.PlayerRegisterInfo).Host
		row.AddCell().Value = player.(postgresql.PlayerRegisterInfo).PlayerName
		row.AddCell().Value = player.(postgresql.PlayerRegisterInfo).RealName
		row.AddCell().Value = player.(postgresql.PlayerRegisterInfo).RegisteredOn.Format(time.RFC3339)
	}
}

func sendFilesToTelegram(filePaths []string, botToken, chatID string) error {
	var allErrors *multierror.Error

	for _, filePath := range filePaths {
		if err := telegram.SendFile(botToken, chatID, filePath); err != nil {
			// Log each error, but continue processing other files
			log.LogInfo(fmt.Sprintf("Failed to send file %s: %v", filePath, err))
			allErrors = multierror.Append(allErrors, fmt.Errorf("failed to send file %s: %w", filePath, err))
		}
	}

	return allErrors.ErrorOrNil() // Returns nil if no errors were added
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
