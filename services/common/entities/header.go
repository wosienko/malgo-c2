package entities

import (
	"github.com/google/uuid"
	"google.golang.org/protobuf/types/known/timestamppb"
	"time"
)

func NewHeader() *Header {
	return &Header{
		ID:             uuid.NewString(),
		PublishedAt:    timestamppb.New(time.Now()),
		IdempotencyKey: uuid.NewString(),
	}
}

func NewHeaderWithIdempotencyKey(idempotencyKey string) *Header {
	return &Header{
		ID:             uuid.NewString(),
		PublishedAt:    timestamppb.New(time.Now()),
		IdempotencyKey: idempotencyKey,
	}
}
