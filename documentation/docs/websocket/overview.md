# Websocket Server

Websocket Server is the server handling processing information to and from the operator.

## Adding New Commands / Events

### Create a New Protobuf Definition

In `services/common/entities`, add a new Protobuf definition to either `comands.proto` or `events.proto`. See the existing definitions for examples.

!!! note
    Do not forget to regenerate the Protobuf files by running `task services:proto` in the root directory.

### Create New Repository

If necessary, create a new repository in `services/malgo-websocket/db/` that will handle database operations for the new command/event.

!!! note
    If the new command/event does not require any database operations or you can reuse existing functionality, you can skip this step.

!!! warning
    If you want to emit an event that something has happened (e.g. C2 command was created), you need to use Outbox Pattern.

### Create New Handler

Depending on the command or event, create a new handler in `services/malgo-websocket/messages/commands` or `services/malgo-websocket/messages/events`.

Exemplary handler for a new command:

```go
func (h *Handler) SendModifiedSessionKeyValueToWebsocket(ctx context.Context, keyValue *entities.SessionKeyValueModified) error {
    log.FromContext(ctx).Info("Sending modified key-value through websockets")

    payload, err := proto.Marshal(keyValue)
    if err != nil {
        return fmt.Errorf("could not marshal key-value into protobuf: %v", err)
    }

    msg := message.NewMessage(
        uuid.NewString(),
        payload,
    )
    return h.pubSub.Publish(ws.SessionKeyValueModifiedTopic, msg)
}
```

As you might notice, the handler arguments are always context.Context and, when using Protobuf, a pointer to the Protobuf message.
Watermill library automatically handles message serialization and deserialization.

!!! note
    Should handler fail, the event will be requeued and retried.

### Add New Handler to the Router

Do not forget to add the new handler to the router in `services/malgo-websocket/messages/router.go`.

```go
err = eventProcessor.AddHandlers(
    cqrs.NewEventHandler(
        "SendNewCommandToWebsocket",
        eventHandler.SendNewCommandsToWebsocket,
    ),
)

err = commandProcessor.AddHandlers(
    cqrs.NewCommandHandler(
        "CreateCommand",
        commandHandler.CreateCommand,
    ),
)
```

## Handling Websocket Messages

### Reading Messages

Reading messages happens in `services/malgo-websocket/ws/readFromWebsocket.go`. Over there, based on the `type` field in JSON, the message is routed to appropriate websocket handler.

### Writing Messages

Writing messages is a bit more complex. To send a message through websocket, Watermill Router Handlers need to publish a message to PubSub queue.

```go hl_lines="13"
func (h *Handler) SendNewCommandsToWebsocket(ctx context.Context, command *entities.CommandCreated) error {
    log.FromContext(ctx).Info("Sending new command through websockets")

    payload, err := proto.Marshal(command)
    if err != nil {
        return fmt.Errorf("could not marshal command into protobuf: %v", err)
    }

    msg := message.NewMessage(
        uuid.NewString(),
        payload,
    )
    return h.pubSub.Publish(ws.NewCommandsTopic, msg)
}
```

It is then received in `services/malgo-websocket/ws/writeToWebsocket.go`, handled and sent to the operator.

```go
newCommandsChannel, err := h.pubSub.Subscribe(ctx, NewCommandsTopic)
if err != nil {
    log.FromContext(context.Background()).Errorf("error listening on new commands channel")
}
for {
    select {
    case message := <-newCommandsChannel:
        if err := h.handleNewCommands(message.Payload); err != nil {
            log.FromContext(context.Background()).Errorf("Error handling new command: %v", err)
        }
        message.Ack()
    }
}
```

```go
func (h *Handler) handleNewCommands(input []byte) error {
    var command entities.CommandCreated
    err := proto.Unmarshal(input, &command)
    if err != nil {
        return fmt.Errorf("error unmarshalling command: %v", err)

    }
    if command.SessionId != h.subscribedSession {
        return nil
    }

    response := internalEntities.CommandSentToOperator{
        MessageType: "new-command",
        ID:          command.CommandId,
        Type:        command.Type,
        Status:      command.Status,
        Command:     command.Command,
        Operator:    command.OperatorName,
        CreatedAt:   command.CreatedAt.AsTime().UTC().String(),
    }

    payload, err := json.Marshal(response)
    if err != nil {
        return fmt.Errorf("error marshalling response: %v", err)
    }

    if err := h.conn.WriteMessage(websocket.TextMessage, payload); err != nil {
        return fmt.Errorf("error sending message through websocket: %v", err)
    }
    return nil
}
```
