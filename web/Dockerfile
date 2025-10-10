FROM node:24-alpine AS builder
RUN corepack enable && corepack prepare pnpm@latest --activate

WORKDIR /app
COPY package.json pnpm-lock.yaml ./
RUN pnpm install --frozen-lockfile
COPY . .
RUN pnpm build

FROM node:24-alpine AS runner
RUN corepack enable && corepack prepare pnpm@latest --activate  # penting!
WORKDIR /app
COPY --from=builder /app ./
EXPOSE 5173
CMD ["pnpm", "start"]
