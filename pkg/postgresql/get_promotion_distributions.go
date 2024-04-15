package postgresql

import (
	"database/sql"
	"fmt"
	"tgs-automation/internal/log"
	"tgs-automation/internal/util"
	"time"

	_ "github.com/lib/pq"
)

type GetPromotionDistributionsInterface interface {
	GetData(brand string, startDate string, endDate string) ([]PromotionDistribute, error)
	Close()
}

type GetPromotionDistributionService struct {
	Db     *sql.DB
	Config util.PostgresqlConfig
}

func NewPromotionDistributionInterface(config util.PostgresqlConfig) GetPromotionDistributionsInterface {
	db, err := NewDataAccessLayer(config)

	if err != nil {
		panic(err)
	}
	return &GetPromotionDistributionService{
		Db:     db,
		Config: config,
	}
}

func (s *GetPromotionDistributionService) Close() {
	s.Db.Close()
}

func (s *GetPromotionDistributionService) GetData(brandCode string, startDate string, endDate string) ([]PromotionDistribute, error) {
	defer s.Db.Close()
	query := "select * from report.get_promotion_distributions($1, $2, $3);"

	// Execute the query
	rows, err := s.Db.Query(query, brandCode, startDate, endDate)
	if err != nil {
		log.LogError(fmt.Sprintf("Error executing query: %v", err))
		return nil, err
	}
	defer rows.Close()

	var promotionDistributes []PromotionDistribute
	for rows.Next() {
		var r PromotionDistribute
		// Assuming the date in the database is stored in a compatible format; adjust the scan accordingly if it's not.
		if err := rows.Scan(&r.Username, &r.PromotionName, &r.PromotionType, &r.CreatedOn, &r.BonusAmount, &r.SentOn); err != nil {
			log.LogFatal(err.Error())
			return nil, err
		}
		promotionDistributes = append(promotionDistributes, r)
	}

	// Check for errors from iterating over rows.
	if err := rows.Err(); err != nil {
		log.LogFatal(err.Error())
		return nil, err
	}

	// Example of printing out the reports
	for _, player := range promotionDistributes {
		fmt.Printf("%+v\n", player)
	}
	return promotionDistributes, nil
}

type PromotionDistribute struct {
	Username      string    `json:"username"`
	PromotionName string    `json:"promotion_name"`
	PromotionType string    `json:"promotion_type"`
	CreatedOn     time.Time `json:"created_on"`
	BonusAmount   float64   `json:"bonus_amount"`
	SentOn        time.Time `json:"sent_on"`
}
