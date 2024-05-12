package db

import (
	"context"
	"database/sql"
	"math"
	"time"

	"github.com/VipWW/malgo-c2/services/common/entities"
	"github.com/VipWW/malgo-c2/services/common/log"
	internalEntities "github.com/VipWW/malgo-c2/services/malgo-gateway/internal/entities"
	"github.com/VipWW/malgo-c2/services/malgo-gateway/internal/messages/events"
	"github.com/VipWW/malgo-c2/services/malgo-gateway/internal/messages/outbox"
	"github.com/jmoiron/sqlx"
	"google.golang.org/protobuf/types/known/timestamppb"
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
	return updateInTx(
		ctx,
		r.db,
		sql.LevelReadCommitted,
		func(ctx2 context.Context, tx *sqlx.Tx) error {
			_, err := r.db.ExecContext(
				ctx2,
				`UPDATE c2_commands
				 SET result_size = $1
				 WHERE id = $2`,
				length,
				commandId,
			)
			if err != nil {
				return err
			}

			if length != 0 {
				return nil
			}

			row := tx.QueryRowxContext(
				ctx2,
				`UPDATE c2_commands
				 SET status = 'completed'
				 WHERE id = $1
				 RETURNING session_id
				 `,
				commandId,
			)
			var sessionID string
			err = row.Scan(&sessionID)
			if err != nil {
				return err
			}

			outboxPublisher, err := outbox.NewPublisherForDb(ctx, tx)
			if err != nil {
				return err
			}

			bus := events.NewBus(outboxPublisher)
			err = bus.Publish(ctx2, &entities.CommandStatusModified{
				Header:    entities.NewHeader(),
				CommandId: commandId,
				Status:    "completed",
				SessionId: sessionID,
			})
			if err != nil {
				return err
			}

			return nil
		},
	)
}

func (r *ResultRepository) AddResultChunk(ctx context.Context, chunk internalEntities.ResultChunk) error {
	return updateInTx(
		ctx,
		r.db,
		sql.LevelReadCommitted,
		func(ctx2 context.Context, tx *sqlx.Tx) error {
			row := tx.QueryRowxContext(
				ctx2,
				`INSERT INTO c2_result_chunks (command_id, chunk_offset, result_chunk)
    				 VALUES ($1, $2, $3)
    				 ON CONFLICT DO NOTHING
    				 RETURNING created_at`,
				chunk.CommandId,
				chunk.Offset,
				chunk.Chunk,
			)
			var createdAt time.Time
			err := row.Scan(&createdAt)
			if err != nil {
				if isErrorUniqueViolation(err) {
					return nil
				}
				return err
			}

			// check the total expected length of the result
			row = tx.QueryRowxContext(
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

			outboxPublisher, err := outbox.NewPublisherForDb(ctx, tx)
			if err != nil {
				return err
			}

			bus := events.NewBus(outboxPublisher)
			// Publish information about inserting a new chunk
			err = bus.Publish(ctx2, &entities.ResultChunkInserted{
				Header:    entities.NewHeader(),
				SessionId: sessionID,
				CommandId: chunk.CommandId,
				CreatedAt: timestamppb.New(createdAt),
				Progress:  int64(math.Round(float64(chunk.Offset+len(chunk.Chunk)) / float64(resultSize) * 100)),
			})
			if err != nil {
				return err
			}

			// Publish information about the first and last chunk
			if (resultSize != chunk.Offset+len(chunk.Chunk)) && chunk.Offset != 0 {
				log.FromContext(ctx2).Infof("Not the first or last chunk, skipping")
				return nil
			}

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
