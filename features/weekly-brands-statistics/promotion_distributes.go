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

	for _, d := range result {
		log.LogInfo(fmt.Sprintf("PromotionDistribute: %v", d))
	}

	setHeaderAndFillData(file, result, fileName, populatePromotionDistributionSheetHeader, "活動派發列表", "")
	return nil

}

func populatePromotionDistributionSheetHeader(headerRow *xlsx.Row, boldStyle *xlsx.Style, sheet *xlsx.Sheet, players []interface{}, dataType string) {
	headerTitles := []string{"用戶名", "活动名称", "活动类型", "活动子类型", "创建时间", "派发时间", "领取金额", "状态"}
	for _, title := range headerTitles {
		cell := headerRow.AddCell()
		cell.Value = title
		cell.SetStyle(boldStyle)
	}

	// Load the +8 time zone location
	loc, err := time.LoadLocation("Asia/Shanghai") // or "Asia/Singapore", "Australia/Perth" depending on exact location
	if err != nil {
		fmt.Println("Error loading location:", err)
		return
	}

	// Populating data
	for _, player := range players {
		row := sheet.AddRow()
		row.AddCell().Value = player.(postgresql.PromotionDistribute).Username
		row.AddCell().Value = player.(postgresql.PromotionDistribute).PromotionName
		row.AddCell().Value = player.(postgresql.PromotionDistribute).PromotionType
		row.AddCell().Value = player.(postgresql.PromotionDistribute).PromotionSubType
		createdOn := player.(postgresql.PromotionDistribute).CreatedOn
		createdOn = createdOn.In(loc)

		row.AddCell().Value = createdOn.Format("15:04:05 02/01/2006")
		row.AddCell().Value = fmt.Sprintf("%v", player.(postgresql.PromotionDistribute).BonusAmount)

		sentOn := player.(postgresql.PromotionDistribute).SentOn

		if !sentOn.IsZero() {
			sentOn := sentOn.In(loc)
			row.AddCell().Value = sentOn.Format("15:04:05 02/01/2006")
			row.AddCell().Value = "已派发"
		} else {
			row.AddCell().Value = ""
			row.AddCell().Value = "未派发"
		}
	}
}
