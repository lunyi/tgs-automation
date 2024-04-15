package postgresql

import (
	"fmt"
	"tgs-automation/internal/log"
	"time"

	_ "github.com/lib/pq"
)

func (s *GetMomoDataService) GetMomoRegisteredPlayers(brandCode string, startDate string, endDate string, timezoneOffset string) ([]PlayerRegisterInfo, error) {
	query := `
select 
    coalesce(a.username, '') AS agent,
    coalesce(PIR.host,'') as host,
	p.username as player,
	coalesce(p.real_name, '') as real_name,
	p.registered_on
	from dbo.players p
	left join dbo.agents a on a.id = p.agent_id
	left join dbo.player_ip_records AS PIR ON p.player_code = PIR.player_code
			AND PIR.ip_type = 1
	where p.brand_id = (select id from dbo.brands where code = $1) 
	and p.registered_on >= $2 and p.registered_on < $3
	order by 1 nulls first,2,3`

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
