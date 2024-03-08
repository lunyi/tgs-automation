package googlesheet

import (
	"bytes"
	"cdnetwork/internal/log"
	"cdnetwork/internal/util"
	"cdnetwork/pkg/postgresql"
	"context"
	"fmt"
	"html/template"
	"io/ioutil"
	"time"

	"golang.org/x/oauth2/google"
	"golang.org/x/oauth2/jwt"
	"google.golang.org/api/sheets/v4"
)

type GoogleSheetServiceInterface interface {
	CreateSheetsService(key string) *sheets.Service
	CreateSheet(sheetsService *sheets.Service, spreadsheetId string, sheetName string)
	PlaceTextCenter(title string, values [][]interface{})
	WriteData(
		spreadsheetId string,
		domains []postgresql.DomainForExcel,
		writeRangeFunc func() string,
		valueRangeFunc func(domains []postgresql.DomainForExcel) *sheets.ValueRange) error
}

type GoogleSheetService struct {
	SpreadsheetId string
	GoogleApiKey  string
	SheetService  *sheets.Service
}

func New(config util.GoogleSheetConfig) (GoogleSheetServiceInterface, error) {
	return &GoogleSheetService{
		GoogleApiKey:  config.GoogleApiKey,
		SpreadsheetId: config.SheetId,
	}, nil
}

func CreateExpiredDomainExcel(gs *GoogleSheetService, sheetName string, domains []postgresql.DomainForExcel) {
	gs.SheetService = gs.CreateSheetsService(gs.GoogleApiKey)
	gs.CreateSheet(gs.SheetService, gs.SpreadsheetId, sheetName)

	gs.WriteData(
		gs.SpreadsheetId,
		domains,
		func() string {
			return fmt.Sprintf("%s!A1:F%d", sheetName, len(domains)+1)
		},
		func(domains []postgresql.DomainForExcel) *sheets.ValueRange {
			valueRange := createValueRangeForDomain(domains)
			gs.PlaceTextCenter(sheetName, valueRange.Values)
			return valueRange
		},
	)

	gs.WriteData(
		gs.SpreadsheetId,
		domains,
		func() string {
			return fmt.Sprintf("%s!A%d:F%d", sheetName, len(domains)+3, len(domains)+3)
		},
		func(domains []postgresql.DomainForExcel) *sheets.ValueRange {
			return createValueRangeForMessage(domains)
		},
	)

	fmt.Println("Data successfully written to Google Sheet.")
}

func (gs *GoogleSheetService) CreateSheetsService(key string) *sheets.Service {
	creds, err := ioutil.ReadFile(key)
	if err != nil {
		message := fmt.Sprintf("Unable to read client secret file: %v", err)
		log.LogFatal(message)
		fmt.Println(message)
	}

	conf, err := google.JWTConfigFromJSON(creds, sheets.SpreadsheetsScope)
	if err != nil {
		message := fmt.Sprintf("Unable to parse client secret file to config: %v", err)
		log.LogFatal(message)
		fmt.Println(message)
	}

	jwtConf := &jwt.Config{
		Email:      conf.Email,
		PrivateKey: conf.PrivateKey,
		Scopes:     []string{sheets.SpreadsheetsScope},
		TokenURL:   conf.TokenURL,
	}

	client := jwtConf.Client(context.Background())

	sheetsService, err := sheets.New(client)
	if err != nil {
		log.LogFatal(fmt.Sprintf("Unable to retrieve Sheets client: %v", err))
	}
	return sheetsService
}

func (gs *GoogleSheetService) PlaceTextCenter(title string, values [][]interface{}) {
	sheetId, err := getSheetID(gs.SheetService, gs.SpreadsheetId, title)

	if err != nil {
		log.LogFatal(fmt.Sprintf("Fail to get sheet id: %v", err))
	}

	requests := []*sheets.Request{
		{
			RepeatCell: &sheets.RepeatCellRequest{
				Range: &sheets.GridRange{
					SheetId:          sheetId,               // Provide the sheet ID here
					StartRowIndex:    0,                     // Start from the first row (header row)
					EndRowIndex:      int64(len(values)),    // End at the last row
					StartColumnIndex: 0,                     // Start from the first column
					EndColumnIndex:   int64(len(values[0])), // End at the last column
				},
				Cell: &sheets.CellData{
					UserEnteredFormat: &sheets.CellFormat{
						HorizontalAlignment: "CENTER",
					},
				},
				Fields: "userEnteredFormat.horizontalAlignment", // Specify the field to update
			},
		},
	}

	batchUpdate := &sheets.BatchUpdateSpreadsheetRequest{
		Requests: requests,
	}

	_, err = gs.SheetService.Spreadsheets.BatchUpdate(gs.SpreadsheetId, batchUpdate).Do()
	if err != nil {
		log.LogFatal(fmt.Sprintf("Unable to set alignment: %v", err))
	}
}

func (gs *GoogleSheetService) WriteData(
	spreadsheetId string,
	domains []postgresql.DomainForExcel,
	writeRangeFunc func() string,
	valueRangeFunc func(domains []postgresql.DomainForExcel) *sheets.ValueRange) error {

	valueRange := valueRangeFunc(domains)
	writeRange := writeRangeFunc()

	_, err := gs.SheetService.Spreadsheets.Values.Update(spreadsheetId, writeRange, valueRange).
		ValueInputOption("RAW").Do()
	if err != nil {
		log.LogFatal(fmt.Sprintf("Unable to update cell value: %v", err))
		return err
	}
	return nil
}

func (gs *GoogleSheetService) CreateSheet(sheetsService *sheets.Service, spreadsheetId string, sheetName string) {
	if checkIfSheetExists(sheetsService, spreadsheetId, sheetName) {
		return
	}

	newSheetProperties := &sheets.SheetProperties{
		Title: sheetName,
	}

	batchUpdateRequest := &sheets.BatchUpdateSpreadsheetRequest{
		Requests: []*sheets.Request{
			{
				AddSheet: &sheets.AddSheetRequest{
					Properties: newSheetProperties,
				},
			},
		},
	}

	_, err := sheetsService.Spreadsheets.BatchUpdate(spreadsheetId, batchUpdateRequest).Do()
	if err != nil {
		log.LogFatal(fmt.Sprintf("Unable to create new sheet: %v", err))
	}

	fmt.Printf("New sheet '%s' created successfully.\n", sheetName)
}

func createValueRangeForDomain(domains []postgresql.DomainForExcel) *sheets.ValueRange {
	values := [][]interface{}{
		{"平台", "代理", "網域", "建立日期", "到期日", "續約價格"},
	}

	for _, domain := range domains {
		row := []interface{}{
			domain.Brand,
			domain.Agent,
			domain.Domain,
			domain.Created,
			domain.Expires,
			"16.06",
		}
		values = append(values, row)
	}

	return &sheets.ValueRange{
		Values: values,
	}
}

func createValueRangeForMessage(domains []postgresql.DomainForExcel) *sheets.ValueRange {
	expirationDate := getResponseDate() + "(五)"

	domainNames := mapStrings(domains, func(domain postgresql.DomainForExcel) string {
		return domain.Domain
	})

	// Prepare data to be passed to the template
	data := struct {
		ExpirationDate string
		Domains        []string
	}{
		ExpirationDate: expirationDate,
		Domains:        domainNames,
	}

	const msgTemplate = `
貴司您好，以下{{ len .Domains }}個域名即將到期，綁定代理、到期日、續約價格請參考上圖；
請於{{ .ExpirationDate }}之前回覆是否需要續約，我司將優先為您處理，謝謝您!
{{- range .Domains }}
{{ . }}
{{- end }}
	`

	tmpl, err := template.New("message").Parse(msgTemplate)
	if err != nil {
		log.LogFatal(fmt.Sprintf("Error parsing message template: %v", err))
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		log.LogFatal(fmt.Sprintf("Error executing template: %v", err))
	}

	valueRange := &sheets.ValueRange{
		Values: [][]interface{}{
			{buf.String()}, // Call the function to compose the message
		},
	}
	return valueRange
}

func mapStrings(domains []postgresql.DomainForExcel, f func(postgresql.DomainForExcel) string) []string {
	result := make([]string, len(domains))
	for i, v := range domains {
		result[i] = f(v)
	}
	return result
}

func checkIfSheetExists(sheetsService *sheets.Service, spreadsheetId string, sheetName string) bool {
	resp, err := sheetsService.Spreadsheets.Get(spreadsheetId).Do()
	if err != nil {
		log.LogFatal(fmt.Sprintf("Unable to retrieve spreadsheet: %v", err))
	}

	for _, sheet := range resp.Sheets {
		if sheet.Properties.Title == sheetName {
			return true
		}
	}
	return false
}

func getLastFridayOfMonth(now time.Time) time.Time {
	// Get the first day of the next month
	firstDayOfNextMonth := now.AddDate(0, 1, 0)
	firstDayOfNextMonth = time.Date(firstDayOfNextMonth.Year(), firstDayOfNextMonth.Month(), 1, 0, 0, 0, 0, time.Local)

	// Subtract one day to get the last day of the current month
	lastDayOfCurrentMonth := firstDayOfNextMonth.Add(-24 * time.Hour)

	// Determine the day of the week for the last day of the month
	dayOfWeek := lastDayOfCurrentMonth.Weekday()

	// Calculate the number of days to subtract to get to the last Friday
	daysToSubtract := int(dayOfWeek - time.Friday)
	if daysToSubtract < 0 {
		daysToSubtract += 7 // Ensure positive result
	}

	// Subtract the appropriate number of days to find the last Friday
	return lastDayOfCurrentMonth.Add(-time.Duration(daysToSubtract) * 24 * time.Hour)
}

func getResponseDate() string {
	lastFridayOfMonth := getLastFridayOfMonth(time.Now())

	formattedLastFriday := lastFridayOfMonth.Format("01/02")

	fmt.Println("Last Friday of the month (formatted as {month/day}):", formattedLastFriday)
	return formattedLastFriday
}

func getSheetID(srv *sheets.Service, spreadsheetID string, sheetTitle string) (int64, error) {
	// Retrieve the spreadsheet
	spreadsheet, err := srv.Spreadsheets.Get(spreadsheetID).Do()
	if err != nil {
		return 0, err
	}

	// Search for the sheet by title
	var sheetID int64
	for _, sheet := range spreadsheet.Sheets {
		if sheet.Properties.Title == sheetTitle {
			sheetID = sheet.Properties.SheetId
			break
		}
	}

	if sheetID == 0 {
		return 0, fmt.Errorf("Sheet with title '%s' not found", sheetTitle)
	}

	return sheetID, nil
}
