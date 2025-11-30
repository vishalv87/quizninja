# Configuration

## Overview

This folder contains application configuration modules for environment variables and site metadata. Configuration is centralized here to ensure type safety and easy access throughout the application.

## Files

| File | Purpose |
|------|---------|
| `env.ts` | Environment variable validation and access |
| `site.ts` | Site metadata (name, description, URLs) |

## Environment Variables (`env.ts`)

### Purpose

Provides type-safe access to environment variables with validation. Uses getter functions to defer access until runtime, preventing build-time errors.

### Required Variables

| Variable | Description |
|----------|-------------|
| `NEXT_PUBLIC_SUPABASE_URL` | Supabase project URL |
| `NEXT_PUBLIC_SUPABASE_ANON_KEY` | Supabase anonymous key |
| `NEXT_PUBLIC_API_BASE_URL` | Backend API base URL |
| `NEXT_PUBLIC_APP_URL` | Frontend application URL |

### Optional Variables

| Variable | Default | Description |
|----------|---------|-------------|
| `NEXT_PUBLIC_APP_NAME` | "QuizNinja" | Application name |
| `NEXT_PUBLIC_SENTRY_DSN` | "" | Sentry error tracking DSN |
| `NEXT_PUBLIC_GA_TRACKING_ID` | "" | Google Analytics ID |

### Usage

```tsx
import { env } from "@/config/env";

// Access environment variables
const supabaseUrl = env.supabase.url;
const apiBaseUrl = env.api.baseUrl;
const appName = env.app.name;
```

### Validation

Call `validateEnv()` at application startup to ensure all required variables are present:

```tsx
import { validateEnv } from "@/config/env";

// Will throw if required vars are missing
validateEnv();
```

### Why Getters?

Environment variables are accessed via getters to:
1. Defer access until actually needed (not at import time)
2. Allow for different values in different environments
3. Prevent "undefined" errors during build

```tsx
// Getter pattern - defers access
export const env = {
  supabase: {
    get url() {
      return getEnvVar("NEXT_PUBLIC_SUPABASE_URL");
    },
  },
} as const;
```

## Site Configuration (`site.ts`)

### Purpose

Centralized site metadata used for SEO, Open Graph, and consistent branding across the application.

### Configuration

```tsx
import { siteConfig } from "@/config/site";

// Available properties
siteConfig.name        // "QuizNinja"
siteConfig.description // "Test Your Knowledge, Compete with Friends"
siteConfig.url         // Application URL
siteConfig.ogImage     // "/og-image.png"
siteConfig.creator     // "QuizNinja Team"
siteConfig.links.github // GitHub repository URL
```

### Usage in Metadata

The root layout uses `siteConfig` for page metadata:

```tsx
// src/app/layout.tsx
import { siteConfig } from "@/config/site";

export const metadata: Metadata = {
  title: {
    default: siteConfig.name,
    template: `%s | ${siteConfig.name}`,
  },
  description: siteConfig.description,
  // ... OpenGraph, Twitter cards, etc.
};
```

## Setting Up Environment

### Development

1. Copy the example file:
   ```bash
   cp .env.example .env.local
   ```

2. Fill in required values:
   ```env
   NEXT_PUBLIC_SUPABASE_URL=https://your-project.supabase.co
   NEXT_PUBLIC_SUPABASE_ANON_KEY=your-anon-key
   NEXT_PUBLIC_API_BASE_URL=http://localhost:8080
   NEXT_PUBLIC_APP_URL=http://localhost:3001
   ```

### Production

Set environment variables in your deployment platform (Vercel, etc.):

```env
NEXT_PUBLIC_SUPABASE_URL=https://your-project.supabase.co
NEXT_PUBLIC_SUPABASE_ANON_KEY=your-anon-key
NEXT_PUBLIC_API_BASE_URL=https://api.yourdomain.com
NEXT_PUBLIC_APP_URL=https://yourdomain.com
```

## Adding New Configuration

### New Environment Variable

1. Add to `env.ts`:
   ```tsx
   export const env = {
     // ... existing
     myService: {
       get apiKey() {
         return getEnvVar("NEXT_PUBLIC_MY_SERVICE_KEY");
       },
     },
   } as const;
   ```

2. If required, add to validation:
   ```tsx
   export function validateEnv() {
     const required = [
       // ... existing
       "NEXT_PUBLIC_MY_SERVICE_KEY",
     ];
   }
   ```

3. Add to `.env.example`:
   ```env
   NEXT_PUBLIC_MY_SERVICE_KEY=
   ```

### New Site Config Property

Add to `site.ts`:

```tsx
export const siteConfig = {
  // ... existing
  myNewProperty: "value",
};
```

## Common Pitfalls

### Missing NEXT_PUBLIC_ Prefix

Client-side environment variables must be prefixed with `NEXT_PUBLIC_`:

```env
# WRONG - not accessible in browser
API_KEY=secret

# CORRECT - accessible in browser
NEXT_PUBLIC_API_KEY=secret
```

### Build-Time Access

Don't access environment variables at the top level of modules imported at build time:

```tsx
// WRONG - may fail at build
const url = process.env.NEXT_PUBLIC_API_URL;

// CORRECT - use getter
export const env = {
  get url() {
    return process.env.NEXT_PUBLIC_API_URL;
  },
};
```

## Related Documentation

- [Parent: Source Overview](../README.md)
- [API Client](../lib/api/README.md) - Uses `env.api.baseUrl`
- [Supabase Client](../lib/supabase/README.md) - Uses `env.supabase.*`
