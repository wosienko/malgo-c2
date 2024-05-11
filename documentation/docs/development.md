# Development

## Required Tools

- [Go](https://go.dev/)
- [gRPC for Go](https://grpc.io/docs/languages/go/quickstart/)
- [Node](https://nodejs.org/en)
- [Bun](https://bun.sh/)
- [Task](https://taskfile.dev/)
- [Docker](https://www.docker.com/)

## Task

Task is an alternative to Make. It is a task runner and build tool that aims to be simpler and easier to use.
To see available tasks, run `task`.

## Project layout

```plaintext
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
```

## Running the project

It is recommended that developers launch every service separately. This way, they can see logs and debug the services more effectively.

### Launching PostgreSQL, Redis and Jaeger

```bash
docker run -d --name some-jaeger \
  -e COLLECTOR_ZIPKIN_HTTP_PORT=9411 \
  -p 5775:5775/udp \
  -p 6831:6831/udp \
  -p 6832:6832/udp \
  -p 5778:5778 \
  -p 16686:16686 \
  -p 14268:14268 \
  -p 9411:9411 \
  jaegertracing/all-in-one:1.6

docker run -d --name some-redis \
  -p 6379:6379 \
  -p 8001:8001 \
  redis/redis-stack:latest

docker run -d --name some-postgres \
  -p5433:5432 \
  -e POSTGRES_PASSWORD=postgres \
  -e POSTGRES_USER=postgres \
  -e POSTGRES_DB=database \
  postgres
```

Feel free to tweak ports, names and environment variables.

### Operator Panel

1. Ensure that PostgreSQL is running.
1. Go to the `malgo-operator` directory.
1. Copy `.env.example` to `.env` and fill in the required fields.
1. Run `task dev` to start the operator panel in development mode.

!!! info
    In the future, you can launch `Task` tasks directly from the root directory.
    Tasks related to Operator Panel are prefixed with `operator:`.
    For example, to run the Operator Panel in development mode, you can run `task operator:dev`.

### Services

1. Ensure that PostgreSQL, Redis and Jaeger are running.
1. Go to the `services` directory.
1. For each service, copy `.env.example` to `.env` and fill in the required fields.
1. For each service, run appropriate Task:
    - `task ws-run` for Websocket service
    - `task gw-run` for gRPC service
    - `task redirector-run` for Redirector service

!!! info
    `common` is not a service, you dummy.

!!! info
    In the future, you can launch `Task` tasks directly from the root directory.
    Tasks related to services are prefixed with `services:`.
    For example, to run the Websocket service, you can run `task services:ws-run`.

!!! warning
    If you leave `ENV=production` in the `.env` file, you won't see traces in Jaeger instantly.

## Implants

### Windows

To build implants on Windows, you need MSVC compiler.

### Linux

To build implants on Linux, you need MinGW compiler.

## Code generation

Some of the functionalities require code generation (e.g. DB schemas, gRPC services).

!!! info
    Following examples assume that you are in the root directory of the project.
    Remember that you may omit the prefix depending on the directory you are in.

### Generate DB schemas

```bash
task operator:dbgen
```

!!! info
    To manually migrate the database, run `task operator:migrate`.

### Generate gRPC services and Protobuf structures

```bash
task services:proto
```
