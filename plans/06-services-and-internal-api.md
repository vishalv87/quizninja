> **For Claude:** REQUIRED SUB-SKILL: Use superpowers:executing-plans to implement this plan task-by-task. Follow TDD: write failing test → run → implement → run → commit.

**Goal:** Port the achievement service (trigger-based checking) and the internal API (5 handlers, auth middleware, HTTP client, routes).
**Architecture:** Achievement service with trigger-based logic for 10+ achievement keys. Internal API with X-Internal-API-Key authentication, separate route group, and internal HTTP client for cross-service calls.
**Tech Stack:** TypeScript, Express v5, Vitest

---

# Plan 06: Services & Internal API

**Priority**: Medium
**Estimated Complexity**: Medium
**Dependencies**: Plan 04 (repository layer), Plan 05 (handlers and routes setup)
**Impact**: Enables background job processing and cross-service communication for quiz scoring and statistics.

---

## Problem Statement

The Go backend has a service layer (achievement checking) and an internal API used by background jobs. The internal API has 6 endpoints protected by a shared secret, used for operations like quiz scoring, statistics updates, and achievement checking that are triggered by background processes.

---

## Context from Parent Plan

**Parent Plan:** `plans/go-to-typescript-rewrite-plan.md`
**Overall Goal:** Full feature parity including internal services.

Internal API endpoints (from `api/internal/routes/routes.go`):
- `POST /internal/v1/attempts/:attemptId/validate` — Validate attempt
- `PUT /internal/v1/attempts/:attemptId` — Update attempt
- `GET /internal/v1/quizzes/:quizId/questions` — Get questions with answers
- `POST /internal/v1/scoring/calculate` — Calculate quiz score
- `POST /internal/v1/users/:userId/statistics` — Update user statistics
- `POST /internal/v1/users/:userId/achievements/check` — Check achievements

Auth: `X-Internal-API-Key` header validated against `INTERNAL_API_SECRET` env var.

---

## Current State Analysis

### Relevant Files

| File | Description |
|------|-------------|
| `api/services/achievement_service.go` | Trigger-based achievement checking with 10+ keys |
| `api/internal/middleware/auth.go` | X-Internal-API-Key validation |
| `api/internal/handlers/attempt_handler.go` | Attempt validation and update |
| `api/internal/handlers/quiz_handler.go` | Get questions with answers |
| `api/internal/handlers/scoring_handler.go` | Score calculation |
| `api/internal/handlers/statistics_handler.go` | Statistics update |
| `api/internal/handlers/achievement_handler.go` | Achievement check trigger |
| `api/internal/client/client.go` | Internal HTTP client (uses net/http) |
| `api/internal/routes/routes.go` | Internal route definitions (51 lines) |
| `api/internal/models/models.go` | Internal request/response models |

### Key Code Patterns

```go
// api/internal/middleware/auth.go — Simple shared secret validation
func InternalAuthMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        apiKey := c.GetHeader("X-Internal-API-Key")
        if apiKey == "" || apiKey != cfg.InternalAPISecret {
            c.JSON(401, gin.H{"error": "Unauthorized"})
            c.Abort()
            return
        }
        c.Next()
    }
}
```

```go
// api/internal/routes/routes.go — 6 internal endpoints
internal := r.Group("/internal/v1")
internal.Use(InternalAuthMiddleware())
internal.POST("/attempts/:attemptId/validate", attemptHandler.ValidateAttempt)
internal.PUT("/attempts/:attemptId", attemptHandler.UpdateAttempt)
internal.GET("/quizzes/:quizId/questions", quizHandler.GetQuestionsWithAnswers)
internal.POST("/scoring/calculate", scoringHandler.CalculateScore)
internal.POST("/users/:userId/statistics", statisticsHandler.UpdateStatistics)
internal.POST("/users/:userId/achievements/check", achievementHandler.CheckAchievements)
```

---

## Implementation

### Overview

Implement the achievement service first (used by handlers and internal API), then the internal API components: auth middleware, handlers, HTTP client, and routes.

### Task 1: Achievement Service

**Files:**
- Create: `api-ts/src/services/achievement.service.ts`
- Create: `api-ts/src/services/index.ts`
- Test: `api-ts/tests/services/achievement.test.ts`

**Step 1: Write failing test**
```typescript
// tests/services/achievement.test.ts
describe('AchievementService', () => {
  it('should check and unlock quiz_first_complete achievement', async () => {
    const service = new AchievementService(mockRepo);
    const unlocked = await service.checkAchievements(userId, 'quiz_complete');
    expect(unlocked).toContainEqual(
      expect.objectContaining({ key: 'quiz_first_complete' })
    );
  });
});
```

**Step 2: Implement trigger-based achievement checking**
```typescript
// src/services/achievement.service.ts
// Trigger types: quiz_complete, streak_update, friend_add, discussion_create, etc.
// Each trigger checks relevant achievement keys against user's progress
```

### Task 2: Internal Auth Middleware

**Files:**
- Create: `api-ts/src/internal/middleware/auth.ts`
- Test: `api-ts/tests/internal/middleware/auth.test.ts`

```typescript
// src/internal/middleware/auth.ts
export function internalAuthMiddleware(config: Config): RequestHandler {
  return (req, res, next) => {
    const apiKey = req.headers['x-internal-api-key'];
    if (!apiKey || apiKey !== config.internalApiSecret) {
      return res.status(401).json({ error: 'Unauthorized' });
    }
    next();
  };
}
```

### Task 3: Internal Handlers (5 files)

**Files:**
- Create: `api-ts/src/internal/handlers/attempt.handler.ts`
- Create: `api-ts/src/internal/handlers/quiz.handler.ts`
- Create: `api-ts/src/internal/handlers/scoring.handler.ts`
- Create: `api-ts/src/internal/handlers/statistics.handler.ts`
- Create: `api-ts/src/internal/handlers/achievement.handler.ts`

### Task 4: Internal HTTP Client

**Files:**
- Create: `api-ts/src/internal/client/client.ts`

```typescript
// src/internal/client/client.ts — uses native fetch
export class InternalClient {
  constructor(
    private baseUrl: string,
    private apiSecret: string,
  ) {}

  async post<T>(path: string, body: unknown): Promise<T> {
    const res = await fetch(`${this.baseUrl}${path}`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
        'X-Internal-API-Key': this.apiSecret,
      },
      body: JSON.stringify(body),
    });
    if (!res.ok) throw new Error(`Internal API error: ${res.status}`);
    return res.json() as Promise<T>;
  }
  // get, put methods similarly
}
```

### Task 5: Internal Routes

**Files:**
- Create: `api-ts/src/internal/routes/routes.ts`
- Modify: `api-ts/src/app.ts` (add internal route setup)

Wire all 6 internal endpoints under `/internal/v1/` with internal auth middleware.

**Commit**
```bash
git add api-ts/src/services/ api-ts/src/internal/ api-ts/tests/
git commit -m "feat(api-ts): add achievement service and internal API (6 endpoints)"
```

---

## Testing Strategy

### Unit Tests
- [ ] Achievement service triggers correct achievements for each trigger type
- [ ] Internal auth middleware rejects missing/wrong API key (401)
- [ ] Internal auth middleware accepts correct API key
- [ ] Score calculation produces correct results
- [ ] Statistics update modifies user stats correctly
- [ ] Internal HTTP client includes correct headers

---

## Acceptance Criteria

- [ ] Achievement service checks 10+ achievement keys based on trigger types
- [ ] All 6 internal API endpoints wired and functional
- [ ] Internal auth validates `X-Internal-API-Key` against `INTERNAL_API_SECRET`
- [ ] Internal HTTP client uses native `fetch` with correct headers
- [ ] Internal routes match Go's `api/internal/routes/routes.go` exactly
- [ ] `npm run typecheck` passes
- [ ] `npm test` passes all tests
- [ ] No TypeScript errors introduced

---

## Dependencies & Execution Order

**Depends on:** Plan 04 (repository methods for achievement/statistics queries), Plan 05 (handler patterns and route setup)
**Blocks:** Plan 07 (integration testing needs internal API working)

The achievement service is also called from the public achievement handler (Plan 05), but can be stubbed there and implemented here.
