package entities

import (
	"github.com/google/uuid"
	"time"
)

type Command struct {
	ID         uuid.UUID `db:"id" json:"id"`
	SessionId  uuid.UUID `db:"session_id" json:"session_id"`
	Type       string    `db:"type" json:"type"`
	Status     string    `db:"status" json:"status"`
	Command    string    `db:"command" json:"command"`
	OperatorId uuid.UUID `db:"operator_id" json:"operator_id"`
	ResultSize int       `db:"result_size" json:"result_size"`
	CreatedAt  time.Time `db:"created_at" json:"created_at"`
}

type CommandSentToOperator struct {
	MessageType string `json:"message_type"`
	ID          string `json:"id"`
	Type        string `json:"type"`
	Status      string `json:"status"`
	Command     string `json:"command"`
	Operator    string `json:"operator"`
	CreatedAt   string `json:"created_at"`
}
