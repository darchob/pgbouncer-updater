package databases

const DefaultQuery = "select rolname,rolpassword from pg_authid where rolpassword is not null order by rolname asc"

type Databases interface {
	// Query results to a map
	ToMap(query string) (map[string]string, error)
	// Exec query with no results excepted
	ToVoid(query string) error
	Close()
}

func NewQuery(dsn string) (Databases, error) {
	cred := Postgres{
		dsn: dsn,
	}

	if err := cred.connect(); err != nil {
		return nil, err
	}

	return &cred, nil
}
