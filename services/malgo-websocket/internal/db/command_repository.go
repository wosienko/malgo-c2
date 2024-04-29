package db

import (
	"context"
	"database/sql"
	"fmt"
	externalEntities "github.com/VipWW/malgo-c2/services/common/entities"
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
		sql.LevelSerializable,
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
