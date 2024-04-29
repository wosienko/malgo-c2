package outbox

import (
	"context"
	"fmt"
	watermillSQL "github.com/ThreeDotsLabs/watermill-sql/v2/pkg/sql"
	"github.com/ThreeDotsLabs/watermill/components/forwarder"
	"github.com/ThreeDotsLabs/watermill/message"
	log2 "github.com/VipWW/malgo-c2/services/common/log"
	"github.com/VipWW/malgo-c2/services/malgo-websocket/internal/observability"
	"github.com/jmoiron/sqlx"
)

func NewPublisherForDb(ctx context.Context, db *sqlx.Tx) (message.Publisher, error) {
	var publisher message.Publisher

	logger := log2.NewWatermill(log2.FromContext(ctx))

	publisher, err := watermillSQL.NewPublisher(
		db,
		watermillSQL.PublisherConfig{
			SchemaAdapter: watermillSQL.DefaultPostgreSQLSchema{},
		},
		logger,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create outbox publisher: %w", err)
	}
	publisher = log2.CorrelationPublisherDecorator{publisher}
	publisher = observability.TracingPublisherDecorator{publisher}

	publisher = forwarder.NewPublisher(publisher, forwarder.PublisherConfig{
		ForwarderTopic: outboxTopic,
	})
	publisher = log2.CorrelationPublisherDecorator{publisher}
	publisher = observability.TracingPublisherDecorator{publisher}

	return publisher, nil
}
