package postgresql

import (
	"cdnetwork/internal/log"
	"cdnetwork/internal/util"
	"database/sql"
	"fmt"
	"time"

	_ "github.com/lib/pq"
)

type GetBrandsIncomeInterface interface {
	GetBrandsIncome() ([]BrandsIncomeModel, error)
	Close()
}

type GetBrandsIncomeService struct {
	Db     *sql.DB
	Config util.PostgresqlConfig
}

func NewBrandsIncomeInterface(config util.PostgresqlConfig) GetBrandsIncomeInterface {
	db, err := NewDataAccessLayer(config)

	if err != nil {
		panic(err)
	}
	return &GetBrandsIncomeService{
		Db:     db,
		Config: config,
	}
}

func (s *GetBrandsIncomeService) Close() {
	s.Db.Close()
}

func (s *GetBrandsIncomeService) GetBrandsIncome() ([]BrandsIncomeModel, error) {
	defer s.Db.Close()
	query := `
	WITH ExchangeRates AS (
		SELECT 
			'PHP' AS currency_code, 
			0.0177 AS rate_to_usd
		UNION ALL 
		SELECT 'HKD', 0.1277
		UNION ALL 
		SELECT 'VND_1000', 0.0398
	), SubQuery AS (
		SELECT 
			b.code as 平台,
			bc.currency_code as 幣別,
			((p.report_date::date + INTERVAL '1 day' + INTERVAL '8 hours')::timestamp)::date as 日期,
			count(distinct p.player_code) as 活躍人數, 
			SUM(p.round_count) as 當日訂單數量_raw,
			ROUND(SUM(p.win_loss_amount * er.rate_to_usd), 2) as 當日營收美金_raw,
			sum(ROUND(sum(p.win_loss_amount * er.rate_to_usd),2)) 
				OVER (PARTITION BY b.code ORDER BY ((p.report_date::date + INTERVAL '1 day' + INTERVAL '8 hours')::timestamp)::date ASC) 
				as 當月累計營收美金
		FROM report.player_aggregates p
		JOIN dbo.brands b ON b.id = p.brand_id
		JOIN dbo.brand_currencies bc ON b.id = bc.brand_id
		LEFT JOIN ExchangeRates er ON bc.currency_code = er.currency_code
		WHERE 
			p.report_date >= DATE_TRUNC('month', CURRENT_DATE) + INTERVAL '8 hours'  
			AND p.report_date < DATE_TRUNC('month', CURRENT_DATE) + INTERVAL '1 month' + INTERVAL '8 hours'  
			AND p.bet_amount > 0
			AND b.code NOT IN ('Sky8', 'GPI')
		GROUP BY 
			b.code,
			bc.currency_code,
			((p.report_date::date + INTERVAL '1 day' + INTERVAL '8 hours')::timestamp)::date
	)

	SELECT 
		'Total' as 平台,
		'' as 幣別,
		CURRENT_DATE - INTERVAL '1 day' as 日期,
		SUM(活躍人數) as 活躍人數,
		TO_CHAR(SUM(當日訂單數量_raw), 'FM999,999,999,999') as 當日訂單數量,
		TO_CHAR(SUM(當日營收美金_raw), 'FM999,999,999,999.99') as 當日營收美金,
		TO_CHAR(SUM(當月累計營收美金), 'FM999,999,999,999.99') as 當月累計營收美金
	FROM SubQuery
	WHERE 日期 = CURRENT_DATE - INTERVAL '1 day'
	UNION ALL
	SELECT 
		平台,
		幣別,
		日期,
		活躍人數,
		TO_CHAR(當日訂單數量_raw, 'FM999,999,999,999') as 當日訂單數量,
		TO_CHAR(當日營收美金_raw, 'FM999,999,999,999.99') as 當日營收美金,
		TO_CHAR(當月累計營收美金, 'FM999,999,999,999.99') as 當月累計營收美金
	FROM SubQuery
	WHERE 日期 = CURRENT_DATE - INTERVAL '1 day';
	
`

	// Execute the query
	rows, err := s.Db.Query(query)
	if err != nil {
		log.LogFatal(fmt.Sprintf("Error executing query: %v", err))
		return nil, err
	}
	defer rows.Close()

	var brands []BrandsIncomeModel
	for rows.Next() {
		var r BrandsIncomeModel
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

type BrandsIncomeModel struct {
	PlatformCode         string    `json:"platform_code"`
	CurrencyCode         string    `json:"currency_code"`
	Date                 time.Time `json:"date"`
	ActiveUserCount      int       `json:"active_user_count"`
	DailyOrderCount      string    `json:"daily_order_count"`
	DailyRevenueUSD      string    `json:"daily_revenue_usd"`
	CumulativeRevenueUSD string    `json:"cumulative_revenue_usd"`
}
