# Bookmark Service

A RESTful API service for managing bookmarks and URL shortening, built with Go using Domain-Driven Design (DDD) principles.

## Features

- **Bookmark Management** — Create, read, update, and delete bookmarks per user
- **URL Shortening** — Generate shortened URLs with expiration
- **Caching** — Redis-backed caching for improved performance
- **Rate Limiting** — Per-IP (public endpoints) and per-user (authenticated endpoints)
- **JWT Authentication** — RSA-based JWT validation for protected routes
- **Pagination & Sorting** — Flexible bookmark listing with query params
- **API Documentation** — Auto-generated Swagger/OpenAPI docs
- **Database Migrations** — Version-controlled schema management

## Architecture

The project follows **Domain-Driven Design (DDD)**:

```
bookmark-service/
├── cmd/
│   ├── api/            # API server entry point
│   └── migrate/        # Database migration entry point
├── internal/
│   ├── api/            # Route registration, middleware
│   ├── infrastructure/ # DB, Redis, JWT initialization
│   └── app/
│       ├── handler/    # HTTP request/response layer
│       ├── service/    # Business logic layer
│       ├── repository/ # Data access layer
│       └── model/      # Domain models
├── migrations/         # SQL migration files
└── docs/               # Swagger documentation
```

## Prerequisites

- Go 1.21+
- PostgreSQL
- Redis
- Docker & Docker Compose (optional)

## Getting Started

### Run with Docker Compose

```bash
make dev-up     # Start PostgreSQL and Redis
make migrate    # Run database migrations
make dev-run    # Start the API server
```

### Run locally

```bash
# Start dependencies
make dev-up

# Run migrations
make migrate

# Start server
go run cmd/api/main.go
```

The server starts on `:8080` by default.

## Configuration

Configure via environment variables (or `.env` file):

| Variable | Default | Description |
|---|---|---|
| `APP_PORT` | `:8080` | Server port |
| `APP_HOST_NAME` | `localhost:8080` | Host for Swagger docs |
| `SERVICE_NAME` | `bookmark_service` | Service identifier |
| `DB_NAME` | `bookmark` | PostgreSQL database name |

## API Endpoints

### Public

| Method | Path | Description |
|---|---|---|
| `GET` | `/health-check` | Health check |
| `POST` | `/v1/links/shorten` | Shorten a URL |
| `GET` | `/v1/links/redirect/:code` | Redirect to original URL |

### Protected (JWT required)

| Method | Path | Description |
|---|---|---|
| `POST` | `/v1/bookmarks` | Create a bookmark |
| `GET` | `/v1/bookmarks` | List bookmarks (with pagination) |
| `PUT` | `/v1/bookmarks/:id` | Update a bookmark |
| `DELETE` | `/v1/bookmarks/:id` | Delete a bookmark |

### API Documentation

Swagger UI is available at: `http://localhost:8080/swagger/index.html`

To regenerate docs:
```bash
make swag-gen
```

## Testing

```bash
make test           # Run tests with coverage (80% threshold required)
make docker-test    # Run tests inside Docker
```

## Database Migrations

```bash
make migrate                        # Run pending migrations
make new-schema name=<schema_name>  # Create a new migration file
```

## Docker

```bash
make docker-build   # Build Docker image
make docker-up      # Start containers
make docker-down    # Stop containers
make docker-release # Push image to Docker Hub
```

## Makefile Targets

```bash
make dev-up         # Start dev dependencies (docker-compose)
make dev-down       # Stop dev dependencies
make dev-run        # Run API server locally
make test           # Run tests with coverage
make migrate        # Run database migrations
make swag-gen       # Generate Swagger docs
make mock-gen       # Generate mocks for testing
make redis-cli      # Access Redis CLI
make redis-monitor  # Monitor Redis commands
make clean          # Clean test cache and coverage
```

## CI/CD

GitHub Actions workflow (`.github/workflows/ci-pure.yaml`) triggers on pull requests and pushes to `main`:

1. Run tests with coverage
2. Build Docker image
3. Push to Docker Hub

Image tags:
- `temporary` — feature branches
- `dev` — `main` branch
- `v*.*.*` — version tags