FROM node:22-slim as builder
WORKDIR /app

COPY ./migrations .
COPY ./drizzle.config.ts .

RUN npm install drizzle-kit

CMD [ "npx", "drizzle-kit", "push:pg", "--config", "drizzle.config.ts" ]