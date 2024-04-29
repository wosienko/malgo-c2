package db

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/VipWW/malgo-c2/services/common/entities"
	"github.com/VipWW/malgo-c2/services/common/log"
	internalEntities "github.com/VipWW/malgo-c2/services/malgo-gateway/internal/entities"
	"github.com/VipWW/malgo-c2/services/malgo-gateway/internal/messages/events"
	"github.com/VipWW/malgo-c2/services/malgo-gateway/internal/messages/outbox"
	"github.com/jmoiron/sqlx"
	"google.golang.org/protobuf/types/known/timestamppb"
	"time"
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

func (r *SessionRepository) UpdateSessionHeartbeat(ctx context.Context, sessionId string) error {
	return updateInTx(
		ctx,
		r.db,
		sql.LevelSerializable,
		func(ctx context.Context, tx *sqlx.Tx) error {
			row := tx.QueryRowContext(
				ctx,
				`UPDATE c2_sessions
    				SET heartbeat_at = NOW()
    				WHERE id = $1
    				RETURNING heartbeat_at, project_id
					`,
				sessionId,
			)
			var heartbeatAt time.Time
			var projectId string
			err := row.Scan(&heartbeatAt, &projectId)
			if err != nil {
				log.FromContext(ctx).Warnf("session not found: %s", sessionId)
			}

			outboxPublisher, err := outbox.NewPublisherForDb(ctx, tx)
			if err != nil {
				return fmt.Errorf("could not create event bus: %w", err)
			}

			err = events.NewBus(outboxPublisher).Publish(ctx, &entities.SessionHeartbeatUpdated{
				Header:      entities.NewHeader(),
				SessionId:   sessionId,
				HeartbeatAt: timestamppb.New(heartbeatAt),
				ProjectId:   projectId,
			})
			if err != nil {
				return fmt.Errorf("could not publish event: %w", err)
			}

			return nil
		},
	)
}

func (r *SessionRepository) RegisterNewSession(ctx context.Context, session internalEntities.RegisterNewSession) error {
	return updateInTx(
		ctx,
		r.db,
		sql.LevelSerializable,
		func(ctx context.Context, tx *sqlx.Tx) error {
			sessionName := fmt.Sprintf("session-%s", session.SessionId)
			row := tx.QueryRowxContext(
				ctx,
				`INSERT INTO c2_sessions (id, project_id, name)
					VALUES ($1, $2, $3)
					ON CONFLICT DO NOTHING
					RETURNING created_at, heartbeat_at
					`,
				session.SessionId,
				session.ProjectId,
				sessionName,
			)
			var createdAt time.Time
			var heartbeatAt time.Time
			err := row.Scan(&createdAt, &heartbeatAt)
			if err != nil {
				return fmt.Errorf("could not insert session: %w", err)
			}

			outboxPublisher, err := outbox.NewPublisherForDb(ctx, tx)
			if err != nil {
				return fmt.Errorf("could not create event bus: %w", err)
			}

			err = events.NewBus(outboxPublisher).Publish(ctx, &entities.SessionRegistered{
				Header:      entities.NewHeader(),
				SessionId:   session.SessionId,
				ProjectId:   session.ProjectId,
				Name:        sessionName,
				CreatedAt:   timestamppb.New(createdAt),
				HeartbeatAt: timestamppb.New(heartbeatAt),
			})
			if err != nil {
				return fmt.Errorf("could not publish event: %w", err)
			}

			return nil
		},
	)
}
