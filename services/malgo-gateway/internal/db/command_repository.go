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

type CommandRepository struct {
	db *sqlx.DB
}

func NewCommandRepository(db *sqlx.DB) *CommandRepository {
	if db == nil {
		panic("db is nil")
	}

	return &CommandRepository{db: db}
}

func (r *CommandRepository) GetCommandInfo(ctx context.Context, sessionId string) (*internalEntities.CommandInfo, error) {
	var commandInfo internalEntities.CommandInfo

	err := updateInTx(
		ctx,
		r.db,
		sql.LevelSerializable,
		func(ctx2 context.Context, tx *sqlx.Tx) error {
			row := tx.QueryRowxContext(
				ctx2,
				`SELECT
			id,
			type,
			octet_length(command) AS length
			FROM c2_commands
			WHERE session_id = $1
			AND status = 'created'
			ORDER BY created_at
			LIMIT 1
			`,
				sessionId,
			)
			err := row.StructScan(&commandInfo)
			if err != nil {
				return err
			}

			_, err = tx.ExecContext(
				ctx2,
				`UPDATE c2_commands
    			SET status = 'queried'
    			WHERE id = $1
    			AND status = 'created'`,
				commandInfo.ID,
			)
			if err != nil {
				return err
			}

			outboxPublisher, err := outbox.NewPublisherForDb(ctx, tx)
			if err != nil {
				return err
			}

			err = events.NewBus(outboxPublisher).Publish(
				ctx,
				&entities.CommandStatusModified{
					Header:    entities.NewHeader(),
					CommandId: commandInfo.ID,
					Status:    "queried",
					SessionId: sessionId,
				},
			)

			return nil
		},
	)
	if err != nil {
		return nil, err
	}

	return &commandInfo, nil
}

func (r *CommandRepository) GetCommandChunk(ctx context.Context, query *internalEntities.CommandChunkQuery) (*internalEntities.CommandChunk, error) {
	var commandChunk internalEntities.CommandChunk

	err := updateInTx(
		ctx,
		r.db,
		sql.LevelSerializable,
		func(ctx2 context.Context, tx *sqlx.Tx) error {
			row := tx.QueryRowxContext(
				ctx2,
				`SELECT
					id,
					session_id,
					substring(command from $1::INTEGER for $2::INTEGER) AS data,
					(octet_length(command) < $1::INTEGER + $2::INTEGER) AS is_last
					FROM c2_commands
					WHERE id = $3
			`,
				query.Offset+1,
				query.Length,
				query.CommandID,
			)
			err := row.StructScan(&commandChunk)
			if err != nil {
				return err
			}

			outboxPublisher, err := outbox.NewPublisherForDb(ctx, tx)
			if err != nil {
				return err
			}

			bus := events.NewBus(outboxPublisher)

			if query.Offset == 0 {
				_, err = tx.ExecContext(
					ctx2,
					`UPDATE c2_commands
    					SET status = 'sending'
    					WHERE id = $1
    					AND status = 'queried'`,
					query.CommandID,
				)
				if err != nil {
					log.FromContext(ctx2).Warnf("failed to update command status: %v", err)
				}

				err = bus.Publish(
					ctx2,
					&entities.CommandStatusModified{
						Header:    entities.NewHeader(),
						CommandId: query.CommandID,
						SessionId: commandChunk.SessionID,
						Status:    "sending",
					},
				)
				if err != nil {
					return err
				}
			}

			if commandChunk.IsLast {
				_, err = tx.ExecContext(
					ctx2,
					`UPDATE c2_commands
    					SET status = 'sent'
    					WHERE id = $1
    					AND status = 'sending'`,
					query.CommandID,
				)
				if err != nil {
					log.FromContext(ctx2).Warnf("failed to update command status: %v", err)
				}

				err = bus.Publish(
					ctx2,
					&entities.CommandStatusModified{
						Header:    entities.NewHeader(),
						CommandId: query.CommandID,
						SessionId: commandChunk.SessionID,
						Status:    "sent",
					},
				)
				if err != nil {
					return err
				}

			}

			commandChunk.Length = len(commandChunk.Data)

			return nil
		},
	)
	if err != nil {
		return nil, err
	}

	return &commandChunk, nil
}
