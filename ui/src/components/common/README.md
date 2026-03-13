# Common Components

## Overview

Shared utility components used throughout the application. These provide consistent patterns for loading states, errors, empty content, and reusable UI patterns.

## Components

| Component | File | Purpose |
|-----------|------|---------|
| GlassCard | `GlassCard.tsx` | Glassmorphic card container |
| LoadingSpinner | `LoadingSpinner.tsx` | Loading indicator |
| EmptyState | `EmptyState.tsx` | Empty content placeholder |
| ErrorBoundary | `ErrorBoundary.tsx` | Error handling wrapper |
| StatsCard | `StatsCard.tsx` | Statistics display card |
| StatsGrid | `StatsGrid.tsx` | Grid of statistics cards |

## GlassCard

A card component with glassmorphism effect (blurred, translucent background).

### Props

| Prop | Type | Default | Description |
|------|------|---------|-------------|
| `children` | `ReactNode` | - | Card content |
| `hover` | `boolean` | `false` | Enable hover lift effect |
| `padding` | `"none" \| "sm" \| "md" \| "lg"` | `"md"` | Padding size |
| `rounded` | `"xl" \| "2xl" \| "3xl"` | `"2xl"` | Border radius |
| `className` | `string` | - | Additional classes |

### Usage

```tsx
import { GlassCard } from "@/components/common/GlassCard";

// Basic
<GlassCard>
  <h2>Card Title</h2>
  <p>Card content</p>
</GlassCard>

// With hover effect
<GlassCard hover>
  <QuizPreview quiz={quiz} />
</GlassCard>

// Custom padding and rounded
<GlassCard padding="lg" rounded="3xl">
  <ProfileHeader user={user} />
</GlassCard>

// No padding (for images)
<GlassCard padding="none">
  <img src={image} className="rounded-2xl" />
</GlassCard>
```

### Styling

```css
/* Applied classes */
.glass-card {
  background: rgba(255, 255, 255, 0.4); /* Light mode */
  background: rgba(0, 0, 0, 0.4);       /* Dark mode */
  backdrop-filter: blur(12px);
  border: 1px solid rgba(255, 255, 255, 0.2);
  box-shadow: 0 10px 15px rgba(0, 0, 0, 0.05);
}
```

---

## LoadingSpinner

Animated loading indicator with size variants.

### Props

| Prop | Type | Default | Description |
|------|------|---------|-------------|
| `size` | `"sm" \| "md" \| "lg"` | `"md"` | Spinner size |
| `className` | `string` | - | Additional classes |

### Usage

```tsx
import { LoadingSpinner } from "@/components/common/LoadingSpinner";

// In loading states
if (isLoading) {
  return <LoadingSpinner />;
}

// Centered in container
<div className="flex justify-center py-12">
  <LoadingSpinner size="lg" />
</div>

// Inline with text
<div className="flex items-center gap-2">
  <LoadingSpinner size="sm" />
  <span>Loading...</span>
</div>
```

---

## EmptyState

Placeholder for empty content with optional icon, title, description, and action.

### Props

| Prop | Type | Default | Description |
|------|------|---------|-------------|
| `icon` | `ReactNode` | - | Icon component |
| `title` | `string` | - | Empty state title |
| `description` | `string` | - | Explanatory text |
| `action` | `ReactNode` | - | Call-to-action button |
| `className` | `string` | - | Additional classes |

### Usage

```tsx
import { EmptyState } from "@/components/common/EmptyState";
import { FileQuestion } from "lucide-react";
import { Button } from "@/components/ui/button";

// Basic
<EmptyState
  icon={<FileQuestion className="h-12 w-12" />}
  title="No quizzes found"
  description="Try adjusting your filters or search terms."
/>

// With action
<EmptyState
  icon={<Users className="h-12 w-12" />}
  title="No friends yet"
  description="Search for users to add as friends."
  action={
    <Button onClick={() => setSearchOpen(true)}>
      Find Friends
    </Button>
  }
/>

// In conditional rendering
{quizzes.length === 0 ? (
  <EmptyState
    title="No quizzes"
    description="Check back later for new quizzes."
  />
) : (
  <QuizGrid quizzes={quizzes} />
)}
```

---

## ErrorBoundary

React error boundary for graceful error handling.

### Props

| Prop | Type | Default | Description |
|------|------|---------|-------------|
| `children` | `ReactNode` | - | Content to wrap |
| `fallback` | `ReactNode` | - | Custom error UI |

### Usage

```tsx
import { ErrorBoundary } from "@/components/common/ErrorBoundary";

// Wrap error-prone components
<ErrorBoundary>
  <QuizContent quiz={quiz} />
</ErrorBoundary>

// With custom fallback
<ErrorBoundary
  fallback={
    <div className="text-center py-12">
      <p>Something went wrong loading this quiz.</p>
      <Button onClick={() => window.location.reload()}>
        Refresh Page
      </Button>
    </div>
  }
>
  <QuizTakingView />
</ErrorBoundary>

// Page-level boundary
export default function QuizPage() {
  return (
    <ErrorBoundary>
      <QuizContent />
    </ErrorBoundary>
  );
}
```

---

## StatsCard

Card displaying a single statistic with label and icon.

### Props

| Prop | Type | Default | Description |
|------|------|---------|-------------|
| `title` | `string` | - | Stat label |
| `value` | `string \| number` | - | Stat value |
| `icon` | `ReactNode` | - | Icon component |
| `description` | `string` | - | Additional info |
| `trend` | `"up" \| "down" \| "neutral"` | - | Trend indicator |
| `className` | `string` | - | Additional classes |

### Usage

```tsx
import { StatsCard } from "@/components/common/StatsCard";
import { Trophy, Target, Clock } from "lucide-react";

<StatsCard
  title="Total Score"
  value={1250}
  icon={<Trophy className="h-5 w-5" />}
  trend="up"
  description="+15% this week"
/>

<StatsCard
  title="Accuracy"
  value="78%"
  icon={<Target className="h-5 w-5" />}
/>

<StatsCard
  title="Time Spent"
  value="12h 30m"
  icon={<Clock className="h-5 w-5" />}
/>
```

---

## StatsGrid

Grid layout for multiple StatsCard components.

### Props

| Prop | Type | Default | Description |
|------|------|---------|-------------|
| `children` | `ReactNode` | - | StatsCard components |
| `columns` | `2 \| 3 \| 4` | `4` | Number of columns |
| `className` | `string` | - | Additional classes |

### Usage

```tsx
import { StatsGrid, StatsCard } from "@/components/common";

<StatsGrid columns={4}>
  <StatsCard title="Quizzes" value={42} icon={<BookOpen />} />
  <StatsCard title="Score" value={1250} icon={<Trophy />} />
  <StatsCard title="Accuracy" value="78%" icon={<Target />} />
  <StatsCard title="Streak" value={7} icon={<Flame />} />
</StatsGrid>

// 3 columns for smaller grids
<StatsGrid columns={3}>
  <StatsCard title="Friends" value={15} />
  <StatsCard title="Achievements" value={8} />
  <StatsCard title="Rank" value="#42" />
</StatsGrid>
```

---

## Common Patterns

### Loading → Error → Empty → Data

```tsx
function DataDisplay() {
  const { data, isLoading, error } = useMyData();

  // Loading state
  if (isLoading) {
    return (
      <div className="flex justify-center py-12">
        <LoadingSpinner size="lg" />
      </div>
    );
  }

  // Error state
  if (error) {
    return (
      <EmptyState
        icon={<AlertCircle className="h-12 w-12 text-destructive" />}
        title="Error loading data"
        description={error.message}
        action={
          <Button onClick={() => refetch()}>Try Again</Button>
        }
      />
    );
  }

  // Empty state
  if (!data || data.length === 0) {
    return (
      <EmptyState
        icon={<FileX className="h-12 w-12" />}
        title="No data found"
        description="There's nothing here yet."
      />
    );
  }

  // Data display
  return <DataList items={data} />;
}
```

### Consistent Card Styling

```tsx
// Use GlassCard for featured/highlighted items
<GlassCard hover>
  <FeaturedQuiz quiz={quiz} />
</GlassCard>

// Use regular Card for list items
<Card>
  <QuizListItem quiz={quiz} />
</Card>
```

## Related Documentation

- [Parent: Components Overview](../README.md)
- [UI Primitives](../ui/README.md) - Base components
- [Dashboard Components](../dashboard/README.md) - Uses StatsCard/StatsGrid
