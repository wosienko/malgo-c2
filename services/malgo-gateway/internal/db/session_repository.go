package db

import (
	"github.com/jmoiron/sqlx"
)

type SessionRepository struct {
	db *sqlx.DB
}

func NewSessionRepository(db *sqlx.DB) *SessionRepository {
	if db == nil {
		panic("db is nil")
	}
	return &SessionRepository{db: db}
}
