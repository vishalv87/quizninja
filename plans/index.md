# Go-to-TypeScript Backend Rewrite - Implementation Sub-Plans

## Overview

This directory contains **7 implementation plans** split from:
`plans/go-to-typescript-rewrite-plan.md`

Each plan is designed to be implementable independently, with clear dependencies and acceptance criteria. The full rewrite covers 94 endpoints, 24 database tables, 100+ repository methods, 13 handler files, and 6 middleware layers.

---

## Execution Order

The recommended order for implementing these plans:

```
Plan 01: Project Scaffolding & Foundation
    ├──→ Plan 02: Drizzle Schema & Zod Models  ─┐
    │                                             ├──→ Plan 04: Repository Layer
    └──→ Plan 03: Utilities & Middleware ─────────┘           │
                                                              ├──→ Plan 05: Handlers & Routes
                                                              │           │
                                                              └──→ Plan 06: Services & Internal API
                                                                              │
                                                                    Plan 07: Docker, Testing & Integration
```

**Parallelism:** Plans 02 and 03 can be implemented simultaneously after Plan 01.

---

## Quick Reference

| # | Plan | Priority | Complexity | Dependencies | Status |
|---|------|----------|------------|--------------|--------|
| 01 | [Project Scaffolding & Foundation](./01-project-scaffolding-and-foundation.md) | High | Medium | None | ⬜ |
| 02 | [Drizzle Schema & Zod Models](./02-drizzle-schema-and-zod-models.md) | High | High | 01 | ⬜ |
| 03 | [Utilities & Middleware](./03-utilities-and-middleware.md) | High | Medium | 01 | ⬜ |
| 04 | [Repository Layer](./04-repository-layer.md) | High | High | 02, 03 | ⬜ |
| 05 | [Handlers & Routes](./05-handlers-and-routes.md) | High | High | 02, 03, 04 | ⬜ |
| 06 | [Services & Internal API](./06-services-and-internal-api.md) | Medium | Medium | 04, 05 | ⬜ |
| 07 | [Docker, Testing & Integration](./07-docker-testing-and-integration.md) | Medium | Medium | All | ⬜ |

---

## Plan Summaries

### Plan 01: Project Scaffolding & Foundation
Create the `api-ts/` directory with all project configuration (package.json, tsconfig, tsup, vitest), Zod-validated config loading for 30+ env vars, Pino logger, postgres.js database connection with exponential backoff retry, and Express v5 server with health endpoint and graceful shutdown matching Go's `main.go`.

### Plan 02: Drizzle Schema & Zod Models
Define Drizzle ORM schema for all 24 database tables (descriptive only — SQL migrations handle table creation) with correct type mappings for TEXT[], JSONB, DECIMAL, and all indexes/constraints. Create 15 Zod model files for request validation and response types with exact JSON field name parity to Go's `json:` struct tags.

### Plan 03: Utilities & Middleware
Port all 10 Go utility files (error classes, validation with XSS/SQLi detection, Supabase auth, idempotency, password hashing, mock auth) and 6 middleware layers (logging, error handler, security headers, CORS, request size limits, 4-tier rate limiting) with middleware applied in the exact same order as Go's `main.go`.

### Plan 04: Repository Layer
Port all 8 repository interfaces and 11 implementation files with 100+ methods using Drizzle query builder. Covers user CRUD, quiz management with complex filters, friend requests, notifications with soft delete, leaderboard rankings, achievements, discussions, and ratings. Preserves Go's transaction behavior exactly.

### Plan 05: Handlers & Routes
Port all 13 handler files with ~80 methods and wire 88+ public routes to exactly match `api/routes/routes.go`. Includes concurrent query patterns (Promise.allSettled for notifications), idempotent registration, and correct route group structure (public, auth-rate-limited, protected with per-user rate limiting).

### Plan 06: Services & Internal API
Port the achievement service with trigger-based checking for 10+ achievement keys, and the internal API with 5 handlers, X-Internal-API-Key auth middleware, internal HTTP client using native fetch, and 6 internal endpoints for cross-service communication.

### Plan 07: Docker, Testing & Integration
Create multi-stage Docker build for Cloud Run, comprehensive test suites for all layers (utils, middleware, repository, handlers, edge cases), test helpers with mock auth and DB isolation, and monorepo integration (root `dev:ts` script, Makefile targets, Docker Compose).

---

## Getting Started

1. Start with Plan 01 (no dependencies — can begin immediately)
2. After Plan 01, start Plans 02 and 03 in parallel
3. Mark each plan's status as you complete it: ⬜ → 🔄 → ✅
4. Refer to individual plan files for detailed implementation steps

## Execution

Each plan is formatted for direct execution via:
- **superpowers:executing-plans** (open a new Claude session in this directory, then use the executing-plans skill)
- **superpowers:subagent-driven-development** (dispatch a fresh subagent per plan in the current session)

Each plan file begins with a "For Claude" directive that automatically loads the correct skills.

---

## High-Risk Areas (from parent plan)

These must be verified across all plans:

1. **Empty arrays vs null (HIGH)** — Go returns `null` for nil slices, `[]` for empty. Check per-method.
2. **JSONB field shapes (HIGH)** — `notifications.data`, `quiz_attempts.answers` need exact TypeScript types.
3. **`omitempty` semantics (HIGH)** — Use `undefined` (not `null`) for empty optional fields.
4. **Response format inconsistency (MEDIUM)** — Some handlers use wrapper, others don't. Check each.
5. **Timestamp format (MEDIUM)** — Ensure consistent ISO 8601 output.
6. **Decimal precision (MEDIUM)** — Cast `numeric()` to `Number()` in repository layer.

---

## Generated

- **Source:** `plans/go-to-typescript-rewrite-plan.md`
- **Generated:** 2026-03-14
- **Total plans:** 7
- **Files analyzed:** 50+ Go source files
