FROM node:24 AS builder

RUN corepack enable && corepack prepare pnpm@latest --activate

WORKDIR /app

COPY package.json pnpm-lock.yaml ./
RUN pnpm install --frozen-lockfile

COPY . .
RUN pnpm run build

FROM node:24

WORKDIR /app

COPY --from=builder /app/dist ./dist

RUN npm install -g serve@latest

EXPOSE 5173

CMD ["serve", "-s", "dist", "-l", "5173", "--no-port-switching"]