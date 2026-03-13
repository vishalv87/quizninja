# Enums

## Overview

Type-safe enum constants using the "const object" pattern. Each enum provides a constant object, TypeScript type, array of values, and type guard function.

## Files

| File | Enums | Domain |
|------|-------|--------|
| `quiz.ts` | QuizDifficulty, QuestionType, AttemptStatus | Quiz system |
| `user.ts` | Theme, ProfileVisibility, FriendRequestStatus, FriendshipStatus, NotificationFrequency | User settings |
| `notification.ts` | NotificationType | Notification system |
| `discussion.ts` | DiscussionType, DiscussionSort | Discussion forum |
| `achievement.ts` | AchievementCategory, AchievementRequirementType | Achievements |
| `leaderboard.ts` | LeaderboardPeriod | Leaderboard |
| `common.ts` | SortOrder | Shared |

## Pattern

Every enum follows the same pattern for consistency:

```tsx
// 1. Const object - the actual values
export const MyEnum = {
  VALUE_ONE: 'value_one',
  VALUE_TWO: 'value_two',
} as const;

// 2. Type - union of all values
export type MyEnum = typeof MyEnum[keyof typeof MyEnum];
// Result: "value_one" | "value_two"

// 3. Array - all values for iteration/validation
export const MY_ENUMS = Object.values(MyEnum);
// Result: ["value_one", "value_two"]

// 4. Type guard - runtime validation
export function isMyEnum(value: unknown): value is MyEnum {
  return typeof value === 'string' && MY_ENUMS.includes(value as MyEnum);
}
```

## Quiz Enums (`quiz.ts`)

### QuizDifficulty

Difficulty levels for quizzes.

| Constant | Value | Description |
|----------|-------|-------------|
| `BEGINNER` | "beginner" | Entry-level difficulty |
| `INTERMEDIATE` | "intermediate" | Medium difficulty |
| `ADVANCED` | "advanced" | Expert-level difficulty |

```tsx
import { QuizDifficulty, isQuizDifficulty } from "@/constants";

// Use constant
const difficulty = QuizDifficulty.INTERMEDIATE;

// Validate unknown input
if (isQuizDifficulty(input)) {
  // input is typed as QuizDifficulty
}
```

### QuestionType

Types of quiz questions.

| Constant | Value | Description |
|----------|-------|-------------|
| `MULTIPLE_CHOICE` | "multiple_choice" | Select one from options |
| `TRUE_FALSE` | "true_false" | Binary choice |
| `SHORT_ANSWER` | "short_answer" | Text input |

### AttemptStatus

Status of a quiz attempt.

| Constant | Value | Description |
|----------|-------|-------------|
| `IN_PROGRESS` | "in_progress" | Quiz being taken |
| `COMPLETED` | "completed" | Quiz finished |
| `ABANDONED` | "abandoned" | Quiz left incomplete |

## User Enums (`user.ts`)

### Theme

UI theme preferences.

| Constant | Value | Description |
|----------|-------|-------------|
| `LIGHT` | "light" | Light mode |
| `DARK` | "dark" | Dark mode |
| `SYSTEM` | "system" | Follow system preference |

### ProfileVisibility

User profile privacy settings.

| Constant | Value | Description |
|----------|-------|-------------|
| `PUBLIC` | "public" | Visible to everyone |
| `FRIENDS_ONLY` | "friends_only" | Visible to friends |
| `PRIVATE` | "private" | Only visible to self |

### FriendRequestStatus

State of a friend request.

| Constant | Value | Description |
|----------|-------|-------------|
| `PENDING` | "pending" | Awaiting response |
| `ACCEPTED` | "accepted" | Request accepted |
| `REJECTED` | "rejected" | Request declined |

### FriendshipStatus

Relationship state between two users.

| Constant | Value | Description |
|----------|-------|-------------|
| `NONE` | "none" | Not connected |
| `PENDING_SENT` | "pending_sent" | Request sent, awaiting |
| `PENDING_RECEIVED` | "pending_received" | Request received |
| `FRIENDS` | "friends" | Connected as friends |

### NotificationFrequency

How often to send notifications.

| Constant | Value | Description |
|----------|-------|-------------|
| `INSTANT` | "instant" | Immediately |
| `DAILY` | "daily" | Once per day |
| `WEEKLY` | "weekly" | Once per week |
| `NEVER` | "never" | Disabled |

### FriendRequestAction

Actions on friend requests.

| Constant | Value | Description |
|----------|-------|-------------|
| `ACCEPT` | "accept" | Accept request |
| `DECLINE` | "decline" | Decline request |

## Notification Enums (`notification.ts`)

### NotificationType

Categories of notifications.

| Constant | Value | Description |
|----------|-------|-------------|
| `FRIEND_REQUEST` | "friend_request" | New friend request |
| `FRIEND_ACCEPTED` | "friend_accepted" | Request accepted |
| `ACHIEVEMENT_UNLOCKED` | "achievement_unlocked" | Achievement earned |
| `QUIZ_REMINDER` | "quiz_reminder" | Quiz reminder |
| `DISCUSSION_REPLY` | "discussion_reply" | Reply to discussion |
| `SYSTEM` | "system" | System notification |

## Discussion Enums (`discussion.ts`)

### DiscussionType

Types of discussion posts.

| Constant | Value | Description |
|----------|-------|-------------|
| `QUESTION` | "question" | Asking a question |
| `GENERAL` | "general" | General discussion |
| `BUG_REPORT` | "bug_report" | Reporting a bug |
| `FEATURE_REQUEST` | "feature_request" | Suggesting a feature |
| `DISCUSSION` | "discussion" | Open discussion |

### DiscussionSort

Sorting options for discussions.

| Constant | Value | Description |
|----------|-------|-------------|
| `RECENT` | "recent" | Most recent first |
| `POPULAR` | "popular" | Most popular first |

## Achievement Enums (`achievement.ts`)

### AchievementCategory

Groups of achievements.

| Constant | Value | Description |
|----------|-------|-------------|
| `QUIZ_MASTER` | "quiz_master" | Quiz completion |
| `SOCIAL` | "social" | Social interactions |
| `STREAK` | "streak" | Consistency |
| `KNOWLEDGE` | "knowledge" | Learning milestones |
| `COMPETITOR` | "competitor" | Competition wins |

### AchievementRequirementType

Types of achievement requirements.

| Constant | Value | Description |
|----------|-------|-------------|
| `QUIZZES_COMPLETED` | "quizzes_completed" | Number of quizzes |
| `TOTAL_POINTS` | "total_points" | Points earned |
| `ACCURACY_PERCENTAGE` | "accuracy_percentage" | Accuracy rate |
| `STREAK_REACHED` | "streak_reached" | Streak days |
| `FRIENDS_ADDED` | "friends_added" | Friends count |
| `DISCUSSIONS_STARTED` | "discussions_started" | Discussions created |

## Leaderboard Enums (`leaderboard.ts`)

### LeaderboardPeriod

Time periods for leaderboard filtering.

| Constant | Value | Description |
|----------|-------|-------------|
| `ALL_TIME` | "all_time" | All time rankings |
| `MONTHLY` | "monthly" | This month |
| `WEEKLY` | "weekly" | This week |

## Usage Examples

### In Type Definitions

```tsx
// types/quiz.ts
import type { QuizDifficulty, QuestionType } from "@/constants";

interface Quiz {
  id: string;
  title: string;
  difficulty: QuizDifficulty;
}

interface Question {
  id: string;
  question_type: QuestionType;
}
```

### In Components

```tsx
import { AttemptStatus, ATTEMPT_STATUSES } from "@/constants";

function StatusBadge({ status }: { status: AttemptStatus }) {
  const colors = {
    [AttemptStatus.IN_PROGRESS]: "bg-yellow-500",
    [AttemptStatus.COMPLETED]: "bg-green-500",
    [AttemptStatus.ABANDONED]: "bg-red-500",
  };

  return <Badge className={colors[status]}>{status}</Badge>;
}
```

### In API Responses

```tsx
import { isQuizDifficulty } from "@/constants";

async function fetchQuiz(id: string) {
  const data = await api.get(`/quizzes/${id}`);

  // Validate response
  if (!isQuizDifficulty(data.difficulty)) {
    throw new Error("Invalid difficulty from API");
  }

  return data;
}
```

### In Forms

```tsx
import { QuizDifficulty } from "@/constants";
import { DIFFICULTY_OPTIONS } from "@/constants";

const form = useForm({
  defaultValues: {
    difficulty: QuizDifficulty.BEGINNER,
  },
});

return (
  <Select {...form.register("difficulty")}>
    {DIFFICULTY_OPTIONS.map(opt => (
      <SelectItem key={opt.value} value={opt.value}>
        {opt.label}
      </SelectItem>
    ))}
  </Select>
);
```

## Adding New Enums

1. Create or edit the appropriate domain file:

```tsx
// constants/enums/myDomain.ts

export const MyNewEnum = {
  OPTION_A: 'option_a',
  OPTION_B: 'option_b',
} as const;

export type MyNewEnum = typeof MyNewEnum[keyof typeof MyNewEnum];
export const MY_NEW_ENUMS = Object.values(MyNewEnum);

export function isMyNewEnum(value: unknown): value is MyNewEnum {
  return typeof value === 'string' && MY_NEW_ENUMS.includes(value as MyNewEnum);
}
```

2. Export from `index.ts`:

```tsx
export * from './myDomain';
```

3. Add labels in `../labels.ts` if needed for UI display.

## Related Documentation

- [Parent: Constants](../README.md)
- [Labels](../labels.ts) - Human-readable labels
- [Types](../../types/README.md) - TypeScript definitions using these enums
