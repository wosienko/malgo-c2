FROM oven/bun:1.1.6 as builder
WORKDIR /app

COPY . .

RUN bun install

CMD [ "bun", "run", "migrate" ]