> **For Claude:** REQUIRED SUB-SKILL: Use superpowers:executing-plans to implement this plan task-by-task. Follow TDD: write failing test → run → implement → run → commit.

**Goal:** Port all utility functions and 6 middleware layers from Go to TypeScript.
**Architecture:** Utilities (errors, validation, auth, idempotency, password) as standalone modules. Middleware (logging, error handler, security, CORS, rate limiting, auth) as Express middleware functions applied in exact Go ordering.
**Tech Stack:** TypeScript, Express v5, Pino, Helmet, cors, express-rate-limit, bcryptjs, jsonwebtoken

---

# Plan 03: Utilities & Middleware

**Priority**: High
**Estimated Complexity**: Medium
**Dependencies**: Plan 01 (project foundation)
**Impact**: Enables handlers (Plan 05) — all request processing flows through middleware. Utilities are used by every other layer.

---

## Problem Statement

The Go backend has 10 utility files and 6 middleware layers. These cross-cutting concerns (error handling, validation, authentication, rate limiting) must be ported to TypeScript with identical behavior. The middleware must be applied in the exact same order as Go's `main.go`.

---

## Context from Parent Plan

**Parent Plan:** `plans/go-to-typescript-rewrite-plan.md`
**Overall Goal:** Exact API compatibility.

Critical constraints:
- Middleware order must match Go's `main.go`: Logger → ErrorHandler → SecurityHeaders → CORS → RequestSize → GlobalRateLimit
- Auth middleware must: extract Bearer token → validate (Supabase or mock) → lookup DB user by Supabase ID → set `req.userId` (database UUID, NOT Supabase UUID)
- Error response format: `{ "error": { "code": 401, "message": "...", "details": "..." } }`
- Success response format: `{ "data": ..., "message": "..." }`
- 4-tier rate limiting: global (100/min IP), auth (5/min IP), write (20/min IP), per-user (90/min user)

---

## Current State Analysis

### Relevant Files

| File | Description |
|------|-------------|
| `api/utils/errors.go` | 4 error classes + ErrorResponse, HandleError, SuccessResponse, CreatedResponse |
| `api/utils/validation.go` | Email, name, URL, message validation + XSS/SQLi detection + sanitization |
| `api/utils/idempotency.go` | In-memory Map with 10min TTL, setInterval cleanup every 5min |
| `api/utils/supabase_auth.go` | HTTP validation via fetch to Supabase `/auth/v1/user` |
| `api/utils/authorization.go` | getUserIDFromContext, ownership checks |
| `api/utils/password.go` | bcryptjs hash/compare |
| `api/utils/mapper.go` | Data transformation helpers |
| `api/utils/mock_jwt.go` | Mock JWT generation/validation for testing |
| `api/utils/mock_auth_manager.go` | Mock auth manager |
| `api/middleware/auth.go` | Bearer token → Supabase/mock validation → DB user lookup |
| `api/middleware/cors.go` | CORS with configurable origins |
| `api/middleware/logging.go` | Structured request logging |
| `api/middleware/rate_limiter.go` | 4-tier rate limiting (global, auth, write, per-user) |
| `api/middleware/request_size.go` | Body size limits (3 tiers) |
| `api/middleware/security_headers.go` | X-Frame-Options, CSP, etc. |
| `api/main.go` | Middleware ordering (lines 71-84) |

### Key Code Patterns

```go
// api/utils/errors.go — Error types and response helpers
type ValidationError struct { Field, Message string }
type NotFoundError struct { Resource, ID string }
type ForbiddenError struct { Message string }
type APIError struct { Code int; Message string; Details string }

func ErrorResponse(c *gin.Context, statusCode int, message string, details ...string) {
    error := APIError{Code: statusCode, Message: message}
    if len(details) > 0 { error.Details = details[0] }
    c.JSON(statusCode, gin.H{"error": error})
}

func SuccessResponse(c *gin.Context, data interface{}, message ...string) {
    response := gin.H{"data": data}
    if len(message) > 0 { response["message"] = message[0] }
    c.JSON(http.StatusOK, response)
}
```

```go
// api/middleware/auth.go — Auth flow
// 1. Extract Bearer token from Authorization header
// 2. Validate with Supabase HTTP API or mock JWT
// 3. Convert Supabase UUID to database UUID
// 4. Lookup user in DB by Supabase ID
// 5. Set c.Set("user_id", dbUser.ID) — database UUID, not Supabase UUID
```

```go
// api/main.go — Middleware order (lines 71-84)
r.Use(middleware.Logger())
r.Use(middleware.ErrorHandler())
r.Use(middleware.SecurityHeaders())
r.Use(middleware.CORS(cfg.AllowedOrigins))
if cfg.RequestSizeLimitEnabled { r.Use(middleware.DefaultRequestSizeLimit()) }
if cfg.RateLimitEnabled { r.Use(middleware.GlobalRateLimit()) }
```

---

## Implementation

### Overview

Port utilities first (they're standalone), then middleware (which depends on utilities). Apply middleware in app.ts in exact Go order.

### Task 1: Error Utilities

**Files:**
- Create: `api-ts/src/utils/errors.ts`
- Test: `api-ts/tests/utils/errors.test.ts`

**Step 1: Write the failing test**
```typescript
// tests/utils/errors.test.ts
import { describe, it, expect } from 'vitest';
import { ValidationError, NotFoundError, ForbiddenError } from '@/utils/errors';

describe('error classes', () => {
  it('ValidationError has correct message', () => {
    const err = new ValidationError('Invalid email format');
    expect(err.message).toBe('Invalid email format');
    expect(err instanceof Error).toBe(true);
  });

  it('NotFoundError formats message correctly', () => {
    const err = new NotFoundError('Quiz', '123');
    expect(err.message).toBe('Quiz with ID 123 not found');
  });
});
```
Run: `cd api-ts && npx vitest run tests/utils/errors`
Expected: FAIL

**Step 2: Implement**
```typescript
// src/utils/errors.ts
export class ValidationError extends Error {
  field: string;
  constructor(message: string, field = '') {
    super(message);
    this.name = 'ValidationError';
    this.field = field;
  }
}

export class NotFoundError extends Error {
  resource: string;
  id: string;
  constructor(resource: string, id: string) {
    super(`${resource} with ID ${id} not found`);
    this.name = 'NotFoundError';
    this.resource = resource;
    this.id = id;
  }
}

export class ForbiddenError extends Error {
  constructor(message: string) {
    super(message);
    this.name = 'ForbiddenError';
  }
}

// Response helpers — match Go's exact JSON shape
export function errorResponse(res: Response, statusCode: number, message: string, details?: string) {
  const error: { code: number; message: string; details?: string } = { code: statusCode, message };
  if (details) error.details = details;
  res.status(statusCode).json({ error });
}

export function successResponse(res: Response, data: unknown, message?: string) {
  const response: Record<string, unknown> = { data };
  if (message) response.message = message;
  res.status(200).json(response);
}

export function createdResponse(res: Response, data: unknown, message?: string) {
  const response: Record<string, unknown> = { data };
  response.message = message ?? 'Resource created successfully';
  res.status(201).json(response);
}
```

**Step 3: Verify & Commit**

### Task 2: Validation Utilities

**Files:**
- Create: `api-ts/src/utils/validation.ts`
- Test: `api-ts/tests/utils/validation.test.ts`

Port regex patterns from Go's `validation.go`:
- `validateEmail()` — RFC 5322 regex
- `validateName()` — 2-100 chars, alphanumeric + punctuation
- `validateURL()` — HTTP/HTTPS only
- `validateMessage()` — Max 500 chars, XSS/SQLi checks
- `sanitizeString()`, `sanitizeHTML()`, `stripHTML()`
- `containsXSS()`, `containsSQLInjection()` — pattern matching
- `sanitizeEmail()`, `sanitizeName()`

### Task 3: Auth Utilities

**Files:**
- Create: `api-ts/src/utils/supabase-auth.ts`
- Create: `api-ts/src/utils/authorization.ts`
- Create: `api-ts/src/utils/mock-auth.ts`
- Create: `api-ts/src/utils/mock-jwt.ts`
- Create: `api-ts/src/utils/password.ts`

Key implementations:
- `validateSupabaseTokenHTTP(token, url, anonKey)` — HTTP call to `${SUPABASE_URL}/auth/v1/user`
- `getUserIdFromRequest(req)` → `req.userId` (typed via Express declaration merging)
- `hashPassword(password)` / `comparePassword(password, hash)` — bcryptjs wrapper
- Mock JWT with jsonwebtoken for testing

### Task 4: Idempotency & Mapper

**Files:**
- Create: `api-ts/src/utils/idempotency.ts`
- Create: `api-ts/src/utils/mapper.ts`
- Create: `api-ts/src/utils/index.ts`

Idempotency store: `Map<string, { result, timestamp }>` with 10min TTL, cleanup every 5min via `setInterval`.

### Task 5: Express Type Extensions

**Files:**
- Create: `api-ts/src/types/express.d.ts`

```typescript
// Extend Express Request to include userId from auth middleware
declare global {
  namespace Express {
    interface Request {
      userId?: string;  // Database UUID set by auth middleware
      authMethod?: 'supabase' | 'mock';
    }
  }
}
```

### Task 6: Middleware — Logging

**Files:**
- Create: `api-ts/src/middleware/logging.ts`

Use `pino-http` with request ID, method, path, status, latency, user ID.

### Task 7: Middleware — Error Handler

**Files:**
- Create: `api-ts/src/middleware/error-handler.ts`
- Test: `api-ts/tests/middleware/error-handler.test.ts`

Express v5 async error handler — dispatches by error type:
```typescript
// src/middleware/error-handler.ts
export function errorHandler(): ErrorRequestHandler {
  return (err, req, res, next) => {
    if (err instanceof ValidationError) {
      return errorResponse(res, 400, 'Validation error', err.message);
    }
    if (err instanceof NotFoundError) {
      return errorResponse(res, 404, 'Resource not found', err.message);
    }
    if (err instanceof ForbiddenError) {
      return errorResponse(res, 403, 'Access forbidden', err.message);
    }
    return errorResponse(res, 500, 'Internal server error', err.message);
  };
}
```

### Task 8: Middleware — Security, CORS, Request Size

**Files:**
- Create: `api-ts/src/middleware/security-headers.ts` — `helmet` or manual headers
- Create: `api-ts/src/middleware/cors.ts` — `cors` npm pkg with `credentials: true`
- Create: `api-ts/src/middleware/request-size.ts` — `express.json({ limit })` with 3 tiers

### Task 9: Middleware — Rate Limiting (4 tiers)

**Files:**
- Create: `api-ts/src/middleware/rate-limiter.ts`
- Test: `api-ts/tests/middleware/rate-limiter.test.ts`

```typescript
// src/middleware/rate-limiter.ts
import rateLimit from 'express-rate-limit';

export function globalRateLimit(config: Config) {
  return rateLimit({
    windowMs: 60 * 1000,
    max: config.rateLimitGlobal, // 100
    keyGenerator: (req) => req.ip,
    standardHeaders: true,
    legacyHeaders: false,
  });
}

export function authRateLimit(config: Config) { /* 5/min per IP */ }
export function writeRateLimit(config: Config) { /* 20/min per IP */ }
export function perUserRateLimit(config: Config) {
  return rateLimit({
    windowMs: 60 * 1000,
    max: config.rateLimitPerUser, // 90
    keyGenerator: (req) => `user:${req.userId}`,
  });
}
```

### Task 10: Middleware — Auth

**Files:**
- Create: `api-ts/src/middleware/auth.ts`
- Test: `api-ts/tests/middleware/auth.test.ts`
- Create: `api-ts/src/middleware/index.ts`

Must match Go's auth flow exactly:
1. Extract Bearer token from `Authorization` header
2. If mock auth enabled → `validateMockJWT(token)`
3. Else → `validateSupabaseTokenHTTP(token, url, anonKey)`
4. Convert Supabase ID to UUID
5. Look up database user by Supabase ID
6. Set `req.userId = dbUser.id` (DATABASE UUID, not Supabase UUID)
7. Set `req.authMethod = 'supabase' | 'mock'`

### Task 11: Wire Middleware in app.ts

**Files:**
- Modify: `api-ts/src/app.ts`

Apply middleware in exact Go order (from `main.go` lines 71-84):
```typescript
app.use(loggingMiddleware(logger));
app.use(errorHandler());
app.use(securityHeaders());
app.use(corsMiddleware(config.allowedOrigins));
if (config.requestSizeLimitEnabled) app.use(express.json({ limit: config.requestSizeDefault }));
if (config.rateLimitEnabled) app.use(globalRateLimit(config));
```

**Commit**
```bash
git add api-ts/src/utils/ api-ts/src/middleware/ api-ts/src/types/ api-ts/tests/
git commit -m "feat(api-ts): add all utilities and 6 middleware layers"
```

---

## Testing Strategy

### Unit Tests
- [ ] All 4 error classes construct correctly and have proper messages
- [ ] Validation functions reject invalid input (bad email, XSS patterns, SQL injection)
- [ ] Validation functions accept valid input
- [ ] Error handler middleware maps error types to correct HTTP status codes
- [ ] Auth middleware rejects requests without Bearer token (401)
- [ ] Auth middleware rejects invalid tokens (401)
- [ ] Rate limiter returns 429 when limit exceeded
- [ ] Idempotency store returns cached result within TTL

---

## Acceptance Criteria

- [ ] All utility functions ported from 10 Go files
- [ ] Error response JSON matches Go format: `{ "error": { "code", "message", "details?" } }`
- [ ] Success response JSON matches Go format: `{ "data": ..., "message?" }`
- [ ] Middleware applied in exact Go order
- [ ] Auth middleware sets `req.userId` as database UUID (not Supabase UUID)
- [ ] 4-tier rate limiting works (global, auth, write, per-user)
- [ ] Mock auth mode works for testing
- [ ] Security: mock auth rejected in production mode
- [ ] `npm run typecheck` passes
- [ ] `npm test` passes all tests
- [ ] No TypeScript errors introduced

---

## Dependencies & Execution Order

**Depends on:** Plan 01 (project foundation, config, logger)
**Blocks:** Plan 04 (repositories need error utilities), Plan 05 (handlers need middleware + utilities)

Can be done **in parallel** with Plan 02 (schema/models) since they don't depend on each other.
