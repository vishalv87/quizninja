# QuizNinja API

> A Go-based REST API backend for the QuizNinja quiz application

## What is this?

QuizNinja API is the backend service that powers the QuizNinja quiz application. It handles everything from user authentication to quiz management, scoring, achievements, social features (friends, discussions), and leaderboards.

**Key Features:**
- Quiz creation, management, and attempt tracking
- User authentication (JWT and Supabase support)
- Achievement and gamification system
- Friends and social features
- Real-time leaderboards
- Discussion forums per quiz
- Notification system

## Quick Start

### Prerequisites

- Go 1.23.1 or later
- PostgreSQL database (or Supabase account)
- Docker (optional, for containerized deployment)

### 1. Clone and Setup

```bash
# Clone the repository
git clone <repository-url>
cd quizninja-api

# Copy environment template
cp .env.example .env
```

### 2. Configure Environment

Edit `.env` with your database credentials:

```bash
# For PostgreSQL
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=your_password
DB_NAME=quizninja

# OR for Supabase
USE_SUPABASE=true
SUPABASE_URL=https://your-project.supabase.co
SUPABASE_ANON_KEY=your-anon-key
SUPABASE_SERVICE_KEY=your-service-key
```

### 3. Run the Application

```bash
# Install dependencies
go mod download

# Run the server
go run main.go

# Or with hot reload (requires air)
air
```

The API will be available at `http://localhost:8080`

### 4. Verify Installation

```bash
curl http://localhost:8080/health
# Expected: {"status": "ok"}
```

## Architecture Overview

```
┌─────────────────────────────────────────────────────────────────┐
│                        HTTP Request                              │
└─────────────────────────────────────────────────────────────────┘
                               │
                               ▼
┌─────────────────────────────────────────────────────────────────┐
│                      Middleware Layer                            │
│  ┌─────────┐ ┌──────────┐ ┌────────┐ ┌──────┐ ┌────────────┐   │
│  │ Logging │ │   CORS   │ │  Auth  │ │ Rate │ │  Security  │   │
│  │         │ │          │ │  JWT   │ │Limit │ │  Headers   │   │
│  └─────────┘ └──────────┘ └────────┘ └──────┘ └────────────┘   │
└─────────────────────────────────────────────────────────────────┘
                               │
                               ▼
┌─────────────────────────────────────────────────────────────────┐
│                        Routes Layer                              │
│         /api/v1/* (Public)    │    /internal/v1/* (Private)     │
└─────────────────────────────────────────────────────────────────┘
                               │
                               ▼
┌─────────────────────────────────────────────────────────────────┐
│                       Handlers Layer                             │
│  ┌──────┐ ┌──────┐ ┌─────────┐ ┌─────────┐ ┌───────────────┐   │
│  │ Auth │ │ Quiz │ │ Friends │ │ Leaders │ │ Achievements  │   │
│  └──────┘ └──────┘ └─────────┘ └─────────┘ └───────────────┘   │
└─────────────────────────────────────────────────────────────────┘
                               │
                               ▼
┌─────────────────────────────────────────────────────────────────┐
│                      Repository Layer                            │
│            Data Access Objects (DAOs) for each entity            │
└─────────────────────────────────────────────────────────────────┘
                               │
                               ▼
┌─────────────────────────────────────────────────────────────────┐
│                        Database                                  │
│              PostgreSQL / Supabase PostgreSQL                    │
└─────────────────────────────────────────────────────────────────┘
```

## Project Structure

```
quizninja-api/
├── main.go                 # Application entry point
├── config/                 # Configuration management
│   └── config.go          # Environment loading & validation
├── database/              # Database connection & migrations
│   ├── database.go        # Connection management
│   ├── migrate.go         # Migration runner
│   ├── schema.sql         # Generated schema
│   └── migrations/        # SQL migration files (67 files)
├── handlers/              # HTTP request handlers (public API)
│   ├── auth_handler.go    # Authentication endpoints
│   ├── quiz_handler.go    # Quiz operations
│   ├── friends_handler.go # Friend management
│   └── ...                # Other handlers
├── internal/              # Internal service-to-service API
│   ├── handlers/          # Internal handlers
│   ├── middleware/        # Internal auth middleware
│   └── routes/            # Internal route setup
├── middleware/            # HTTP middleware
│   ├── auth.go            # JWT authentication
│   ├── rate_limiter.go    # Rate limiting
│   └── ...                # Other middleware
├── models/                # Data models & structs
│   ├── user.go            # User models
│   ├── quiz.go            # Quiz models
│   └── ...                # Other models
├── repository/            # Data access layer
│   ├── interfaces.go      # Repository interfaces
│   ├── user_repository.go # User data access
│   └── ...                # Other repositories
├── routes/                # Route definitions
│   └── routes.go          # All route setup
├── services/              # Business logic
│   └── achievement_service.go
├── utils/                 # Utility functions
│   ├── logger.go          # Logging utilities
│   ├── auth_interfaces.go # Auth abstractions
│   └── ...                # Other utilities
├── cmd/                   # CLI commands
│   └── generate-schema/   # Schema generator
├── Dockerfile             # Docker build config
├── docker-compose.yml     # Docker Compose setup
└── .env.example           # Environment template
```

## API Overview

### Public API (`/api/v1/`)

| Category | Endpoints | Auth Required |
|----------|-----------|---------------|
| Health | `GET /health` | No |
| Auth | `POST /auth/register`, `POST /auth/login`, `POST /auth/logout` | Partial |
| Quizzes | `GET /quizzes`, `GET /quizzes/:id`, `GET /quizzes/featured` | No |
| Quiz Attempts | `POST /quizzes/:id/attempts`, `POST /quizzes/:id/attempts/:id/submit` | Yes |
| Users | `GET /users/:id`, `GET /auth/profile` | Yes |
| Friends | `POST /friends/requests`, `GET /friends` | Yes |
| Leaderboard | `GET /leaderboard`, `GET /leaderboard/rank` | Yes |
| Achievements | `GET /achievements`, `GET /achievements/progress` | Yes |
| Notifications | `GET /notifications`, `PUT /notifications/:id/read` | Yes |
| Discussions | `GET /discussions`, `POST /discussions` | Yes |
| Favorites | `POST /favorites`, `GET /favorites` | Yes |
| Ratings | `POST /quizzes/:id/ratings`, `GET /quizzes/:id/ratings` | Yes |

### Internal API (`/internal/v1/`)

Used for service-to-service communication. Requires `X-Internal-API-Key` header.

| Endpoint | Purpose |
|----------|---------|
| `POST /attempts/:id/validate` | Validate quiz attempt |
| `POST /scoring/calculate` | Calculate quiz score |
| `POST /users/:id/statistics` | Update user statistics |
| `POST /users/:id/achievements/check` | Check and unlock achievements |

## Technology Stack

| Component | Technology |
|-----------|------------|
| Language | Go 1.23.1 |
| Web Framework | Gin |
| Database | PostgreSQL / Supabase |
| Authentication | JWT (golang-jwt) |
| Rate Limiting | ulule/limiter |
| Logging | Logrus |
| Containerization | Docker |
| Deployment | Google Cloud Run |

## Deployment

### Docker

```bash
# Build image
docker build -t quizninja-api .

# Run container
docker run -p 8080:8080 --env-file .env quizninja-api
```

### Docker Compose

```bash
docker-compose up -d
```

### Google Cloud Run

The application is optimized for Cloud Run with:
- Multi-stage Docker build
- Nginx reverse proxy
- Supervisor process management
- Health check endpoint

## Folder Documentation

Each folder has its own detailed documentation:

| Folder | Documentation | Description |
|--------|--------------|-------------|
| `/config` | [config/README.md](./config/README.md) | Configuration management |
| `/database` | [database/README.md](./database/README.md) | Database & migrations |
| `/handlers` | [handlers/README.md](./handlers/README.md) | Public API handlers |
| `/internal` | [internal/README.md](./internal/README.md) | Internal API |
| `/middleware` | [middleware/README.md](./middleware/README.md) | HTTP middleware |
| `/models` | [models/README.md](./models/README.md) | Data models |
| `/repository` | [repository/README.md](./repository/README.md) | Data access layer |
| `/routes` | [routes/README.md](./routes/README.md) | API routing |
| `/services` | [services/README.md](./services/README.md) | Business logic |
| `/utils` | [utils/README.md](./utils/README.md) | Utilities |
| `/cmd` | [cmd/README.md](./cmd/README.md) | CLI commands |

## Common Tasks

### Running in Development

```bash
# With hot reload
air

# Without hot reload
go run main.go
```

### Running Migrations

Migrations run automatically on startup. To generate a new migration:

1. Create a new SQL file in `database/migrations/` with format `XXX_description.sql`
2. Restart the application

### Running Tests

```bash
go test ./...
```

### Building for Production

```bash
CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o main .
```

## Environment Variables

See [config/README.md](./config/README.md) for a complete list of environment variables.

Key variables:
- `PORT` - Server port (default: 8080)
- `DB_HOST`, `DB_PORT`, `DB_USER`, `DB_PASSWORD`, `DB_NAME` - PostgreSQL connection
- `USE_SUPABASE` - Enable Supabase mode
- `RATE_LIMIT_ENABLED` - Enable rate limiting
- `LOG_LEVEL` - Logging verbosity (DEBUG, INFO, WARN, ERROR)

## Contributing

1. Read the folder-specific documentation before making changes
2. Follow existing code patterns and conventions
3. Write tests for new functionality
4. Update documentation when adding features

## License

[Add your license here]
