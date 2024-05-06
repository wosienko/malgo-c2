# Malgo - Developer Documentation

This is a developer documentation - a document not only describing how Malgo Command & Control works, but also how to contribute to it.

## What is Malgo

Malgo is a Command & Control framework designed by Red Team Operators for Red Team Operators. The main goal is to render starting new projects easier than ever. There is no need for hassle with setting up new servers from the ground up each time.

## Why is Malgo

Thus far, every single C2 framework is widely fingerprinted and needs severe obfucation and reconfiguration to make it somewhat usable in the field.
Moreover, every solution requires that new server be started for each new engagement, rendering preparations more tedious.
As if that was not enough, advanced solutions are very pricey. There is no need for these expenses.

## Development

### Required Tools

- [Go](https://go.dev/)
- [gRPC for Go](https://grpc.io/docs/languages/go/quickstart/)
- [Node](https://nodejs.org/en)
- [Bun](https://bun.sh/)
- [Task](https://taskfile.dev/)
- [Docker](https://www.docker.com/)

### Project layout

    .env.example                 # Example of the .env file for Docker Compose startup
    docker-compose.yml
    Taskfile.yml                 # Main config for the Task command
    documentation/               # This developer documentation
    malgo-operator/              # Red Team Operators' panel
    docker/                      # Configuration files for 3rd party docker containers
    services/                    # Folder with all services written in Go
             malgo-websocket/    # Websocket service
             malgo-gateway/      # gRPC service
             malgo-redirector/   # Redirectors
             common/             # Common components for all services (libraries, etc.)
