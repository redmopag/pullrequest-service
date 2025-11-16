FROM golang:1.25-alpine AS builder
WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o pr-service ./cmd/pr-service

FROM alpine:3.19

WORKDIR /root/
COPY --from=builder /app/pr-service .

CMD ["./pr-service"]
