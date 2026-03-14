> **For Claude:** REQUIRED SUB-SKILL: Use superpowers:executing-plans to implement this plan task-by-task. Follow TDD: write failing test → run → implement → run → commit.

**Goal:** Type-safe Drizzle schema for all 24 tables and Zod request/response models for all 15 model files.
**Architecture:** Drizzle schema files describe existing DB tables (no migrations — SQL handles that). Zod schemas validate requests and define response types with exact JSON field name parity to Go.
**Tech Stack:** TypeScript, Drizzle ORM, Zod, Vitest

---

# Plan 02: Drizzle Schema & Zod Models

**Priority**: High
**Estimated Complexity**: High
**Dependencies**: Plan 01 (project scaffolding must exist)
**Impact**: Enables the repository layer (Plan 04) and handlers (Plan 05) to be type-safe.

---

## Problem Statement

The Go backend has 24 database tables and 15 model files with Go struct tags (`json:`, `binding:`, `omitempty`). These must be faithfully represented in TypeScript:
1. **Drizzle schema** — descriptive only (matches existing PostgreSQL tables, does NOT create them)
2. **Zod models** — request validation + response types with exact JSON field names matching Go's `json:` tags

---

## Context from Parent Plan

**Parent Plan:** `plans/go-to-typescript-rewrite-plan.md`
**Overall Goal:** Exact API compatibility — same request/response shapes.

Critical constraints:
- JSON field names must match Go `json:` tags exactly (snake_case)
- `TEXT[]` → `text().array()`, `JSONB` → `jsonb()`, `DECIMAL` → `numeric()` (returns string, cast to number)
- All indexes and constraints from `schema.sql` must be declared in Drizzle
- `omitempty` → use `undefined` (not `null`) for empty optional fields
- `*string` (Go pointer) → `string | null` (`.nullable()`)

---

## Current State Analysis

### Relevant Files

| File | Description |
|------|-------------|
| `api/database/schema.sql` | Full DB schema (808 lines, 24 tables, 64 indexes) |
| `api/models/user.go` | User + UserPreferences with custom StringArray type |
| `api/models/quiz.go` | Quiz + Question models with array fields |
| `api/models/quiz_attempt.go` | QuizAttempt with JSONB answers field |
| `api/models/notification.go` | Notification with JSONB data field |
| `api/models/discussion.go` | Discussion, Reply, and Like models |
| `api/models/achievement.go` | Achievement + UserAchievement models |
| `api/models/friends.go` | FriendRequest, Friendship, Friend models |
| `api/models/leaderboard.go` | LeaderboardEntry + UserRankInfo |
| `api/models/rating.go` | QuizRating model |
| `api/models/favorites.go` | UserQuizFavorite model |
| `api/models/statistics.go` | UserStatistics + QuizStatistics |
| `api/models/settings.go` | AppSettings model |
| `api/models/onboarding.go` | OnboardingStatus |
| `api/models/auth.go` | Auth request/response structs |
| `api/models/user_profile.go` | UserProfile response model |
| `api/repository/interfaces.go` | Repository interfaces (define what queries return) |

### Key Code Patterns

```go
// api/models/user.go — Custom StringArray type for PostgreSQL arrays
type StringArray []string

type User struct {
    ID                    uuid.UUID  `json:"id"`
    Email                 string     `json:"email"`
    PasswordHash          string     `json:"-"` // Never serialized
    Name                  string     `json:"name"`
    AvatarURL             *string    `json:"avatar_url,omitempty"` // Nullable + omitempty
    Level                 string     `json:"level"`
    TotalPoints           int        `json:"total_points"`
    // ...
}

type UserPreferences struct {
    SelectedCategories    StringArray            `json:"selected_categories"` // TEXT[]
    NotificationTypes     map[string]interface{} `json:"notification_types"` // JSONB
    OnboardingCompletedAt *time.Time             `json:"onboarding_completed_at,omitempty"`
}
```

```sql
-- api/database/schema.sql — Table with arrays, JSONB, timestamps
CREATE TABLE quizzes (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    title VARCHAR(255) NOT NULL,
    tags TEXT[] DEFAULT '{}',
    time_limit INTEGER NOT NULL DEFAULT 300,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    -- ...
);
CREATE INDEX idx_quizzes_category ON quizzes(category);
```

---

## Implementation

### Overview

First define all 24 Drizzle schema tables (grouped by domain), then create 15 Zod model files with request validation schemas and response types. Verify with `drizzle-kit introspect` against a live DB.

### Task 1: Core Drizzle Schema Tables

**Files:**
- Create: `api-ts/src/db/schema/users.ts`
- Create: `api-ts/src/db/schema/quizzes.ts`
- Create: `api-ts/src/db/schema/questions.ts`
- Create: `api-ts/src/db/schema/quiz-attempts.ts`
- Create: `api-ts/src/db/schema/quiz-statistics.ts`
- Create: `api-ts/src/db/schema/quiz-ratings.ts`
- Create: `api-ts/src/db/schema/categories.ts`
- Create: `api-ts/src/db/schema/difficulty-levels.ts`
- Create: `api-ts/src/db/schema/notification-frequencies.ts`
- Create: `api-ts/src/db/schema/app-settings.ts`
- Create: `api-ts/src/db/schema/user-preferences.ts`
- Create: `api-ts/src/db/schema/achievements.ts`
- Create: `api-ts/src/db/schema/user-achievements.ts`
- Create: `api-ts/src/db/schema/user-category-performance.ts`
- Create: `api-ts/src/db/schema/user-quiz-favorites.ts`
- Create: `api-ts/src/db/schema/user-rank-history.ts`
- Create: `api-ts/src/db/schema/leaderboard-snapshots.ts`
- Create: `api-ts/src/db/schema/friend-requests.ts`
- Create: `api-ts/src/db/schema/friendships.ts`
- Create: `api-ts/src/db/schema/notifications.ts`
- Create: `api-ts/src/db/schema/discussions.ts`
- Create: `api-ts/src/db/schema/discussion-replies.ts`
- Create: `api-ts/src/db/schema/discussion-likes.ts`
- Create: `api-ts/src/db/schema/discussion-reply-likes.ts`
- Create: `api-ts/src/db/schema/index.ts`
- Create: `api-ts/drizzle.config.ts`

**Step 1: Write the failing test**
```typescript
// tests/db/schema.test.ts
import { describe, it, expect } from 'vitest';
import * as schema from '@/db/schema';

describe('Drizzle schema', () => {
  it('should export all 24 tables', () => {
    expect(schema.users).toBeDefined();
    expect(schema.quizzes).toBeDefined();
    expect(schema.questions).toBeDefined();
    // ... all 24 tables
  });

  it('should have correct column types for users', () => {
    // Verify users table has expected columns
  });
});
```
Run: `cd api-ts && npx vitest run tests/db/schema`
Expected: FAIL

**Step 2: Implement schema files**

Each schema file follows this pattern:
```typescript
// src/db/schema/users.ts
import { pgTable, uuid, varchar, text, boolean, integer, numeric,
         timestamp, index, uniqueIndex } from 'drizzle-orm/pg-core';

export const users = pgTable('users', {
  id: uuid('id').primaryKey().defaultRandom(),
  email: varchar('email', { length: 255 }).notNull().unique(),
  passwordHash: text('password_hash').notNull(),
  name: varchar('name', { length: 100 }).notNull(),
  avatarUrl: text('avatar_url'),
  level: varchar('level', { length: 50 }).notNull().default('beginner'),
  totalPoints: integer('total_points').notNull().default(0),
  currentStreak: integer('current_streak').notNull().default(0),
  bestStreak: integer('best_streak').notNull().default(0),
  totalQuizzesCompleted: integer('total_quizzes_completed').notNull().default(0),
  averageScore: numeric('average_score').notNull().default('0'),
  isOnline: boolean('is_online').notNull().default(false),
  lastActive: timestamp('last_active', { withTimezone: true }),
  authMethod: varchar('auth_method', { length: 50 }).notNull().default('supabase'),
  supabaseId: text('supabase_id').unique(),
  lastAuthMethod: varchar('last_auth_method', { length: 50 }),
  migratedAt: timestamp('migrated_at', { withTimezone: true }),
  createdAt: timestamp('created_at', { withTimezone: true }).notNull().defaultNow(),
  updatedAt: timestamp('updated_at', { withTimezone: true }).notNull().defaultNow(),
}, (table) => [
  index('idx_users_email').on(table.email),
  index('idx_users_supabase_id').on(table.supabaseId),
]);
```

**Special type handling:**
- `TEXT[]` → `text('col').array()` (for `quizzes.tags`, `questions.options`, `user_preferences.selected_categories`)
- `JSONB` → `jsonb('col')` (for `quiz_attempts.answers`, `notifications.data`, `user_preferences.notification_types`)
- `DECIMAL/NUMERIC` → `numeric('col')` (returns string from DB, cast to `Number()` in repository)

**Step 3: Verify**
Run: `cd api-ts && npx vitest run tests/db/schema`
Expected: PASS

**Step 4: Commit**
```bash
git add api-ts/src/db/schema/ api-ts/drizzle.config.ts api-ts/tests/db/
git commit -m "feat(api-ts): add Drizzle schema for all 24 database tables"
```

### Task 2: Zod Models (15 files)

**Files:**
- Create: `api-ts/src/models/user.ts`
- Create: `api-ts/src/models/auth.ts`
- Create: `api-ts/src/models/quiz.ts`
- Create: `api-ts/src/models/quiz-attempt.ts`
- Create: `api-ts/src/models/notification.ts`
- Create: `api-ts/src/models/achievement.ts`
- Create: `api-ts/src/models/friends.ts`
- Create: `api-ts/src/models/leaderboard.ts`
- Create: `api-ts/src/models/discussion.ts`
- Create: `api-ts/src/models/rating.ts`
- Create: `api-ts/src/models/favorites.ts`
- Create: `api-ts/src/models/statistics.ts`
- Create: `api-ts/src/models/settings.ts`
- Create: `api-ts/src/models/onboarding.ts`
- Create: `api-ts/src/models/user-profile.ts`
- Create: `api-ts/src/models/index.ts`

**Step 1: Write failing tests**
```typescript
// tests/models/auth.test.ts
import { describe, it, expect } from 'vitest';
import { registerRequestSchema, loginRequestSchema } from '@/models/auth';

describe('auth schemas', () => {
  it('should validate register request with required fields', () => {
    const result = registerRequestSchema.safeParse({
      email: 'test@example.com',
      password: 'password123',
      name: 'Test User',
    });
    expect(result.success).toBe(true);
  });

  it('should reject register without email', () => {
    const result = registerRequestSchema.safeParse({
      password: 'password123',
      name: 'Test User',
    });
    expect(result.success).toBe(false);
  });
});
```
Run: `cd api-ts && npx vitest run tests/models/`
Expected: FAIL

**Step 2: Implement models**

Go-to-Zod mapping for each model file:
```typescript
// src/models/auth.ts — matches api/models/auth.go
import { z } from 'zod';

// Go: RegisterRequest with binding:"required" → Zod non-optional
export const registerRequestSchema = z.object({
  email: z.string().email(),
  password: z.string().min(1),
  name: z.string().min(1),
});

// Go: LoginRequest
export const loginRequestSchema = z.object({
  email: z.string().email(),
  password: z.string().min(1),
});

// Go: AuthResponse — response type (no validation needed)
export type AuthResponse = {
  user: UserResponse;
  token: string;
  message: string;
};

export type RegisterRequest = z.infer<typeof registerRequestSchema>;
export type LoginRequest = z.infer<typeof loginRequestSchema>;
```

**Key Go→Zod mappings to apply consistently across all 15 files:**
- `binding:"required"` → field is non-optional (Zod default)
- `binding:"email"` → `.email()`
- `binding:"oneof=..."` → `.enum([...])`
- `json:"snake_case"` → Use exact same field names in Zod schema
- `*string` → `.nullable()`
- `time.Time` → `z.string()` (ISO 8601)
- `uuid.UUID` → `z.string().uuid()`
- `omitempty` → field uses `undefined` (omit from Zod output schemas)

**Step 3: Verify**
Run: `cd api-ts && npx vitest run tests/models/`
Expected: PASS

**Step 4: Commit**
```bash
git add api-ts/src/models/ api-ts/tests/models/
git commit -m "feat(api-ts): add Zod request/response models for all 15 model files"
```

---

## Testing Strategy

### Unit Tests
- [ ] All 24 Drizzle schema tables export correctly with expected columns
- [ ] Zod request schemas reject invalid input (missing required fields, bad email format)
- [ ] Zod request schemas accept valid input
- [ ] JSON field names in response types match Go `json:` tags exactly
- [ ] Nullable fields accept both value and null
- [ ] Array fields default to empty array `[]` not `null`

---

## Acceptance Criteria

- [ ] All 24 tables defined in `src/db/schema/` matching `api/database/schema.sql`
- [ ] All indexes and constraints from schema.sql declared
- [ ] All 15 model files created in `src/models/` with Zod schemas
- [ ] JSON field names match Go `json:` tags exactly (snake_case)
- [ ] `TEXT[]` columns use `text().array()`, `JSONB` uses `jsonb()`, `DECIMAL` uses `numeric()`
- [ ] `npm run typecheck` passes
- [ ] `npm test` passes all model/schema tests
- [ ] No TypeScript errors introduced

---

## Dependencies & Execution Order

**Depends on:** Plan 01 (project must exist with package.json, tsconfig, drizzle dependency)
**Blocks:** Plan 04 (repository layer needs schema), Plan 05 (handlers need models)

Depends on Plan 01 because the Drizzle and Zod dependencies must be installed, and the project build system must be working.
