
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

#FROM alpine:3.22 # if you want more debug tools
FROM gcr.io/distroless:nonroot

COPY --from=builder /tmp/app/dist/build /etc/abramad/irbankmock/server

WORKDIR /etc/abramad/irbankmock
CMD ["./server"]
