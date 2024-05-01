package events

import (
	"context"
	"fmt"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/VipWW/malgo-c2/services/common/entities"
	"github.com/VipWW/malgo-c2/services/malgo-websocket/internal/ws"
	"github.com/google/uuid"
	"google.golang.org/protobuf/proto"
)

func (h *Handler) SendResultToWebsocket(ctx context.Context, command *entities.CommandStatusModified) error {
	if command.Status != "completed" {
		return nil
	}

	cmdFromDB, err := h.commandRepo.GetCommandByID(ctx, command.CommandId)
	if err != nil {
		return fmt.Errorf("could not get command from DB: %v", err)
	}

	if cmdFromDB.Type != "command" {
		return nil
	}

	result, err := h.resultRepo.GetResultForCommand(ctx, command.CommandId)
	if err != nil {
		return fmt.Errorf("could not get result for command from DB: %v", err)
	}

	payload, _ := proto.Marshal(&entities.CommandResultRetrieved{
		Header:    command.Header,
		CommandId: command.CommandId,
		SessionId: command.SessionId,
		Result:    result,
	})

	msg := message.NewMessage(
		uuid.NewString(),
		payload,
	)
	return h.pubSub.Publish(ws.CommandResultRetrievedTopic, msg)
}
