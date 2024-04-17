package main

import (
	"fmt"
	"tgs-automation/internal/log"
	"tgs-automation/pkg/postgresql"

	"github.com/tealeg/xlsx"
)

func exportPlayerCount(pcp postgresql.PlayerCountProvider, file *xlsx.File, filename string, brand string, startDate string, endDate string, columnName string) error {

	count, err := pcp.GetPlayerCount(brand, startDate, endDate)

	var players []interface{}
	players = append(players, count)

	if err != nil {
		err = setHeaderAndFillData(file, players, filename, populatePlayerCount, columnName, columnName)
		if err != nil {
			log.LogFatal(err.Error())
		}
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
