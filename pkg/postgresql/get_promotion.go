package postgresql

import (
	"database/sql"
	"tgs-automation/internal/util"

	_ "github.com/lib/pq"
)

type GetPromotionInterface interface {
	GetPromotionDistributions(brand string, startDate string, endDate string) ([]PromotionDistribute, error)
	GetPromotionTypes() ([]Category, error)
	Close()
}

type GetPromotionService struct {
	Db     *sql.DB
	Config util.PostgresqlConfig
}

func (s *GetPromotionService) Close() {
	s.Db.Close()
}
func NewPromotionInterface(config util.PostgresqlConfig) GetPromotionInterface {
	db, err := NewDataAccessLayer(config)

	if err != nil {
		panic(err)
	}
	return &GetPromotionService{
		Db:     db,
		Config: config,
	}
}
