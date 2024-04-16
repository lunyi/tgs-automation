package postgresql

import (
	"database/sql"
	"fmt"
	"tgs-automation/internal/log"
	"tgs-automation/internal/util"
	"time"

	_ "github.com/lib/pq"
)

type GetMomoDataInterface interface {
	GetFirstDepositedPlayers(brandCode string, startDate string, endDate string, timezoneOffset string) ([]PlayerFirstDeposit, error)
	GetRegisteredPlayers(brandCode string, startDate string, endDate string, timezoneOffset string) ([]PlayerRegisterInfo, error)
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

func (s *GetMomoDataService) GetFirstDepositedPlayers(brandCode string, startDate string, endDate string, timezoneOffset string) ([]PlayerFirstDeposit, error) {
	query := "select * from report.get_first_deposite_players_by_brand($1, $2, $3);"

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
	Agent              string
	Host               string
	PlayerName         string
	DailyDepositAmount float64
	DailyDepositCount  int
	FirstDepositOn     time.Time // This will hold either the first deposit or registered date based on your choice
}
