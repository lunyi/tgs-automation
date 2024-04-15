package main

import (
	"cdnetwork/internal/log"
	"cdnetwork/internal/util"
	"cdnetwork/pkg/postgresql"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/tealeg/xlsx"
)

func main() {
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)

	config := util.GetConfig()

	app1 := postgresql.NewPromotionTypesInterface(config.Postgresql)
	defer app1.Close()
	app1.GetPromotionTypes()

	return
	// Initialize the database connection or service interface
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

func exportPlayerAdjustFile(app postgresql.GetPlayersAdjustAmountInterface, brand string) error {
	file := xlsx.NewFile()
	start := time.Now().AddDate(0, 0, -8).Format("0102")
	end := time.Now().AddDate(0, 0, -2).Format("0102")
	filename := fmt.Sprintf("%s-%s_%s.xlsx", start, end, brand)

	startDate := time.Now().AddDate(0, 0, -8).Format("20060102+8")
	endDate := time.Now().AddDate(0, 0, -1).Format("20060102+8")

	adjustType := []struct {
		key    int
		column string
		sheet  string
	}{
		{4, "反水金額", "反水"},
		{5400, "優惠調帳", "優惠調帳"},
		{2, "公司調帳", "公司調帳"},
	}

	for _, value := range adjustType {
		playerAdjustAmounts, err := app.GetData(brand, startDate, endDate, value.key)

		if err != nil {
			log.LogFatal(err.Error())
		}

		var data []interface{}
		for _, p := range playerAdjustAmounts {
			data = append(data, p)
		}

		err = createSheet(file, data, filename, populateSheetHeader, value.column, value.sheet)
		if err != nil {
			log.LogFatal(err.Error())
		}
	}

	return nil
}

type PopulatorFunc func(*xlsx.Row, *xlsx.Style, *xlsx.Sheet, []interface{}, string)

func createSheet(file *xlsx.File, players []interface{}, excelFilename string, populate PopulatorFunc, dataType string, sheetName string) error {

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

func populateSheetHeader(headerRow *xlsx.Row, boldStyle *xlsx.Style, sheet *xlsx.Sheet, players []interface{}, dataType string) {
	headerTitles := []string{"玩家用戶名", dataType, "派發前餘額", "派發後餘額", "執行時間", "執行者", "描述"}
	for _, title := range headerTitles {
		cell := headerRow.AddCell()
		cell.Value = title
		cell.SetStyle(boldStyle)
	}

	// Populating data
	for _, player := range players {
		row := sheet.AddRow()
		row.AddCell().Value = player.(postgresql.PlayerAdjustAmountData).PlayerName
		row.AddCell().Value = fmt.Sprintf("%.2f", player.(postgresql.PlayerAdjustAmountData).Amount)
		row.AddCell().Value = fmt.Sprintf("%.2f", player.(postgresql.PlayerAdjustAmountData).BeforeBalance)
		row.AddCell().Value = fmt.Sprintf("%.2f", player.(postgresql.PlayerAdjustAmountData).AfterBalance)
		row.AddCell().Value = player.(postgresql.PlayerAdjustAmountData).ExecutionTime.Format(time.RFC3339)
		row.AddCell().Value = player.(postgresql.PlayerAdjustAmountData).Executor
		row.AddCell().Value = player.(postgresql.PlayerAdjustAmountData).Description
	}
}
