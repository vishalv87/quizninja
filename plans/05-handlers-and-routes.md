> **For Claude:** REQUIRED SUB-SKILL: Use superpowers:executing-plans to implement this plan task-by-task. Follow TDD: write failing test → run → implement → run → commit.

**Goal:** Port all 13 handler files (~80 methods) and wire up all 88+ public routes to exactly match `api/routes/routes.go`.
**Architecture:** Express route handlers calling repository methods through dependency injection. Each handler file mirrors its Go counterpart. Route groups match Go's structure exactly (public, auth-rate-limited, protected with per-user rate limiting).
**Tech Stack:** TypeScript, Express v5, Zod, Vitest, supertest

---

# Plan 05: Handlers & Routes

**Priority**: High
**Estimated Complexity**: High
**Dependencies**: Plan 02 (Zod models), Plan 03 (middleware + utilities), Plan 04 (repository layer)
**Impact**: This is the API surface — all 88+ public endpoints that the frontend calls.

---

## Problem Statement

The Go backend has 13 handler files with ~80 methods covering all API endpoints. Each handler method validates input (Gin binding → Zod), calls repository methods, and formats responses. The route wiring in `routes/routes.go` (222 lines) must be replicated exactly — same paths, methods, middleware groups.

---

## Context from Parent Plan

**Parent Plan:** `plans/go-to-typescript-rewrite-plan.md`
**Overall Goal:** Exact API compatibility — same routes, same request/response shapes.

Critical constraints:
- Response format varies per handler: some use `SuccessResponse` wrapper (`{"data": ...}`), others use `c.JSON` directly (no wrapper). **Check each handler individually.**
- Concurrent queries in notification handler: Go uses `sync.WaitGroup` + goroutines → TS uses `Promise.allSettled`
- Auth handler has idempotency for registration
- Route paths and HTTP methods must match Go exactly

---

## Current State Analysis

### Relevant Files

| File | Description |
|------|-------------|
| `api/routes/routes.go` | All 88+ route definitions (222 lines) — THE reference |
| `api/handlers/quiz_handler.go` | 14 methods, largest handler (27KB) |
| `api/handlers/notification_handler.go` | 12 methods, concurrent queries (11KB) |
| `api/handlers/discussion_handler.go` | 12 methods (18KB) |
| `api/handlers/achievement_handler.go` | 10 methods (13KB) |
| `api/handlers/friends_handler.go` | 10 methods (14KB) |
| `api/handlers/auth_handler.go` | 6 methods with idempotency (11KB) |
| `api/handlers/preferences_handler.go` | 7 methods (7.3KB) |
| `api/handlers/leaderboard_handler.go` | 4 methods (6.9KB) |
| `api/handlers/rating_handler.go` | 6 methods (6.8KB) |
| `api/handlers/user_handler.go` | 1 method (5.7KB) |
| `api/handlers/favorites_handler.go` | 4 methods (3.5KB) |
| `api/handlers/categories_handler.go` | 2 methods (1.1KB) |
| `api/handlers/app_settings_handler.go` | 2 methods (1.6KB) |

### Key Code Patterns

```go
// Gin → Express pattern mapping
c.ShouldBindJSON(&req)     → const result = schema.safeParse(req.body)
c.Param("id")              → req.params.id
c.DefaultQuery("page","1") → (req.query.page as string) ?? '1'
c.Get("user_id")           → req.userId
c.JSON(200, gin.H{...})    → res.status(200).json({...})
c.GetHeader("X-...")       → req.headers['x-...']
c.Abort()                  → return

// Go handler constructor pattern
type QuizHandler struct {
    repo   *repository.Repository
    config *config.Config
}
func NewQuizHandler(cfg *config.Config) *QuizHandler {
    return &QuizHandler{repo: repository.NewRepository(), config: cfg}
}
```

```go
// Concurrent query pattern (notification_handler.go)
var wg sync.WaitGroup
wg.Add(2)
go func() {
    defer wg.Done()
    defer func() { if r := recover(); r != nil { /* handle panic */ } }()
    notifications, total, notifErr = h.repo.Notification.GetNotifications(...)
}()
go func() {
    defer wg.Done()
    defer func() { if r := recover(); r != nil { /* handle panic */ } }()
    unreadCount, unreadErr = h.repo.Notification.GetUnreadNotificationCount(...)
}()
wg.Wait()
```

```go
// Route groups from routes.go
// Public: GET /api/v1/quizzes, /featured, /category/:category, /categories
// Auth (rate limited): POST /api/v1/auth/register, /login
// Protected (auth + per-user rate limit): all other endpoints
```

---

## Implementation

### Overview

Implement handlers from simplest to most complex (matching the Go plan order), then wire all routes. Each handler gets its own test file using supertest.

### Task 1: Simple Handlers (categories, app-settings, preferences)

**Files:**
- Create: `api-ts/src/handlers/categories.handler.ts` (2 methods)
- Create: `api-ts/src/handlers/app-settings.handler.ts` (2 methods)
- Create: `api-ts/src/handlers/preferences.handler.ts` (7 methods)
- Test: `api-ts/tests/handlers/categories.test.ts`

**Step 1: Write failing test**
```typescript
// tests/handlers/categories.test.ts
import { describe, it, expect } from 'vitest';
import request from 'supertest';
import { createApp } from '@/app';

describe('GET /api/v1/categories', () => {
  it('should return list of categories', async () => {
    const app = createApp(testConfig);
    const res = await request(app).get('/api/v1/categories');
    expect(res.status).toBe(200);
    expect(res.body.data).toBeInstanceOf(Array);
  });
});
```
Run: `cd api-ts && npx vitest run tests/handlers/categories`
Expected: FAIL

**Step 2: Implement handler**
```typescript
// src/handlers/categories.handler.ts
import type { Repository } from '@/repository/interfaces';
import { successResponse } from '@/utils/errors';
import type { Request, Response, NextFunction } from 'express';

export function createCategoriesHandler(repo: Repository) {
  return {
    getCategories: async (req: Request, res: Response, next: NextFunction) => {
      try {
        const categories = await repo.categories.getCategories();
        successResponse(res, categories);
      } catch (err) {
        next(err);
      }
    },
    getCategoryGroups: async (req: Request, res: Response, next: NextFunction) => {
      try {
        const groups = await repo.categories.getCategoryGroups();
        successResponse(res, groups);
      } catch (err) {
        next(err);
      }
    },
  };
}
```

### Task 2: Auth Handler (6 methods — register with idempotency)

**Files:**
- Create: `api-ts/src/handlers/auth.handler.ts`
- Test: `api-ts/tests/handlers/auth.test.ts`

Methods: Register (idempotent), Login, Logout, GetProfile, UpdateProfile, GetUserStats

**Critical:** Register must use idempotency key from `X-Idempotency-Key` header.

### Task 3: User Handler (1 method)

**Files:**
- Create: `api-ts/src/handlers/user.handler.ts`

Method: GetUserProfile (view another user's public profile)

### Task 4: Quiz Handler (14 methods — largest)

**Files:**
- Create: `api-ts/src/handlers/quiz.handler.ts`
- Test: `api-ts/tests/handlers/quiz.test.ts`

Methods:
- GetQuizzes, GetQuizByID, GetQuizQuestions, GetFeaturedQuizzes, GetQuizzesByCategory
- GetUserQuizzes, StartQuizAttempt, SubmitQuizAttempt, UpdateQuizAttempt, AbandonQuizAttempt
- GetUserAttempts, GetAttemptDetails, GetUserQuizAttempt, GetUserLatestCompletedAttempt

### Task 5: Favorites & Rating Handlers

**Files:**
- Create: `api-ts/src/handlers/favorites.handler.ts` (4 methods)
- Create: `api-ts/src/handlers/rating.handler.ts` (6 methods)

### Task 6: Friends Handler (7 methods)

**Files:**
- Create: `api-ts/src/handlers/friends.handler.ts`

Methods: SendFriendRequest, GetFriendRequests, RespondToFriendRequest, CancelFriendRequest, GetFriends, RemoveFriend, SearchUsers

### Task 7: Notification Handler (12 methods — concurrent queries)

**Files:**
- Create: `api-ts/src/handlers/notification.handler.ts`
- Test: `api-ts/tests/handlers/notification.test.ts`

**Critical pattern — concurrent queries:**
```typescript
// TS equivalent of Go's WaitGroup + goroutines
const [notifResult, unreadResult] = await Promise.allSettled([
  repo.notification.getNotifications(userId, filters),
  repo.notification.getUnreadNotificationCount(userId),
]);

// Handle each result independently (matches Go's error-per-goroutine pattern)
if (notifResult.status === 'rejected') { /* handle */ }
if (unreadResult.status === 'rejected') { /* handle */ }
```

### Task 8: Leaderboard & Achievement Handlers

**Files:**
- Create: `api-ts/src/handlers/leaderboard.handler.ts` (4 methods)
- Create: `api-ts/src/handlers/achievement.handler.ts` (8 methods)

### Task 9: Discussion Handler (12 methods)

**Files:**
- Create: `api-ts/src/handlers/discussion.handler.ts`
- Create: `api-ts/src/handlers/index.ts`

Methods: CRUD for discussions, replies, likes, stats

### Task 10: Route Wiring

**Files:**
- Create: `api-ts/src/routes/routes.ts`
- Modify: `api-ts/src/app.ts` (add route setup)

Must exactly match `api/routes/routes.go` (222 lines):

```typescript
// src/routes/routes.ts
export function setupRoutes(app: Express, config: Config, repo: Repository) {
  // Health (already exists from Plan 01)

  const api = Router();
  app.use('/api/v1', api);

  api.get('/ping', (req, res) => res.json({ message: 'pong' }));

  // Public endpoints
  const quizzes = Router();
  api.use('/quizzes', quizzes);
  quizzes.get('', quizHandler.getQuizzes);
  quizzes.get('/featured', quizHandler.getFeaturedQuizzes);
  quizzes.get('/category/:category', quizHandler.getQuizzesByCategory);
  quizzes.get('/categories', categoriesHandler.getCategoryGroups);

  // ... all route groups from routes.go

  // Auth (with auth rate limit)
  const auth = Router();
  if (config.rateLimitEnabled) auth.use(authRateLimit(config));
  api.use('/auth', auth);
  auth.post('/register', authHandler.register);
  auth.post('/login', authHandler.login);

  // Protected (with auth middleware + per-user rate limit)
  const protected = Router();
  protected.use(authMiddleware(config));
  if (config.rateLimitEnabled) protected.use(perUserRateLimit(config));
  api.use('/', protected);
  // ... all 60+ protected endpoints
}
```

**Commit**
```bash
git add api-ts/src/handlers/ api-ts/src/routes/ api-ts/tests/handlers/
git commit -m "feat(api-ts): add all 13 handlers (80+ methods) and route wiring"
```

---

## Testing Strategy

### Unit Tests
- [ ] Each public endpoint returns correct status code
- [ ] Request validation rejects malformed input with 400
- [ ] Protected endpoints return 401 without auth token
- [ ] Notification handler handles concurrent query failures gracefully
- [ ] Register endpoint respects idempotency key
- [ ] Quiz listing supports pagination and filters

---

## Acceptance Criteria

- [ ] All 88+ public routes wired matching `api/routes/routes.go` exactly
- [ ] All 13 handler files created with ~80 methods total
- [ ] Route audit: extract all routes from both Go and TS servers, diff should be empty
- [ ] Request validation uses Zod schemas from Plan 02
- [ ] Response format matches Go per-handler (SuccessResponse wrapper vs direct JSON)
- [ ] Concurrent queries in notification handler use `Promise.allSettled`
- [ ] Auth rate limiting on `/auth/register` and `/auth/login`
- [ ] Per-user rate limiting on all protected endpoints
- [ ] `npm run typecheck` passes
- [ ] `npm test` passes all handler tests
- [ ] No TypeScript errors introduced

---

## Dependencies & Execution Order

**Depends on:** Plan 02 (Zod models for request validation), Plan 03 (middleware for auth/rate limiting, error utilities), Plan 04 (repository methods handlers call)
**Blocks:** Plan 06 (services and internal API), Plan 07 (end-to-end testing)

This is the largest plan by number of methods. Cannot start until the repository layer (Plan 04) is complete.
