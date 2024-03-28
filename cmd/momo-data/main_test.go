package main

import (
	"testing"

	"github.com/stretchr/testify/mock"
)

func TestGetMomoDataExcels(t *testing.T) {
	record := map[string]string{
		"Field": "exampleField",
		"File":  "exampleFile",
	}

	app := &mockGetMomoDataInterface{}
	brand := "exampleBrand"
	yesterday := "exampleYesterday"
	today := "exampleToday"
	filenames := []string{}

	expectedFilename := "exampleFilename"

	app.On("GetMomodata", brand, record["Field"], yesterday, today, "+08:00").Return([]string{"player1", "player2"}, nil)
	CreateExcelMock.On("CreateExcel", []string{"player1", "player2"}, brand, record["Column"], record["File"]).Return(expectedFilename, nil)

	filenames = getMomoDataExcels(record, app, brand, yesterday, today, filenames)

	if len(filenames) != 1 {
		t.Errorf("Expected 1 filename, got %d", len(filenames))
	}

	if filenames[0] != expectedFilename {
		t.Errorf("Expected filename %s, got %s", expectedFilename, filenames[0])
	}

	app.AssertExpectations(t)
	CreateExcelMock.AssertExpectations(t)
}

type mockGetMomoDataInterface struct {
	mock.Mock
}

func (m *mockGetMomoDataInterface) GetMomodata(brand, field, yesterday, today, timezone string) ([]string, error) {
	args := m.Called(brand, field, yesterday, today, timezone)
	return args.Get(0).([]string), args.Error(1)
}

type mockCreateExcel struct {
	mock.Mock
}

func (m *mockCreateExcel) CreateExcel(players []string, brand, column, file string) (string, error) {
	args := m.Called(players, brand, column, file)
	return args.String(0), args.Error(1)
}

var CreateExcelMock = &mockCreateExcel{}
