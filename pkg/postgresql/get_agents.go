package postgresql

import (
	"cdnetwork/internal/log"
	"cdnetwork/internal/util"
	"cdnetwork/pkg/namecheap"
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

type DomainForExcel struct {
	Brand   string
	Agent   string
	Domain  string
	Created string
	Expires string
}

type GetAgentServiceInterface interface {
	GetAgentDomains(domains []namecheap.FilteredDomain) ([]DomainForExcel, error)
}

type GetAgentService struct {
	Db     *sql.DB
	Config util.PostgresqlConfig
}

func New(config util.PostgresqlConfig) GetAgentServiceInterface {
	db, err := NewDataAccessLayer(config)

	if err != nil {
		panic(err)
	}
	return &GetAgentService{
		Db:     db,
		Config: config,
	}
}

func (s *GetAgentService) GetAgentDomains(domains []namecheap.FilteredDomain) ([]DomainForExcel, error) {

	createTempTable(s.Db, domains)
	defer deleteTempTable(s.Db)

	query := `
	SELECT b.code, a.username, td.name, td.created, td.expires
	FROM dbo.agents a
	JOIN dbo.brands b ON a.brand_id = b.id
	JOIN dbo.agent_domains bd ON a.id = bd.agent_id
	JOIN temp_domains td ON bd.domain_name = td.name
	ORDER BY 1,2
	`
	rows, err := s.Db.Query(query)
	fmt.Println("start query")

	if err != nil {
		log.LogError(err.Error())
		return nil, err
	}
	defer rows.Close()

	var result []DomainForExcel
	for rows.Next() {
		var (
			code     string
			username string
			name     string
			created  string
			expires  string
		)
		// Scan the values from the row into variables
		err := rows.Scan(&code, &username, &name, &created, &expires)
		if err != nil {
			//panic(err)
			return nil, err
		}
		// Print the values
		fmt.Printf("brand: %s, agent: %s, domain: %s, created: %s, expired: %s\n", code, username, name, created, expires)

		result = append(result, DomainForExcel{
			Brand:   code,
			Agent:   username,
			Domain:  name,
			Created: created,
			Expires: expires,
		})
	}
	// Check for errors from iterating over rows
	if err = rows.Err(); err != nil {
		log.LogFatal("get domains: " + err.Error())
		//panic(err)
		return nil, err
	}
	return result, nil
}

func createTempTable(db *sql.DB, domains []namecheap.FilteredDomain) {
	fmt.Println("start create temp table")
	createTempSql := `
	CREATE TEMPORARY TABLE temp_domains (
		name VARCHAR(255),
		created VARCHAR(255),
		expires VARCHAR(255)
	);
	`
	_, err := db.Exec(createTempSql)
	if err != nil {
		log.LogFatal(err.Error())
		panic(err)
	}

	fmt.Println("end create temp table")

	// Insert domain data into the temporary table
	for _, domain := range domains {
		fmt.Printf("domain: %s, created: %s, expired: %s\n", domain.Name, domain.Created, domain.Expires)
		_, err := db.Exec(`INSERT INTO temp_domains (name, created, expires) VALUES ($1, $2, $3);`, domain.Name, domain.Created, domain.Expires)
		if err != nil {
			log.LogFatal(err.Error())
			panic(err)
		}
	}

	fmt.Println("end insert temp table")
}

func deleteTempTable(db *sql.DB) {
	fmt.Println("drop temp table")
	sql := `DROP TABLE temp_domains;`
	_, err := db.Exec(sql)
	if err != nil {
		log.LogFatal(err.Error())
		panic(err)
	}
	defer db.Close()
}
