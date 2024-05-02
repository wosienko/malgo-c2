ARG GRPC_PORT=8080

FROM golang:1.22.2 as builder

WORKDIR /app

COPY . .

RUN go mod download

RUN go build -o /app/malgo-gateway ./cmd/server/main.go

FROM golang:1.22.2-nanoserver as runner

WORKDIR /app

COPY --from=builder /app/malgo-gateway /app

RUN chown -R go:go /app

USER go

EXPOSE ${GRPC_PORT}
ENV GRPC_ADDR="0.0.0.0:${GRPC_PORT}"

CMD ["/app/malgo-websocket"]