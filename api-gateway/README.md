# API Gateway

Gin HTTP gateway for the React frontend. It owns REST payload shape, JWT validation through `UserService.ValidateToken`, and frontend-facing error/status mapping while existing services keep the gRPC business contracts.

## Run

```bash
go run ./cmd/api-gateway
```

Environment:

| Variable | Default |
| --- | --- |
| `HTTP_ADDR` | `:8080` |
| `USER_GRPC_ADDR` | `localhost:50051` |
| `MOVIE_GRPC_ADDR` | `localhost:50052` |
| `BOOKING_GRPC_ADDR` | `localhost:50053` |
| `REQUEST_TIMEOUT` | `5s` |
| `CORS_ALLOWED_ORIGINS` | `http://localhost:5173` |
| `IP_RATE_PER_SEC` / `IP_RATE_BURST` | `20` / `40` |
| `USER_RATE_PER_SEC` / `USER_RATE_BURST` | `10` / `20` |

Rate limiting is in-memory per process. Use a shared store before running multiple gateway replicas behind one public endpoint.
