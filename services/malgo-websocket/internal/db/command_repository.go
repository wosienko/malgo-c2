package db

import (
	"context"
	"database/sql"
	"fmt"
	externalEntities "github.com/VipWW/malgo-c2/services/common/entities"
	"github.com/VipWW/malgo-c2/services/common/log"
	"github.com/VipWW/malgo-c2/services/malgo-websocket/internal/entities"
	"github.com/VipWW/malgo-c2/services/malgo-websocket/internal/messages/events"
	"github.com/VipWW/malgo-c2/services/malgo-websocket/internal/messages/outbox"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"google.golang.org/protobuf/types/known/timestamppb"
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

func (c *CommandRepository) AddCommand(ctx context.Context, command entities.Command) error {
	return updateInTx(
		ctx,
		c.db,
		sql.LevelReadCommitted,
		func(ctx context.Context, tx *sqlx.Tx) error {
			if command.ID == uuid.Nil {
				command.ID = uuid.New()
			}
			row := tx.QueryRowxContext(
				ctx,
				`INSERT INTO c2_commands (id, session_id, command, operator_id)
				VALUES ($1, $2, $3, $4)
				RETURNING id, session_id, type, status, command, result_size, created_at, operator_id`,
				command.ID.String(),
				command.SessionId.String(),
				command.Command,
				command.OperatorId.String(),
			)
			err := row.StructScan(&command)
			if err != nil {
				return fmt.Errorf("could not insert command: %w", err)
			}

			outboxPublisher, err := outbox.NewPublisherForDb(ctx, tx)
			if err != nil {
				return fmt.Errorf("could not create event bus: %w", err)
			}

			userNameSurnameRow := tx.QueryRowContext(
				ctx,
				`SELECT name, surname FROM users WHERE id = $1`,
				command.OperatorId.String(),
			)

			var name, surname string
			err = userNameSurnameRow.Scan(&name, &surname)
			if err != nil {
				return fmt.Errorf("could not get operator name and surname: %w", err)
			}

			err = events.NewBus(outboxPublisher).Publish(ctx, &externalEntities.CommandCreated{
				Header:       externalEntities.NewHeader(),
				CommandId:    command.ID.String(),
				SessionId:    command.SessionId.String(),
				Type:         command.Type,
				Status:       command.Status,
				Command:      command.Command,
				ResultSize:   int64(command.ResultSize),
				CreatedAt:    timestamppb.New(command.CreatedAt),
				OperatorId:   command.OperatorId.String(),
				OperatorName: fmt.Sprintf("%s %s", name, surname),
			})
			if err != nil {
				return fmt.Errorf("could not publish event: %w", err)
			}

			return nil
		},
	)
}

func (c *CommandRepository) GetCommandByID(ctx context.Context, id string) (*entities.Command, error) {
	var command entities.Command
	row := c.db.QueryRowxContext(
		ctx,
		`SELECT id, session_id, type, status, command, result_size, created_at, operator_id
		FROM c2_commands
		WHERE id = $1`,
		id,
	)
	err := row.StructScan(&command)
	if err != nil {
		return nil, fmt.Errorf("could not get command by id: %w", err)
	}

	return &command, nil
}

func (c *CommandRepository) CancelCommand(ctx context.Context, id string) error {
	return updateInTx(
		ctx,
		c.db,
		sql.LevelReadCommitted,
		func(ctx context.Context, tx *sqlx.Tx) error {
			row := tx.QueryRowxContext(
				ctx,
				`UPDATE c2_commands
    				SET status = 'canceled'
    				WHERE id = $1
    				AND status = 'created'
    				RETURNING session_id
					`,
				id,
			)
			var sessionId string
			err := row.Scan(&sessionId)
			if err != nil {
				log.FromContext(ctx).Warnf("could not cancel command: %v", err)
				return nil
			}

			outboxPublisher, err := outbox.NewPublisherForDb(ctx, tx)
			if err != nil {
				return fmt.Errorf("could not create event bus: %w", err)
			}

			err = events.NewBus(outboxPublisher).Publish(ctx, &externalEntities.CommandStatusModified{
				Header:    externalEntities.NewHeader(),
				CommandId: id,
				Status:    "canceled",
				SessionId: sessionId,
			})
			if err != nil {
				return fmt.Errorf("could not publish event: %w", err)
			}

			return nil
		},
	)
}
