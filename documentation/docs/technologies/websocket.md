# Websocket Server

## Event-Driven Development

While the notion "Event-Driven Development" is not a technology per sé, it had to be mentioned as early as possible.
It significantly influences all the choices made with regards to the Websocket Server and the Gateway Server.

Websocket Server is fully asynchronous and event-driven. Operators emit commands and events through Operator Panel, which are routed and handled by the Websocket Server. Gateway Server may also emit events and commands.

## Go

Websocket Server is written in Go. Go is a statically typed, compiled language that is designed for simplicity and efficiency. It is a great choice for writing high-performance servers. It is a perfect balance between performance and ease of development.

### Gorilla Websocket

Gorilla Websocket is a popular Websocket library for Go. It is used in the Websocket Server to handle Websocket connections. Gorilla libraries, albeit slightly dated, are still among most popular HTTP-oriented libraries for Go.

### sqlx

sqlx is a library that provides a set of extensions on top of the standard `database/sql` library. It is used in the Websocket Server to interact with the database.

### Watermill

Watermill is a Go library for working efficiently with message streams. It is used in the Websocket Server to handle messages and events. Watermill is a great choice for working with message streams in Go.

Watermill is the core library that enables the Websocket Server to be fully asynchronous and event-driven. Advanced functionalities abstract away the complexity of working with message streams.

### Protobuf

Protobuf is a language-agnostic serialization format developed by Google. Websocket Server uses Protobuf throughout events and commands. Protobuf is known for its' efficiency and simplicity.
