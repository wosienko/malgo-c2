package db

import (
	"context"
	"database/sql"
	"github.com/VipWW/malgo-c2/services/common/entities"
	"github.com/VipWW/malgo-c2/services/common/log"
	internalEntities "github.com/VipWW/malgo-c2/services/malgo-gateway/internal/entities"
	"github.com/VipWW/malgo-c2/services/malgo-gateway/internal/messages/events"
	"github.com/VipWW/malgo-c2/services/malgo-gateway/internal/messages/outbox"
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

func (r *ResultRepository) SetResultLength(ctx context.Context, commandId string, length int) error {
	_, err := r.db.ExecContext(
		ctx,
		`UPDATE c2_commands
		 SET result_size = $1
		 WHERE id = $2`,
		length,
		commandId,
	)
	return err
}

func (r *ResultRepository) AddResultChunk(ctx context.Context, chunk internalEntities.ResultChunk) error {
	return updateInTx(
		ctx,
		r.db,
		sql.LevelReadCommitted,
		func(ctx2 context.Context, tx *sqlx.Tx) error {
			_, err := tx.NamedExecContext(
				ctx2,
				`INSERT INTO c2_result_chunks (command_id, chunk_offset, result_chunk)
    				 VALUES (:command_id, :offset, :chunk)
    				 ON CONFLICT DO NOTHING`,
				chunk,
			)
			if err != nil {
				return err
			}

			// check the total expected length of the result
			row := tx.QueryRowxContext(
				ctx2,
				`SELECT result_size, session_id
				FROM c2_commands
				WHERE id = $1`,
				chunk.CommandId,
			)
			var resultSize int
			var sessionID string
			err = row.Scan(&resultSize, &sessionID)
			if err != nil {
				return err
			}

			log.FromContext(ctx2).Infof("Result size: %d, chunk offset: %d, chunk size: %d, sum: %d", resultSize, chunk.Offset, len(chunk.Chunk), chunk.Offset+len(chunk.Chunk))

			if (resultSize != chunk.Offset+len(chunk.Chunk)) && chunk.Offset != 0 {
				log.FromContext(ctx2).Infof("Not the first or last chunk, skipping")
				return nil
			}

			outboxPublisher, err := outbox.NewPublisherForDb(ctx, tx)
			if err != nil {
				return err
			}

			bus := events.NewBus(outboxPublisher)

			if resultSize == chunk.Offset+len(chunk.Chunk) {
				_, err := tx.ExecContext(
					ctx2,
					`UPDATE c2_commands
    					SET status = 'completed'
    					WHERE id = $1
    					AND status = 'retrieving' OR status = 'sent'`,
					chunk.CommandId,
				)
				if err != nil {
					return err
				}
				err = bus.Publish(ctx2, &entities.CommandStatusModified{
					Header:    entities.NewHeader(),
					CommandId: chunk.CommandId,
					Status:    "completed",
					SessionId: sessionID,
				})
				if err != nil {
					return err
				}
			} else if chunk.Offset == 0 {
				_, err := tx.ExecContext(
					ctx2,
					`UPDATE c2_commands
					SET status = 'retrieving'
					WHERE id = $1
					AND status = 'sent'`,
					chunk.CommandId,
				)
				if err != nil {
					return err
				}

				err = bus.Publish(ctx2, &entities.CommandStatusModified{
					Header:    entities.NewHeader(),
					CommandId: chunk.CommandId,
					Status:    "retrieving",
					SessionId: sessionID,
				})
				if err != nil {
					return err
				}
			}

			return nil
		},
	)
}
