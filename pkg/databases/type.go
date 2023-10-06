package databases

import (
	"database/sql"

	_ "github.com/lib/pq"
	log "github.com/sirupsen/logrus"
)

type Postgres struct {
	dsn  string
	conn *sql.DB
}

func (p *Postgres) ToMap(query string) (map[string]string, error) {
	defer p.Close()

	// Execute query from provider connection
	rows, err := p.execQuery(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Parse query into Map[string]string
	data := make(map[string]string)
	for rows.Next() {
		var key, value string

		if err := rows.Scan(&key, &value); err != nil {
			return nil, err
		}
		data[key] = value
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return data, nil
}

func (p *Postgres) ToVoid(query string) error {
	defer p.Close()

	// Execute query from provider connection
	res, err := p.execQuery(query)
	if err != nil {
		return err
	}
	defer res.Close()

	return nil
}

func (p *Postgres) Close() {
	p.conn.Close()
}

func (p *Postgres) execQuery(query string) (*sql.Rows, error) {
	err := p.conn.Ping()
	if err != nil {
		return nil, err
	}

	// Execute a simple query
	rows, err := p.conn.Query(query)
	if err != nil {
		return nil, err
	}

	// Iterate over the rows
	return rows, nil

}

func (p *Postgres) connect() error {
	var err error

	p.conn, err = sql.Open("postgres", p.dsn)
	if err := p.conn.Ping(); err != nil {
		log.Error("PGBouncer SQL Ping ", p.conn.Ping())
	}

	return err
}
