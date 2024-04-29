package db

import (
	"github.com/jmoiron/sqlx"
)

type CommandRepository struct {
	db *sqlx.DB
}

func NewCommandRepository(db *sqlx.DB) *CommandRepository {
	if db == nil {
		panic("db is nil")
	}

	return &CommandRepository{db: db}
}
