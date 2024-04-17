package main

import (
	"fmt"
	"tgs-automation/internal/log"
	"tgs-automation/pkg/postgresql"

	"github.com/tealeg/xlsx"
)

func exportPlayerCount(pcp postgresql.PlayerCountProvider, params BrandStatParams, columnName string) error {

	count, err := pcp.GetPlayerCount(params.Brand, params.StartDate, params.EndDate)

	if err != nil {
		log.LogFatal(err.Error())
		return err
	}
	var players []interface{}
	players = append(players, count)

	err = setHeaderAndFillData(params.File, players, params.Filename, populatePlayerCount, columnName, columnName)
	if err != nil {
		log.LogFatal(err.Error())
		return err
	}
	return nil
}

func populatePlayerCount(headerRow *xlsx.Row, boldStyle *xlsx.Style, sheet *xlsx.Sheet, players []interface{}, dataType string) {
	headerTitles := []string{dataType}
	for _, title := range headerTitles {
		cell := headerRow.AddCell()
		cell.Value = title
		cell.SetStyle(boldStyle)
	}
	row := sheet.AddRow()
	row.AddCell().Value = fmt.Sprintf("%v", players[0])
}
