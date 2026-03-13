# Dashboard Components

## Overview

Components for the main dashboard page, displaying user statistics, featured quizzes, and recent activity.

## Components

| Component | File | Purpose |
|-----------|------|---------|
| FeaturedQuizzesDashboard | `FeaturedQuizzesDashboard.tsx` | Featured quiz carousel/grid |
| DashboardStatsCard | `DashboardStatsCard.tsx` | User statistics overview |
| RecentActivity | `RecentActivity.tsx` | Recent quiz attempts |

## DashboardStatsCard

Displays user statistics in a grid layout.

### Props

| Prop | Type | Description |
|------|------|-------------|
| `stats` | `UserStats` | User statistics data |
| `isLoading` | `boolean` | Loading state |

### Statistics Displayed

| Stat | Icon | Description |
|------|------|-------------|
| Quizzes Completed | BookOpen | Total completed quizzes |
| Total Score | Trophy | Cumulative points |
| Accuracy | Target | Average accuracy percentage |
| Current Streak | Flame | Consecutive days active |

### Usage

```tsx
import { DashboardStatsCard } from "@/components/dashboard/DashboardStatsCard";
import { useUserStats } from "@/hooks/useUserStats";

function Dashboard() {
  const { data: stats, isLoading } = useUserStats();

  return <DashboardStatsCard stats={stats} isLoading={isLoading} />;
}
```

### Layout

```tsx
<StatsGrid columns={4}>
  <StatsCard
    title="Quizzes"
    value={stats.quizzes_completed}
    icon={<BookOpen />}
  />
  <StatsCard
    title="Total Score"
    value={stats.total_points}
    icon={<Trophy />}
  />
  <StatsCard
    title="Accuracy"
    value={`${stats.accuracy}%`}
    icon={<Target />}
  />
  <StatsCard
    title="Streak"
    value={stats.current_streak}
    icon={<Flame />}
  />
</StatsGrid>
```

---

## FeaturedQuizzesDashboard

Displays featured quizzes for quick access.

### Props

| Prop | Type | Description |
|------|------|-------------|
| `quizzes` | `Quiz[]` | Featured quiz list |
| `isLoading` | `boolean` | Loading state |

### Features

- Horizontal scroll on mobile
- Grid layout on desktop
- Quick-start button
- Category badges
- Difficulty indicators

### Usage

```tsx
import { FeaturedQuizzesDashboard } from "@/components/dashboard/FeaturedQuizzesDashboard";
import { useFeaturedQuizzes } from "@/hooks/useFeaturedQuizzes";

function Dashboard() {
  const { data: quizzes, isLoading } = useFeaturedQuizzes();

  return (
    <section>
      <h2>Featured Quizzes</h2>
      <FeaturedQuizzesDashboard quizzes={quizzes} isLoading={isLoading} />
    </section>
  );
}
```

---

## RecentActivity

Shows user's recent quiz attempts.

### Props

| Prop | Type | Description |
|------|------|-------------|
| `attempts` | `QuizAttempt[]` | Recent attempts |
| `limit` | `number` | Max items to show |

### Display

- Quiz title
- Score achieved
- Date/time
- Result (pass/fail)
- Link to full results

### Usage

```tsx
<RecentActivity attempts={recentAttempts} limit={5} />
```

---

## Dashboard Page Layout

```tsx
// app/(dashboard)/dashboard/page.tsx
export default function DashboardPage() {
  return (
    <div className="space-y-8">
      {/* Page header */}
      <div>
        <h1 className="text-2xl font-bold">Dashboard</h1>
        <p className="text-muted-foreground">
          Welcome back! Here's your progress.
        </p>
      </div>

      {/* Stats grid */}
      <DashboardStatsCard />

      {/* Featured quizzes */}
      <section>
        <div className="flex justify-between items-center mb-4">
          <h2 className="text-xl font-semibold">Featured Quizzes</h2>
          <Link href="/quizzes">View All</Link>
        </div>
        <FeaturedQuizzesDashboard />
      </section>

      {/* Recent activity */}
      <section>
        <h2 className="text-xl font-semibold mb-4">Recent Activity</h2>
        <RecentActivity />
      </section>
    </div>
  );
}
```

## Data Hooks

```tsx
// User statistics
const { data: stats } = useUserStats();

// Featured quizzes
const { data: featured } = useFeaturedQuizzes();

// Recent attempts (from user stats or separate hook)
const { data: attempts } = useUserAttempts({ limit: 5 });
```

## Related Documentation

- [Parent: Components Overview](../README.md)
- [Common Components](../common/README.md) - StatsCard, StatsGrid
- [Quiz Components](../quiz/README.md) - QuizCard
- [useUserStats Hook](../../hooks/README.md)
