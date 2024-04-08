package postgresql

import (
	"cdnetwork/internal/log"
	"cdnetwork/internal/util"
	"database/sql"
	"fmt"
	"time"

	_ "github.com/lib/pq"
)

type GetMomoDataInterface interface {
	GetMomoFirstDepositePlayers(brandCode string, startDate string, endDate string, timezoneOffset string) ([]PlayerFirstDeposit, error)
	GetMomoRegisteredPlayers(brandCode string, startDate string, endDate string, timezoneOffset string) ([]PlayerRegisterInfo, error)
	Close()
}

type GetMomoDataService struct {
	Db     *sql.DB
	Config util.PostgresqlConfig
}

func NewMomoDataInterface(config util.PostgresqlConfig) GetMomoDataInterface {
	db, err := NewDataAccessLayer(config)

	if err != nil {
		panic(err)
	}
	return &GetMomoDataService{
		Db:     db,
		Config: config,
	}
}

func (s *GetMomoDataService) Close() {
	s.Db.Close()
}

func (s *GetMomoDataService) GetMomoFirstDepositePlayers(brandCode string, startDate string, endDate string, timezoneOffset string) ([]PlayerFirstDeposit, error) {
	query := `
SELECT 
    a.username AS agent,
    PIR.host,
    p.username AS playername, 
    pdp.daily_deposit_amount, 
    pdp.daily_deposit_count,
    pdp.first_deposit_on
FROM 
    dbo.player_daily_payment_statistics pdp
JOIN dbo.players p ON p.player_code = pdp.player_code
LEFT JOIN
    dbo.player_ip_records AS PIR ON P.player_code = PIR.player_code
        AND PIR.ip_type = 1
LEFT JOIN dbo.agents a ON a.id = p.agent_id
WHERE p.brand_id = (SELECT id FROM dbo.brands WHERE code=$1)
AND pdp.first_deposit_on >= $2 AND pdp.first_deposit_on < $3 order by 1 nulls first,2,3
`

	log.LogInfo(fmt.Sprintf("Query: %s", query))
	log.LogInfo(fmt.Sprintf("Brand: %s, StartDate: %s, EndDate: %s", brandCode, startDate, endDate))
	// Execute the query
	rows, err := s.Db.Query(query, brandCode, startDate+timezoneOffset, endDate+timezoneOffset)
	if err != nil {
		log.LogFatal(fmt.Sprintf("Error executing query: %v", err))
		return nil, err
	}
	defer rows.Close()

	var players []PlayerFirstDeposit
	for rows.Next() {
		var pi PlayerFirstDeposit
		// Adjust the Scan based on the date field used
		if err := rows.Scan(&pi.Agent, &pi.Host, &pi.PlayerName, &pi.DailyDepositAmount, &pi.DailyDepositCount, &pi.FirstDepositOn); err != nil {
			log.LogFatal(fmt.Sprintf("Error scanning row: %v", err))
		}
		if !pi.Agent.Valid {
			pi.Agent.String = ""
		}
		players = append(players, pi)
		//fmt.Printf("%+v\n", pi)
	}

	// Check for errors from iterating over rows
	if err = rows.Err(); err != nil {
		log.LogFatal(fmt.Sprintf("Error iterating over rows: %v", err))
		return nil, err
	}
	return players, nil
}

type PlayerFirstDeposit struct {
	Agent              sql.NullString
	Host               string // Use sql.NullString for nullable fields
	PlayerName         string
	DailyDepositAmount float64
	DailyDepositCount  int
	FirstDepositOn     time.Time // This will hold either the first deposit or registered date based on your choice
}
