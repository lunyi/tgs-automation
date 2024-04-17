package main

import (
	"fmt"
	"tgs-automation/internal/log"
	"tgs-automation/pkg/postgresql"

	"github.com/tealeg/xlsx"
)

func exportPromotionDistributes(app postgresql.GetPromotionInterface, params BrandStatParams) error {
	promotionTypes, err := app.GetPromotionTypes()
	defer app.Close()
	if err != nil {
		log.LogFatal(err.Error())
	}

	log.LogInfo(fmt.Sprintf("Promotion types: %v", promotionTypes))
	data, err := app.GetPromotionDistributions(params.Brand, params.StartDate, params.EndDate)

	if err != nil {
		log.LogFatal(err.Error())
		return err
	}

	var result []interface{}
	for _, d := range data {
		for _, p := range promotionTypes {
			for _, s := range p.PromotionType {
				if s.Name == d.PromotionSubType {
					d.PromotionType = p.Trans.Zh
					d.PromotionSubType = s.Trans.Zh
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

	setHeaderAndFillData(params.File, result, params.Filename, populatePromotionDistributionSheetHeader, "活動派發列表", "")
	return nil
}

func populatePromotionDistributionSheetHeader(headerRow *xlsx.Row, boldStyle *xlsx.Style, sheet *xlsx.Sheet, players []interface{}, dataType string) {
	headerTitles := []string{"用戶名", "活动名称", "活动类型", "活动子类型", "创建时间", "派发时间", "领取金额", "状态"}
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

		row.AddCell().Value = player.(postgresql.PromotionDistribute).CreatedOn

		senton := player.(postgresql.PromotionDistribute).SentOn
		row.AddCell().Value = senton
		row.AddCell().Value = fmt.Sprintf("%.2f", player.(postgresql.PromotionDistribute).BonusAmount)

		if senton == "" {
			row.AddCell().Value = "未派发"
		} else {
			row.AddCell().Value = "已派发"
		}
	}
}
