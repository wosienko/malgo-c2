FROM golang:1.22.2-alpine as builder
ARG WS_PORT=8080

WORKDIR /app

COPY . .

RUN go mod download

RUN go build -o /app/executable ./malgo-websocket/cmd/server/main.go

FROM golang:1.22.2-alpine as runner

WORKDIR /app

COPY --from=builder /app/executable /app/executable

EXPOSE ${WS_PORT}

CMD ["/app/executable"]