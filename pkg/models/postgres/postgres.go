package postgres

import (
	"github.com/jmoiron/sqlx"
)

type Postgres struct{}

func New(db *sqlx.DB) *Postgres {
	return &Postgres{}
}
