package driver

import (
	"database/sql"
	"time"

	_ "github.com/jackc/pgx/v5"
	_ "github.com/jackc/pgx/v5/pgconn"
	_ "github.com/jackc/pgx/v5/stdlib"
)

// Holds database connection pull
// Needs if we want to use some other DB than PostgreSQL
type DB struct {
	SQL *sql.DB
}

// Variable to
var dbConn = &DB{}

// Declaring some connection pool properties
const maxOpenDbConns = 10
const maxIdleDbConns = 5
const maxDbLifeTime = 5 * time.Minute

// ConnectSQL creates database  pool for Postgres
func ConnectSQL(dsn string) (*DB, error) {
	d, err := NewDataBase(dsn)
	if err != nil {
		panic(err)
	}
	d.SetMaxOpenConns(maxOpenDbConns)
	d.SetMaxIdleConns(maxIdleDbConns)
	d.SetConnMaxLifetime(maxDbLifeTime)

	dbConn.SQL = d

	err = testDB(d)
	if err != nil {
		return nil, err
	}

	return dbConn, nil
}

// testDB pings database
func testDB(d *sql.DB) error {
	err := d.Ping()
	if err != nil {
		return err
	}
	return nil
}

// NewDataBase creates a new database for the application
func NewDataBase(dsn string) (*sql.DB, error) {
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, err
	}

	if err = db.Ping(); err != nil {
		return nil, err

	}
	return db, nil
}
