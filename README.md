# MALGO C2 - Operator & Services

Short overview

- MALGO is a C2 (command & control) project with a SvelteKit operator UI (malgo-operator) and several Go services (gateway, websocket, redirector).
- Main parts:
  - Operator UI: [malgo-operator](malgo-operator/) — SvelteKit frontend + API.
  - Gateway service: gRPC gateway in Go — [services/malgo-gateway/cmd/server/main.go](services/malgo-gateway/cmd/server/main.go).
  - Websocket forwarder: real-time notifications — [services/malgo-websocket/internal/ws/writeToWebsocket.go](services/malgo-websocket/internal/ws/writeToWebsocket.go).
  - DNS redirector (implant comms): [services/malgo-redirector/internal/dnsproxy/handler.go](services/malgo-redirector/internal/dnsproxy/handler.go).

Quick links

- Repo compose: [docker-compose.yml](docker-compose.yml)
- Operator Dockerfile: [malgo-operator/Dockerfile](malgo-operator/Dockerfile)
- Operator env example: [malgo-operator/.env.example](malgo-operator/.env.example)
- Top-level env example: [.env.example](.env.example)
- DB config for migrations: [malgo-operator/drizzle.config.ts](malgo-operator/drizzle.config.ts)
- DB bootstrap & schema wiring: [malgo-operator/src/lib/db/db.server.ts](malgo-operator/src/lib/db/db.server.ts)
- Key DB schema examples: [malgo-operator/src/lib/db/schema/c2_sessions.ts](malgo-operator/src/lib/db/schema/c2_sessions.ts), [malgo-operator/src/lib/db/schema/c2_commands.ts](malgo-operator/src/lib/db/schema/c2_commands.ts)

Architecture summary

- Frontend (malgo-operator) uses SvelteKit and Bun; router + API endpoints are colocated under [malgo-operator/src/routes](malgo-operator/src/routes).
  - Example API endpoints: sessions listing at [`getAllSessionsForProject`](malgo-operator/src/lib/services/c2-sessions-service.ts) and session commands at [`getCommandsForSession`](malgo-operator/src/lib/services/c2-commands-service.ts).
  - App-level layout and auth use: [malgo-operator/src/routes/+layout.svelte](malgo-operator/src/routes/+layout.svelte) and Lucia auth in [malgo-operator/src/lib/auth.server.ts](malgo-operator/src/lib/auth.server.ts).
- Backend services:
  - Gateway (gRPC) exposes implant APIs; server entry: [services/malgo-gateway/cmd/server/main.go](services/malgo-gateway/cmd/server/main.go).
  - Websocket service forwards events to operator UI; relevant code: [services/malgo-websocket/internal/ws/writeToWebsocket.go](services/malgo-websocket/internal/ws/writeToWebsocket.go).
  - Redirector handles DNS/TXT-based implant comms (DNS chunking); handler: [services/malgo-redirector/internal/dnsproxy/handler.go](services/malgo-redirector/internal/dnsproxy/handler.go).

Development - local (operator)

1. Prepare environment
   - Copy env examples:
     - Root: cp .env.example .env
     - Operator: cp malgo-operator/.env.example malgo-operator/.env
2. Start dev server (operator)
   - Using Bun (project configured for Bun):
     - cd malgo-operator
     - bun install
     - bun run dev
   - The project expects the DB via DATABASE_URL; see [malgo-operator/drizzle.config.ts](malgo-operator/drizzle.config.ts).
3. Useful files
   - package scripts & metadata: [malgo-operator/package.json](malgo-operator/package.json)
   - Playwright tests: [malgo-operator/playwright.config.ts](malgo-operator/playwright.config.ts), example test: [malgo-operator/tests/test.ts](malgo-operator/tests/test.ts)

Development - services (Go)

- Each service contains a Dockerfile; examples:
  - Gateway build: [services/gateway.Dockerfile](services/gateway.Dockerfile)
  - Gateway server entry: [services/malgo-gateway/cmd/server/main.go](services/malgo-gateway/cmd/server/main.go)
- Local dev:
  - Install Go 1.22.x (as used in Dockerfiles).
  - Build & run service with `go build` / `go run` or use Dockerfiles in CI/dev.

Run everything (Docker)

- The repo has a compose file: [docker-compose.yml](docker-compose.yml).
- Typical flow:
  - Fill environment variables in `.env` files.
  - docker compose up --build
- Individual Dockerfiles for components are located at service folders (see links above).

Database & migrations

- Drizzle is used for schema/migrations:
  - Schema glob: [malgo-operator/drizzle.config.ts](malgo-operator/drizzle.config.ts)
  - DB access: [malgo-operator/src/lib/db/db.server.ts](malgo-operator/src/lib/db/db.server.ts)
  - Schema examples: [malgo-operator/src/lib/db/schema/c2_sessions.ts](malgo-operator/src/lib/db/schema/c2_sessions.ts), [malgo-operator/src/lib/db/schema/c2_commands.ts](malgo-operator/src/lib/db/schema/c2_commands.ts)
- Running migrations:
  - Use drizzle-kit command configured in [malgo-operator/drizzle.Dockerfile](malgo-operator/drizzle.Dockerfile) or run locally with npx/npm.

API & important code pointers

- Session list for a project: [`getAllSessionsForProject`](malgo-operator/src/lib/services/c2-sessions-service.ts)
- Commands pagination & count helpers: [`getCommandsForSession`](malgo-operator/src/lib/services/c2-commands-service.ts), [`getCountOfCommandsForSession`](malgo-operator/src/lib/services/c2-commands-service.ts)
- Project service helpers: [`getProjects`](malgo-operator/src/lib/services/project-service.ts), [`getProjectById`](malgo-operator/src/lib/services/project-service.ts), [`getCountOfProjects`](malgo-operator/src/lib/services/project-service.ts)
- Websocket store (client side) referenced in many components: [malgo-operator/src/lib/stores/Websocket.js](malgo-operator/src/lib/stores/Websocket.js)
- Key UI routes:
  - Operator UI root/layout: [malgo-operator/src/routes/+layout.svelte](malgo-operator/src/routes/+layout.svelte)
  - Project pages: [malgo-operator/src/routes/(authorized)/projects](<malgo-operator/src/routes/(authorized)/projects>)
  - Admin pages: [malgo-operator/src/routes/(authorized)/admin](<malgo-operator/src/routes/(authorized)/admin>)
