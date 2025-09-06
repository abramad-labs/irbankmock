
ARG IMAGE_TAG

FROM golang:1.22.2 AS builder
WORKDIR /tmp/app
ARG IMAGE_TAG

COPY ./go.mod ./go.sum ./

RUN --mount=type=cache,target=/go/pkg/mod/ go mod download -x

COPY . .

RUN --mount=type=cache,target=/go/pkg/mod/ \
    --mount=type=cache,target=/root/.cache \
    CGO_ENABLED=0 GOOS=linux go build -C . -ldflags "-X 'github.com/abramad-labs/irbankmock/internal/version.ServerVersion=${IMAGE_TAG:-development}'" -o dist/build ./cmd/server/main.go

FROM oven/bun:1.2.4 AS frontbuilder

ENV BUN_INSTALL_CACHE_DIR=/usr/.bun/install/cache

WORKDIR /app

COPY web/app/bun.lock web/app/package.json ./


RUN --mount=type=cache,target=/app/.next/cache \
    --mount=type=cache,target=/root/.cache \
    --mount=type=cache,target=${BUN_INSTALL_CACHE_DIR} \
    bun install --frozen-lockfile

COPY web/app .

RUN --mount=type=cache,target=/app/.next/cache \
    --mount=type=cache,target=${BUN_INSTALL_CACHE_DIR} \
    bun run build

#FROM alpine:3.22 # if you want more debug tools
FROM gcr.io/distroless/static-debian12:nonroot

WORKDIR /opt/irbankmock/data
WORKDIR /opt/irbankmock/web/app

COPY --from=frontbuilder /app/out ./

WORKDIR /etc/abramad/irbankmock

COPY --from=builder /tmp/app/dist/build ./server


ENV IRBANKMOCK_DATA_PATH=/opt/irbankmock/data
ENV IRBANKMOCK_WEBAPP_PATH=/opt/irbankmock/web/app

LABEL org.opencontainers.image.source=https://github.com/abramad-labs/irbankmock
LABEL org.opencontainers.image.description="service for testing the internet payment gateways of Iranian banks"
LABEL org.opencontainers.image.licenses=MIT

CMD ["./server"]
