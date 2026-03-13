# QuizNinja UI - Performance Optimization Plan

> **Last Updated**: 2025-11-16
> **Status**: Planning Phase
> **Estimated Timeline**: 12 weeks (3 months)

---

## Executive Summary

This document outlines a comprehensive, phased approach to improving the performance of QuizNinja UI. The plan addresses both **actual performance improvements** (reducing load times, bundle sizes, render times) and **perceived performance enhancements** (making the app feel faster through UX patterns).

### Current State Analysis
- **Build size**: 427 MB (.next folder)
- **Client components**: 103 "use client" directives
- **Code splitting**: None (beyond Next.js defaults)
- **Component memoization**: Minimal to none
- **Image optimization**: Not using Next.js Image component
- **Bundle analysis**: Not configured

### Goals
- Reduce initial bundle size by 60-70%
- Achieve Lighthouse Performance Score of 90+
- Make all interactions feel instant (<100ms perceived response)
- Improve Core Web Vitals to "Good" thresholds
- Enhance user satisfaction and engagement

---

## Table of Contents

1. [Phase 1: Quick Wins & Foundation](#phase-1-quick-wins--foundation)
2. [Phase 2: Code Splitting & Bundle Optimization](#phase-2-code-splitting--bundle-optimization)
3. [Phase 3: Server Components & Data Fetching](#phase-3-server-components--data-fetching)
4. [Phase 4: Component Optimization & Memoization](#phase-4-component-optimization--memoization)
5. [Phase 5: Perceived Performance & UX Enhancements](#phase-5-perceived-performance--ux-enhancements)
6. [Phase 6: Quiz Taking Experience](#phase-6-quiz-taking-experience)
7. [Phase 7: Image & Asset Optimization](#phase-7-image--asset-optimization)
8. [Phase 8: Advanced Optimizations](#phase-8-advanced-optimizations)
9. [Phase 9: Monitoring & Continuous Improvement](#phase-9-monitoring--continuous-improvement)
10. [Success Criteria & KPIs](#success-criteria--kpis)
11. [Implementation Guidelines](#implementation-guidelines)

---

## Phase 1: Quick Wins & Foundation

**Timeline**: Week 1
**Type**: Real Performance + Perceived Performance
**Risk Level**: Low
**Effort**: Low

### Objectives
- Achieve immediate, visible performance improvements
- Establish baseline metrics for comparison
- Set up tooling for ongoing optimization
- Improve perceived performance during navigation

### Tasks

#### 1.1 Remove Production Console Logs
**Type**: Real Performance
**Impact**: Medium
**Files Affected**:
- `src/middleware.ts` (multiple instances)
- `src/app/(dashboard)/layout.tsx`
- `src/app/layout.tsx`
- Debug/test files

**Action**:
```typescript
// Before
console.log('User authenticated:', user)

// After - use logger utility with environment check
import { logger } from '@/lib/logger'
logger.debug('User authenticated:', user) // Only logs in development
```

**Expected Outcome**:
- Reduced JavaScript execution overhead
- Cleaner production console
- No debugging information leaked to users

---

#### 1.2 Install & Configure Bundle Analyzer
**Type**: Real Performance (Tooling)
**Impact**: High (Visibility)

**Installation**:
```bash
npm install --save-dev @next/bundle-analyzer
```

**Configuration** (`next.config.js`):
```javascript
const withBundleAnalyzer = require('@next/bundle-analyzer')({
  enabled: process.env.ANALYZE === 'true',
})

module.exports = withBundleAnalyzer({
  reactStrictMode: true,
  swcMinify: true,
})
```

**Usage**:
```bash
ANALYZE=true npm run build
```

**Expected Outcome**:
- Visual representation of bundle composition
- Identify largest dependencies
- Document baseline bundle sizes

---

#### 1.3 Optimize Font Loading
**Type**: Real Performance
**Impact**: Medium
**File**: `src/app/layout.tsx`

**Current Implementation**:
```typescript
const inter = Inter({ subsets: ["latin"] })
```

**Optimize**:
```typescript
const inter = Inter({
  subsets: ["latin"],
  display: 'swap', // Prevent invisible text during font load
  preload: true,
  variable: '--font-inter',
})
```

**Expected Outcome**:
- Faster text rendering (no FOIT - Flash of Invisible Text)
- Improved First Contentful Paint

---

#### 1.4 Gate Development Code
**Type**: Real Performance
**Impact**: Low-Medium

**Action**:
- Wrap all development-only code with environment checks
- Remove or disable React Query DevTools in production
- Remove debug routes from production build

**Example**:
```typescript
// React Query DevTools
const queryClient = new QueryClient({
  defaultOptions: {
    queries: {
      staleTime: 1000 * 60, // 60 seconds
      retry: 1,
      refetchOnWindowFocus: false,
    },
  },
})

// Only include DevTools in development
const DevTools =
  process.env.NODE_ENV === 'development'
    ? ReactQueryDevtools
    : () => null
```

---

#### 1.5 Production Build Analysis
**Type**: Baseline Measurement

**Action**:
1. Run production build: `npm run build`
2. Document build output sizes
3. Run Lighthouse on production build
4. Document current Web Vitals

**Metrics to Record**:
- Total bundle size
- Page-specific bundle sizes
- First Load JS shared by all
- Lighthouse scores (Performance, Accessibility, Best Practices, SEO)
- Core Web Vitals (LCP, FID, CLS, FCP, TTFB)

---

#### 1.6 Add Global Loading Indicator
**Type**: Perceived Performance
**Impact**: High
**Effort**: Low

**Installation**:
```bash
npm install nprogress
npm install --save-dev @types/nprogress
```

**Implementation** (`src/app/layout.tsx` or new `LoadingProvider`):
```typescript
'use client'
import { useEffect } from 'react'
import { usePathname, useSearchParams } from 'next/navigation'
import NProgress from 'nprogress'
import 'nprogress/nprogress.css'

export function NavigationProgress() {
  const pathname = usePathname()
  const searchParams = useSearchParams()

  useEffect(() => {
    NProgress.done()
  }, [pathname, searchParams])

  useEffect(() => {
    const handleStart = () => NProgress.start()
    const handleComplete = () => NProgress.done()

    // Listen to route changes
    window.addEventListener('beforeunload', handleStart)

    return () => {
      window.removeEventListener('beforeunload', handleStart)
    }
  }, [])

  return null
}
```

**Styling** (`globals.css`):
```css
#nprogress .bar {
  background: hsl(var(--primary)) !important;
  height: 3px;
}

#nprogress .peg {
  box-shadow: 0 0 10px hsl(var(--primary)), 0 0 5px hsl(var(--primary));
}
```

**Expected Outcome**:
- Users see immediate visual feedback during navigation
- Perceived performance improvement
- Reduced frustration during page transitions

---

#### 1.7 Add Skeleton Loaders
**Type**: Perceived Performance
**Impact**: High
**Files to Create/Modify**:
- `components/quiz/QuizListSkeleton.tsx` (new)
- `components/leaderboard/LeaderboardListSkeleton.tsx` (new)
- `components/leaderboard/LeaderboardSkeleton.tsx` (new)

**Example Implementation** (`QuizListSkeleton.tsx`):
```typescript
import { Skeleton } from '@/components/ui/skeleton'
import { Card, CardContent, CardHeader } from '@/components/ui/card'

export function QuizListSkeleton({ count = 6 }: { count?: number }) {
  return (
    <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
      {Array.from({ length: count }).map((_, i) => (
        <Card key={i}>
          <CardHeader>
            <Skeleton className="h-6 w-3/4" />
            <Skeleton className="h-4 w-1/2 mt-2" />
          </CardHeader>
          <CardContent>
            <Skeleton className="h-4 w-full mb-2" />
            <Skeleton className="h-4 w-full mb-2" />
            <Skeleton className="h-10 w-full mt-4" />
          </CardContent>
        </Card>
      ))}
    </div>
  )
}
```

**Usage in Components**:
```typescript
function QuizzesPage() {
  const { data: quizzes, isLoading } = useQuizzes()

  if (isLoading) {
    return <QuizListSkeleton count={6} />
  }

  return <QuizList quizzes={quizzes} />
}
```

**Expected Outcome**:
- Users see content structure immediately
- Reduced perceived loading time
- Better user experience during data fetching

---

### Phase 1 Success Metrics

- [ ] Zero console.logs in production build
- [ ] Bundle analyzer installed and baseline documented
- [ ] Font loading optimized (display: swap)
- [ ] Development code gated behind environment checks
- [ ] Production build metrics documented
- [ ] Global loading indicator visible during navigation
- [ ] Skeleton loaders implemented for 3+ list components
- [ ] Lighthouse score baseline established

**Baseline Metrics to Record**:
```
Performance Score: ___
FCP: ___ ms
LCP: ___ ms
TBT: ___ ms
CLS: ___
Total Bundle Size: ___ KB
First Load JS: ___ KB
```

---

## Phase 2: Code Splitting & Bundle Optimization

**Timeline**: Week 2-3
**Type**: Real Performance
**Risk Level**: Medium
**Effort**: Medium-High

### Objectives
- Reduce initial bundle size by 40-50%
- Implement dynamic imports for heavy components
- Enable route-based code splitting
- Optimize dependency tree

### Tasks

#### 2.1 Dynamic Imports for Modals and Dialogs
**Type**: Real Performance
**Impact**: High
**Effort**: Medium

**Components to Lazy Load** (10+ components):
- All Dialog/AlertDialog usage
- Sheet components (mobile menus)
- Popover components with heavy content
- Complex forms

**Implementation Pattern**:
```typescript
// Before
import { CreateQuizDialog } from '@/components/quiz/CreateQuizDialog'

function QuizzesPage() {
  return <CreateQuizDialog />
}

// After
import dynamic from 'next/dynamic'
import { Skeleton } from '@/components/ui/skeleton'

const CreateQuizDialog = dynamic(
  () => import('@/components/quiz/CreateQuizDialog').then(mod => ({
    default: mod.CreateQuizDialog
  })),
  {
    loading: () => <Skeleton className="h-96 w-full" />,
    ssr: false // Dialogs don't need SSR
  }
)

function QuizzesPage() {
  return <CreateQuizDialog />
}
```

**Files to Update**:
- Components using Dialog/AlertDialog
- Components using Sheet (mobile navigation)
- Components with complex forms (edit profile, settings)
- Achievement detail modals
- Discussion modals

**Expected Outcome**:
- Initial bundle reduced by 15-20%
- Faster initial page load
- Modals load on-demand

---

#### 2.2 Route-Based Code Splitting
**Type**: Real Performance
**Impact**: High
**Effort**: Medium

**Routes to Split**:
```typescript
// src/app/(dashboard)/dashboard/page.tsx
import dynamic from 'next/dynamic'

// Lazy load dashboard widgets
const QuickStats = dynamic(() => import('@/components/dashboard/QuickStats'))
const ActivityFeed = dynamic(() => import('@/components/dashboard/ActivityFeed'))
const FeaturedQuizzes = dynamic(() => import('@/components/dashboard/FeaturedQuizzes'))
const FriendActivity = dynamic(() => import('@/components/dashboard/FriendActivity'))
const AchievementShowcase = dynamic(() => import('@/components/dashboard/AchievementShowcase'))

// Load above-fold content immediately, below-fold lazily
const BelowFoldContent = dynamic(
  () => import('@/components/dashboard/BelowFoldContent'),
  { ssr: false }
)
```

**Routes to Optimize**:
- `/dashboard` - Split into above/below fold
- `/achievements` - Heavy component, lazy load
- `/discussions` - Lazy load discussion list
- `/settings` - Load settings panels on-demand
- `/leaderboard` - Lazy load leaderboard table
- `/profile/edit` - Lazy load form sections

**Expected Outcome**:
- Each route loads only necessary code
- Faster navigation between sections
- Reduced Time to Interactive

---

#### 2.3 Component-Level Lazy Loading
**Type**: Real Performance
**Impact**: Medium
**Effort**: Low-Medium

**Below-Fold Components to Defer**:
- Achievement cards (not visible on load)
- Activity feed (below fold on dashboard)
- Comments section (discussion pages)
- Related quizzes (quiz detail page)
- User statistics charts (profile page)

**Implementation with Intersection Observer**:
```typescript
'use client'
import { useEffect, useState, useRef } from 'react'
import dynamic from 'next/dynamic'

const HeavyComponent = dynamic(() => import('./HeavyComponent'))

export function LazyLoadWrapper({ children }) {
  const [isVisible, setIsVisible] = useState(false)
  const ref = useRef(null)

  useEffect(() => {
    const observer = new IntersectionObserver(
      ([entry]) => {
        if (entry.isIntersecting) {
          setIsVisible(true)
          observer.disconnect()
        }
      },
      { rootMargin: '100px' } // Load 100px before visible
    )

    if (ref.current) {
      observer.observe(ref.current)
    }

    return () => observer.disconnect()
  }, [])

  return <div ref={ref}>{isVisible ? children : <div className="h-64" />}</div>
}
```

**Expected Outcome**:
- Faster initial render
- Progressive content loading
- Improved perceived performance

---

#### 2.4 Dependency Audit & Optimization
**Type**: Real Performance
**Impact**: Medium
**Effort**: High

**Action Items**:

1. **Analyze Current Dependencies** (15+ Radix UI packages)
   ```bash
   npm run build -- --profile
   npx depcheck # Find unused dependencies
   ```

2. **Potential Optimizations**:
   - Ensure Radix UI packages are tree-shakeable
   - Check if any Radix components can be replaced with native HTML + CSS
   - Verify axios is properly tree-shaken (or consider fetch API)
   - Check if Supabase client bundles unused auth providers

3. **Remove Unused Dependencies**:
   ```bash
   # Run depcheck to find unused deps
   npx depcheck
   npm uninstall <unused-deps>
   ```

4. **Optimize Imports**:
   ```typescript
   // Before - imports entire library
   import * as RadixDialog from '@radix-ui/react-dialog'

   // After - imports only needed parts
   import { Dialog, DialogContent, DialogTrigger } from '@radix-ui/react-dialog'
   ```

**Expected Outcome**:
- 5-10% bundle size reduction
- Faster build times
- Cleaner dependency tree

---

#### 2.5 Optimize Radix UI Imports
**Type**: Real Performance
**Impact**: Low-Medium
**Effort**: Low

**Verify Tree-Shaking**:
- Check that shadcn/ui components use named imports
- Ensure no barrel imports that prevent tree-shaking
- Verify production build excludes unused Radix components

**Files to Audit**:
- `components/ui/*.tsx` (20+ shadcn components)

**Optimization**:
```typescript
// Ensure all imports are specific
import { Dialog, DialogContent } from '@radix-ui/react-dialog'
// NOT: import * as Dialog from '@radix-ui/react-dialog'
```

---

#### 2.6 Configure Manual Chunk Splitting
**Type**: Real Performance
**Impact**: Medium
**Effort**: Medium

**Add to `next.config.js`**:
```javascript
module.exports = {
  webpack: (config, { isServer }) => {
    if (!isServer) {
      config.optimization.splitChunks = {
        chunks: 'all',
        cacheGroups: {
          default: false,
          vendors: false,
          // Vendor chunk for react and react-dom
          react: {
            name: 'react',
            test: /[\\/]node_modules[\\/](react|react-dom|scheduler)[\\/]/,
            priority: 20,
          },
          // UI libraries chunk
          ui: {
            name: 'ui',
            test: /[\\/]node_modules[\\/](@radix-ui|@hookform|zod)[\\/]/,
            priority: 10,
          },
          // Supabase chunk
          supabase: {
            name: 'supabase',
            test: /[\\/]node_modules[\\/](@supabase)[\\/]/,
            priority: 10,
          },
          // Common chunk for shared code
          commons: {
            name: 'commons',
            minChunks: 2,
            priority: 5,
          },
        },
      }
    }
    return config
  },
}
```

**Expected Outcome**:
- Better caching (vendor code changes less frequently)
- Parallel chunk downloads
- Faster subsequent page loads

---

### Phase 2 Success Metrics

- [ ] Initial bundle size reduced by 40%+
- [ ] All modals/dialogs dynamically imported
- [ ] Route-based code splitting implemented
- [ ] Below-fold components lazy loaded
- [ ] Dependency audit completed
- [ ] Manual chunk splitting configured
- [ ] First Contentful Paint improved by 30%
- [ ] Time to Interactive improved by 25%

**Metrics to Record**:
```
Performance Score: ___
Bundle Size Reduction: ___% (from ___ KB to ___ KB)
FCP Improvement: ___% (from ___ ms to ___ ms)
TTI Improvement: ___% (from ___ ms to ___ ms)
Lighthouse Score: ___
```

---

## Phase 3: Server Components & Data Fetching

**Timeline**: Week 4-5
**Type**: Real Performance + Perceived Performance
**Risk Level**: Medium-High
**Effort**: High

### Objectives
- Convert 60-70% of client components to server components
- Reduce client bundle by additional 30%
- Implement parallel data fetching on server
- Enable streaming and suspense for faster perceived loads

### Tasks

#### 3.1 Client Component Audit
**Type**: Planning
**Impact**: High
**Effort**: Medium

**Action**:
1. Create spreadsheet of all 103 "use client" components
2. Categorize each component:
   - ✅ **Can be server component** (static UI, no interactivity)
   - ⚠️ **Partial conversion** (server shell, client interactive parts)
   - ❌ **Must remain client** (uses hooks, event handlers, state)

3. Prioritize conversion by impact:
   - High: List components, cards, layouts
   - Medium: Detail pages, profile pages
   - Low: Highly interactive components

**Expected Output**:
- Conversion plan spreadsheet
- Estimated bundle reduction per component
- Risk assessment

---

#### 3.2 Convert List & Card Components
**Type**: Real Performance
**Impact**: High
**Effort**: High

**Components to Convert** (Server Shell + Client Interactive Parts):

**Pattern**:
```typescript
// QuizCard.tsx - Server Component (default export)
import { QuizCardClient } from './QuizCardClient'
import type { Quiz } from '@/types'

interface QuizCardProps {
  quiz: Quiz
}

export function QuizCard({ quiz }: QuizCardProps) {
  // Server component - renders static parts
  return (
    <div className="card">
      <h3>{quiz.title}</h3>
      <p>{quiz.description}</p>
      {/* Client component for interactive parts */}
      <QuizCardClient quizId={quiz.id} initialFavoriteState={quiz.isFavorite} />
    </div>
  )
}

// QuizCardClient.tsx - Client Component
'use client'
import { useIsFavorite, useToggleFavorite } from '@/hooks'

interface QuizCardClientProps {
  quizId: string
  initialFavoriteState: boolean
}

export function QuizCardClient({ quizId, initialFavoriteState }: QuizCardClientProps) {
  const { data: isFavorite } = useIsFavorite(quizId, initialFavoriteState)
  const { mutate: toggleFavorite } = useToggleFavorite()

  return (
    <button onClick={() => toggleFavorite(quizId)}>
      {isFavorite ? '❤️' : '🤍'}
    </button>
  )
}
```

**Components to Split**:
- `QuizCard` → Server + Client parts
- `FriendCard` → Server + Client parts
- `DiscussionCard` → Server + Client parts
- `LeaderboardRow` → Server + Client parts
- `AchievementCard` → Server + Client parts
- `UserProfileCard` → Server + Client parts

**Expected Outcome**:
- 50-60% reduction in client JavaScript for list pages
- Faster initial render
- SEO improvements

---

#### 3.3 Convert Layout Components (Partial)
**Type**: Real Performance
**Impact**: Medium
**Effort**: Medium

**Layouts to Optimize**:

1. **Header** - Split navigation and user menu
   ```typescript
   // Header.tsx - Server Component
   export function Header() {
     return (
       <header>
         <Logo /> {/* Server component */}
         <Navigation /> {/* Server component */}
         <UserMenu /> {/* Client component (has dropdown) */}
       </header>
     )
   }
   ```

2. **Sidebar** - Keep navigation links server, client interactions separate
   ```typescript
   // Sidebar.tsx - Server Component
   export function Sidebar() {
     return (
       <aside>
         <SidebarLinks /> {/* Server */}
         <ThemeToggle /> {/* Client - has state */}
       </aside>
     )
   }
   ```

**Expected Outcome**:
- Reduced layout hydration cost
- Faster navigation rendering

---

#### 3.4 Parallel Data Fetching on Server
**Type**: Real Performance
**Impact**: High
**Effort**: Medium

**Current Problem**:
Dashboard makes sequential API calls, causing waterfall:
```
Load page → Fetch stats → Fetch quizzes → Fetch activity
```

**Solution - Parallel Fetching**:
```typescript
// app/(dashboard)/dashboard/page.tsx (Server Component)
import { getQuizzes } from '@/lib/api/quiz'
import { getUserStats } from '@/lib/api/user'
import { getFriendActivity } from '@/lib/api/friends'

export default async function DashboardPage() {
  // Parallel data fetching
  const [stats, quizzes, activity] = await Promise.all([
    getUserStats(),
    getQuizzes({ featured: true, limit: 6 }),
    getFriendActivity({ limit: 3 }),
  ])

  return (
    <div>
      <QuickStats data={stats} />
      <FeaturedQuizzes quizzes={quizzes} />
      <FriendActivity activity={activity} />
    </div>
  )
}
```

**Pages to Optimize**:
- `/dashboard` - Parallel fetch stats, quizzes, activity
- `/quizzes/[id]` - Parallel fetch quiz details and user attempt
- `/leaderboard` - Parallel fetch global and friend leaderboards
- `/profile/[userId]` - Parallel fetch profile, stats, achievements

**Expected Outcome**:
- 2-3x faster data loading
- Reduced Time to First Byte for page content
- Better perceived performance

---

#### 3.5 Implement Streaming & Suspense
**Type**: Perceived Performance
**Impact**: High
**Effort**: Medium

**Pattern - Stream Slow Data**:
```typescript
// app/(dashboard)/dashboard/page.tsx
import { Suspense } from 'react'
import { QuickStatsSkeleton } from '@/components/dashboard/QuickStatsSkeleton'

export default function DashboardPage() {
  return (
    <div>
      {/* Fast content renders immediately */}
      <h1>Dashboard</h1>

      {/* Slow content streams in */}
      <Suspense fallback={<QuickStatsSkeleton />}>
        <QuickStats /> {/* Server Component that fetches data */}
      </Suspense>

      <Suspense fallback={<FeaturedQuizzesSkeleton />}>
        <FeaturedQuizzes />
      </Suspense>

      <Suspense fallback={<ActivityFeedSkeleton />}>
        <ActivityFeed />
      </Suspense>
    </div>
  )
}

// components/dashboard/QuickStats.tsx (Server Component)
async function QuickStats() {
  const stats = await getUserStats() // This runs on server
  return <QuickStatsDisplay stats={stats} />
}
```

**Pages to Add Streaming**:
- Dashboard (multiple data sources)
- Leaderboard (slow query)
- Profile pages (stats, achievements)
- Quiz results page (analysis data)

**Expected Outcome**:
- Users see content progressively
- Page appears to load 2-3x faster
- Reduced perceived loading time
- Better UX during slow API responses

---

#### 3.6 Prefetch Critical Routes
**Type**: Perceived Performance
**Impact**: Medium
**Effort**: Low

**Implementation**:
```typescript
// components/layout/Sidebar.tsx
import Link from 'next/link'

export function Sidebar() {
  return (
    <nav>
      {/* prefetch defaults to true for Link */}
      <Link href="/dashboard" prefetch={true}>Dashboard</Link>
      <Link href="/quizzes" prefetch={true}>Quizzes</Link>
      <Link href="/friends" prefetch={true}>Friends</Link>

      {/* Disable prefetch for less common routes */}
      <Link href="/settings" prefetch={false}>Settings</Link>
    </nav>
  )
}
```

**Routes to Prefetch**:
- Dashboard (from any page)
- Quizzes list (high traffic)
- Active quiz (from dashboard)
- Leaderboard (common navigation)

**Routes NOT to Prefetch**:
- Settings (infrequent)
- Edit profile (infrequent)
- Admin pages

**Expected Outcome**:
- Instant navigation to common routes
- Improved perceived performance
- Better user experience

---

#### 3.7 Server-Side Data Prefetching
**Type**: Real Performance
**Impact**: Medium
**Effort**: Medium

**Use React Query on Server**:
```typescript
// app/(dashboard)/dashboard/page.tsx
import { HydrationBoundary, QueryClient, dehydrate } from '@tanstack/react-query'
import { getQuizzes } from '@/lib/api/quiz'

export default async function DashboardPage() {
  const queryClient = new QueryClient()

  // Prefetch data on server
  await queryClient.prefetchQuery({
    queryKey: ['quizzes', { featured: true }],
    queryFn: () => getQuizzes({ featured: true }),
  })

  return (
    <HydrationBoundary state={dehydrate(queryClient)}>
      <DashboardClient /> {/* Client component has data immediately */}
    </HydrationBoundary>
  )
}
```

**Expected Outcome**:
- Data available immediately on client
- No loading states on initial render
- Faster perceived performance

---

### Phase 3 Success Metrics

- [ ] 60%+ of components converted to server components
- [ ] "use client" directives reduced from 103 to <40
- [ ] Parallel data fetching implemented on 5+ pages
- [ ] Streaming & Suspense added to dashboard and 3+ pages
- [ ] Critical routes prefetched
- [ ] Client bundle reduced by additional 30%
- [ ] Dashboard loads 2x faster

**Metrics to Record**:
```
Client Components: ___ (from 103)
Server Components: ___
Client Bundle Reduction: ___% (additional)
Dashboard Load Time: ___ ms (from ___ ms)
FCP: ___ ms
LCP: ___ ms
Lighthouse Score: ___
```

---

## Phase 4: Component Optimization & Memoization

**Timeline**: Week 6
**Type**: Real Performance
**Risk Level**: Low-Medium
**Effort**: Medium

### Objectives
- Eliminate unnecessary component re-renders
- Optimize list rendering performance
- Reduce memory usage
- Achieve 60fps scrolling in all lists

### Tasks

#### 4.1 Memoize List Item Components
**Type**: Real Performance
**Impact**: High
**Effort**: Low-Medium

**Components to Memoize**:

```typescript
// Before: QuizCard.tsx
export function QuizCard({ quiz, onFavorite }: QuizCardProps) {
  // Component re-renders when parent re-renders
  return <div>...</div>
}

// After: QuizCard.tsx
import { memo } from 'react'

export const QuizCard = memo(function QuizCard({ quiz, onFavorite }: QuizCardProps) {
  return <div>...</div>
}, (prevProps, nextProps) => {
  // Custom comparison - only re-render if quiz data changed
  return prevProps.quiz.id === nextProps.quiz.id &&
         prevProps.quiz.isFavorite === nextProps.quiz.isFavorite
})
```

**Components to Memoize**:
- `QuizCard` (used in lists of 10-50 items)
- `FriendCard` (used in lists)
- `LeaderboardRow` (used in lists of 100+ items)
- `DiscussionCard` (used in lists)
- `AchievementCard` (used in lists)
- `NotificationItem` (used in lists)
- `FriendCard` (used in lists)

**Expected Outcome**:
- 70-80% reduction in re-renders
- Smoother scrolling
- Lower CPU usage

---

#### 4.2 Memoize Callback Functions
**Type**: Real Performance
**Impact**: Medium
**Effort**: Low

**Pattern**:
```typescript
// Before: QuizList.tsx
function QuizList({ quizzes }: QuizListProps) {
  const handleFavorite = (quizId: string) => {
    // This function is recreated on every render
    toggleFavorite(quizId)
  }

  return quizzes.map(quiz => (
    <QuizCard key={quiz.id} quiz={quiz} onFavorite={handleFavorite} />
  ))
}

// After: QuizList.tsx
import { useCallback } from 'react'

function QuizList({ quizzes }: QuizListProps) {
  const handleFavorite = useCallback((quizId: string) => {
    toggleFavorite(quizId)
  }, []) // Empty deps - function doesn't change

  return quizzes.map(quiz => (
    <QuizCard key={quiz.id} quiz={quiz} onFavorite={handleFavorite} />
  ))
}
```

**Files to Update**:
- All list components passing callbacks to children
- Event handlers in forms
- Components with inline arrow functions as props

**Expected Outcome**:
- Prevents unnecessary child re-renders (when combined with memo)
- More stable references
- Better performance with memoized components

---

#### 4.3 Memoize Computed Values
**Type**: Real Performance
**Impact**: Medium
**Effort**: Medium

**Pattern**:
```typescript
// Before: Dashboard.tsx
function Dashboard() {
  const { data: stats } = useUserStats()

  // This is recalculated on every render
  const progressPercentage = (stats.completedQuizzes / stats.totalQuizzes) * 100
  const averageScore = stats.scores.reduce((a, b) => a + b, 0) / stats.scores.length

  return <div>{progressPercentage}% - Avg: {averageScore}</div>
}

// After: Dashboard.tsx
import { useMemo } from 'react'

function Dashboard() {
  const { data: stats } = useUserStats()

  const progressPercentage = useMemo(() => {
    if (!stats) return 0
    return (stats.completedQuizzes / stats.totalQuizzes) * 100
  }, [stats?.completedQuizzes, stats?.totalQuizzes])

  const averageScore = useMemo(() => {
    if (!stats?.scores.length) return 0
    return stats.scores.reduce((a, b) => a + b, 0) / stats.scores.length
  }, [stats?.scores])

  return <div>{progressPercentage}% - Avg: {averageScore}</div>
}
```

**Computations to Memoize**:
- Filtered/sorted lists (quiz filters, leaderboard sorting)
- Statistical calculations (averages, percentages)
- Derived state (form validation, complex conditions)
- Heavy transformations (data mapping, grouping)

**Files to Audit**:
- Dashboard components (stats calculations)
- Leaderboard (ranking calculations)
- Quiz results (score calculations)
- Profile stats (achievement progress)

**Expected Outcome**:
- Reduced CPU usage
- Faster renders
- Better performance on low-end devices

---

#### 4.4 Implement Virtual Scrolling for Long Lists
**Type**: Real Performance
**Impact**: High (for long lists)
**Effort**: Medium

**Installation**:
```bash
npm install @tanstack/react-virtual
```

**Implementation** (Leaderboard with 100+ rows):
```typescript
// Before: LeaderboardTable.tsx
function LeaderboardTable({ entries }: { entries: LeaderboardEntry[] }) {
  return (
    <div className="overflow-auto h-screen">
      {entries.map(entry => (
        <LeaderboardRow key={entry.userId} entry={entry} />
      ))}
    </div>
  )
}

// After: LeaderboardTable.tsx
import { useVirtualizer } from '@tanstack/react-virtual'
import { useRef } from 'react'

function LeaderboardTable({ entries }: { entries: LeaderboardEntry[] }) {
  const parentRef = useRef<HTMLDivElement>(null)

  const virtualizer = useVirtualizer({
    count: entries.length,
    getScrollElement: () => parentRef.current,
    estimateSize: () => 60, // Estimated row height
    overscan: 5, // Render 5 extra items outside viewport
  })

  return (
    <div ref={parentRef} className="overflow-auto h-screen">
      <div
        style={{
          height: `${virtualizer.getTotalSize()}px`,
          width: '100%',
          position: 'relative',
        }}
      >
        {virtualizer.getVirtualItems().map(virtualRow => {
          const entry = entries[virtualRow.index]
          return (
            <div
              key={entry.userId}
              style={{
                position: 'absolute',
                top: 0,
                left: 0,
                width: '100%',
                height: `${virtualRow.size}px`,
                transform: `translateY(${virtualRow.start}px)`,
              }}
            >
              <LeaderboardRow entry={entry} />
            </div>
          )
        })}
      </div>
    </div>
  )
}
```

**Lists to Virtualize** (if they exceed 20 items):
- Leaderboard (100+ items)
- Notifications list (50+ items)
- Quiz search results (100+ items)
- Discussion list (if long)

**Expected Outcome**:
- Smooth 60fps scrolling even with 1000+ items
- Constant memory usage regardless of list length
- Dramatically improved performance for long lists

---

#### 4.5 Optimize React Query Configuration
**Type**: Real Performance
**Impact**: Medium
**Effort**: Low

**Current Configuration** (`src/app/layout.tsx`):
```typescript
const queryClient = new QueryClient({
  defaultOptions: {
    queries: {
      staleTime: 1000 * 60, // 60 seconds
      retry: 1,
      refetchOnWindowFocus: false,
    },
  },
})
```

**Optimized Configuration**:
```typescript
const queryClient = new QueryClient({
  defaultOptions: {
    queries: {
      staleTime: 1000 * 60 * 5, // 5 minutes (reduce refetching)
      cacheTime: 1000 * 60 * 10, // 10 minutes (keep in cache longer)
      retry: 1,
      refetchOnWindowFocus: false,
      refetchOnReconnect: false, // Don't refetch on reconnect
      refetchOnMount: false, // Don't refetch if data is fresh
    },
  },
})
```

**Per-Query Optimizations**:
```typescript
// User stats - update frequently
export function useUserStats() {
  return useQuery({
    queryKey: ['user-stats'],
    queryFn: getUserStats,
    staleTime: 1000 * 30, // 30 seconds (fresh)
  })
}

// Quiz list - update less frequently
export function useQuizzes(filters?: QuizFilters) {
  return useQuery({
    queryKey: ['quizzes', filters],
    queryFn: () => getQuizzes(filters),
    staleTime: 1000 * 60 * 5, // 5 minutes
  })
}

// Leaderboard - can be stale
export function useLeaderboard() {
  return useQuery({
    queryKey: ['leaderboard'],
    queryFn: getLeaderboard,
    staleTime: 1000 * 60 * 10, // 10 minutes
    cacheTime: 1000 * 60 * 30, // Cache for 30 minutes
  })
}
```

**Implement Query Prefetching**:
```typescript
// Prefetch quiz details on hover
function QuizCard({ quiz }: QuizCardProps) {
  const queryClient = useQueryClient()

  const handleMouseEnter = () => {
    queryClient.prefetchQuery({
      queryKey: ['quiz', quiz.id],
      queryFn: () => getQuizDetails(quiz.id),
    })
  }

  return <div onMouseEnter={handleMouseEnter}>...</div>
}
```

**Expected Outcome**:
- Reduced API calls
- Faster perceived performance (data from cache)
- Lower server load

---

#### 4.6 Audit and Fix Re-render Issues
**Type**: Real Performance
**Impact**: Medium
**Effort**: Medium

**Install React DevTools Profiler** and audit:
1. Dashboard page - identify components re-rendering unnecessarily
2. Quiz list page - check list re-renders
3. Quiz taking page - ensure only active question re-renders

**Common Issues to Fix**:
```typescript
// Issue 1: Inline object creation
// Before
<Component style={{ color: 'red' }} /> // New object every render

// After
const style = { color: 'red' } // Stable reference
<Component style={style} />

// Issue 2: Inline array/object in deps
// Before
useEffect(() => {
  doSomething()
}, [{ id: quiz.id }]) // New object every render

// After
useEffect(() => {
  doSomething()
}, [quiz.id]) // Primitive value

// Issue 3: Anonymous functions
// Before
<button onClick={() => handleClick(id)}>Click</button>

// After
const handleButtonClick = useCallback(() => handleClick(id), [id])
<button onClick={handleButtonClick}>Click</button>
```

**Expected Outcome**:
- Identify all unnecessary re-renders
- Fix React anti-patterns
- Improve overall performance

---

### Phase 4 Success Metrics

- [ ] All list item components memoized
- [ ] All callback functions in lists memoized
- [ ] Computed values in 10+ components memoized
- [ ] Virtual scrolling implemented for 2+ long lists
- [ ] React Query configuration optimized
- [ ] Query prefetching implemented for quiz cards
- [ ] Re-render audit completed and issues fixed
- [ ] 60fps scrolling in all lists
- [ ] Re-render count reduced by 70%

**Metrics to Record**:
```
Re-render Count (Dashboard): ___ (before) → ___ (after)
List Scroll Performance: ___ fps
Memory Usage: ___ MB (before) → ___ MB (after)
Largest Component Re-renders: ___
React Query Cache Hit Rate: ___%
```

---

## Phase 5: Perceived Performance & UX Enhancements

**Timeline**: Week 7-8
**Type**: Perceived Performance
**Risk Level**: Low
**Effort**: Medium

### Objectives
- Make all user interactions feel instant (<100ms)
- Eliminate blank screen moments
- Implement optimistic UI updates
- Add smooth loading transitions

### Tasks

#### 5.1 Optimistic UI Updates
**Type**: Perceived Performance
**Impact**: High
**Effort**: Medium

**Pattern - Favorite/Unfavorite**:
```typescript
// Before: useToggleFavorite.ts
export function useToggleFavorite() {
  const queryClient = useQueryClient()

  return useMutation({
    mutationFn: (quizId: string) => toggleFavoriteApi(quizId),
    onSuccess: () => {
      // Wait for server response before updating UI
      queryClient.invalidateQueries({ queryKey: ['quizzes'] })
    },
  })
}

// After: useToggleFavorite.ts (Optimistic)
export function useToggleFavorite() {
  const queryClient = useQueryClient()

  return useMutation({
    mutationFn: (quizId: string) => toggleFavoriteApi(quizId),

    // Update UI immediately (before server responds)
    onMutate: async (quizId) => {
      // Cancel outgoing refetches
      await queryClient.cancelQueries({ queryKey: ['quizzes'] })

      // Snapshot previous value
      const previousQuizzes = queryClient.getQueryData(['quizzes'])

      // Optimistically update
      queryClient.setQueryData(['quizzes'], (old: Quiz[]) => {
        return old.map(quiz =>
          quiz.id === quizId
            ? { ...quiz, isFavorite: !quiz.isFavorite }
            : quiz
        )
      })

      return { previousQuizzes }
    },

    // Rollback on error
    onError: (err, quizId, context) => {
      queryClient.setQueryData(['quizzes'], context.previousQuizzes)
    },

    // Always refetch after error or success
    onSettled: () => {
      queryClient.invalidateQueries({ queryKey: ['quizzes'] })
    },
  })
}
```

**Interactions to Make Optimistic**:
- ✅ Toggle favorite (quiz)
- ✅ Like/unlike discussion
- ✅ Send friend request
- ✅ Accept/decline friend request
- ✅ Mark notification as read
- ✅ Select quiz answer (during quiz)

**Expected Outcome**:
- Instant visual feedback
- Actions feel immediate
- Better UX even on slow connections

---

#### 5.2 Enhanced Loading States
**Type**: Perceived Performance
**Impact**: High
**Effort**: Medium

**Implement Progressive Loading**:
```typescript
// Before: QuizzesPage.tsx
function QuizzesPage() {
  const { data: quizzes, isLoading } = useQuizzes()

  if (isLoading) {
    return <Spinner /> // Blank screen until all data loads
  }

  return <QuizList quizzes={quizzes} />
}

// After: QuizzesPage.tsx (Progressive)
function QuizzesPage() {
  const { data: quizzes, isLoading, isFetching } = useQuizzes()

  return (
    <div>
      {/* Show skeleton on initial load */}
      {isLoading && <QuizListSkeleton count={6} />}

      {/* Show content as soon as available */}
      {quizzes && (
        <div className={isFetching ? 'opacity-50' : ''}>
          <QuizList quizzes={quizzes} />
        </div>
      )}

      {/* Subtle loading indicator during background refetch */}
      {isFetching && !isLoading && (
        <div className="fixed top-4 right-4">
          <Spinner size="sm" />
        </div>
      )}
    </div>
  )
}
```

**Enhanced Skeleton Patterns**:
```typescript
// Shimmer effect for skeletons
// globals.css
@keyframes shimmer {
  0% {
    background-position: -1000px 0;
  }
  100% {
    background-position: 1000px 0;
  }
}

.skeleton {
  animation: shimmer 2s infinite linear;
  background: linear-gradient(
    to right,
    #f0f0f0 4%,
    #e8e8e8 25%,
    #f0f0f0 36%
  );
  background-size: 1000px 100%;
}
```

**Pages to Enhance**:
- Dashboard (progressive widget loading)
- Quiz list (skeleton → content)
- Leaderboard (skeleton → progressive rows)
- Profile (skeleton → content)

**Expected Outcome**:
- Zero blank screens
- Content appears progressively
- Users always see something

---

#### 5.3 Prefetching Strategies
**Type**: Perceived Performance
**Impact**: High
**Effort**: Medium

**1. Hover Prefetch**:
```typescript
// QuizCard.tsx
function QuizCard({ quiz }: QuizCardProps) {
  const queryClient = useQueryClient()
  const router = useRouter()
  const [prefetched, setPrefetched] = useState(false)

  const handleMouseEnter = () => {
    if (!prefetched) {
      // Prefetch quiz details
      queryClient.prefetchQuery({
        queryKey: ['quiz', quiz.id],
        queryFn: () => getQuizDetails(quiz.id),
      })

      // Prefetch Next.js route
      router.prefetch(`/quizzes/${quiz.id}`)

      setPrefetched(true)
    }
  }

  return (
    <div onMouseEnter={handleMouseEnter}>
      <Link href={`/quizzes/${quiz.id}`}>
        {quiz.title}
      </Link>
    </div>
  )
}
```

**2. Predictive Prefetch (Quiz Taking)**:
```typescript
// QuizTakingPage.tsx
function QuizTakingPage({ quiz, currentQuestionIndex }: Props) {
  const queryClient = useQueryClient()

  // Prefetch next question data when user is on current question
  useEffect(() => {
    const nextQuestionIndex = currentQuestionIndex + 1
    if (nextQuestionIndex < quiz.questions.length) {
      const nextQuestion = quiz.questions[nextQuestionIndex]
      // Preload any data needed for next question
    }

    // If on last question, prefetch results page data
    if (currentQuestionIndex === quiz.questions.length - 1) {
      queryClient.prefetchQuery({
        queryKey: ['quiz-results', attemptId],
        queryFn: () => getQuizResults(attemptId),
      })
    }
  }, [currentQuestionIndex])

  return <QuestionDisplay question={quiz.questions[currentQuestionIndex]} />
}
```

**3. Route-Based Prefetch**:
```typescript
// Dashboard.tsx - Prefetch likely next routes
function Dashboard() {
  useEffect(() => {
    const router = useRouter()

    // User is likely to go to quizzes next
    router.prefetch('/quizzes')

    // Also prefetch friends page
    router.prefetch('/friends')
  }, [])

  return <DashboardContent />
}
```

**Expected Outcome**:
- Near-instant navigation on hover
- Quiz questions load instantly
- Results appear immediately after quiz completion

---

#### 5.4 Smart Caching Strategies
**Type**: Perceived Performance
**Impact**: Medium
**Effort**: Low-Medium

**Stale-While-Revalidate Pattern**:
```typescript
// Leaderboard - show cached data, update in background
export function useLeaderboard() {
  return useQuery({
    queryKey: ['leaderboard'],
    queryFn: getLeaderboard,
    staleTime: 1000 * 60, // Consider fresh for 1 minute
    cacheTime: 1000 * 60 * 30, // Keep in cache for 30 minutes

    // Show stale data while fetching fresh data
    keepPreviousData: true,
  })
}

// Usage in component
function Leaderboard() {
  const { data, isLoading, isFetching } = useLeaderboard()

  return (
    <div>
      {data && <LeaderboardTable data={data} />}

      {/* Subtle indicator during background refresh */}
      {isFetching && !isLoading && (
        <div className="text-xs text-muted">Updating...</div>
      )}
    </div>
  )
}
```

**Background Refresh for Notifications**:
```typescript
// Automatically refresh notifications every 30 seconds
export function useNotifications() {
  return useQuery({
    queryKey: ['notifications'],
    queryFn: getNotifications,
    refetchInterval: 30000, // 30 seconds
    refetchIntervalInBackground: true, // Even when tab not focused
  })
}
```

**Expected Outcome**:
- Users always see data (even if slightly stale)
- Background updates don't interrupt UX
- Fresh data without loading states

---

#### 5.5 Micro-interactions & Feedback
**Type**: Perceived Performance
**Impact**: Medium
**Effort**: Low-Medium

**Button Loading States**:
```typescript
// Before
<Button onClick={handleSubmit}>Submit</Button>

// After
<Button onClick={handleSubmit} disabled={isSubmitting}>
  {isSubmitting ? (
    <>
      <Spinner className="mr-2" size="sm" />
      Submitting...
    </>
  ) : (
    'Submit'
  )}
</Button>
```

**Success Animations**:
```typescript
// Toast notifications for actions
import { toast } from 'sonner'

function useToggleFavorite() {
  return useMutation({
    mutationFn: toggleFavoriteApi,
    onSuccess: (data) => {
      toast.success(
        data.isFavorite ? 'Added to favorites!' : 'Removed from favorites',
        { duration: 2000 }
      )
    },
  })
}

// Confetti for achievements
import confetti from 'canvas-confetti'

function AchievementUnlocked({ achievement }: Props) {
  useEffect(() => {
    confetti({
      particleCount: 100,
      spread: 70,
      origin: { y: 0.6 }
    })
  }, [])

  return <div>Achievement Unlocked: {achievement.name}</div>
}
```

**Smooth Transitions**:
```typescript
// Add transitions to state changes
// globals.css
.fade-in {
  animation: fadeIn 0.3s ease-in;
}

@keyframes fadeIn {
  from { opacity: 0; transform: translateY(10px); }
  to { opacity: 1; transform: translateY(0); }
}

// Component
function QuizList({ quizzes }: Props) {
  return (
    <div className="grid gap-4">
      {quizzes.map((quiz, index) => (
        <div
          key={quiz.id}
          className="fade-in"
          style={{ animationDelay: `${index * 0.05}s` }}
        >
          <QuizCard quiz={quiz} />
        </div>
      ))}
    </div>
  )
}
```

**Expected Outcome**:
- Users get immediate feedback for actions
- Actions feel satisfying and complete
- Better perceived responsiveness

---

#### 5.6 Smooth Page Transitions
**Type**: Perceived Performance
**Impact**: Medium
**Effort**: Low

**Framer Motion for Page Transitions** (optional - adds dependency):
```bash
npm install framer-motion
```

**Alternative - CSS-only transitions**:
```typescript
// app/layout.tsx
export default function RootLayout({ children }: Props) {
  return (
    <html>
      <body>
        <div className="page-transition">
          {children}
        </div>
      </body>
    </html>
  )
}

// globals.css
.page-transition {
  animation: pageEnter 0.2s ease-out;
}

@keyframes pageEnter {
  from {
    opacity: 0;
    transform: translateY(8px);
  }
  to {
    opacity: 1;
    transform: translateY(0);
  }
}
```

**Expected Outcome**:
- Smooth transitions between pages
- Professional feel
- Reduced jarring page changes

---

### Phase 5 Success Metrics

- [ ] Optimistic updates implemented for 6+ interactions
- [ ] Enhanced loading states on all data-heavy pages
- [ ] Progressive loading implemented
- [ ] Hover prefetch on quiz cards
- [ ] Predictive prefetch in quiz taking flow
- [ ] Stale-while-revalidate for leaderboard
- [ ] Button loading states on all forms
- [ ] Success feedback for all actions
- [ ] Smooth page transitions

**Metrics to Record**:
```
Perceived Response Time: < 100ms for ___ interactions
Loading State Coverage: ___%
Prefetch Success Rate: ___%
User Satisfaction Score: ___
```

---

## Phase 6: Quiz Taking Experience

**Timeline**: Week 9
**Type**: Real Performance + Perceived Performance
**Risk Level**: Medium
**Effort**: Medium

### Objectives
- Optimize core quiz-taking flow
- Achieve <50ms question transitions
- Eliminate lag during answer selection
- Instant results display

### Tasks

#### 6.1 Preload Quiz Data
**Type**: Real Performance
**Impact**: High
**Effort**: Low

**Current Implementation**: Questions loaded one at a time
**Optimized**: Load all questions upfront

```typescript
// Before: QuizTakingPage.tsx
function QuizTakingPage({ quizId }: Props) {
  const { data: currentQuestion } = useQuestion(quizId, currentIndex)
  // Loads one question at a time - causes delays
}

// After: QuizTakingPage.tsx
function QuizTakingPage({ quizId }: Props) {
  // Load entire quiz data upfront
  const { data: quiz } = useQuiz(quizId, {
    includeQuestions: true, // Fetch all questions at once
  })

  // No network requests during quiz - instant transitions
  const currentQuestion = quiz?.questions[currentIndex]
}
```

**API Change**:
```typescript
// lib/api/quiz.ts
export async function getQuizWithQuestions(quizId: string): Promise<QuizWithQuestions> {
  const response = await apiClient.get<QuizWithQuestions>(
    `${API_ENDPOINTS.QUIZZES.DETAIL(quizId)}?include_questions=true`
  )
  return response
}

// Hook
export function useQuiz(quizId: string) {
  return useQuery({
    queryKey: ['quiz', quizId, 'full'],
    queryFn: () => getQuizWithQuestions(quizId),
    staleTime: Infinity, // Quiz data doesn't change during session
  })
}
```

**Expected Outcome**:
- First question shows immediately
- All subsequent questions load instantly (no network requests)
- Better offline experience

---

#### 6.2 Optimize Question Transitions
**Type**: Real Performance + Perceived Performance
**Impact**: High
**Effort**: Medium

**Eliminate Layout Shifts**:
```typescript
// QuizTakingPage.tsx
function QuizTakingPage({ quiz }: Props) {
  const [currentIndex, setCurrentIndex] = useState(0)
  const currentQuestion = quiz.questions[currentIndex]

  // Preload next question's images (if any)
  useEffect(() => {
    if (currentIndex + 1 < quiz.questions.length) {
      const nextQuestion = quiz.questions[currentIndex + 1]
      // Preload image if exists
      if (nextQuestion.imageUrl) {
        const img = new Image()
        img.src = nextQuestion.imageUrl
      }
    }
  }, [currentIndex, quiz.questions])

  return (
    <div className="quiz-container">
      {/* Fixed height to prevent layout shift */}
      <div className="min-h-[400px]">
        <QuestionDisplay
          question={currentQuestion}
          onAnswer={handleAnswer}
        />
      </div>
    </div>
  )
}
```

**Smooth Transitions**:
```typescript
// QuestionDisplay.tsx
import { motion, AnimatePresence } from 'framer-motion'

function QuestionDisplay({ question, onAnswer }: Props) {
  return (
    <AnimatePresence mode="wait">
      <motion.div
        key={question.id}
        initial={{ opacity: 0, x: 20 }}
        animate={{ opacity: 1, x: 0 }}
        exit={{ opacity: 0, x: -20 }}
        transition={{ duration: 0.15 }} // Fast transition
      >
        <h3>{question.text}</h3>
        <AnswerOptions options={question.options} onSelect={onAnswer} />
      </motion.div>
    </AnimatePresence>
  )
}
```

**Expected Outcome**:
- Question transitions <50ms
- No layout shifts
- Smooth, professional animations

---

#### 6.3 Instant Answer Feedback
**Type**: Perceived Performance
**Impact**: High
**Effort**: Low

```typescript
// AnswerOption.tsx
function AnswerOption({ option, onSelect }: Props) {
  const [isSelected, setIsSelected] = useState(false)

  const handleClick = () => {
    // Immediate visual feedback
    setIsSelected(true)

    // Delay actual submission slightly for animation
    setTimeout(() => {
      onSelect(option.id)
    }, 150)
  }

  return (
    <button
      onClick={handleClick}
      className={cn(
        "answer-option",
        isSelected && "ring-2 ring-primary scale-105"
      )}
    >
      {option.text}
      {isSelected && <CheckIcon className="ml-2" />}
    </button>
  )
}
```

**Expected Outcome**:
- Answer selection feels instant
- Clear visual feedback
- Better UX

---

#### 6.4 Background Auto-Save Optimization
**Type**: Real Performance
**Impact**: Medium
**Effort**: Low

**Current**: 30-second interval auto-save (from codebase analysis)

**Optimize**:
```typescript
// hooks/useQuizAutoSave.ts
export function useQuizAutoSave(attemptId: string, answers: Answers) {
  const [isSaving, setIsSaving] = useState(false)
  const debouncedAnswers = useDebounce(answers, 2000) // Debounce 2 seconds

  useEffect(() => {
    // Only save if answers changed
    if (debouncedAnswers) {
      saveQuizProgress(attemptId, debouncedAnswers)
        .then(() => setIsSaving(false))
        .catch(err => logger.error('Auto-save failed:', err))
    }
  }, [debouncedAnswers, attemptId])

  return { isSaving }
}

// Usage
function QuizTakingPage() {
  const [answers, setAnswers] = useState({})
  const { isSaving } = useQuizAutoSave(attemptId, answers)

  return (
    <div>
      {isSaving && <SaveIndicator />}
      <QuestionDisplay />
    </div>
  )
}
```

**Expected Outcome**:
- Non-blocking saves
- Saves only when needed (debounced)
- Clear save status indicator

---

#### 6.5 Timer Optimization (Web Worker)
**Type**: Real Performance
**Impact**: Medium
**Effort**: Medium-High

**Problem**: Timer running on main thread can cause jank

**Solution**: Move timer to Web Worker

```typescript
// workers/quiz-timer.worker.ts
let timeRemaining = 0
let isPaused = false
let interval: NodeJS.Timeout | null = null

self.addEventListener('message', (e) => {
  const { type, payload } = e.data

  switch (type) {
    case 'START':
      timeRemaining = payload.duration
      interval = setInterval(() => {
        if (!isPaused && timeRemaining > 0) {
          timeRemaining--
          self.postMessage({ type: 'TICK', payload: { timeRemaining } })
        }
        if (timeRemaining === 0) {
          self.postMessage({ type: 'TIMEOUT' })
          if (interval) clearInterval(interval)
        }
      }, 1000)
      break

    case 'PAUSE':
      isPaused = true
      break

    case 'RESUME':
      isPaused = false
      break

    case 'STOP':
      if (interval) clearInterval(interval)
      break
  }
})

// hooks/useQuizTimer.ts
export function useQuizTimer(duration: number) {
  const [timeRemaining, setTimeRemaining] = useState(duration)
  const workerRef = useRef<Worker>()

  useEffect(() => {
    workerRef.current = new Worker(
      new URL('../workers/quiz-timer.worker.ts', import.meta.url)
    )

    workerRef.current.addEventListener('message', (e) => {
      if (e.data.type === 'TICK') {
        setTimeRemaining(e.data.payload.timeRemaining)
      } else if (e.data.type === 'TIMEOUT') {
        // Handle timeout
      }
    })

    workerRef.current.postMessage({ type: 'START', payload: { duration } })

    return () => {
      workerRef.current?.postMessage({ type: 'STOP' })
      workerRef.current?.terminate()
    }
  }, [duration])

  return { timeRemaining }
}
```

**Expected Outcome**:
- Timer doesn't block UI
- Smoother quiz experience
- Better performance on lower-end devices

---

#### 6.6 Preload Results Page
**Type**: Perceived Performance
**Impact**: High
**Effort**: Low

```typescript
// QuizTakingPage.tsx
function QuizTakingPage({ quiz, attemptId }: Props) {
  const queryClient = useQueryClient()
  const router = useRouter()
  const isLastQuestion = currentIndex === quiz.questions.length - 1

  // When user is on last question, prefetch results
  useEffect(() => {
    if (isLastQuestion) {
      // Prefetch results data
      queryClient.prefetchQuery({
        queryKey: ['quiz-results', attemptId],
        queryFn: () => getQuizResults(attemptId),
      })

      // Prefetch Next.js route
      router.prefetch(`/quizzes/${quiz.id}/results/${attemptId}`)
    }
  }, [isLastQuestion])

  const handleSubmitQuiz = async () => {
    await submitQuiz(attemptId)

    // Instant navigation (data already loaded)
    router.push(`/quizzes/${quiz.id}/results/${attemptId}`)
  }

  return <QuestionDisplay />
}
```

**Expected Outcome**:
- Results appear instantly after last answer
- No waiting after quiz completion
- Better completion experience

---

### Phase 6 Success Metrics

- [ ] All quiz questions preloaded on quiz start
- [ ] Question transitions <50ms
- [ ] Answer selection feels instant (<100ms)
- [ ] Auto-save optimized with debouncing
- [ ] Timer moved to Web Worker
- [ ] Results page preloaded on last question
- [ ] Zero lag during quiz taking

**Metrics to Record**:
```
Quiz Load Time: ___ ms
Question Transition Time: ___ ms
Answer Selection Response: ___ ms
Results Display Time: ___ ms (after submission)
User Satisfaction (Quiz Experience): ___
```

---

## Phase 7: Image & Asset Optimization

**Timeline**: Week 10
**Type**: Real Performance
**Risk Level**: Low
**Effort**: Medium

### Objectives
- Migrate to Next.js Image component
- Reduce image sizes by 60%+
- Implement lazy loading for all images
- Improve Largest Contentful Paint (LCP)

### Tasks

#### 7.1 Migrate to next/image
**Type**: Real Performance
**Impact**: High
**Effort**: Medium

**Pattern**:
```typescript
// Before
<img src="/images/quiz-banner.jpg" alt="Quiz" />

// After
import Image from 'next/image'

<Image
  src="/images/quiz-banner.jpg"
  alt="Quiz"
  width={800}
  height={400}
  quality={85} // Good quality/size balance
  placeholder="blur" // Blur placeholder while loading
  blurDataURL="data:image/..." // Base64 blur
  priority={isAboveFold} // Load immediately if above fold
/>
```

**Find all images**:
```bash
# Search for <img> tags
grep -r "<img" src/components
grep -r "<img" src/app
```

**Components to Update**:
- Quiz card images
- User avatars
- Category images
- Achievement badges
- Landing page images
- Profile banners

**Expected Outcome**:
- Automatic WebP/AVIF conversion
- Responsive images (different sizes for different screens)
- Lazy loading by default
- Better LCP scores

---

#### 7.2 Configure Image Optimization
**Type**: Real Performance
**Impact**: Medium
**Effort**: Low

**Update `next.config.js`**:
```javascript
module.exports = {
  images: {
    formats: ['image/avif', 'image/webp'], // Modern formats
    deviceSizes: [640, 750, 828, 1080, 1200, 1920], // Breakpoints
    imageSizes: [16, 32, 48, 64, 96, 128, 256, 384], // Smaller sizes
    minimumCacheTTL: 60 * 60 * 24 * 30, // Cache for 30 days

    // If using external images
    remotePatterns: [
      {
        protocol: 'https',
        hostname: '**.supabase.co', // Supabase storage
        pathname: '/storage/v1/object/public/**',
      },
    ],
  },
}
```

**Expected Outcome**:
- Optimal image formats served
- Responsive images at all breakpoints
- Long-term caching

---

#### 7.3 Implement Lazy Loading
**Type**: Real Performance
**Impact**: Medium
**Effort**: Low

**Pattern**:
```typescript
// Above fold - load immediately
<Image
  src={quiz.coverImage}
  alt={quiz.title}
  priority={true} // Disable lazy loading
  width={800}
  height={400}
/>

// Below fold - lazy load
<Image
  src={quiz.coverImage}
  alt={quiz.title}
  loading="lazy" // Default behavior
  width={800}
  height={400}
/>
```

**Blur Placeholders** (Better UX):
```typescript
import { getPlaiceholder } from 'plaiceholder'

// Generate blur at build time
export async function getStaticProps() {
  const { base64, img } = await getPlaiceholder('/images/quiz.jpg')

  return {
    props: {
      imageProps: {
        ...img,
        blurDataURL: base64,
      },
    },
  }
}

// Component
<Image
  {...imageProps}
  placeholder="blur"
  alt="Quiz"
/>
```

**Expected Outcome**:
- Only visible images load initially
- Better initial page load
- Smoother scrolling (images load as you scroll)

---

#### 7.4 CDN Configuration (if applicable)
**Type**: Real Performance
**Impact**: High
**Effort**: Low-High (depends on hosting)

**If self-hosting**:
- Use Vercel (built-in CDN for Next.js)
- Or configure Cloudflare CDN

**If using Supabase Storage**:
```typescript
// Ensure images are served from Supabase CDN
const imageUrl = supabase.storage
  .from('quiz-images')
  .getPublicUrl(imagePath)

// Use with next/image loader
// next.config.js
module.exports = {
  images: {
    loader: 'custom',
    loaderFile: './lib/image-loader.ts',
  },
}

// lib/image-loader.ts
export default function supabaseLoader({ src, width, quality }) {
  return `${src}?width=${width}&quality=${quality || 75}`
}
```

**Expected Outcome**:
- Images served from edge locations
- Faster image load times globally
- Reduced origin server load

---

#### 7.5 Icon Optimization
**Type**: Real Performance
**Impact**: Low
**Effort**: Low

**Verify tree-shaking works**:
```typescript
// Good - named imports (tree-shakeable)
import { User, Settings, LogOut } from 'lucide-react'

// Bad - default import (includes all icons)
import * as Icons from 'lucide-react'
```

**Check bundle**:
```bash
ANALYZE=true npm run build
# Verify lucide-react size is reasonable (<50kb)
```

**Expected Outcome**:
- Only used icons in bundle
- Smaller overall bundle size

---

#### 7.6 Font Optimization (Additional)
**Type**: Real Performance
**Impact**: Low
**Effort**: Low

**Subset fonts** (if using many characters):
```typescript
// app/layout.tsx
const inter = Inter({
  subsets: ['latin'],
  display: 'swap',
  preload: true,

  // Only include weights you use
  weight: ['400', '500', '600', '700'],

  // Optionally subset specific characters (if not using full latin)
  // axes: ['wght'],
})
```

**Expected Outcome**:
- Smaller font files
- Faster font loading

---

### Phase 7 Success Metrics

- [ ] All <img> tags migrated to next/image
- [ ] Image optimization configured
- [ ] Lazy loading implemented for below-fold images
- [ ] Blur placeholders added to key images
- [ ] CDN configured (if applicable)
- [ ] Icon tree-shaking verified
- [ ] Image size reduced by 60%+
- [ ] LCP improved by 25%

**Metrics to Record**:
```
Total Image Size: ___ KB (before) → ___ KB (after)
Image Format Distribution: JPEG ___%, WebP ___%, AVIF ___%
LCP: ___ ms (before) → ___ ms (after)
Images Lazy Loaded: ___%
```

---

## Phase 8: Advanced Optimizations

**Timeline**: Week 11-12
**Type**: Real Performance + Advanced
**Risk Level**: Medium-High
**Effort**: High

### Objectives
- Implement edge runtime for critical paths
- Add ISR for static content
- Implement service worker for offline support
- Optimize server-side performance

### Tasks

#### 8.1 Edge Runtime for Middleware
**Type**: Real Performance
**Impact**: Medium
**Effort**: Low

**Update Middleware**:
```typescript
// src/middleware.ts
export const config = {
  matcher: [
    '/((?!_next/static|_next/image|favicon.ico).*)',
  ],
  runtime: 'edge', // Run on edge for faster global response
}

export async function middleware(request: NextRequest) {
  // Runs on edge - closer to users
  const session = await validateSession(request)

  if (!session) {
    return NextResponse.redirect(new URL('/login', request.url))
  }

  return NextResponse.next()
}
```

**Expected Outcome**:
- Faster auth checks globally
- Reduced latency for protected routes
- Better performance for international users

---

#### 8.2 Incremental Static Regeneration (ISR)
**Type**: Real Performance
**Impact**: Medium
**Effort**: Medium

**Landing Page** (Static with periodic updates):
```typescript
// app/page.tsx
export const revalidate = 3600 // Revalidate every hour

export default async function LandingPage() {
  // This page is statically generated at build time
  // and regenerated every hour
  const featuredQuizzes = await getQuizzes({ featured: true })

  return (
    <div>
      <Hero />
      <FeaturedQuizzes quizzes={featuredQuizzes} />
    </div>
  )
}
```

**Category Pages** (Static):
```typescript
// app/(dashboard)/categories/page.tsx
export const revalidate = 1800 // Revalidate every 30 minutes

export default async function CategoriesPage() {
  const categories = await getCategories()

  return <CategoryGrid categories={categories} />
}
```

**Quiz Detail Pages** (ISR):
```typescript
// app/(dashboard)/quizzes/[id]/page.tsx
export const dynamicParams = true // Enable for pages not generated at build
export const revalidate = 300 // Revalidate every 5 minutes

export async function generateStaticParams() {
  // Generate popular quizzes at build time
  const popularQuizzes = await getQuizzes({ popular: true, limit: 100 })

  return popularQuizzes.map(quiz => ({
    id: quiz.id,
  }))
}

export default async function QuizDetailPage({ params }: { params: { id: string } }) {
  const quiz = await getQuizDetails(params.id)

  return <QuizDetail quiz={quiz} />
}
```

**Expected Outcome**:
- Static performance for dynamic content
- Faster page loads
- Reduced server load

---

#### 8.3 Service Worker & PWA
**Type**: Real Performance + Offline
**Impact**: Medium-High
**Effort**: High

**Install next-pwa**:
```bash
npm install next-pwa
```

**Configure**:
```javascript
// next.config.js
const withPWA = require('next-pwa')({
  dest: 'public',
  register: true,
  skipWaiting: true,
  disable: process.env.NODE_ENV === 'development',
})

module.exports = withPWA({
  // existing config
})
```

**Add Manifest** (`public/manifest.json`):
```json
{
  "name": "QuizNinja",
  "short_name": "QuizNinja",
  "description": "Take quizzes, compete with friends, earn achievements",
  "start_url": "/dashboard",
  "display": "standalone",
  "background_color": "#ffffff",
  "theme_color": "#000000",
  "icons": [
    {
      "src": "/icons/icon-192.png",
      "sizes": "192x192",
      "type": "image/png"
    },
    {
      "src": "/icons/icon-512.png",
      "sizes": "512x512",
      "type": "image/png"
    }
  ]
}
```

**Update Layout**:
```typescript
// app/layout.tsx
export const metadata = {
  manifest: '/manifest.json',
  themeColor: '#000000',
}
```

**Expected Outcome**:
- Offline support for previously visited pages
- Install prompt for app-like experience
- Faster repeat visits (cached assets)
- Better mobile experience

---

#### 8.4 HTTP/2 & Compression
**Type**: Real Performance
**Impact**: Medium
**Effort**: Low (if using Vercel) / High (if self-hosting)

**Vercel**: Enabled by default (HTTP/2, gzip, brotli)

**Self-hosting** (Next.js standalone):
```javascript
// server.js
const express = require('express')
const next = require('next')
const compression = require('compression')

const dev = process.env.NODE_ENV !== 'production'
const app = next({ dev })
const handle = app.getRequestHandler()

app.prepare().then(() => {
  const server = express()

  // Enable compression (gzip/brotli)
  server.use(compression())

  // Serve Next.js
  server.all('*', (req, res) => {
    return handle(req, res)
  })

  server.listen(3000, err => {
    if (err) throw err
    console.log('> Ready on http://localhost:3000')
  })
})
```

**Expected Outcome**:
- Smaller response sizes (50-70% reduction)
- Faster data transfer
- Better performance on slow connections

---

#### 8.5 Database Query Optimization (Backend)
**Type**: Real Performance
**Impact**: High
**Effort**: High (Backend work)

**Action items for backend team**:

1. **Add database indexes**:
   - Quiz queries by category, difficulty
   - User stats queries
   - Leaderboard queries (sorted by score)

2. **Optimize N+1 queries**:
   - Use joins instead of multiple queries
   - Example: Get quizzes with user favorites in single query

3. **Add database caching**:
   - Cache leaderboard (Redis)
   - Cache user stats
   - Cache category lists

4. **Pagination**:
   - Limit quiz list queries
   - Cursor-based pagination for leaderboard

**Frontend changes**:
```typescript
// Request paginated data
export function useQuizzes(page = 1, limit = 20) {
  return useQuery({
    queryKey: ['quizzes', page, limit],
    queryFn: () => getQuizzes({ page, limit }),
    keepPreviousData: true, // Smooth pagination
  })
}
```

**Expected Outcome**:
- 2-5x faster API responses
- Reduced database load
- Better scalability

---

#### 8.6 Background Sync (PWA)
**Type**: Perceived Performance
**Impact**: Medium
**Effort**: Medium-High

**Service Worker** (background sync for quiz progress):
```javascript
// public/sw.js
self.addEventListener('sync', event => {
  if (event.tag === 'sync-quiz-progress') {
    event.waitUntil(syncQuizProgress())
  }
})

async function syncQuizProgress() {
  const pendingProgress = await getFromIndexedDB('pending-progress')

  for (const progress of pendingProgress) {
    try {
      await fetch('/api/quiz/progress', {
        method: 'POST',
        body: JSON.stringify(progress),
      })
      await removeFromIndexedDB('pending-progress', progress.id)
    } catch (err) {
      console.error('Sync failed:', err)
    }
  }
}
```

**Expected Outcome**:
- Quiz progress syncs even when offline
- Better reliability
- Enhanced offline experience

---

### Phase 8 Success Metrics

- [ ] Middleware running on edge
- [ ] ISR implemented for 3+ pages
- [ ] Service worker configured
- [ ] PWA manifest added
- [ ] Compression enabled
- [ ] Backend query optimization plan created
- [ ] Background sync implemented

**Metrics to Record**:
```
Middleware Response Time: ___ ms (edge) vs ___ ms (serverless)
ISR Cache Hit Rate: ___%
Offline Support: Yes/No
Compression Ratio: ___% reduction
API Response Time: ___ ms (before) → ___ ms (after)
```

---

## Phase 9: Monitoring & Continuous Improvement

**Timeline**: Ongoing
**Type**: Monitoring & Metrics
**Risk Level**: Low
**Effort**: Medium (setup) + Low (ongoing)

### Objectives
- Track real-world performance
- Monitor Core Web Vitals
- Set performance budgets
- Implement automated testing

### Tasks

#### 9.1 Web Vitals Monitoring
**Type**: Monitoring
**Impact**: High (Visibility)
**Effort**: Low

**Implement in App**:
```typescript
// app/layout.tsx
import { Analytics } from '@vercel/analytics/react'
import { SpeedInsights } from '@vercel/speed-insights/next'

export default function RootLayout({ children }: Props) {
  return (
    <html>
      <body>
        {children}
        <Analytics />
        <SpeedInsights />
      </body>
    </html>
  )
}

// Optional: Custom reporting
// app/web-vitals.ts
export function reportWebVitals(metric: NextWebVitalsMetric) {
  // Send to analytics
  if (metric.label === 'web-vital') {
    console.log(metric) // Or send to your analytics service

    // Example: Send to Google Analytics
    window.gtag?.('event', metric.name, {
      value: Math.round(metric.value),
      event_label: metric.id,
      non_interaction: true,
    })
  }
}
```

**Track Core Web Vitals**:
- LCP (Largest Contentful Paint)
- FID (First Input Delay) / INP (Interaction to Next Paint)
- CLS (Cumulative Layout Shift)
- FCP (First Contentful Paint)
- TTFB (Time to First Byte)

**Expected Outcome**:
- Real-world performance data
- Identify performance regressions
- Track improvements over time

---

#### 9.2 Error Tracking (Sentry)
**Type**: Monitoring
**Impact**: High
**Effort**: Low

**Install Sentry**:
```bash
npm install @sentry/nextjs
npx @sentry/wizard@latest -i nextjs
```

**Configure**:
```javascript
// sentry.client.config.js
import * as Sentry from '@sentry/nextjs'

Sentry.init({
  dsn: process.env.NEXT_PUBLIC_SENTRY_DSN,
  tracesSampleRate: 0.1, // 10% of transactions
  environment: process.env.NODE_ENV,

  // Performance monitoring
  integrations: [
    new Sentry.BrowserTracing({
      tracePropagationTargets: [process.env.NEXT_PUBLIC_API_BASE_URL],
    }),
  ],
})
```

**Expected Outcome**:
- Production error tracking
- Performance transaction tracking
- User session replay (optional)

---

#### 9.3 Performance Budgets
**Type**: Enforcement
**Impact**: High
**Effort**: Low

**Create Budget File** (`.lighthouserc.js`):
```javascript
module.exports = {
  ci: {
    collect: {
      url: ['http://localhost:3000/dashboard', 'http://localhost:3000/quizzes'],
      numberOfRuns: 3,
    },
    assert: {
      assertions: {
        'categories:performance': ['error', { minScore: 0.9 }],
        'categories:accessibility': ['warn', { minScore: 0.9 }],

        // Performance budgets
        'first-contentful-paint': ['error', { maxNumericValue: 1200 }],
        'largest-contentful-paint': ['error', { maxNumericValue: 2000 }],
        'total-blocking-time': ['warn', { maxNumericValue: 200 }],
        'cumulative-layout-shift': ['error', { maxNumericValue: 0.1 }],

        // Bundle size budgets
        'total-byte-weight': ['error', { maxNumericValue: 500000 }], // 500kb
      },
    },
  },
}
```

**Add to package.json**:
```json
{
  "scripts": {
    "perf": "lhci autorun"
  }
}
```

**Expected Outcome**:
- Prevent performance regressions
- Enforce standards
- CI/CD integration

---

#### 9.4 Lighthouse CI
**Type**: Automation
**Impact**: High
**Effort**: Medium

**GitHub Actions** (`.github/workflows/lighthouse.yml`):
```yaml
name: Lighthouse CI
on: [pull_request]

jobs:
  lighthouse:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
      - uses: actions/setup-node@v3
        with:
          node-version: 18

      - name: Install dependencies
        run: npm ci

      - name: Build
        run: npm run build

      - name: Run Lighthouse CI
        run: |
          npm install -g @lhci/cli
          lhci autorun
        env:
          LHCI_GITHUB_APP_TOKEN: ${{ secrets.LHCI_GITHUB_APP_TOKEN }}
```

**Expected Outcome**:
- Automated performance testing
- Performance scores on every PR
- Prevent regressions before merge

---

#### 9.5 Bundle Size Monitoring
**Type**: Monitoring
**Impact**: Medium
**Effort**: Low

**Add to CI** (GitHub Actions):
```yaml
- name: Analyze bundle size
  run: |
    npm run build
    npx -p nextjs-bundle-analysis report
```

**Use Bundlephobia Bot**:
- Automatically comments on PRs with bundle size changes
- Tracks bundle size over time

**Expected Outcome**:
- Visibility into bundle changes
- Prevent bundle bloat
- Track over time

---

#### 9.6 Real User Monitoring (RUM)
**Type**: Monitoring
**Impact**: High
**Effort**: Medium

**Options**:
1. **Vercel Analytics** (if using Vercel)
2. **Google Analytics 4** with Web Vitals
3. **Custom RUM**

**Custom RUM Implementation**:
```typescript
// lib/rum.ts
export function trackPerformance() {
  if (typeof window === 'undefined') return

  // Track page load
  window.addEventListener('load', () => {
    const perfData = performance.getEntriesByType('navigation')[0] as PerformanceNavigationTiming

    const metrics = {
      dns: perfData.domainLookupEnd - perfData.domainLookupStart,
      tcp: perfData.connectEnd - perfData.connectStart,
      ttfb: perfData.responseStart - perfData.requestStart,
      download: perfData.responseEnd - perfData.responseStart,
      domInteractive: perfData.domInteractive - perfData.fetchStart,
      domComplete: perfData.domComplete - perfData.fetchStart,
      loadComplete: perfData.loadEventEnd - perfData.fetchStart,
    }

    // Send to analytics
    sendToAnalytics('page_performance', metrics)
  })
}

// app/layout.tsx
useEffect(() => {
  trackPerformance()
}, [])
```

**Expected Outcome**:
- Real user data (not synthetic)
- Geographic performance insights
- Device/browser performance breakdown

---

#### 9.7 A/B Testing Performance Improvements
**Type**: Validation
**Impact**: Medium
**Effort**: Medium-High

**Setup**:
```bash
npm install @vercel/flags
```

**Example - Test Optimizations**:
```typescript
// Test virtual scrolling vs regular scrolling
import { useFlag } from '@vercel/flags/react'

function Leaderboard() {
  const useVirtualScrolling = useFlag('virtual-scrolling')

  if (useVirtualScrolling) {
    return <VirtualLeaderboard />
  }

  return <RegularLeaderboard />
}
```

**Track Results**:
- Engagement metrics
- Bounce rate
- Time on site
- User satisfaction

**Expected Outcome**:
- Data-driven decisions
- Validate improvements
- Measure real impact

---

### Phase 9 Success Metrics

- [ ] Web Vitals monitoring implemented
- [ ] Sentry error tracking configured
- [ ] Performance budgets defined
- [ ] Lighthouse CI integrated
- [ ] Bundle size monitoring automated
- [ ] RUM implemented
- [ ] A/B testing framework set up

**Ongoing Tracking**:
```
Weekly Reports:
- LCP P75: ___ ms
- FID/INP P75: ___ ms
- CLS P75: ___
- Error Rate: ___%
- Bundle Size: ___ KB
```

---

## Success Criteria & KPIs

### Real Performance Targets

#### Core Web Vitals (95th percentile)
- ✅ **Largest Contentful Paint (LCP)**: < 2.0s (Good: < 2.5s)
- ✅ **Interaction to Next Paint (INP)**: < 200ms (Good: < 200ms)
- ✅ **Cumulative Layout Shift (CLS)**: < 0.1 (Good: < 0.1)
- ✅ **First Contentful Paint (FCP)**: < 1.2s (Good: < 1.8s)
- ✅ **Time to First Byte (TTFB)**: < 600ms (Good: < 800ms)

#### Additional Metrics
- ✅ **Time to Interactive (TTI)**: < 3.0s
- ✅ **Total Blocking Time (TBT)**: < 200ms
- ✅ **Speed Index**: < 2.5s

#### Bundle Size
- ✅ **Initial bundle (gzipped)**: < 200 KB
- ✅ **Total JavaScript**: < 500 KB
- ✅ **First Load JS**: < 150 KB

#### Lighthouse Scores (Mobile)
- ✅ **Performance**: 90+
- ✅ **Accessibility**: 95+
- ✅ **Best Practices**: 95+
- ✅ **SEO**: 95+

---

### Perceived Performance Targets

#### User Experience
- ✅ **All interactions**: Feel instant (< 100ms perceived)
- ✅ **Route transitions**: < 200ms
- ✅ **Loading state coverage**: 100% of async operations
- ✅ **Skeleton screen coverage**: All data-heavy pages
- ✅ **Zero blank screens**: Always show something to user

#### Engagement Metrics
- ✅ **Bounce rate**: Reduce by 30%
- ✅ **Time on site**: Increase by 20%
- ✅ **Pages per session**: Increase by 15%
- ✅ **Quiz completion rate**: Increase by 25%
- ✅ **User satisfaction**: 4.5+ / 5.0 rating

---

### Per-Phase Targets

| Phase | Metric | Target | Baseline | Actual |
|-------|--------|--------|----------|--------|
| Phase 1 | Baseline established | ✅ | - | - |
| Phase 2 | Bundle reduction | 40% | TBD | - |
| Phase 2 | FCP improvement | 30% | TBD | - |
| Phase 3 | Client components | < 40 | 103 | - |
| Phase 3 | Dashboard load | 2x faster | TBD | - |
| Phase 4 | Re-renders | 70% reduction | TBD | - |
| Phase 4 | List scrolling | 60 fps | TBD | - |
| Phase 5 | Interaction response | < 100ms | TBD | - |
| Phase 6 | Question transition | < 50ms | TBD | - |
| Phase 7 | Image size | 60% reduction | TBD | - |
| Phase 7 | LCP improvement | 25% | TBD | - |

---

## Implementation Guidelines

### Testing Between Phases

After each phase:
1. **Run Lighthouse** (mobile & desktop)
2. **Measure Web Vitals** (lab & field data)
3. **Test on slow 3G** (throttled connection)
4. **User testing** (5-10 users for subjective feedback)
5. **Compare vs baseline** (document improvements)

### Rollback Plan

- Each phase = separate git branch
- Merge only after validation
- Tag releases for easy rollback
- Keep monitoring after merge

### Documentation

- Update this plan with actual metrics after each phase
- Document any blockers or issues
- Share learnings with team
- Update timeline if needed

### Dependencies

**New dependencies** (minimal):
- `@next/bundle-analyzer` (dev) - Phase 1
- `nprogress` - Phase 1
- `@tanstack/react-virtual` - Phase 4 (optional)
- `framer-motion` - Phase 5 (optional, if not using CSS)
- `next-pwa` - Phase 8
- `@vercel/analytics` - Phase 9
- `@sentry/nextjs` - Phase 9

**To remove** (after audit):
- Unused dependencies identified in Phase 2

---

## Timeline Summary

| Week | Phase | Focus |
|------|-------|-------|
| 1 | Phase 1 | Quick wins, baseline |
| 2-3 | Phase 2 | Code splitting, bundle optimization |
| 4-5 | Phase 3 | Server components, data fetching |
| 6 | Phase 4 | Memoization, re-render optimization |
| 7-8 | Phase 5 | Perceived performance, UX |
| 9 | Phase 6 | Quiz taking optimization |
| 10 | Phase 7 | Image & asset optimization |
| 11-12 | Phase 8 | Advanced optimizations |
| Ongoing | Phase 9 | Monitoring, continuous improvement |

**Total Estimated Timeline**: 12 weeks (3 months)

---

## Appendix

### Useful Commands

```bash
# Build analysis
ANALYZE=true npm run build

# Production build
npm run build
npm run start

# Lighthouse
npx lighthouse http://localhost:3000 --view

# Bundle size
npx next build --profile

# Find unused dependencies
npx depcheck

# Check bundle contents
npx webpack-bundle-analyzer .next/stats.json
```

### Resources

- [Next.js Performance Docs](https://nextjs.org/docs/advanced-features/measuring-performance)
- [Web Vitals](https://web.dev/vitals/)
- [React Performance](https://react.dev/learn/render-and-commit)
- [Bundle Analysis](https://bundlephobia.com/)

---

## Revision History

| Date | Version | Changes | Author |
|------|---------|---------|--------|
| 2025-11-16 | 1.0 | Initial plan created | - |
| | | | |
| | | | |

---

**Next Steps**:
1. Review and approve this plan
2. Set up tracking spreadsheet for metrics
3. Begin Phase 1 implementation
4. Schedule weekly progress reviews