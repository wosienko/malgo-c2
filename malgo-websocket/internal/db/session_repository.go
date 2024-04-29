package db

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/jmoiron/sqlx"
	"malgo-websocket/internal/entities"
	"malgo-websocket/internal/messages/events"
	"malgo-websocket/internal/messages/outbox"
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

func (s *SessionRepository) AddKeyValue(ctx context.Context, value entities.SessionKeyValue) error {
	return updateInTx(
		ctx,
		s.db,
		sql.LevelSerializable,
		func(ctx context.Context, tx *sqlx.Tx) error {
			row := tx.QueryRowxContext(
				ctx,
				`SELECT 1 FROM c2_sessions WHERE id = $1 AND data::jsonb ? $2`,
				value.SessionId,
				value.Key,
			)
			var disposable int
			err := row.Scan(&disposable)
			if err == nil {
				fmt.Printf("Key already exists\n")
				return nil
			}
			if !errors.Is(err, sql.ErrNoRows) {
				return fmt.Errorf("could not check if key exists: %w", err)
			}

			_, err = tx.ExecContext(
				ctx,
				`UPDATE c2_sessions SET data = data || jsonb_build_object($2::varchar, $3::varchar) WHERE id = $1`,
				value.SessionId,
				value.Key,
				value.Value,
			)
			if err != nil {
				return fmt.Errorf("could not update session: %w", err)
			}

			outboxPublisher, err := outbox.NewPublisherForDb(ctx, tx)
			if err != nil {
				return fmt.Errorf("could not create event bus: %w", err)
			}

			err = events.NewBus(outboxPublisher).Publish(ctx, &entities.SessionKeyValueModified{
				Header:    entities.NewHeader(),
				SessionId: value.SessionId,
				Key:       value.Key,
				Value:     value.Value,
			})
			if err != nil {
				return fmt.Errorf("could not publish event: %w", err)
			}
			return nil
		},
	)
}

func (s *SessionRepository) DeleteKeyValue(ctx context.Context, value entities.SessionKeyValue) error {
	return updateInTx(
		ctx,
		s.db,
		sql.LevelSerializable,
		func(ctx context.Context, tx *sqlx.Tx) error {
			_, err := tx.ExecContext(
				ctx,
				`UPDATE c2_sessions SET data = data - $2 WHERE id = $1`,
				value.SessionId,
				value.Key,
			)
			if err != nil {
				return fmt.Errorf("could not update session: %w", err)
			}

			outboxPublisher, err := outbox.NewPublisherForDb(ctx, tx)
			if err != nil {
				return fmt.Errorf("could not create event bus: %w", err)
			}

			err = events.NewBus(outboxPublisher).Publish(ctx, &entities.SessionKeyValueDeleted{
				Header:    entities.NewHeader(),
				SessionId: value.SessionId,
				Key:       value.Key,
			})
			if err != nil {
				return fmt.Errorf("could not publish event: %w", err)
			}
			return nil
		},
	)
}

func (s *SessionRepository) ModifyKeyValue(ctx context.Context, value entities.SessionKeyValue) error {
	return updateInTx(
		ctx,
		s.db,
		sql.LevelSerializable,
		func(ctx context.Context, tx *sqlx.Tx) error {
			_, err := tx.ExecContext(
				ctx,
				`UPDATE c2_sessions SET data = data || jsonb_build_object($2::varchar, $3::varchar) WHERE id = $1`,
				value.SessionId,
				value.Key,
				value.Value,
			)
			if err != nil {
				return fmt.Errorf("could not update session: %w", err)
			}

			outboxPublisher, err := outbox.NewPublisherForDb(ctx, tx)
			if err != nil {
				return fmt.Errorf("could not create event bus: %w", err)
			}

			err = events.NewBus(outboxPublisher).Publish(ctx, &entities.SessionKeyValueModified{
				Header:    entities.NewHeader(),
				SessionId: value.SessionId,
				Key:       value.Key,
				Value:     value.Value,
			})
			if err != nil {
				return fmt.Errorf("could not publish event: %w", err)
			}
			return nil
		},
	)
}

func (s *SessionRepository) RenameSession(ctx context.Context, value entities.SessionName) error {
	return updateInTx(
		ctx,
		s.db,
		sql.LevelSerializable,
		func(ctx context.Context, tx *sqlx.Tx) error {
			_, err := tx.ExecContext(
				ctx,
				`UPDATE c2_sessions SET name = $2 WHERE id = $1`,
				value.SessionId,
				value.Name,
			)
			if err != nil {
				return fmt.Errorf("could not update session: %w", err)
			}

			outboxPublisher, err := outbox.NewPublisherForDb(ctx, tx)
			if err != nil {
				return fmt.Errorf("could not create event bus: %w", err)
			}

			projectId, err := s.GetProjectForSession(ctx, value.SessionId)
			if err != nil {
				return fmt.Errorf("could not get project for session: %w", err)
			}

			err = events.NewBus(outboxPublisher).Publish(ctx, &entities.SessionNameModified{
				Header:    entities.NewHeader(),
				SessionId: value.SessionId,
				Name:      value.Name,
				ProjectId: projectId,
			})
			if err != nil {
				return fmt.Errorf("could not publish event: %w", err)
			}
			return nil
		},
	)
}

func (s *SessionRepository) GetProjectForSession(ctx context.Context, sessionId string) (string, error) {
	var project string
	err := s.db.GetContext(ctx, &project, `SELECT project_id FROM c2_sessions WHERE id = $1`, sessionId)
	if err != nil {
		return "", fmt.Errorf("could not get project for session: %w", err)
	}
	return project, nil
}
