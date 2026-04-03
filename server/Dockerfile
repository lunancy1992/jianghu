## Build stage
FROM golang:1.22-alpine AS builder

RUN apk add --no-cache gcc musl-dev sqlite-dev

WORKDIR /build

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=1 GOOS=linux go build -o jianghu-server ./cmd/server

## Runtime stage
FROM alpine:3.19

RUN apk add --no-cache sqlite-libs ca-certificates tzdata

ENV TZ=Asia/Shanghai

WORKDIR /app

COPY --from=builder /build/jianghu-server .
COPY --from=builder /build/config.example.yaml ./config.yaml

RUN mkdir -p data

EXPOSE 8080

ENTRYPOINT ["./jianghu-server"]
