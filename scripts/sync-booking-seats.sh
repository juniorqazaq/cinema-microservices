#!/usr/bin/env bash
# Copy seat IDs from movie_db into booking DB (for legacy booking-service or compat seats table).
set -euo pipefail
CONTAINER="${POSTGRES_CONTAINER:-cinema-microservices-postgres-1}"

docker exec "$CONTAINER" psql -U postgres -d booking -v ON_ERROR_STOP=1 <<'SQL'
CREATE TABLE IF NOT EXISTS seats (
    id UUID PRIMARY KEY,
    is_available BOOLEAN NOT NULL DEFAULT true
);
SQL

ids=$(docker exec "$CONTAINER" psql -U postgres -d movie_db -t -A -c "SELECT id FROM seats")
count=0
for id in $ids; do
  docker exec "$CONTAINER" psql -U postgres -d booking -q -c \
    "INSERT INTO seats (id, is_available) VALUES ('$id', true) ON CONFLICT (id) DO NOTHING;"
  count=$((count + 1))
done
echo "Synced $count seat IDs into booking.seats"
