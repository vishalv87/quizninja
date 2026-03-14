> **For Claude:** REQUIRED SUB-SKILL: Use superpowers:executing-plans to implement this plan task-by-task. Follow TDD: write failing test → run → implement → run → commit.

**Goal:** Port all 11 repository files with 100+ methods from Go to TypeScript using Drizzle ORM.
**Architecture:** TypeScript interfaces mirroring Go's `interfaces.go`, with Drizzle query builder for standard queries and `db.execute(sql)` for complex ones. Repository aggregator pattern preserved.
**Tech Stack:** TypeScript, Drizzle ORM, postgres.js, Vitest

---

# Plan 04: Repository Layer

**Priority**: High
**Estimated Complexity**: High
**Dependencies**: Plan 02 (Drizzle schema + Zod models), Plan 03 (error utilities)
**Impact**: Enables all handlers (Plan 05) and services (Plan 06) — the entire data access layer.

---

## Problem Statement

The Go backend has 8 repository interfaces with 100+ methods across 11 implementation files. These must be ported to TypeScript using Drizzle ORM as the query builder. The largest file is `quiz.repository.go` at 1238 lines with complex filter/pagination queries.

**Critical:** Do NOT add transactions where Go doesn't have them. Match Go behavior exactly for parity, even where it has race conditions.

---

## Context from Parent Plan

**Parent Plan:** `plans/go-to-typescript-rewrite-plan.md`
**Overall Goal:** Exact query behavior parity with Go.

Key constraints:
- Drizzle query builder for standard queries, `db.execute(sql`...`)` for complex ones
- `numeric()` columns return string — cast to `Number()` in repository
- `null` vs `[]` handling must match Go per-method (Go is inconsistent)
- No new transactions where Go doesn't have them
- Repository aggregator pattern: single `Repository` class combining all sub-repositories

---

## Current State Analysis

### Relevant Files

| File | Description |
|------|-------------|
| `api/repository/interfaces.go` | 8 interfaces, 100+ method signatures (191 lines) |
| `api/repository/user_repository.go` | User CRUD, preferences, statistics (18 methods) |
| `api/repository/quiz_repository.go` | Quiz, attempts, favorites (1238 lines, 25 methods) |
| `api/repository/friends_repository.go` | Friend requests, friendships, search (16 methods) |
| `api/repository/notification_repository.go` | Notifications with soft delete, JSONB (14 methods) |
| `api/repository/leaderboard_repository.go` | Leaderboard, rankings, scores (8 methods) |
| `api/repository/achievement_repository.go` | Achievements, progress (9 methods) |
| `api/repository/discussion_repository.go` | Discussions, replies, likes |
| `api/repository/rating_repository.go` | Quiz ratings CRUD |
| `api/repository/categories_repository.go` | Categories listing |
| `api/repository/app_settings_repository.go` | App settings cache |
| `api/repository/preferences_repository.go` | User preferences |

### Key Code Patterns

```go
// api/repository/interfaces.go — Repository aggregator pattern
type Repository struct {
    User         UserRepositoryInterface
    Quiz         QuizRepositoryInterface
    Friends      FriendsRepositoryInterface
    Leaderboard  LeaderboardRepositoryInterface
    Achievement  AchievementRepositoryInterface
    Notification NotificationRepositoryInterface
    Discussion   DiscussionRepositoryInterface
    Rating       *RatingRepository
}

func NewRepository() *Repository {
    return &Repository{
        User:         NewUserRepository(),
        Quiz:         NewQuizRepository(),
        // ... all 8
    }
}
```

```go
// Go query pattern → Drizzle equivalent
// Go: database.DB.Query(sql, args) → Drizzle: db.select().from(table).where(...)
// Go: database.DB.QueryRow(...).Scan(...) → Drizzle: const [row] = await db.select()...
// Go: tx.Begin()/Commit()/Rollback() → Drizzle: db.transaction(async (tx) => { ... })
// Go: pq.StringArray → Drizzle text().array() handles natively
// Go: sql.NullString → Drizzle returns null for nullable columns
// Go: pq.Array(values) in ANY($n) → inArray() from drizzle-orm
```

---

## Implementation

### Overview

Define TypeScript interfaces matching Go's `interfaces.go`, then implement each repository file in dependency order. Use Drizzle query builder wherever possible, raw SQL for complex queries.

### Task 1: Repository Interfaces

**Files:**
- Create: `api-ts/src/repository/interfaces.ts`

```typescript
// src/repository/interfaces.ts — mirrors api/repository/interfaces.go exactly
export interface UserRepository {
  createUser(user: NewUser): Promise<User>;
  getUserByID(id: string): Promise<User | null>;
  getUserByEmail(email: string): Promise<User | null>;
  updateUser(user: Partial<User> & { id: string }): Promise<void>;
  deleteUser(id: string): Promise<void>;
  // ... all 18 methods from Go interface
}

export interface QuizRepository {
  getQuizByID(id: string): Promise<Quiz | null>;
  getQuizzes(filters: QuizFilters): Promise<{ quizzes: Quiz[]; total: number }>;
  // ... all 25 methods from Go interface
}

// ... all 8 interfaces

export interface Repository {
  user: UserRepository;
  quiz: QuizRepository;
  friends: FriendsRepository;
  leaderboard: LeaderboardRepository;
  achievement: AchievementRepository;
  notification: NotificationRepository;
  discussion: DiscussionRepository;
  rating: RatingRepository;
}
```

### Task 2: User Repository (foundational, needed by auth)

**Files:**
- Create: `api-ts/src/repository/user.repository.ts`
- Test: `api-ts/tests/repository/user.repository.test.ts`

**Step 1: Write failing test**
```typescript
describe('UserRepository', () => {
  it('should create and retrieve a user', async () => {
    const repo = new UserRepositoryImpl(db);
    const user = await repo.createUser({ email: 'test@example.com', ... });
    expect(user.id).toBeDefined();
    const found = await repo.getUserByID(user.id);
    expect(found?.email).toBe('test@example.com');
  });
});
```
Run: `cd api-ts && npx vitest run tests/repository/user`

**Step 2: Implement** (18 methods)
- CRUD: createUser, getUserByID, getUserByEmail, updateUser, deleteUser
- Preferences: createUserPreferences, getUserPreferences, updateUserPreferences, deleteUserPreferences
- Combined: getUserWithPreferences
- Status: updateUserOnlineStatus, updateUserLastActive
- Statistics: getUserStatistics, updateUserStatistics
- Supabase: getUserBySupabaseID (needed by auth middleware)

### Task 3: Simple Repositories (categories, app-settings, preferences)

**Files:**
- Create: `api-ts/src/repository/categories.repository.ts`
- Create: `api-ts/src/repository/app-settings.repository.ts`
- Create: `api-ts/src/repository/preferences.repository.ts`

These are small and have minimal dependencies.

### Task 4: Quiz Repository (largest — 1238 lines in Go, 25 methods)

**Files:**
- Create: `api-ts/src/repository/quiz.repository.ts`
- Test: `api-ts/tests/repository/quiz.repository.test.ts`

This is the most complex repository. Key methods:
- `getQuizzes(filters)` — complex dynamic WHERE, pagination, sorting
- `startQuizAttempt` / `submitQuizAttempt` — JSONB answers handling
- `getUserFavorites` — join with quiz table
- `getAttemptWithDetails` — multi-table join

**Drizzle pattern for complex filters:**
```typescript
const conditions: SQL[] = [];
if (filters.category) conditions.push(eq(quizzes.category, filters.category));
if (filters.difficulty) conditions.push(eq(quizzes.difficulty, filters.difficulty));
if (filters.search) conditions.push(ilike(quizzes.title, `%${filters.search}%`));

const results = await db
  .select()
  .from(quizzes)
  .where(and(...conditions))
  .orderBy(desc(quizzes.createdAt))
  .limit(filters.pageSize)
  .offset((filters.page - 1) * filters.pageSize);
```

### Task 5: Friends Repository

**Files:**
- Create: `api-ts/src/repository/friends.repository.ts`

16 methods: friend requests, friendships, user search, friend notifications.

### Task 6: Notification Repository

**Files:**
- Create: `api-ts/src/repository/notification.repository.ts`

14 methods with JSONB data handling and soft delete support.

### Task 7: Remaining Repositories

**Files:**
- Create: `api-ts/src/repository/leaderboard.repository.ts`
- Create: `api-ts/src/repository/achievement.repository.ts`
- Create: `api-ts/src/repository/discussion.repository.ts`
- Create: `api-ts/src/repository/rating.repository.ts`

### Task 8: Repository Aggregator

**Files:**
- Create: `api-ts/src/repository/index.ts`

```typescript
// src/repository/index.ts
import type { Repository } from './interfaces';

export function createRepository(db: DrizzleDb): Repository {
  return {
    user: new UserRepositoryImpl(db),
    quiz: new QuizRepositoryImpl(db),
    friends: new FriendsRepositoryImpl(db),
    leaderboard: new LeaderboardRepositoryImpl(db),
    achievement: new AchievementRepositoryImpl(db),
    notification: new NotificationRepositoryImpl(db),
    discussion: new DiscussionRepositoryImpl(db),
    rating: new RatingRepositoryImpl(db),
  };
}
```

**Commit**
```bash
git add api-ts/src/repository/ api-ts/tests/repository/
git commit -m "feat(api-ts): add repository layer with 11 files and 100+ methods"
```

---

## Testing Strategy

### Unit Tests
- [ ] User CRUD operations (create, read, update, delete)
- [ ] Quiz listing with filters and pagination
- [ ] Quiz attempt create/update with JSONB answers
- [ ] Friend request lifecycle (send, accept, reject, cancel)
- [ ] Notification CRUD with soft delete
- [ ] Leaderboard queries return correct ordering
- [ ] Achievement unlock is idempotent

---

## Acceptance Criteria

- [ ] All 8 repository interfaces defined in TypeScript matching Go's `interfaces.go`
- [ ] All 11 repository implementation files created
- [ ] 100+ methods ported with Drizzle query builder
- [ ] `numeric()` columns cast to `Number()` in repository layer
- [ ] No transactions added where Go doesn't have them
- [ ] `null` vs `[]` handling matches Go per-method
- [ ] Repository aggregator creates all sub-repositories
- [ ] `npm run typecheck` passes
- [ ] `npm test` passes all repository tests
- [ ] No TypeScript errors introduced

---

## Dependencies & Execution Order

**Depends on:** Plan 02 (Drizzle schema tables + Zod model types), Plan 03 (error classes used by repositories)
**Blocks:** Plan 05 (handlers call repository methods), Plan 06 (services + internal API use repositories)

The repository layer is the bridge between the database and the business logic. Cannot start until the schema and error utilities exist.
