> **For Claude:** REQUIRED SUB-SKILL: Use superpowers:executing-plans to implement this plan task-by-task. Follow TDD: write failing test → run → implement → run → commit.

**Goal:** Multi-stage Docker build, comprehensive test coverage, and monorepo integration for the TypeScript backend.
**Architecture:** Dockerfile for Cloud Run deployment, test suites for all layers (utils, middleware, repository, handlers), and root package.json integration for parallel development.
**Tech Stack:** TypeScript, Docker, Vitest, supertest, concurrently

---

# Plan 07: Docker, Testing & Integration

**Priority**: Medium
**Estimated Complexity**: Medium
**Dependencies**: All previous plans (01-06)
**Impact**: Production-readiness — deployment, test coverage, and developer experience.

---

## Problem Statement

The TypeScript backend needs: (1) a Docker image for Cloud Run deployment matching the Go backend's container setup, (2) comprehensive test coverage across all layers, and (3) integration into the monorepo so developers can run `npm run dev:ts` from root.

---

## Context from Parent Plan

**Parent Plan:** `plans/go-to-typescript-rewrite-plan.md`
**Overall Goal:** Drop-in replacement for the Go backend.

Verification strategy from parent plan:
1. Route audit — diff all registered routes between Go and TS backends
2. Response snapshot tests — identical requests to both backends, compare JSON
3. Frontend smoke test — point React frontend at TS backend
4. Edge case tests — empty arrays, null fields, float precision, idempotency, rate limiting, expired tokens

---

## Current State Analysis

### Relevant Files

| File | Description |
|------|-------------|
| `api/Dockerfile` | Existing Go Dockerfile (reference for structure) |
| `package.json` (root) | Monorepo root with `concurrently` dev script |
| `Makefile` | Commands: setup, dev, api-dev, ui-dev, docker-up |
| `docker-compose.yml` | Currently only Go API service |
| `api/.env.example` | All env vars (reference for .env.test) |

### Key Code Patterns

```dockerfile
# From parent plan — multi-stage Docker build
FROM node:20-alpine AS builder
WORKDIR /app
COPY package*.json ./
RUN npm ci
COPY tsconfig.json tsup.config.ts ./
COPY src/ src/
RUN npm run build

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

```json
// Root package.json — existing dev script
{
  "scripts": {
    "dev": "concurrently -k -n api,ui -c blue,magenta \"cd api && air\" \"cd ui && npm run dev\""
  }
}
```

---

## Implementation

### Overview

Create Dockerfile, write comprehensive tests for all layers, set up test helpers, and integrate into the monorepo root.

### Task 1: Dockerfile

**Files:**
- Create: `api-ts/Dockerfile`

Multi-stage build:
- Stage 1 (builder): `node:20-alpine`, install deps, build with tsup
- Stage 2 (runtime): `node:20-alpine`, copy dist + node_modules + migrations
- Build context from monorepo root: `docker build -f api-ts/Dockerfile -t quizninja-api-ts .`
- Must copy `api/database/migrations` (shared SQL migrations)
- Health check: `wget --spider http://localhost:8080/health`
- Port 8080

### Task 2: Test Setup & Helpers

**Files:**
- Create: `api-ts/tests/setup.ts`
- Create: `api-ts/tests/helpers/auth.ts`
- Create: `api-ts/tests/helpers/db.ts`

```typescript
// tests/setup.ts
// - Load .env.test
// - Setup test database connection
// - Run migrations before all tests
// - Cleanup after all tests

// tests/helpers/auth.ts
// - Generate mock JWT tokens for test requests
// - Create authenticated supertest agent

// tests/helpers/db.ts
// - Test database connection
// - Transaction-based test isolation (begin → test → rollback)
// - Seed data helpers
```

### Task 3: Utility Tests

**Files:**
- Create/extend: `api-ts/tests/utils/validation.test.ts`
- Create/extend: `api-ts/tests/utils/idempotency.test.ts`
- Create/extend: `api-ts/tests/utils/errors.test.ts`

Tests for:
- Email validation (valid, invalid, edge cases)
- Name validation (length, special chars)
- URL validation (http/https, invalid protocols)
- XSS detection (script tags, event handlers, data URIs)
- SQL injection detection (UNION, DROP, comment patterns)
- Sanitization functions (HTML escape, strip tags)
- Idempotency store (TTL expiry, concurrent access)
- Error class construction and message formatting

### Task 4: Middleware Tests

**Files:**
- Create/extend: `api-ts/tests/middleware/auth.test.ts`
- Create/extend: `api-ts/tests/middleware/error-handler.test.ts`
- Create/extend: `api-ts/tests/middleware/rate-limiter.test.ts`

Tests for:
- Auth middleware: no token (401), invalid token (401), expired token (401), valid token (sets req.userId)
- Error handler: ValidationError → 400, NotFoundError → 404, ForbiddenError → 403, unknown → 500
- Rate limiter: under limit (200), at limit (429), rate limit headers present

### Task 5: Repository Tests (against test DB)

**Files:**
- Create/extend: `api-ts/tests/repository/user.repository.test.ts`
- Create/extend: `api-ts/tests/repository/quiz.repository.test.ts`

Tests against actual PostgreSQL with migrations applied:
- User CRUD lifecycle
- Quiz listing with filters and pagination
- Quiz attempt create → update → submit
- Favorite add/remove/check
- Friend request lifecycle
- Notification create → read → soft delete → restore

### Task 6: Handler Tests (supertest)

**Files:**
- Create/extend: `api-ts/tests/handlers/` (one file per handler)

Tests using supertest against Express app with mock auth:
- Public endpoints accessible without auth
- Protected endpoints return 401 without auth
- Request validation returns 400 for invalid input
- Correct response shapes matching Go

### Task 7: Edge Case Tests

**Files:**
- Create: `api-ts/tests/edge-cases.test.ts`

From parent plan's verification strategy:
- Empty arrays vs null in responses
- Float precision for scores (numeric → Number())
- Idempotency key handling for registration
- Expired token behavior
- Rate limiting across multiple requests
- JSONB field round-trip serialization

### Task 8: Monorepo Integration

**Files:**
- Modify: `package.json` (root) — add `dev:ts` script
- Modify: `Makefile` — add api-ts targets

```json
// Root package.json — add dev:ts script
{
  "scripts": {
    "dev": "concurrently -k -n api,ui -c blue,magenta \"cd api && air\" \"cd ui && npm run dev\"",
    "dev:ts": "concurrently -k -n api-ts,ui -c green,magenta \"cd api-ts && npm run dev\" \"cd ui && npm run dev\""
  }
}
```

Add to Makefile:
```makefile
api-ts-dev:
	cd api-ts && npm run dev

api-ts-build:
	cd api-ts && npm run build

api-ts-test:
	cd api-ts && npm test
```

### Task 9: Docker Compose Update

**Files:**
- Modify: `docker-compose.yml` — add api-ts service (commented out by default)

```yaml
# api-ts:
#   build:
#     context: .
#     dockerfile: api-ts/Dockerfile
#   ports:
#     - "8080:8080"
#   env_file: api-ts/.env
#   networks:
#     - quizninja-net
```

**Commit**
```bash
git add api-ts/Dockerfile api-ts/tests/ package.json Makefile docker-compose.yml
git commit -m "feat(api-ts): add Dockerfile, comprehensive tests, and monorepo integration"
```

---

## Testing Strategy

### Full Test Matrix
- [ ] Utils: validation, errors, idempotency, sanitization (unit tests)
- [ ] Middleware: auth, error handler, rate limiting (unit tests with supertest)
- [ ] Repository: all 8 repositories against test PostgreSQL (integration tests)
- [ ] Handlers: all 13 handlers via supertest (integration tests with mock auth)
- [ ] Edge cases: empty arrays, null fields, float precision, JSONB round-trip
- [ ] Docker: `docker build` succeeds, container starts, health check passes

---

## Acceptance Criteria

- [ ] `docker build -f api-ts/Dockerfile -t quizninja-api-ts .` succeeds
- [ ] Docker container starts and `/health` returns 200
- [ ] `npm run dev:ts` from monorepo root starts both api-ts and ui
- [ ] `cd api-ts && npm test` passes all tests with >80% coverage on critical paths
- [ ] Route audit: all 88+ public + 6 internal routes match Go backend
- [ ] Frontend smoke test: React UI works with TS backend on same endpoints
- [ ] Edge case tests pass (empty arrays, null fields, float precision)
- [ ] Makefile targets work: `make api-ts-dev`, `make api-ts-build`, `make api-ts-test`
- [ ] All existing tests pass
- [ ] No TypeScript errors introduced

---

## Dependencies & Execution Order

**Depends on:** All previous plans (01-06) — the entire backend must be implemented
**Blocks:** None (this is the final plan)

This plan wraps up the rewrite by ensuring production readiness, test coverage, and developer experience.
