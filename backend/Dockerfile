# syntax=docker/dockerfile:1
FROM golang:1.24-alpine AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o /vilivala-server ./cmd/server

FROM alpine:3.20
RUN apk add --no-cache ca-certificates tzdata wget
ENV TZ=Asia/Shanghai

WORKDIR /app
COPY --from=builder /vilivala-server /app/vilivala-server
RUN mkdir -p /app/uploads

EXPOSE 8080
HEALTHCHECK --interval=15s --timeout=3s --start-period=10s --retries=3 \
  CMD wget -qO- http://127.0.0.1:8080/api/health || exit 1

CMD ["/app/vilivala-server"]
