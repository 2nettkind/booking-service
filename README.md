# Cinema Booking System

A seat-booking service with timed holds and a strict "one ticket, one buyer"
guarantee under concurrent load.

** URL Link: ** website: https://booking-w1i7.onrender.com/

## Stack

Go 1.25.1 · Gin · Redis (Upstash) · Docker · Render

## Architecture

\`\`\`
cmd/main.go                  entry point
internal/
├── adapters/redis/          go-redis wrapper, supports rediss:// URLs
└── booking/
    ├── domain.go            Booking model, BookingStore interface, errors
    ├── service.go           business logic (thin)
    ├── handler.go           Gin HTTP handlers
    ├── redis_store.go       production store
    ├── memory_store.go      test store (no deps)
    └── service_test.go      concurrency integration test
static/index.html            single-page frontend
\`\`\`

## API

| Method | Path                                  | Description           |
|--------|---------------------------------------|-----------------------|
| GET    | `/movies`                             | list movies           |
| GET    | `/movies/:movieID/seats`              | booked seats          |
| POST   | `/movies/:movieID/seats/:seatID/hold` | hold a seat (2 min)   |
| PUT    | `/sessions/:sessionID/confirm`        | confirm booking       |
| DELETE | `/sessions/:sessionID`                | release booking       |

## Run locally

\`\`\`bash
# option 1: local Redis
docker run -d -p 6379:6379 redis:7-alpine
go run ./cmd

# option 2: Upstash
export REDIS_URL="rediss://default:password@host.upstash.io:6379"
go run ./cmd
\`\`\`

Open http://localhost:8080.
