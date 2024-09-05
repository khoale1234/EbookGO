package driver

import (
	"database/sql"
	"time"

	_ "github.com/jackc/pgx/v4"
	_ "github.com/jackc/pgx/v4/stdlib"
)

type DB struct {
	SQL *sql.DB
}

var dbConn = &DB{}

const maxOpenDbConn = 10
const maxIdDbConn = 5
const maxDbLifeTime = 5 * time.Minute

func ConnectSQL(dsn string) (*DB, error) {
	db, err := NewDatabase(dsn)
	if err != nil {
		panic(err)
	}
	db.SetMaxIdleConns(maxIdDbConn)
	db.SetConnMaxLifetime(maxDbLifeTime)
	db.SetMaxOpenConns(maxOpenDbConn)
	dbConn.SQL = db
	err = TestDb(db)
	if err != nil {
		return nil, err
	}
	return dbConn, nil
}
func TestDb(db *sql.DB) error {
	err := db.Ping()
	if err != nil {
		return err
	}
	return nil
}
func NewDatabase(dsn string) (*sql.DB, error) {
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, err
	}
	if err = db.Ping(); err != nil {
		return nil, err
	}
	return db, err
}
