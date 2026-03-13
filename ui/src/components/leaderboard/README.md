# Leaderboard Components

## Overview

Components for displaying global and category-specific rankings. Users can see how they compare to others based on quiz performance and achievements.

## Components

| Component | File | Purpose |
|-----------|------|---------|
| LeaderboardTable | `LeaderboardTable.tsx` | Main ranking table |
| UserRankCard | `UserRankCard.tsx` | Current user's rank highlight |

## LeaderboardTable

Sortable table showing user rankings.

### Props

| Prop | Type | Description |
|------|------|-------------|
| `entries` | `LeaderboardEntry[]` | Ranking data |
| `isLoading` | `boolean` | Loading state |
| `error` | `Error \| null` | Error state |

### Columns

| Column | Visibility | Description |
|--------|------------|-------------|
| Rank | Always | Position with medal icons |
| User | Always | Avatar, name, "You" badge |
| Points | Always | Total points earned |
| Quizzes | sm+ screens | Quizzes completed |
| Achievements | md+ screens | Achievement count |

### Rank Icons

```tsx
// Top 3 positions get special icons
const getRankIcon = (rank: number) => {
  switch (rank) {
    case 1: return <Trophy className="text-yellow-500" />;  // Gold
    case 2: return <Medal className="text-gray-400" />;     // Silver
    case 3: return <Award className="text-amber-600" />;    // Bronze
    default: return null;
  }
};
```

### Current User Highlight

```tsx
// Highlight row if it's the current user
<TableRow className={isCurrentUser ? "bg-primary/5 font-medium" : ""}>
  <Badge variant="outline">You</Badge>
</TableRow>
```

### States

**Loading:**
```tsx
<div className="space-y-3">
  {[...Array(10)].map((_, i) => (
    <div className="flex items-center gap-3 p-4 border rounded-lg">
      <Skeleton className="h-10 w-10 rounded-full" />
      <Skeleton className="h-4 w-32" />
    </div>
  ))}
</div>
```

**Error:**
```tsx
<Alert variant="destructive">
  <AlertCircle className="h-4 w-4" />
  <AlertDescription>Failed to load leaderboard</AlertDescription>
</Alert>
```

**Empty:**
```tsx
<div className="text-center py-12">
  <Trophy className="h-12 w-12 text-muted-foreground" />
  <h3>No Leaderboard Data</h3>
  <p>Be the first to compete!</p>
</div>
```

### Usage

```tsx
import { LeaderboardTable } from "@/components/leaderboard/LeaderboardTable";
import { useLeaderboard } from "@/hooks/useLeaderboard";

function LeaderboardPage() {
  const { data: entries, isLoading, error } = useLeaderboard();

  return (
    <LeaderboardTable
      entries={entries}
      isLoading={isLoading}
      error={error}
    />
  );
}
```

---

## UserRankCard

Highlights the current user's ranking.

### Props

| Prop | Type | Description |
|------|------|-------------|
| `rank` | `number` | User's current rank |
| `totalUsers` | `number` | Total users in leaderboard |
| `points` | `number` | User's total points |
| `percentile` | `number` | Top X% percentile |

### Display

- Current rank position
- Percentile indicator
- Points earned
- Progress to next rank

### Usage

```tsx
import { UserRankCard } from "@/components/leaderboard/UserRankCard";

<UserRankCard
  rank={42}
  totalUsers={1000}
  points={2500}
  percentile={5}
/>
```

---

## Leaderboard Page Structure

```tsx
// /leaderboard page
<div className="space-y-8">
  {/* Page Header */}
  <div>
    <h1>Leaderboard</h1>
    <p>See how you rank against other quiz ninjas</p>
  </div>

  {/* Current User Rank */}
  <UserRankCard
    rank={userRank}
    totalUsers={totalUsers}
    points={userPoints}
    percentile={percentile}
  />

  {/* Filter Tabs */}
  <Tabs defaultValue="global">
    <TabsList>
      <TabsTrigger value="global">Global</TabsTrigger>
      <TabsTrigger value="weekly">This Week</TabsTrigger>
      <TabsTrigger value="monthly">This Month</TabsTrigger>
    </TabsList>

    <TabsContent value="global">
      <LeaderboardTable entries={globalEntries} />
    </TabsContent>
  </Tabs>
</div>
```

## Data Types

```typescript
interface LeaderboardEntry {
  user_id: string;
  name: string;
  avatar?: string;
  rank: number;
  points: number;
  quizzes_completed: number;
  achievements?: Achievement[];
  is_current_user?: boolean;
}

interface LeaderboardFilters {
  timeframe: "all" | "weekly" | "monthly";
  category?: string;
  limit?: number;
}
```

## Hooks Used

```typescript
// Get leaderboard entries
const { data, isLoading, error } = useLeaderboard({
  timeframe: "all",
  limit: 100,
});

// Get user's rank
const { data: userRank } = useUserRank();
```

## Ranking Algorithm

Rankings are based on:
1. **Total Points** (primary) - Points from all quizzes
2. **Quiz Count** (tiebreaker) - Number of quizzes completed
3. **Achievement Count** (secondary) - Bonus for achievements

## Related Documentation

- [Parent: Components Overview](../README.md)
- [Dashboard Components](../dashboard/README.md)
- [Profile Components](../profile/README.md)
- [API Types](../../types/README.md)
- [useLeaderboard Hook](../../hooks/README.md)

