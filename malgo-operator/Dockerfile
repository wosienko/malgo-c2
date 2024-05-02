FROM oven/bun:1.1.6 as builder
WORKDIR /app

COPY . .

RUN bun install
RUN bun run build


FROM oven/bun:1.1.6-slim as runner

WORKDIR /app

COPY --from=builder /app/build /app

RUN bun install

RUN chown -R bun:bun /app

USER bun

EXPOSE 3000

CMD ["bun", "run", "start"]
