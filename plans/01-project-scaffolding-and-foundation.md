> **For Claude:** REQUIRED SUB-SKILL: Use superpowers:executing-plans to implement this plan task-by-task. Follow TDD: write failing test → run → implement → run → commit.

**Goal:** Working Express v5 server with health endpoint, config, logging, and DB connection.
**Architecture:** Create `api-ts/` directory with TypeScript project scaffolding, Zod-validated config, Pino logging, postgres.js connection with retry, and graceful shutdown.
**Tech Stack:** TypeScript, Express v5, Zod, Pino, postgres.js, tsup, Vitest

---

# Plan 01: Project Scaffolding & Foundation

**Priority**: High
**Estimated Complexity**: Medium
**Dependencies**: None (can be done first)
**Impact**: Foundation for all other plans — nothing else can start until this is complete.

---

## Problem Statement

The QuizNinja Go backend needs a TypeScript equivalent. This plan creates the project skeleton: package.json, TypeScript config, build tooling, environment config, database connection, logger, and a minimal Express server with health check — matching the Go server's setup in `api/main.go`.

---

## Context from Parent Plan

**Parent Plan:** `plans/go-to-typescript-rewrite-plan.md`
**Overall Goal:** Rewrite the Go backend to TypeScript with exact API compatibility.

Key constraints:
- Must run on port 8080 (same as Go)
- Must use same env vars as Go (30+ vars, see `api/config/config.go`)
- DB connection must use retry with exponential backoff (matches Go's `cenkalti/backoff`)
- Graceful shutdown on SIGINT/SIGTERM with 30s timeout

---

## Current State Analysis

### Relevant Files

| File | Description |
|------|-------------|
| `api/config/config.go` | Go config with 30+ env vars, Zod validation target |
| `api/main.go` | Server setup, middleware ordering, graceful shutdown |
| `api/database/database.go` | DB connection with retry logic |
| `api/.env.example` | All environment variables documented |
| `api/utils/logger.go` | Logrus logger setup (port to Pino) |
| `ui/tsconfig.json` | Existing TS config pattern in monorepo (path alias `@/*`) |
| `package.json` (root) | Monorepo root with `concurrently` |

### Key Code Patterns

```go
// api/config/config.go — Config struct with env loading and defaults
cfg := &Config{
    DBHost:         getEnv("DB_HOST", "localhost"),
    DBPort:         getEnv("DB_PORT", "5432"),
    DBUser:         getEnv("DB_USER", "postgres"),
    DBPassword:     getEnv("DB_PASSWORD", ""),
    DBName:         getEnv("DB_NAME", "quizninja"),
    Port:           getEnv("PORT", "8080"),
    AllowedOrigins: getEnv("ALLOWED_ORIGINS", "http://localhost:3000"),
    // ... 30+ more vars
}
```

```go
// api/main.go — Server with timeouts and graceful shutdown
srv := &http.Server{
    Addr:           ":" + port,
    Handler:        r,
    ReadTimeout:    10 * time.Second,
    WriteTimeout:   10 * time.Second,
    IdleTimeout:    60 * time.Second,
    MaxHeaderBytes: 1 << 20,
}
// Graceful shutdown with 30s timeout on SIGINT/SIGTERM
```

---

## Implementation

### Overview

Create the `api-ts/` directory with all project configuration, a Zod-validated config module, Pino logger, postgres.js connection with retry, and Express v5 server with health endpoint and graceful shutdown.

### Task 1: Project Configuration Files

**Files:**
- Create: `api-ts/package.json`
- Create: `api-ts/tsconfig.json`
- Create: `api-ts/tsup.config.ts`
- Create: `api-ts/vitest.config.ts`
- Create: `api-ts/.env.example`
- Create: `api-ts/.env.test.example`

**Step 1: Create package.json**
```json
{
  "name": "quizninja-api-ts",
  "version": "1.0.0",
  "type": "module",
  "scripts": {
    "dev": "tsx watch src/server.ts",
    "build": "tsup",
    "start": "node dist/server.js",
    "test": "vitest run",
    "test:watch": "vitest",
    "typecheck": "tsc --noEmit",
    "migrate": "tsx src/db/migrate.ts"
  }
}
```

**Step 2: Install dependencies**
```bash
cd api-ts && npm install express@5 drizzle-orm postgres zod pino pino-http pino-pretty cors express-rate-limit helmet bcryptjs uuid dotenv jsonwebtoken
```

```bash
cd api-ts && npm install -D typescript tsup vitest tsx supertest drizzle-kit @types/node @types/express @types/cors @types/bcryptjs @types/uuid @types/jsonwebtoken @types/supertest
```

**Step 3: Create tsconfig.json**
- `strict: true`, `target: ES2022`, `module: NodeNext`
- Path alias: `@/*` → `./src/*`
- ESM output

**Step 4: Create tsup.config.ts**
- ESM format, node20 target, clean build

**Step 5: Create vitest.config.ts**
- Node environment, setup file at `tests/setup.ts`

**Step 6: Create .env.example**
- Document all 30+ env vars matching `api/config/config.go`

### Task 2: Config Module

**Files:**
- Create: `api-ts/src/config/config.ts`
- Test: `api-ts/tests/config/config.test.ts`

**Step 1: Write the failing test**
```typescript
// tests/config/config.test.ts
import { describe, it, expect } from 'vitest';
import { configSchema } from '@/config/config';

describe('configSchema', () => {
  it('should parse valid config with defaults', () => {
    const result = configSchema.safeParse({});
    expect(result.success).toBe(true);
    expect(result.data?.port).toBe('8080');
    expect(result.data?.dbHost).toBe('localhost');
  });

  it('should reject mock auth in release mode', () => {
    // Security check matching Go: cfg.GinMode == "release" && cfg.UseMockAuth
  });
});
```
Run: `cd api-ts && npx vitest run tests/config/`
Expected: FAIL (module not found)

**Step 2: Implement config with Zod**
```typescript
// src/config/config.ts
import { z } from 'zod';
import 'dotenv/config';

export const configSchema = z.object({
  // Database
  dbHost: z.string().default('localhost'),
  dbPort: z.string().default('5432'),
  dbUser: z.string().default('postgres'),
  dbPassword: z.string().default(''),
  dbName: z.string().default('quizninja'),
  port: z.string().default('8080'),
  nodeEnv: z.enum(['development', 'production', 'test']).default('development'),
  allowedOrigins: z.string().default('http://localhost:3000'),
  // Supabase (18 fields)
  useSupabase: z.coerce.boolean().default(false),
  // ... all 30+ fields from config.go
  // Rate limiting
  rateLimitEnabled: z.coerce.boolean().default(true),
  rateLimitGlobal: z.coerce.number().default(100),
  // ... etc
}).transform(cfg => ({
  ...cfg,
  // MB to bytes conversion for request sizes
  requestSizeDefault: cfg.requestSizeDefault * 1024 * 1024,
}));

export type Config = z.infer<typeof configSchema>;
export function loadConfig(): Config { /* parse from process.env */ }
```

**Step 3: Verify**
Run: `cd api-ts && npx vitest run tests/config/`
Expected: PASS

**Step 4: Commit**
```bash
git add api-ts/package.json api-ts/tsconfig.json api-ts/tsup.config.ts api-ts/vitest.config.ts api-ts/.env.example api-ts/src/config/
git commit -m "feat(api-ts): add project scaffolding and Zod config module"
```

### Task 3: Logger

**Files:**
- Create: `api-ts/src/utils/logger.ts`

**Step 1: Implement Pino logger**
```typescript
// src/utils/logger.ts
import pino from 'pino';
import type { Config } from '@/config/config';

export function createLogger(config: Config) {
  return pino({
    level: config.logLevel.toLowerCase(),
    redact: ['req.headers.authorization', 'password', 'passwordHash'],
    transport: config.nodeEnv === 'development'
      ? { target: 'pino-pretty' }
      : undefined,
  });
}
```

### Task 4: Database Connection

**Files:**
- Create: `api-ts/src/db/connection.ts`
- Test: `api-ts/tests/db/connection.test.ts`

**Step 1: Write the failing test**
```typescript
// tests/db/connection.test.ts
describe('createDbConnection', () => {
  it('should export a drizzle db instance', () => {
    // Test that the module exports the expected shape
  });
});
```

**Step 2: Implement with retry**
```typescript
// src/db/connection.ts
import postgres from 'postgres';
import { drizzle } from 'drizzle-orm/postgres-js';
import type { Config } from '@/config/config';

export function createDbPool(config: Config) {
  const connectionString = config.useSupabase
    ? `postgres://${config.supabaseDbUser}:${config.supabaseDbPassword}@${config.supabaseDbHost}:${config.supabaseDbPort}/${config.supabaseDbName}`
    : `postgres://${config.dbUser}:${config.dbPassword}@${config.dbHost}:${config.dbPort}/${config.dbName}`;

  return postgres(connectionString, {
    max: 25,
    idle_timeout: 120,
    max_lifetime: 300,
    prepare: false, // required for Supabase transaction pooler
  });
}

export async function connectWithRetry(pool: ReturnType<typeof postgres>, logger: pino.Logger) {
  // Exponential backoff: 100ms initial, 1.5x multiplier, 10s max, 30s timeout
  // Matches Go's cenkalti/backoff setup
}
```

### Task 5: Migration Runner

**Files:**
- Create: `api-ts/src/db/migrate.ts`

```typescript
// src/db/migrate.ts
// Reads and executes SQL files from ../api/database/migrations/ in order
// Same migration files used by the Go backend
```

### Task 6: Express App & Server

**Files:**
- Create: `api-ts/src/app.ts`
- Create: `api-ts/src/server.ts`
- Test: `api-ts/tests/app.test.ts`

**Step 1: Write the failing test**
```typescript
// tests/app.test.ts
import { describe, it, expect } from 'vitest';
import request from 'supertest';
import { createApp } from '@/app';

describe('Health endpoint', () => {
  it('GET /health returns 200 with status healthy', async () => {
    const app = createApp(/* mock config */);
    const res = await request(app).get('/health');
    expect(res.status).toBe(200);
    expect(res.body).toEqual({
      status: 'healthy',
      message: 'QuizNinja API is running',
    });
  });
});
```

**Step 2: Implement app.ts**
```typescript
// src/app.ts
import express from 'express';
import type { Config } from '@/config/config';

export function createApp(config: Config) {
  const app = express();
  app.get('/health', (req, res) => {
    res.json({ status: 'healthy', message: 'QuizNinja API is running' });
  });
  return app;
}
```

**Step 3: Implement server.ts**
```typescript
// src/server.ts — matches api/main.go structure
// 1. Load config
// 2. Create logger
// 3. Validate config
// 4. Connect to DB with retry
// 5. Create Express app
// 6. Start HTTP server with timeouts (read:10s, write:10s, idle:60s)
// 7. Graceful shutdown on SIGINT/SIGTERM with 30s timeout
```

**Step 4: Verify**
Run: `cd api-ts && npx vitest run`
Expected: PASS

**Step 5: Commit**
```bash
git add api-ts/src/ api-ts/tests/
git commit -m "feat(api-ts): add Express server with health endpoint, DB connection, and graceful shutdown"
```

---

## Testing Strategy

### Unit Tests
- [ ] Config schema parses valid env vars with correct defaults
- [ ] Config rejects mock auth in production mode
- [ ] Health endpoint returns `{ status: "healthy", message: "QuizNinja API is running" }`

---

## Acceptance Criteria

- [ ] `cd api-ts && npm install` succeeds
- [ ] `cd api-ts && npm run typecheck` passes with zero errors
- [ ] `cd api-ts && npm test` passes all tests
- [ ] `cd api-ts && npm run build` produces `dist/server.js`
- [ ] `cd api-ts && npm run dev` starts server on port 8080
- [ ] `GET /health` returns `200 { status: "healthy", message: "QuizNinja API is running" }`
- [ ] `GET /api/v1/ping` returns `200 { message: "pong" }`
- [ ] Server shuts down gracefully on SIGINT
- [ ] All existing tests pass
- [ ] No TypeScript errors introduced

---

## Dependencies & Execution Order

**Depends on:** None
**Blocks:** Plan 02, Plan 03, Plan 04, Plan 05, Plan 06, Plan 07

This is the foundation — all other plans require the project to exist with a working build system, config, and database connection.
