package db

import (
	"context"
	"github.com/jmoiron/sqlx"
)

type ResultRepository struct {
	db *sqlx.DB
}

func NewResultRepository(db *sqlx.DB) *ResultRepository {
	if db == nil {
		panic("db is nil")
	}
	return &ResultRepository{db: db}
}

func (r *ResultRepository) GetResultForCommand(ctx context.Context, commandId string) (string, error) {
	var result string
	rows, err := r.db.QueryxContext(
		ctx,
		`SELECT result_chunk
		 FROM c2_result_chunks
		 WHERE command_id = $1
		 ORDER BY chunk_offset`,
		commandId,
	)
	if err != nil {
		return "", err
	}
	defer rows.Close()

	for rows.Next() {
		// results are in hex, so we need to convert them to string
		var chunk []byte
		err = rows.Scan(&chunk)
		if err != nil {
			return "", err
		}
		result += string(chunk)
	}

	return result, nil
}
