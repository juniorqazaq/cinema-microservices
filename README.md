# Cinema microservices

Monorepository for the cinema booking system services.

## Services

| Service        | Path            | Description                          |
| -------------- | --------------- | ------------------------------------ |
| **user-service** | [`user-service/`](user-service/) | gRPC users, auth, JWT, profiles (port `50051`) |

### user-service

```bash
cd user-service
cp .env.example .env
make build
make test
# or: docker compose up -d --build
```

Details: see [user-service/README.md](user-service/README.md).
