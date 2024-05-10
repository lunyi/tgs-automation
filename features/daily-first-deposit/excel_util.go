// excel_util.go
package main

import (
	"github.com/tealeg/xlsx"
)

type PopulatorFunc func(*xlsx.Row, *xlsx.Style, *xlsx.Sheet, []interface{})

func initializeExcel(sheetName string) (*xlsx.File, *xlsx.Sheet, error) {
	file := xlsx.NewFile()
	sheet, err := file.AddSheet(sheetName)
	if err != nil {
		return nil, nil, err
	}
	return file, sheet, nil
}

func populateExcelData(sheet *xlsx.Sheet, players []interface{}, populate PopulatorFunc) error {
	boldStyle := xlsx.NewStyle()
	boldStyle.Font.Bold = true
	headerRow := sheet.AddRow()

	populate(headerRow, boldStyle, sheet, players)
	return nil
}

func saveExcelFile(file *xlsx.File, filename string) error {
	if err := file.Save(filename); err != nil {
		return err
	}
	return nil
}
