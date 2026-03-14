# Plan: Rewrite QuizNinja Go Backend to TypeScript

## Context

The QuizNinja Go backend (`api/`) needs to be completely rewritten in TypeScript. The Go backend is a substantial application — 13 handler files, 11 repository files, 6 middleware layers, 62+ SQL migrations, 100+ repository methods, and an internal service API. **No feature should be missed.**

The TypeScript backend will live in a new `api-ts/` directory alongside the existing `api/`, with **exact API compatibility** — same routes, same ports, same request/response shapes — so the frontend works with either backend without changes.

## Tech Stack

| Concern | Choice | Rationale |
|---------|--------|-----------|
| Framework | Express v5 | Same as Kinio project |
| ORM | Drizzle ORM | SQL-like, type-safe, close to Go's raw SQL approach |
| Validation | Zod | Same as Kinio |
| Logging | Pino | Same as Kinio |
| Testing | Vitest + supertest | Same as Kinio |
| Build | tsup | Same as Kinio |
| DB Driver | postgres (postgres.js) | Works with Drizzle, handles pooling natively |
| Runtime | Node.js 20+ | ESM, stable |

## Directory Structure

```
api-ts/
├── src/
│   ├── config/
│   │   └── config.ts              # Env loading with Zod validation (30+ vars)
│   ├── db/
│   │   ├── connection.ts          # Drizzle + postgres.js pool setup
│   │   ├── migrate.ts             # Legacy SQL migration runner
│   │   └── schema/                # Drizzle schema (25 tables)
│   │       ├── users.ts
│   │       ├── quizzes.ts
│   │       ├── questions.ts
│   │       ├── quiz-attempts.ts
│   │       ├── quiz-statistics.ts
│   │       ├── quiz-ratings.ts
│   │       ├── categories.ts
│   │       ├── difficulty-levels.ts
│   │       ├── notification-frequencies.ts
│   │       ├── app-settings.ts
│   │       ├── user-preferences.ts
│   │       ├── achievements.ts
│   │       ├── user-achievements.ts
│   │       ├── user-category-performance.ts
│   │       ├── user-quiz-favorites.ts
│   │       ├── user-rank-history.ts
│   │       ├── leaderboard-snapshots.ts
│   │       ├── friend-requests.ts
│   │       ├── friendships.ts
│   │       ├── notifications.ts
│   │       ├── discussions.ts
│   │       ├── discussion-replies.ts
│   │       ├── discussion-likes.ts
│   │       ├── discussion-reply-likes.ts
│   │       └── index.ts
│   ├── models/                     # Zod schemas + TypeScript types
│   │   ├── user.ts
│   │   ├── auth.ts
│   │   ├── quiz.ts
│   │   ├── quiz-attempt.ts
│   │   ├── notification.ts
│   │   ├── achievement.ts
│   │   ├── friends.ts
│   │   ├── leaderboard.ts
│   │   ├── discussion.ts
│   │   ├── rating.ts
│   │   ├── favorites.ts
│   │   ├── statistics.ts
│   │   ├── settings.ts
│   │   ├── onboarding.ts
│   │   ├── user-profile.ts
│   │   └── index.ts
│   ├── middleware/
│   │   ├── auth.ts                 # Supabase JWT + mock auth
│   │   ├── cors.ts                 # CORS (cors npm pkg)
│   │   ├── logging.ts              # Request logging (pino-http)
│   │   ├── rate-limiter.ts         # 4-tier rate limiting (express-rate-limit)
│   │   ├── request-size.ts         # Body size limits
│   │   ├── security-headers.ts     # Security headers (helmet)
│   │   ├── error-handler.ts        # Global error handler
│   │   └── index.ts
│   ├── repository/
│   │   ├── interfaces.ts           # TypeScript interfaces (mirrors Go interfaces.go)
│   │   ├── user.repository.ts
│   │   ├── quiz.repository.ts
│   │   ├── friends.repository.ts
│   │   ├── leaderboard.repository.ts
│   │   ├── achievement.repository.ts
│   │   ├── notification.repository.ts
│   │   ├── discussion.repository.ts
│   │   ├── rating.repository.ts
│   │   ├── preferences.repository.ts
│   │   ├── categories.repository.ts
│   │   ├── app-settings.repository.ts
│   │   └── index.ts
│   ├── handlers/
│   │   ├── auth.handler.ts         # Register, Login, Logout, GetProfile, UpdateProfile, GetUserStats
│   │   ├── quiz.handler.ts         # 14 methods (largest handler)
│   │   ├── user.handler.ts         # GetUserProfile
│   │   ├── friends.handler.ts      # 7 methods
│   │   ├── notification.handler.ts # 13 methods (concurrent queries)
│   │   ├── leaderboard.handler.ts  # 4 methods
│   │   ├── achievement.handler.ts  # 8 methods
│   │   ├── discussion.handler.ts   # 12 methods
│   │   ├── favorites.handler.ts    # 4 methods
│   │   ├── rating.handler.ts       # 6 methods
│   │   ├── categories.handler.ts   # 2 methods
│   │   ├── preferences.handler.ts  # 7 methods
│   │   ├── app-settings.handler.ts # 2 methods
│   │   └── index.ts
│   ├── services/
│   │   ├── achievement.service.ts  # Trigger-based achievement checking
│   │   └── index.ts
│   ├── routes/
│   │   └── routes.ts               # All public + protected routes
│   ├── internal/
│   │   ├── middleware/
│   │   │   └── auth.ts             # X-Internal-API-Key validation
│   │   ├── handlers/
│   │   │   ├── attempt.handler.ts
│   │   │   ├── quiz.handler.ts
│   │   │   ├── scoring.handler.ts
│   │   │   ├── statistics.handler.ts
│   │   │   └── achievement.handler.ts
│   │   ├── client/
│   │   │   └── client.ts           # Internal HTTP client (fetch)
│   │   └── routes/
│   │       └── routes.ts
│   ├── utils/
│   │   ├── supabase-auth.ts        # Token validation via Supabase HTTP API
│   │   ├── authorization.ts        # Ownership checks
│   │   ├── logger.ts               # Pino logger setup
│   │   ├── validation.ts           # Email/name/URL validation + XSS/SQLi detection + sanitization
│   │   ├── idempotency.ts          # In-memory store with TTL (Map + setInterval cleanup)
│   │   ├── password.ts             # bcryptjs wrapper
│   │   ├── errors.ts               # ValidationError, NotFoundError, ForbiddenError, ApiError classes
│   │   ├── mapper.ts               # Data transformation helpers
│   │   ├── mock-auth.ts            # Mock auth for testing
│   │   ├── mock-jwt.ts             # Mock JWT generation/validation
│   │   └── index.ts
│   ├── app.ts                      # Express app setup (middleware stack + routes)
│   └── server.ts                   # HTTP server + graceful shutdown (SIGINT/SIGTERM)
├── tests/
│   ├── setup.ts
│   ├── helpers/
│   │   ├── auth.ts
│   │   └── db.ts
│   ├── handlers/                   # One test file per handler
│   ├── middleware/
│   ├── repository/
│   └── utils/
├── package.json
├── tsconfig.json
├── tsup.config.ts
├── vitest.config.ts
├── drizzle.config.ts
├── Dockerfile
├── .env.example
└── .env.test.example
```

## Implementation Phases

### Phase 1: Project Scaffolding & Foundation

**Goal:** Working Express v5 server with health endpoint, config, logging, and DB connection.

**Files to create:**
- `package.json` — dependencies (see deps list below)
- `tsconfig.json` — strict, ESM, ES2022, path aliases `@/*`
- `tsup.config.ts` — ESM output, node20 target
- `vitest.config.ts` — node environment, setup file
- `.env.example` — all 30+ env vars documented
- `src/config/config.ts` — Zod-validated config loading from env
- `src/utils/logger.ts` — Pino setup with redaction of sensitive fields
- `src/db/connection.ts` — postgres.js pool (max:25, idle:120s, lifetime:300s, prepare:false)
- `src/db/migrate.ts` — runs existing SQL migrations from `../api/database/migrations/`
- `src/app.ts` — Express app with health endpoint
- `src/server.ts` — HTTP server with graceful shutdown

**Key dependencies:**
```
express@5, drizzle-orm, postgres, zod, pino, pino-http, pino-pretty,
cors, express-rate-limit, helmet, bcryptjs, uuid,
dotenv, jsonwebtoken
```

**Dev deps:**
```
typescript, tsup, vitest, tsx, supertest, drizzle-kit,
@types/node, @types/express, @types/cors, @types/bcryptjs,
@types/uuid, @types/jsonwebtoken, @types/supertest
```

**DB connection with retry:** Async retry loop with exponential backoff (100ms initial, 1.5x multiplier, 10s max interval, 30s timeout) — matches Go's `cenkalti/backoff` setup.

### Phase 2: Drizzle Schema Definition (25 tables)

**Goal:** Define Drizzle schema matching the existing PostgreSQL schema exactly.

The schema is **descriptive only** — it does NOT create tables. The 62 SQL migrations handle that. The Drizzle schema exists so queries are type-safe.

**Reference file:** `api/database/schema.sql` (808 lines, auto-generated)

**Special type handling:**
- `TEXT[]` columns (`quizzes.tags`, `questions.options`, `user_preferences.selected_categories`) → `text('col').array()`
- `JSONB` columns (`quiz_attempts.answers`, `notifications.data`, `user_preferences.notification_types`) → `jsonb('col')`
- `DECIMAL` columns (scores) → `numeric('col')` — returns string, cast to number in code
- All indexes and constraints from schema.sql must be declared

**Verification:** Run `drizzle-kit introspect` against the live DB and compare with hand-written schema.

### Phase 3: Models (Zod Schemas + Types)

**Goal:** TypeScript types and Zod request validation schemas for all 15 model files.

**Key mappings from Go:**
- `binding:"required"` → Zod `.min(1)` or non-optional
- `binding:"email"` → `.email()`
- `binding:"oneof=..."` → `.enum([...])`
- `json:"snake_case"` → Use exact same snake_case field names in JSON responses
- `*string` (Go pointer) → `string | null` (Zod `.nullable()`)
- `time.Time` → `string` (ISO 8601 date string)
- `uuid.UUID` → `string` (UUID format)
- `omitempty` → Conditionally omit fields (use `undefined` in objects, which JSON.stringify skips)

**Critical:** All JSON field names must match the Go `json:` tags exactly for frontend compatibility.

### Phase 4: Utilities

**Goal:** Port all utility functions.

| Go File | TS File | Notes |
|---------|---------|-------|
| `errors.go` | `errors.ts` | 4 error classes + `errorResponse`, `handleError`, `successResponse`, `createdResponse` |
| `validation.go` | `validation.ts` | Regex patterns (email, name, URL, XSS, SQLi detection), sanitization functions |
| `idempotency.go` | `idempotency.ts` | In-memory Map with 10min TTL, setInterval cleanup every 5min |
| `supabase_auth.go` | `supabase-auth.ts` | HTTP validation via `fetch` to `{SUPABASE_URL}/auth/v1/user` |
| `authorization.go` | `authorization.ts` | `getUserIDFromContext` → `req.userId`, ownership checks |
| `password.go` | `password.ts` | bcryptjs hash/compare |
| `mapper.go` | `mapper.ts` | Data transformation helpers |
| `mock_auth.go` + `mock_jwt.go` | `mock-auth.ts` + `mock-jwt.ts` | jsonwebtoken for mock JWT |
| `logger.go` (logrus) | `logger.ts` (Pino) | Already done in Phase 1 |

### Phase 5: Middleware Stack (6 layers)

**Middleware order must match Go's `main.go`:**
1. **Logging** — `pino-http` with request ID, method, path, status, latency, user ID
2. **Error handler** — Express v5 async error handler, dispatches by error type
3. **Security headers** — `helmet` or manual headers (X-Frame-Options, CSP, etc.)
4. **CORS** — `cors` npm pkg, credentials:true, all methods/headers
5. **Request size** — `express.json({ limit })` with 3 tiers (default:10MB, auth:1MB, write:5MB)
6. **Global rate limit** — `express-rate-limit` memory store, 100 req/min per IP

**Auth middleware** (applied to protected route group):
- Extract Bearer token → validate via Supabase HTTP or mock JWT → look up DB user by Supabase ID → set `req.userId` (database UUID, NOT Supabase UUID)
- Extend Express `Request` type via declaration merging

**Per-user rate limit** (applied to protected route group):
- 90 req/min keyed by `user:${req.userId}`

**Auth rate limit** (applied to auth route group only):
- 5 req/min per IP

### Phase 6: Repository Layer (11 files, 100+ methods)

**Goal:** Port all repository methods using Drizzle query builder, falling back to `db.execute(sql`...`)` for complex queries.

**Implementation order (by dependency):**
1. `user.repository.ts` — foundational, needed by auth middleware
2. `categories.repository.ts` — simple, no FK deps
3. `app-settings.repository.ts` — simple
4. `preferences.repository.ts` — depends on user
5. `quiz.repository.ts` — largest (1238 lines in Go), complex filters/pagination
6. `friends.repository.ts` — depends on user
7. `notification.repository.ts` — depends on user, JSONB handling
8. `leaderboard.repository.ts` — depends on user, friends
9. `achievement.repository.ts` — depends on user
10. `discussion.repository.ts` — depends on user, quiz
11. `rating.repository.ts` — depends on user, quiz

**Go-to-Drizzle pattern mapping:**
- `database.DB.Query(sql, args)` → `db.select().from(table).where(...)`
- `database.DB.QueryRow(...).Scan(...)` → `const [row] = await db.select()...`
- `tx.Begin()/Commit()/Rollback()` → `db.transaction(async (tx) => { ... })`
- `pq.StringArray` → Drizzle `text().array()` handles natively
- `sql.NullString` etc. → Drizzle returns `null` for nullable columns
- `pq.Array(values)` in `ANY($n)` → `inArray()` from drizzle-orm or `sql` template

**Critical: DO NOT add transactions where Go doesn't have them.** Match Go behavior exactly for parity, even where it has race conditions.

### Phase 7: Handlers (13 files, ~80 methods)

**Goal:** Port all handler methods.

**Implementation order (simplest first):**
1. `categories.handler.ts` — 2 methods, public
2. `app-settings.handler.ts` — 2 methods
3. `preferences.handler.ts` — 7 methods
4. `auth.handler.ts` — 6 methods (critical: register with idempotency)
5. `user.handler.ts` — 1 method
6. `quiz.handler.ts` — 14 methods (largest)
7. `favorites.handler.ts` — 4 methods
8. `rating.handler.ts` — 6 methods
9. `friends.handler.ts` — 7 methods
10. `notification.handler.ts` — 13 methods (concurrent queries via `Promise.allSettled`)
11. `leaderboard.handler.ts` — 4 methods
12. `achievement.handler.ts` — 8 methods
13. `discussion.handler.ts` — 12 methods

**Gin-to-Express pattern mapping:**
| Go (Gin) | TypeScript (Express v5) |
|----------|------------------------|
| `c.ShouldBindJSON(&req)` | `schema.safeParse(req.body)` |
| `c.Param("id")` | `req.params.id` |
| `c.DefaultQuery("page", "1")` | `(req.query.page as string) ?? '1'` |
| `c.Get("user_id")` | `req.userId` |
| `c.JSON(200, gin.H{...})` | `res.status(200).json({...})` |
| `c.GetHeader("X-...")` | `req.headers['x-...']` |
| `c.Abort()` | `return` |

**Concurrent query pattern** (notification handler):
```ts
// Go: sync.WaitGroup + goroutines + panic recovery
// TS: Promise.allSettled (settles independently, handles rejections)
const [notifResult, unreadResult] = await Promise.allSettled([
  repo.notification.getNotifications(userId, filters),
  repo.notification.getUnreadNotificationCount(userId),
]);
```

**Response format parity:** Some handlers use `utils.SuccessResponse` (wraps in `{"data": ...}`), others use `c.JSON` directly (no wrapper). Each handler must be checked individually against the Go source.

### Phase 8: Routes

**Goal:** Wire up all routes to exactly match `api/routes/routes.go` (221 lines).

**Route groups (from Go source):**
- `/health` — public health check
- `/api/v1/ping` — public ping
- `/api/v1/quizzes` — public quiz list (GET), featured, by-category, categories
- `/api/v1/categories` — public flat category list
- `/api/v1/config/app-settings` — public app settings
- `/api/v1/preferences/*` — public preference lookups
- `/api/v1/auth/register|login` — auth rate limited
- `/api/v1/` (protected) — all authenticated endpoints (see routes.go for full list)

### Phase 9: Services

**Goal:** Port `achievement.service.ts` — trigger-based achievement checking with 10+ achievement keys.

### Phase 10: Internal API

**Goal:** Port the 5 internal handlers, internal auth middleware, internal HTTP client, and internal routes.

**Internal routes (from `internal/routes/routes.go`):**
- `POST /internal/v1/attempts/:attemptId/validate`
- `PUT /internal/v1/attempts/:attemptId`
- `GET /internal/v1/quizzes/:quizId/questions`
- `POST /internal/v1/scoring/calculate`
- `POST /internal/v1/users/:userId/statistics`
- `POST /internal/v1/users/:userId/achievements/check`

**Auth:** `X-Internal-API-Key` header validated against `INTERNAL_API_SECRET` env var.

### Phase 11: Dockerfile

**Goal:** Multi-stage Docker build for Cloud Run.

```dockerfile
# Stage 1: Build
FROM node:20-alpine AS builder
WORKDIR /app
COPY package*.json ./
RUN npm ci
COPY tsconfig.json tsup.config.ts ./
COPY src/ src/
RUN npm run build

# Stage 2: Runtime
FROM node:20-alpine
WORKDIR /app
COPY --from=builder /app/dist ./dist
COPY --from=builder /app/node_modules ./node_modules
COPY --from=builder /app/package.json ./
COPY api/database/migrations ./database/migrations
EXPOSE 8080
HEALTHCHECK CMD wget --spider http://localhost:8080/health || exit 1
CMD ["node", "dist/server.js"]
```

Build context from monorepo root: `docker build -f api-ts/Dockerfile -t quizninja-api-ts .`

### Phase 12: Testing

**Goal:** Comprehensive test coverage for all layers.

- **Utils tests:** validation, idempotency, errors, sanitization
- **Middleware tests:** auth flow, rate limiting, error handler
- **Repository tests:** against test DB with migrations applied
- **Handler tests:** supertest against Express app with mock auth

### Phase 13: Monorepo Integration

Update root `package.json`:
```json
{
  "scripts": {
    "dev": "concurrently -k -n api,ui -c blue,magenta \"cd api && air\" \"cd ui && npm run dev\"",
    "dev:ts": "concurrently -k -n api-ts,ui -c green,magenta \"cd api-ts && npm run dev\" \"cd ui && npm run dev\""
  }
}
```

## High-Risk Areas & Mitigations

### 1. Empty arrays vs null (HIGH)
Go returns `null` for nil slices, `[]` for empty slices. Inconsistent across the codebase. Each repository method must be checked.
**Mitigation:** Always use `?? []` in TS to match Go's `make([]T, 0)` pattern, and explicitly return `null` where Go uses nil slices.

### 2. JSONB field shapes (HIGH)
`notifications.data`, `quiz_attempts.answers`, `user_preferences.notification_types` use JSONB. Go marshals/unmarshals manually.
**Mitigation:** Define exact TypeScript types for each JSONB field and test round-trip serialization.

### 3. `omitempty` semantics (HIGH)
Go omits zero-value fields with `omitempty`. In TS, `undefined` is omitted by `JSON.stringify` but `null` is not.
**Mitigation:** Map each Go model's `omitempty` fields and use `undefined` (not `null`) for empty optional fields.

### 4. Response format inconsistency (MEDIUM)
Some Go handlers use `SuccessResponse` wrapper, others use `c.JSON` directly.
**Mitigation:** Check each handler individually against the Go source.

### 5. Timestamp format (MEDIUM)
Go's `time.Time` → RFC 3339 with microseconds. JS `Date.toISOString()` → milliseconds only.
**Mitigation:** Use Drizzle's `timestamp` mode which returns `Date` objects — ensure consistent ISO 8601 output.

### 6. Decimal precision (MEDIUM)
Go uses `float64` for scores. Drizzle's `numeric` returns string.
**Mitigation:** Cast to `Number()` in repository layer, match Go's rounding patterns.

## Verification Strategy

1. **Route audit:** Extract all registered routes from both servers, diff method + path
2. **Response snapshot tests:** For each endpoint, send identical requests to both backends, compare JSON responses
3. **Frontend smoke test:** Point the React frontend at the TS backend, exercise all features
4. **Edge case tests:** Empty arrays, null fields, float precision, idempotency, rate limiting, expired tokens

## Critical Go Source Files

These files are the authoritative references during implementation:

- `api/routes/routes.go` — Complete route map (221 lines)
- `api/internal/routes/routes.go` — Internal API routes
- `api/database/schema.sql` — Full DB schema (808 lines, 25 tables)
- `api/repository/interfaces.go` — All repository interfaces (100+ methods)
- `api/utils/errors.go` — Error types and response helpers
- `api/config/config.go` — All 30+ env vars with defaults
- `api/handlers/notification_handler.go` — Concurrent query pattern
- `api/models/user.go` — Custom StringArray type for PostgreSQL arrays
- `api/middleware/auth.go` — Auth middleware flow
- `api/middleware/rate_limiter.go` — 4-tier rate limiting config
