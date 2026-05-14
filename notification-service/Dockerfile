FROM golang:latest AS builder
WORKDIR /app
COPY . .
RUN rm -f go.sum && \
    go mod tidy && \
    CGO_ENABLED=0 GOOS=linux go build -o notification-service ./cmd/main.go

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /app
COPY --from=builder /app/notification-service .
CMD ["./notification-service"]
