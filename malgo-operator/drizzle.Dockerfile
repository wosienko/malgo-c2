FROM node:22-slim as builder
WORKDIR /app

COPY . .

RUN npm install --force drizzle-kit

CMD [ "npx", "drizzle-kit", "push:pg", "--config", "drizzle.config.ts" ]