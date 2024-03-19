package googlesheet

import (
	"fmt"
	"testing"

	"cdnetwork/internal/util"
	"cdnetwork/pkg/postgresql"

	"github.com/stretchr/testify/assert"
	gomock "go.uber.org/mock/gomock"
)

func TestNew(t *testing.T) {
	config := util.GoogleSheetConfig{
		GoogleApiKey: "your-api-key",
		SheetId:      "your-sheet-id",
	}

	svc, _, err := New(config)
	if err != nil {
		t.Errorf("Failed to create GoogleSheetService: %v", err)
	}

	// Add assertions here to test the created service

	// Example assertion:``
	if svc.(*GoogleSheetService).GoogleApiKey != config.GoogleApiKey {
		t.Errorf("Expected GoogleApiKey to be %s, but got %s", config.GoogleApiKey, svc.(*GoogleSheetService).GoogleApiKey)
	}

	if svc.(*GoogleSheetService).SpreadsheetId != config.SheetId {
		t.Errorf("Expected SheetId to be %s, but got %s", config.SheetId, svc.(*GoogleSheetService).SpreadsheetId)
	}
}

func TestCreateExpiredDomainExcel(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockGS := NewMockGoogleSheetServiceInterface(ctrl)
	mockGSS := &GoogleSheetService{SpreadsheetId: "your-sheet-id", GoogleApiKey: "your-api-key"}

	sheetName := "testSheet"
	domains := []postgresql.DomainForExcel{
		{Domain: "example.com", Expires: "2022-01-01"},
		{Domain: "example.org", Expires: "2022-02-01"},
	}

	mockGS.EXPECT().CreateSheetsService(gomock.Any()).Return(nil, nil)
	mockGS.EXPECT().CreateSheet(gomock.Any(), gomock.Any(), sheetName).Return(nil)
	mockGS.EXPECT().WriteData(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
	mockGS.EXPECT().WriteData(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)

	err := CreateExpiredDomainExcel(mockGS, mockGSS, sheetName, domains)
	assert.NoError(t, err)
}

func TestCreateExpiredDomainExcel_Error(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockGS := NewMockGoogleSheetServiceInterface(ctrl)
	mockGSS := &GoogleSheetService{SpreadsheetId: "your-sheet-id", GoogleApiKey: "your-api-key"}

	sheetName := "testSheet"
	domains := []postgresql.DomainForExcel{
		{Domain: "example.com", Expires: "2022-01-01"},
		{Domain: "example.org", Expires: "2022-02-01"},
	}

	mockGS.EXPECT().CreateSheetsService(gomock.Any()).Return(nil, fmt.Errorf("error creating sheet service"))

	err := CreateExpiredDomainExcel(mockGS, mockGSS, sheetName, domains)
	assert.Error(t, err)
	assert.EqualError(t, err, "error creating sheet service")
}

func TestCreateExpiredDomainExcel_NilService(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	//mockGS := NewMockGoogleSheetServiceInterface(ctrl)
	mockGSS := &GoogleSheetService{SpreadsheetId: "your-sheet-id", GoogleApiKey: "your-api-key"}

	sheetName := "testSheet"
	domains := []postgresql.DomainForExcel{
		{Domain: "example.com", Expires: "2022-01-01"},
		{Domain: "example.org", Expires: "2022-02-01"},
	}

	err := CreateExpiredDomainExcel(nil, mockGSS, sheetName, domains)
	assert.Error(t, err)
	assert.EqualError(t, err, "GoogleSheetServiceInterface or GoogleSheetService is nil")
}

func TestCreateExpiredDomainExcel_NilGoogleSheetService(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockGS := NewMockGoogleSheetServiceInterface(ctrl)
	//mockGSS := &GoogleSheetService{SpreadsheetId: "your-sheet-id", GoogleApiKey: "your-api-key"}

	sheetName := "testSheet"
	domains := []postgresql.DomainForExcel{
		{Domain: "example.com", Expires: "2022-01-01"},
		{Domain: "example.org", Expires: "2022-02-01"},
	}

	err := CreateExpiredDomainExcel(mockGS, nil, sheetName, domains)
	assert.Error(t, err)
	assert.EqualError(t, err, "GoogleSheetServiceInterface or GoogleSheetService is nil")
}
