package dbrepo

import (
	"Ebook/internal/config"
	"database/sql"
)

type postgresDBRepo struct {
	App *config.AppConfig
	DB  *sql.DB
}

func NewPostgresRepo(conn *sql.DB, a *config.AppConfig) *postgresDBRepo {
	return &postgresDBRepo{
		App: a,
		DB:  conn,
	}
}
