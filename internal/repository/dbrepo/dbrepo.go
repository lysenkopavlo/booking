package dbrepo

import (
	"database/sql"

	"github.com/lysenkopavlo/booking/internal/config"
	"github.com/lysenkopavlo/booking/internal/repository"
)

type postgresDbRepo struct {
	DB  *sql.DB
	App *config.AppConfig
}

func NewPostgresRepo(conn *sql.DB, a *config.AppConfig) repository.DataBaseRepo {
	return &postgresDbRepo{
		DB:  conn,
		App: a,
	}
}
