package postgresql

import (
	"database/sql"
	"fmt"
	"tgs-automation/internal/log"
	"tgs-automation/internal/util"
	"time"

	_ "github.com/lib/pq"
)

type GetDailyBrandsRevenueInterface interface {
	GetDailyBrandsRevenue() ([]BrandsRevenueModel, error)
	Close()
}

type GetDailyBrandsRevenueService struct {
	Db     *sql.DB
	Config util.PostgresqlConfig
}

func NewDailyBrandsRevenueInterface(config util.PostgresqlConfig) GetDailyBrandsRevenueInterface {
	db, err := NewDataAccessLayer(config)

	if err != nil {
		panic(err)
	}
	return &GetDailyBrandsRevenueService{
		Db:     db,
		Config: config,
	}
}

func (s *GetDailyBrandsRevenueService) Close() {
	s.Db.Close()
}

func (s *GetDailyBrandsRevenueService) GetDailyBrandsRevenue() ([]BrandsRevenueModel, error) {
	defer s.Db.Close()
	query := "select * from report.report_get_daily_brands_revenue();"

	// Execute the query
	rows, err := s.Db.Query(query)
	if err != nil {
		log.LogFatal(fmt.Sprintf("Error executing query: %v", err))
		return nil, err
	}
	defer rows.Close()

	var brands []BrandsRevenueModel
	for rows.Next() {
		var r BrandsRevenueModel
		// Assuming the date in the database is stored in a compatible format; adjust the scan accordingly if it's not.
		if err := rows.Scan(&r.PlatformCode, &r.CurrencyCode, &r.Date, &r.ActiveUserCount, &r.DailyOrderCount, &r.DailyRevenueUSD, &r.CumulativeRevenueUSD); err != nil {
			log.LogFatal(err.Error())
			return nil, err
		}
		brands = append(brands, r)
	}

	// Check for errors from iterating over rows.
	if err := rows.Err(); err != nil {
		log.LogFatal(err.Error())
		return nil, err
	}

	// Example of printing out the reports
	for _, report := range brands {
		fmt.Printf("%+v\n", report)
	}
	return brands, nil
}

type BrandsRevenueModel struct {
	PlatformCode         string    `json:"platform_code"`
	CurrencyCode         string    `json:"currency_code"`
	Date                 time.Time `json:"date"`
	ActiveUserCount      int       `json:"active_user_count"`
	DailyOrderCount      string    `json:"daily_order_count"`
	DailyRevenueUSD      string    `json:"daily_revenue_usd"`
	CumulativeRevenueUSD string    `json:"cumulative_revenue_usd"`
}
