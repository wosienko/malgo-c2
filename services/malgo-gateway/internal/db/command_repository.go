package db

import (
	"context"
	"database/sql"
	"github.com/VipWW/malgo-c2/services/common/entities"
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
