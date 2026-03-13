# QuizNinja

A gamified quiz platform with a Go REST API backend and Next.js frontend.

## Structure

```
api/   - Go (Gin) REST API with PostgreSQL/Supabase
ui/    - Next.js 14 (App Router) TypeScript frontend
```

## Quick Start

```bash
# Install dependencies for both projects
make setup

# Start API (Terminal 1) - requires air for hot reload
make api-dev

# Start UI (Terminal 2)
make ui-dev
```

The API runs on `http://localhost:8080` and the UI on `http://localhost:3001`.

## Setup

Each project has its own `.env.example` — copy and configure both:

```bash
cp api/.env.example api/.env    # Configure database & Supabase credentials
cp ui/.env.example ui/.env      # Configure Supabase & API URL
```

See [api/README.md](api/README.md) and [ui/README.md](ui/README.md) for detailed setup instructions.

## Available Commands

Run `make help` to see all available commands:

| Command | Description |
|---------|-------------|
| `make setup` | Set up both projects |
| `make api-dev` | Start API with hot reload |
| `make ui-dev` | Start UI dev server |
| `make api-build` | Build API binary |
| `make ui-build` | Build UI for production |
| `make ui-lint` | Lint UI code |
| `make docker-up` | Start services via Docker Compose |
