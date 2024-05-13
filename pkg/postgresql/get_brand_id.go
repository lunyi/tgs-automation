package postgresql

import (
	"fmt"
	"tgs-automation/internal/log"
	"tgs-automation/internal/util"

	_ "github.com/lib/pq"
)

func GetBrandId(config util.PostgresqlConfig, brandCode string) (string, error) {
	db, _ := NewDataAccessLayer(config)
	defer db.Close()
	var brandId string
	query := `SELECT id FROM dbo.brands WHERE code = $1`
	err := db.QueryRow(query, brandCode).Scan(&brandId)

	if err != nil {
		log.LogError(fmt.Sprintf("Error executing query: %v", err))
		return "", err
	}

	return brandId, nil
}
