package main

import (
	"fmt"
	"tgs-automation/internal/log"
	"tgs-automation/internal/util"
	"tgs-automation/pkg/postgresql"
	"time"

	"github.com/tealeg/xlsx"
)

func exportPromotionDistributes(config util.TgsConfig, file *xlsx.File, brand string, fileName string, startDate string, endDate string) error {

	appPromotionTypes := postgresql.NewPromotionTypesInterface(config.Postgresql)
	promotionTypes, err := appPromotionTypes.GetPromotionTypes()
	if err != nil {
		log.LogFatal(err.Error())
	}

	log.LogInfo(fmt.Sprintf("Promotion types: %v", promotionTypes))
	appPromotionTypes.Close()
	appPromotionDistributions := postgresql.NewPromotionDistributionInterface(config.Postgresql)
	data, err := appPromotionDistributions.GetData(brand, startDate, endDate)

	if err != nil {
		log.LogFatal(err.Error())
		return err
	}

	var result []interface{}
	for _, d := range data {
		log.LogInfo(fmt.Sprintf("PromotionDistribute: %v", d.PromotionSubType))
		for _, p := range promotionTypes {
			for _, s := range p.PromotionType {
				log.LogInfo(fmt.Sprintf("PromotionSubTypeName: %v", s.Name))
				if s.Name == d.PromotionSubType {

					log.LogInfo("==========")
					d.PromotionType = p.Trans.Zh
					d.PromotionSubType = s.Trans.Zh

					log.LogInfo(fmt.Sprintf("中文: %v, %v", p.Trans.Zh, s.Trans.Zh))
					log.LogInfo(fmt.Sprintf("轉換後: %v, %v", d.PromotionType, d.PromotionSubType))

					result = append(result, d)
					continue
				}
			}
		}
	}

	log.LogInfo("===================================")

	for _, d := range data {
		log.LogInfo(fmt.Sprintf("PromotionDistribute: %v", d))
	}

	setHeaderAndFillData(file, result, fileName, populatePromotionDistributionSheetHeader, "活動派發列表", "")
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
