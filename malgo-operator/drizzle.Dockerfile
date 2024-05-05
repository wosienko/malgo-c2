FROM node:22-slim as builder
WORKDIR /app

COPY . .
RUN rm ./package.json

RUN npm install drizzle-kit drizzle-orm dotenv

CMD [ "npx", "drizzle-kit", "push:pg", "--config", "drizzle.config.ts" ]