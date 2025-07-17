# Caregiver Microservices Application

This is a Go-based backend application for managing caregiver schedules and visits. It provides APIs for:

- User authentication and management
- Schedule creation and management
- Task tracking for caregivers
- Visit status updates (start/end)

Built with Go, PostgreSQL, and Docker, this application follows Clean Architecture principles.

## Architecture Approach

The application uses a three-layer architecture:

1. **Domain Layer** (`src/domain/`): Core business entities and interfaces
2. **Application Layer** (`src/application/`): Business logic in use cases
3. **Infrastructure Layer** (`src/infrastructure/`): External implementations (database, web, etc.)

This separation ensures:

- Business logic is independent of frameworks
- Components are easily testable
- Dependencies point inward

## Key Technologies

- **Gin**: HTTP web framework
- **GORM**: PostgreSQL ORM
- **JWT**: Authentication
- **Zap**: Structured logging

## Running the Application

### Prerequisites

- Go 1.16+
- PostgreSQL 12+
- Docker (optional)

### Local Setup

1. Copy environment file:

   ```bash
   cp .env.example .env
   ```

2. Install dependencies:

   ```bash
   go mod download
   ```

3. Run the application:
   ```bash
   go run main.go
   ```

### Docker Setup

```bash
# Start containers
docker-compose up -d

# View logs
docker-compose logs -f app

# Stop containers
docker-compose down
```

The API will be available at `http://localhost:8085`

## Testing

Run tests with:

```bash
# All tests
go test ./src/...

# With coverage
go test -cover ./src/...

# Specific package
go test ./src/application/usecases/schedule

# Using scripts
./scripts/run-tests.sh  # Linux/Mac
.\scripts\run-tests.ps1 # Windows
```
