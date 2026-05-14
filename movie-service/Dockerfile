FROM golang:1.22-alpine AS builder
WORKDIR /app
RUN apk add --no-cache git
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
	go build -ldflags="-w -s" -o /bin/movie-service ./cmd/main.go

FROM gcr.io/distroless/static-debian12 AS runtime
COPY --from=builder /bin/movie-service /movie-service
EXPOSE 50052
ENTRYPOINT ["/movie-service"]
