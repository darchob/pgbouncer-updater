package databases

import (
	"database/sql"
	"log"
	"reflect"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	_ "github.com/lib/pq"
	"github.com/stretchr/testify/assert"
)

var pg_authid = struct {
	rolname     string
	rolpassword string
}{
	rolname:     "postgres",
	rolpassword: "postgres",
}

func NewMock() (*sql.DB, sqlmock.Sqlmock) {
	db, mock, err := sqlmock.New()
	if err != nil {
		log.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	return db, mock
}

func TestPostgres_ToMap(t *testing.T) {
	db, mock := NewMock()
	repo := Postgres{
		dsn:  "sqlmock_db_0",
		conn: db,
	}

	query := "select rolname,rolpassword from pg_authid where rolpassword is not null order by rolname asc"

	rows := sqlmock.NewRows([]string{"rolname", "rolpassword"}).AddRow(pg_authid.rolname, pg_authid.rolpassword)

	want := map[string]string{
		"postgres": "postgres",
	}

	mock.ExpectQuery(query).WillReturnRows(rows)

	roles, err := repo.ToMap(query)
	assert.NotNil(t, roles)
	assert.NoError(t, err)

	if !reflect.DeepEqual(roles, want) {
		t.Errorf("Databases.ToMap(string) error = %v", err)
	}
}

func TestPostgres_ToVoid(t *testing.T) {
	db, mock := NewMock()
	repo := Postgres{
		dsn:  "sqlmock_db_0",
		conn: db,
	}

	query := "select rolname,rolpassword from pg_authid where rolpassword is not null order by rolname asc"

	rows := sqlmock.NewRows([]string{"rolname", "rolpassword"}).AddRow(pg_authid.rolname, pg_authid.rolpassword)

	mock.ExpectQuery(query).WillReturnRows(rows)

	 err := repo.ToVoid(query)
	assert.NoError(t, err)
}
