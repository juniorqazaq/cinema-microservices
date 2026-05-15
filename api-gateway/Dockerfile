FROM golang:1.26-alpine AS builder
WORKDIR /src
RUN apk add --no-cache git
COPY . .
RUN cd api-gateway && CGO_ENABLED=0 GOOS=linux \
	go build -ldflags="-w -s" -o /bin/api-gateway ./cmd/api-gateway

FROM gcr.io/distroless/static-debian12 AS runtime
COPY --from=builder /bin/api-gateway /api-gateway
EXPOSE 8080
ENTRYPOINT ["/api-gateway"]
