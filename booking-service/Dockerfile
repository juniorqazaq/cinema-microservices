FROM golang:latest AS builder
WORKDIR /app
COPY . .
RUN find . -name "*_test.go" -delete && \
    go mod edit -droprequire=github.com/testcontainers/testcontainers-go && \
    go mod edit -droprequire=github.com/testcontainers/testcontainers-go/modules/postgres && \
    go mod edit -droprequire=github.com/shirou/gopsutil/v3 && \
    go mod edit -droprequire=github.com/lufia/plan9stats && \
    rm -f go.sum && \
    go mod tidy && \
    CGO_ENABLED=0 GOOS=linux go build -o booking-service ./cmd/main.go

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /app
COPY --from=builder /app/booking-service .
EXPOSE 50053
CMD ["./booking-service"]
