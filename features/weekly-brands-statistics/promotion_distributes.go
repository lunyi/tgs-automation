package main

import (
	"fmt"
	"tgs-automation/internal/log"
	"tgs-automation/internal/util"
	"tgs-automation/pkg/postgresql"
	"time"

	"github.com/tealeg/xlsx"
)

func exportPromotionDistributes(config util.TgsConfig, file *xlsx.File, fileName string) error {

	app1 := postgresql.NewPromotionTypesInterface(config.Postgresql)
	types, err := app1.GetPromotionTypes()
	if err != nil {
		log.LogFatal(err.Error())
	}

	log.LogInfo(fmt.Sprintf("Promotion types: %v", types))
	app1.Close()

	app2 := postgresql.NewPromotionDistributionInterface(config.Postgresql)

	startDate := time.Now().AddDate(0, 0, -8).Format("20060102+8")
	endDate := time.Now().AddDate(0, 0, -1).Format("20060102+8")

	data, err := app2.GetData("MOPH", startDate, endDate)

	if err != nil {
		log.LogFatal(err.Error())
	}

	for _, d := range data {
		log.LogInfo(fmt.Sprintf("PromotionDistribute: %v", d))
	}

	//createSheet(file, data, fileName, populatePromotionDistributionSheetHeader, "活動派發列表", "")
	return nil

}

func populatePromotionDistributionSheetHeader(headerRow *xlsx.Row, boldStyle *xlsx.Style, sheet *xlsx.Sheet, players []interface{}, dataType string) {
	headerTitles := []string{"用戶名", "活動名稱", "活動類型", "活動子類型", "派發時間", "領取金額"}
	for _, title := range headerTitles {
		cell := headerRow.AddCell()
		cell.Value = title
		cell.SetStyle(boldStyle)
	}

	// Populating data
	for _, player := range players {
		row := sheet.AddRow()
		row.AddCell().Value = player.(postgresql.PromotionDistribute).Username
		row.AddCell().Value = player.(postgresql.PromotionDistribute).PromotionName
		row.AddCell().Value = player.(postgresql.PromotionDistribute).PromotionType
		row.AddCell().Value = player.(postgresql.PromotionDistribute).PromotionSubType
		row.AddCell().Value = player.(postgresql.PromotionDistribute).SentOn.Format(time.RFC3339)
		row.AddCell().Value = fmt.Sprintf("%v", player.(postgresql.PromotionDistribute).BonusAmount)
	}
}
