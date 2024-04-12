package postgresql

import (
	"cdnetwork/internal/log"
	"cdnetwork/internal/util"
	"database/sql"
	"fmt"
	"time"

	_ "github.com/lib/pq"
)

type GetPlayersAdjustAmountInterface interface {
	GetData(brandCode string, startDate string, endDate string, transType int) ([]PlayerAdjustAmountData, error)
	Close()
}

type GetPlayersAdjustAmountService struct {
	Db     *sql.DB
	Config util.PostgresqlConfig
}

func NewGetPlayersAdjustAmountInterface(config util.PostgresqlConfig) GetPlayersAdjustAmountInterface {
	db, err := NewDataAccessLayer(config)

	if err != nil {
		panic(err)
	}
	return &GetPlayersAdjustAmountService{
		Db:     db,
		Config: config,
	}
}

func (s *GetPlayersAdjustAmountService) Close() {
	s.Db.Close()
}

func (s *GetPlayersAdjustAmountService) GetData(brandCode string, startDate string, endDate string, transType int) ([]PlayerAdjustAmountData, error) {
	defer s.Db.Close()
	query := "select * from report.get_players_adjust_amount($1, $2, $3, $4);"

	// Execute the query
	rows, err := s.Db.Query(query, brandCode, startDate, endDate, transType)
	if err != nil {
		log.LogError(fmt.Sprintf("Error executing query: %v", err))
		return nil, err
	}
	defer rows.Close()

	var players []PlayerAdjustAmountData
	for rows.Next() {
		var r PlayerAdjustAmountData
		// Assuming the date in the database is stored in a compatible format; adjust the scan accordingly if it's not.
		if err := rows.Scan(&r.PlayerName, &r.Amount, &r.BeforeBalance, &r.ExecutionTime, &r.ExecutionTime, &r.Executor, &r.Description); err != nil {
			log.LogFatal(err.Error())
			return nil, err
		}
		players = append(players, r)
	}

	// Check for errors from iterating over rows.
	if err := rows.Err(); err != nil {
		log.LogFatal(err.Error())
		return nil, err
	}

	// Example of printing out the reports
	for _, player := range players {
		fmt.Printf("%+v\n", player)
	}
	return players, nil
}

type PlayerAdjustAmountData struct {
	PlayerName    string    `json:"玩家用戶名"` // Username of the player
	Amount        float64   `json:"活動紅利"`  // Bonus amount
	BeforeBalance float64   `json:"派發前餘額"` // Balance before distribution
	AfterBalance  float64   `json:"派發後餘額"` // Balance after distribution
	ExecutionTime time.Time `json:"執行時間"`  // Time of record
	Executor      string    `json:"執行者"`   // Executor, who carried out the transaction
	Description   string    `json:"描述"`    // Description of the transaction
}
