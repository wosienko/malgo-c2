# Malgo Command & Control

## Deployment guide

### Tools required for deployment

- [Docker](https://www.docker.com/)
- [Go](https://go.dev/) (depends whether you want to build from sources)

### Main Server

Deployment of the main server is as simple as populating the `.env` in the root folder (see `.env.example` to see the exemplary default values).

Once everything is populated, run: `docker compose up -d`. It will spin up the following services:

- Red Team Operators' panel - Fullstack web application.
- Websocket server - server handling realtime communication with the operator.
- Gateway server - gRPC server designed for communication with redirectors.
- Schema migration container - simple container designed only to perform database migrations.
- PostgreSQL - database for the main server.
- Redis - Redis Stack instance used for Event-Driven Development with streams.
- Jaeger - tracing platform. Enables tracing of each event for easier debugging.
- Grafana - dashboard for monitoring metrics of the services.
- VictoriaMetrics - alternative for Prometheus.
- VictoriaMetrics agent - service scraping the services for new metrics.

### Redirectors

Once the main server is up and running, redirectors need to be set up.
This approach to C2 architecture renders setting up multiple redirectors on different servers
easier than ever.

There are two options for preparing redirectors.

#### Option 1

1. Download binary for your platform from the GitHub release page.
1. create a .env file and populate the following environment variables:
    - `GRPC_ADDR` (e.g. `GRPC_ADDR=localhost:8082`)
    - `DNS_ADDR` (e.g. `DNS_ADDR=127.0.0.1:53`)
1. run the binary.

#### Option 2

1. Navigate to `services/malgo-redirector`.
1. Populate `.env` (see `.env.example` in that directory to see the exemplary default values).
1. Build from sources and run the binary (e.g. `go run ./cmd/server/main.go`).

## Development Guide

### Tools required for development

- [Go](https://go.dev/)
- [gRPC for Go](https://grpc.io/docs/languages/go/quickstart/)
- [Bun](https://bun.sh/)
- [Task](https://taskfile.dev/)
- [Docker](https://www.docker.com/)

### Documentation

Documentation may be found under `/documentation` and is developed using Material for MkDocs.

To launch it, run:

```bash
task docs
```

It will launch a Docker container and bind it [here](http://localhost:8888/).

### Working on Documentation

To launch an interactive, hot-reload version, run the following command:

```bash
task docs-dev
```

This will also launch a Docker container and bind it [here](http://localhost:8889).
