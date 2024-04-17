package main

import (
	"fmt"
	"tgs-automation/internal/log"
	"tgs-automation/pkg/postgresql"
	"time"

	"github.com/tealeg/xlsx"
)

func exportPlayerAdjustFile(app postgresql.GetPlayersAdjustAmountInterface, params BrandStatParams) error {
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
		playerAdjustAmounts, err := app.GetData(params.Brand, params.StartDate, params.EndDate, value.key)

		if err != nil {
			log.LogFatal(err.Error())
		}

		var data []interface{}
		for _, p := range playerAdjustAmounts {
			data = append(data, p)
		}

		err = setHeaderAndFillData(params.File, data, params.Filename, populateSheetHeader, value.sheet, value.column)
		if err != nil {
			log.LogFatal(err.Error())
		}
	}

	return nil
}

type PopulatorFunc func(*xlsx.Row, *xlsx.Style, *xlsx.Sheet, []interface{}, string)

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
