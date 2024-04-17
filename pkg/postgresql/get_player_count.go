package postgresql

import (
	"fmt"
	"tgs-automation/internal/log"
	"tgs-automation/internal/util"

	_ "github.com/lib/pq"
)

type PlayerCountProvider interface {
	GetPlayerCount(brandCode, startDate, endDate string) (int, error)
}

func getPlayerCount(config util.PostgresqlConfig, query string, brandCode string, startDate string, endDate string) (int, error) {
	db, _ := NewDataAccessLayer(config)
	defer db.Close()
	var count int

	err := db.QueryRow(query, brandCode, startDate, endDate).Scan(&count)
	if err != nil {
		log.LogError(fmt.Sprintf("Error executing query: %v", err))
		return -1, err
	}

	return count, nil
}
