package postgresql

import (
	"testing"
	"tgs-automation/internal/util"
	"tgs-automation/pkg/namecheap"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	_ "github.com/lib/pq"
)

func TestGetAgentDomains(t *testing.T) {
	// Create a new mock database
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	// Set up test data
	testDomains := []namecheap.FilteredDomain{
		{
			Name:    "example.com",
			Created: time.Now().Format("2006-01-02"),
			Expires: time.Now().AddDate(1, 0, 0).Format("2006-01-02"),
		},
	}

	// Expect CREATE TEMPORARY TABLE query
	mock.ExpectExec("CREATE TEMPORARY TABLE temp_domains").
		WillReturnResult(sqlmock.NewResult(0, 0))

	// Expect INSERT INTO query for each domain
	for _, domain := range testDomains {
		mock.ExpectExec("INSERT INTO temp_domains").
			WithArgs(domain.Name, domain.Created, domain.Expires).
			WillReturnResult(sqlmock.NewResult(0, 1))
	}

	// Set up expected rows
	rows := sqlmock.NewRows([]string{"code", "username", "name", "created", "expires"}).
		AddRow("brand1", "agent1", "example.com", time.Now().Format("2006-01-02"), time.Now().AddDate(1, 0, 0).Format("2006-01-02"))

		// Mock the query
	mock.ExpectQuery("SELECT").WillReturnRows(rows)

	// Expect DROP TABLE query
	mock.ExpectExec("DROP TABLE temp_domains").
		WillReturnResult(sqlmock.NewResult(0, 0))

	// Create GetAgentService with mock database
	service := GetAgentService{
		Db:     db,
		Config: util.PostgresqlConfig{},
	}

	// Call the method we are testing
	result, err := service.GetAgentDomains(testDomains)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	// Check the result
	if len(result) != 1 {
		t.Errorf("unexpected result length: got %d, expected %d", len(result), 1)
	}

	// Check individual result fields
	expected := DomainForExcel{
		Brand:   "brand1",
		Agent:   "agent1",
		Domain:  "example.com",
		Created: time.Now().Format("2006-01-02"),
		Expires: time.Now().AddDate(1, 0, 0).Format("2006-01-02"),
	}
	if result[0] != expected {
		t.Errorf("unexpected result: got %+v, expected %+v", result[0], expected)
	}

	// Check mock expectations
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}
