package postgresql

import (
	"fmt"
	"tgs-automation/internal/log"
	"time"

	_ "github.com/lib/pq"
)

func (s *GetMomoDataService) GetRegisteredPlayers(brandCode string, startDate string, endDate string, timezoneOffset string) ([]PlayerRegisterInfo, error) {
	query := "select * from report.get_registered_players_by_brand($1, $2, $3);"

	log.LogInfo(fmt.Sprintf("Query: %s", query))
	log.LogInfo(fmt.Sprintf("Brand: %s, StartDate: %s, EndDate: %s", brandCode, startDate, endDate))
	// Execute the query
	rows, err := s.Db.Query(query, brandCode, startDate+timezoneOffset, endDate+timezoneOffset)
	if err != nil {
		log.LogFatal(fmt.Sprintf("Error executing query: %v", err))
		return nil, err
	}
	defer rows.Close()

	var players []PlayerRegisterInfo
	for rows.Next() {
		var pi PlayerRegisterInfo
		// Adjust the Scan based on the date field used
		if err := rows.Scan(&pi.Agent, &pi.Host, &pi.PlayerName, &pi.RealName, &pi.RegisteredOn); err != nil {
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

type PlayerRegisterInfo struct {
	Agent        string
	Host         string // Use sql.NullString for nullable fields
	PlayerName   string
	RealName     string
	RegisteredOn time.Time
}
