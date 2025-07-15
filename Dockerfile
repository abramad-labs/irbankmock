
ARG IMAGE_TAG

FROM golang:1.22.2 AS builder
WORKDIR /tmp/app
ARG IMAGE_TAG

COPY ./go.mod ./go.sum ./

RUN --mount=type=cache,target=/go/pkg/mod/ go mod download -x

COPY . .

RUN --mount=type=cache,target=/go/pkg/mod/ \
    --mount=type=cache,target=/root/.cache \
    CGO_ENABLED=0 GOOS=linux go build -C . -ldflags "-X 'github.com/abramad-labs/irbankmock/internal/version.ServerVersion=$IMAGE_TAG'" -o dist/build ./cmd/server/main.go

FROM oven/bun:1.2 AS frontbuilder

WORKDIR /app

COPY web/app/bun.lockb web/app/package.json ./

RUN --mount=type=cache,target=/root/.bun \
    --mount=type=cache,target=/app/.next/cache \
    --mount=type=cache,target=/root/.cache \
    bun install --frozen-lockfile

COPY web/app .

# Build the Next.js project
RUN --mount=type=cache,target=/app/.next/cache \
    bun --bun run build

#FROM alpine:3.22 # if you want more debug tools
FROM gcr.io/distroless:nonroot

WORKDIR /opt/irbankmock/web/app

COPY --from=frontbuilder /app/out ./

WORKDIR /etc/abramad/irbankmock

COPY --from=builder /tmp/app/dist/build ./server

ENV IRBANKMOCK_WEBAPP_PATH=/opt/irbankmock/web/app

CMD ["./server"]
