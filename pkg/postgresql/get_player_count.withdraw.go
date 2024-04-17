package postgresql

import (
	"tgs-automation/internal/util"

	_ "github.com/lib/pq"
)

// 提款人數
type WithdrawPlayerCountService struct {
	Config util.PostgresqlConfig
}

func (bpc *WithdrawPlayerCountService) GetPlayerCount(brandCode, startDate, endDate string) (int, error) {
	query := "SELECT report.get_player_withdraw_count($1, $2, $3);"
	return getPlayerCount(bpc.Config, query, brandCode, startDate, endDate)
}
