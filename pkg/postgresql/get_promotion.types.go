package postgresql

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"tgs-automation/internal/log"

	_ "github.com/lib/pq"
)

func (s *GetPromotionService) GetPromotionTypes() ([]Category, error) {
	query := "SELECT value FROM dbo.settings where key = 'PromotionTypes'"

	var value string
	err := s.Db.QueryRow(query).Scan(&value)
	if err != nil {
		if err == sql.ErrNoRows {
			// Handle no rows returned
			log.LogError("No rows were returned")
			return nil, nil
		} else {
			// Handle other errors
			log.LogFatal(fmt.Sprintf("Error executing query: %v", err))
			return nil, err
		}
	}

	// Parse JSON data
	var categories []Category
	err = json.Unmarshal([]byte(value), &categories)
	if err != nil {
		fmt.Println("Error parsing JSON:", err)
		log.LogFatal(fmt.Sprintf("Error parsing JSON: %v", err))
		return nil, err
	}

	return categories, nil
}

type Translation struct {
	En string `json:"En"`
	Zh string `json:"Zh"`
	vi string `json:"vi"`
}

type PromotionType struct {
	Name                   string      `json:"Name"`
	SettlementIntervalType string      `json:"SettlementIntervalType,omitempty"`
	Trans                  Translation `json:"Trans"`
	AutoSentPlayerBonus    *Checkbox   `json:"autoSentPlayerBonus,omitempty"`
	Rule                   *Rule       `json:"rule,omitempty"`
	Unlimited              *Checkbox   `json:"unlimited,omitempty"`
}

type Checkbox struct {
	Checked  bool `json:"checked"`
	Disabled bool `json:"disabled"`
}

type Rule struct {
	PlayerAutoJoin struct {
		AdminUserOnly bool `json:"adminUserOnly"`
		EditProps     struct {
			AfterStartProps struct {
				ReadOnly bool `json:"readOnly"`
				Disabled bool `json:"disabled"`
			} `json:"afterStartProps"`
		} `json:"editProps"`
	} `json:"playerAutoJoin"`
}

type Category struct {
	Name          string          `json:"Name"`
	PromotionType []PromotionType `json:"Types"`
	Trans         Translation     `json:"Trans"`
}
