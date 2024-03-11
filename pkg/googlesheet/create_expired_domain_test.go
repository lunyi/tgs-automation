package googlesheet

import (
	"cdnetwork/pkg/postgresql"

	"github.com/stretchr/testify/mock"
	"google.golang.org/api/sheets/v4"
)

type MockGoogleSheetService struct {
	mock.Mock
}

// Mock the CreateSheetsService method
func (m *MockGoogleSheetService) CreateSheetsService(key string) *sheets.Service {
	args := m.Called(key)
	return args.Get(0).(*sheets.Service)
}

// Mock the CreateSheet method
func (m *MockGoogleSheetService) CreateSheet(sheetsService *sheets.Service, spreadsheetId string, sheetName string) {
	m.Called(sheetsService, spreadsheetId, sheetName)
}

// Mock the PlaceTextCenter method
func (m *MockGoogleSheetService) PlaceTextCenter(title string, values [][]interface{}) {
	m.Called(title, values)
}

// Mock the WriteData method
func (m *MockGoogleSheetService) WriteData(spreadsheetId string, domains []postgresql.DomainForExcel, writeRangeFunc func() string, valueRangeFunc func(domains []postgresql.DomainForExcel) *sheets.ValueRange) error {
	args := m.Called(spreadsheetId, domains, writeRangeFunc, valueRangeFunc)
	return args.Error(0)
}

// func TestCreateExpiredDomainExcel(t *testing.T) {
// 	mockService := new(MockGoogleSheetService)

// 	domains := []postgresql.DomainForExcel{
// 		// Populate this slice with test data
// 	}

// 	// Set expectations for calls to mockService methods
// 	mockService.On("CreateSheetsService", "your_google_api_key.json").Return(&sheets.Service{}, nil)
// 	mockService.On("CreateSheet", mock.Anything, "your_spreadsheet_id", "TestSheetName").Return()
// 	mockService.On("WriteData", "your_spreadsheet_id", mock.Anything, mock.Anything, mock.Anything).Return(nil)

// 	// Call the function under test
// 	CreateExpiredDomainExcel(mockService.(*GoogleSheetService), "TestSheetName", domains)

// 	// Assert that the expected methods were called
// 	mockService.AssertExpectations(t)
// }

// func TestCreateExpiredDomainExce1l(t *testing.T) {
// 	ctrl := gomock.NewController(t)
// 	defer ctrl.Finish()

// 	// Mock GoogleSheetServiceInterface
// 	mockGS := mock.NewMockGoogleSheetServiceInterface(ctrl)

// 	// Set expectations
// 	sheetName := "TestSheet"
// 	domains := []postgresql.DomainForExcel{
// 		// create test data here
// 	}
// 	mockGS.EXPECT().CreateSheetsService(gomock.Any()).Return(nil)
// 	mockGS.EXPECT().CreateSheet(gomock.Any(), gomock.Any(), sheetName)
// 	mockGS.EXPECT().WriteData(
// 		gomock.Any(),
// 		domains,
// 		gomock.Any(),
// 		gomock.Any(),
// 	).Times(2)
// 	mockGS.EXPECT().PlaceTextCenter(gomock.Any(), gomock.Any())

// 	// Call the function under test
// 	//googlesheet.CreateExpiredDomainExcel(mockGS, sheetName, domains)
// }
