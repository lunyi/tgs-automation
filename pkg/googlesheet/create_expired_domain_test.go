package googlesheet

import (
	"cdnetwork/pkg/postgresql"
	"testing"

	"go.uber.org/mock/gomock"
)

func TestCreateExpiredDomainExcel(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Mock GoogleSheetServiceInterface
	mockGS := NewMockGoogleSheetServiceInterface(ctrl)

	// Set expectations
	sheetName := "TestSheet"
	domains := []postgresql.DomainForExcel{
		// create test data here
	}
	mockGS.EXPECT().CreateSheetsService(gomock.Any()).Return(nil)
	mockGS.EXPECT().CreateSheet(gomock.Any(), gomock.Any(), sheetName)
	mockGS.EXPECT().WriteData(
		gomock.Any(),
		domains,
		gomock.Any(),
		gomock.Any(),
	).Times(2)

	svc := &GoogleSheetService{
		SpreadsheetId: "test_id",
		GoogleApiKey:  "test_api_key",
	}
	// Call the function under test
	CreateExpiredDomainExcel(mockGS, svc, sheetName, domains)
}
