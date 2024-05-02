ARG WS_PORT=8080

FROM golang:1.22.2 as builder

WORKDIR /app

COPY . .

RUN go mod download

RUN go build -o /app/malgo-websocket ./cmd/server/main.go

FROM golang:1.22.2-nanoserver as runner

WORKDIR /app

COPY --from=builder /app/malgo-websocket /app

RUN chown -R go:go /app

USER go

EXPOSE ${WS_PORT}
ENV WS_ADDR="0.0.0.0:${WS_PORT}"

CMD ["/app/malgo-websocket"]