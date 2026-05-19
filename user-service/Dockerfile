FROM golang:1.23-alpine AS builder
WORKDIR /app
RUN apk add --no-cache git
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
	go build -ldflags="-w -s" -o /bin/user-service ./cmd/main.go

FROM gcr.io/distroless/static-debian12 AS runtime
COPY --from=builder /bin/user-service /user-service
EXPOSE 50051
ENTRYPOINT ["/user-service"]
