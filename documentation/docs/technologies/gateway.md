# Gateway Server

## Event-Driven Development

While the notion "Event-Driven Development" is not a technology per sé, it had to be mentioned as early as possible.
It significantly influences all the choices made with regards to the Websocket Server and the Gateway Server.

Gateway Server is mostly synchronous. However, it heavily interacts with the event-driven Websocket Server.

## Go

Websocket Server is written in Go. Go is a statically typed, compiled language that is designed for simplicity and efficiency. It is a great choice for writing high-performance servers. It is a perfect balance between performance and ease of development.

### gRPC

gRPC is a high-performance, open-source universal RPC framework. Gateway Server is a gRPC server that exposes API to Redirectors.
gRPC uses Protocol Buffers as the interface definition language. The major upside over regular HTTP APIs is that gRPC is faster, more efficient, and more secure.

### sqlx

sqlx is a library that provides a set of extensions on top of the standard `database/sql` library. It is used in the Gateway Server to interact with the database.

### Watermill

Watermill is used in the Gateway Server to abstract away sending messages to the message broker.

### Protobuf

Protobuf is a language-agnostic serialization format developed by Google. Gateway Server uses Protobuf for gRPC and sending messages to the message broker.
