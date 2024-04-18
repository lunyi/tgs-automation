package main

import (
	"context"
	"fmt"
	log2 "log"
	"os"
	"path/filepath"
	"strings"
	"tgs-automation/internal/log"
	"tgs-automation/internal/util"
	"tgs-automation/pkg/postgresql"
	"tgs-automation/pkg/signalhandler"
	"tgs-automation/pkg/telegram"
	"time"

	"github.com/tealeg/xlsx"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/jaeger"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.9.0"
)

func init() {
	// Initialize Jaeger exporter to send traces to
	exporter, err := jaeger.New(jaeger.WithCollectorEndpoint(jaeger.WithEndpoint("http://0.0.0.0:4317")))
	if err != nil {
		log.LogFatal(fmt.Sprintf("Failed to create Jaeger exporter: %v", err))
	}

	// Create a new tracer provider with a batch span processor and the Jaeger exporter
	tp := trace.NewTracerProvider(
		trace.WithBatcher(exporter),
		// Add resource attributes like service name
		trace.WithResource(resource.NewWithAttributes(semconv.SchemaURL, semconv.ServiceNameKey.String("YourServiceName"))),
	)
	otel.SetTracerProvider(tp)

	// Ensure all spans are flushed before the application exits
	defer func() {
		if err := tp.Shutdown(context.Background()); err != nil {
			log2.Fatalf("Failed to shutdown TracerProvider: %v", err)
			log.LogFatal(fmt.Sprintf("Failed to shutdown TracerProvider: %v", err))
		}
	}()
}

// 每日首存人數和註冊玩家資料
func main() {

	signalhandler.StartListening()

	config := util.GetConfig()
	app := postgresql.NewMomoDataInterface(config.Postgresql)
	defer app.Close()

	today := time.Now().Format("2006-01-02") // Today's date in "YYYY-MM-DD" format
	yesterday := time.Now().AddDate(0, 0, -1).Format("2006-01-02")
	prefilename := time.Now().AddDate(0, 0, -1).Format("0102")

	brands := []struct {
		Code   string
		ChatID int64
	}{
		{"MOVN2", config.MomoTelegram.Movn2ChatId},
		{"MOPH", config.MomoTelegram.MophChatId},
	}

	for _, brand := range brands {
		playerFirstDepositFile := createExcelPlayerFirstDeposit(app, brand.Code, yesterday, today, prefilename)
		playerRegisteredFile := createExcelPlayerRegistered(app, brand.Code, yesterday, today, prefilename)
		filenames := []string{playerFirstDepositFile, playerRegisteredFile}

		sendFilesToTelegram(filenames, config.MomoTelegram.Token, fmt.Sprintf("%d", brand.ChatID))
		fmt.Println("-----")
	}
	deleteFiles()
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

type PopulatorFunc func(*xlsx.Row, *xlsx.Style, *xlsx.Sheet, []interface{})

func createExcel(players []interface{}, excelFilename string, populate PopulatorFunc) error {
	file := xlsx.NewFile()
	sheet, err := file.AddSheet("PlayerInfo")
	if err != nil {
		log.LogFatal(fmt.Sprintf("AddSheet failed: %s", err))
		return err
	}

	boldStyle := xlsx.NewStyle()
	boldStyle.Font.Bold = true
	headerRow := sheet.AddRow()
	// Populating data
	populate(headerRow, boldStyle, sheet, players)
	// Save the file to the disk
	err = file.Save(excelFilename)
	if err != nil {
		log.LogFatal(fmt.Sprintf("Save failed:: %s", err))
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

func sendFilesToTelegram(filePaths []string, botToken, chatID string) {
	for _, filePath := range filePaths {
		telegram.SendFile(botToken, chatID, filePath)
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
