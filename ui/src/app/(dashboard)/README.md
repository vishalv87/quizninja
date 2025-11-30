# Dashboard Routes

## Overview

The `(dashboard)` route group contains all protected routes requiring authentication. These pages share a common layout with header, sidebar, and require the user to be logged in.

## Route Group Features

- **Authentication Required**: All routes wrapped with `AuthGuard`
- **Shared Layout**: Header, Sidebar, MobileNav
- **Real-time Notifications**: `DashboardNotificationListener`
- **URL Path**: No `/dashboard` prefix in URL (route groups don't affect URL)

## Routes

### Main Pages

| Route | Purpose |
|-------|---------|
| `/dashboard` | Main dashboard with stats and featured quizzes |
| `/quizzes` | Browse all quizzes with filters |
| `/categories` | Browse quizzes by category |
| `/favorites` | User's favorite quizzes |

### Quiz Flow

| Route | Purpose |
|-------|---------|
| `/quizzes/[id]` | Quiz details page |
| `/quizzes/[id]/take` | Quiz taking interface |
| `/quizzes/[id]/results/[attemptId]` | Quiz results view |
| `/quizzes/category/[categoryId]` | Quizzes in category |

### User & Profile

| Route | Purpose |
|-------|---------|
| `/profile` | Current user's profile |
| `/profile/edit` | Edit profile |
| `/profile/[userId]` | View other user's profile |
| `/settings` | User settings |

### Social & Community

| Route | Purpose |
|-------|---------|
| `/friends` | Friends list |
| `/friends/requests` | Friend requests |
| `/discussions` | Discussion forum |
| `/discussions/[id]` | Single discussion |

### Progress & Competition

| Route | Purpose |
|-------|---------|
| `/achievements` | User achievements |
| `/leaderboard` | Global rankings |
| `/notifications` | All notifications |

## Layout

```tsx
// (dashboard)/layout.tsx
export default function DashboardLayout({ children }: { children: ReactNode }) {
  return (
    <AuthGuard requireAuth={true}>
      {/* Real-time notification toasts */}
      <DashboardNotificationListener />

      <div className="flex h-screen flex-col overflow-hidden">
        {/* Top navigation */}
        <Header />

        <div className="flex flex-1 min-h-0">
          {/* Desktop sidebar */}
          <Sidebar />

          {/* Mobile navigation */}
          <MobileNav />

          {/* Main content area */}
          <main className="flex-1 overflow-y-auto">
            <div className="container mx-auto px-4 py-6 sm:px-6">
              {children}
            </div>
          </main>
        </div>
      </div>
    </AuthGuard>
  );
}
```

## Page Examples

### Dashboard Page

```tsx
// dashboard/page.tsx
"use client";

import { useFeaturedQuizzes } from "@/hooks/useFeaturedQuizzes";
import { useUserStats } from "@/hooks/useUserStats";
import { FeaturedQuizzesDashboard } from "@/components/dashboard/FeaturedQuizzesDashboard";
import { DashboardStatsCard } from "@/components/dashboard/DashboardStatsCard";

export default function DashboardPage() {
  const { data: featured } = useFeaturedQuizzes();
  const { data: stats } = useUserStats();

  return (
    <div className="space-y-8">
      <h1 className="text-2xl font-bold">Dashboard</h1>

      <DashboardStatsCard stats={stats} />

      <FeaturedQuizzesDashboard quizzes={featured} />
    </div>
  );
}
```

### Quiz Detail Page (Dynamic Route)

```tsx
// quizzes/[id]/page.tsx
"use client";

import { useParams } from "next/navigation";
import { useQuiz } from "@/hooks/useQuiz";
import { QuizDetail } from "@/components/quiz/QuizDetail";

export default function QuizPage() {
  const { id } = useParams<{ id: string }>();
  const { data: quiz, isLoading, error } = useQuiz(id);

  if (isLoading) return <Skeleton />;
  if (error) return <ErrorState error={error} />;
  if (!quiz) return <NotFound />;

  return <QuizDetail quiz={quiz} />;
}
```

### Quiz Taking Page

```tsx
// quizzes/[id]/take/page.tsx
"use client";

import { useParams, useRouter } from "next/navigation";
import { useQuizStore } from "@/store/quizStore";
import { useQuizAttempt } from "@/hooks/useQuizAttempt";
import { QuestionCard } from "@/components/quiz/QuestionCard";
import { QuizTimer } from "@/components/quiz/QuizTimer";
import { QuizProgress } from "@/components/quiz/QuizProgress";

export default function TakeQuizPage() {
  const { id } = useParams<{ id: string }>();
  const router = useRouter();
  const { currentQuiz, currentQuestionIndex, answers, setAnswer } = useQuizStore();
  const { submitAttempt, isSubmitting } = useQuizAttempt(id);

  const handleSubmit = async () => {
    const result = await submitAttempt(Object.values(answers));
    router.push(`/quizzes/${id}/results/${result.attempt.id}`);
  };

  return (
    <div className="max-w-3xl mx-auto space-y-6">
      <div className="flex justify-between items-center">
        <QuizProgress current={currentQuestionIndex} total={currentQuiz.question_count} />
        <QuizTimer />
      </div>

      <QuestionCard
        question={currentQuiz.questions[currentQuestionIndex]}
        answer={answers[currentQuiz.questions[currentQuestionIndex].id]}
        onAnswer={setAnswer}
      />

      <QuizNavigation onSubmit={handleSubmit} isSubmitting={isSubmitting} />
    </div>
  );
}
```

### Results Page (Nested Dynamic)

```tsx
// quizzes/[id]/results/[attemptId]/page.tsx
"use client";

import { useParams } from "next/navigation";

export default function ResultsPage() {
  const { id, attemptId } = useParams<{ id: string; attemptId: string }>();

  return <QuizResults quizId={id} attemptId={attemptId} />;
}
```

### Profile Page (View Others)

```tsx
// profile/[userId]/page.tsx
"use client";

import { useParams } from "next/navigation";
import { useProfile } from "@/hooks/useProfile";
import { UserProfileCard } from "@/components/profile/UserProfileCard";

export default function UserProfilePage() {
  const { userId } = useParams<{ userId: string }>();
  const { data: profile, isLoading } = useProfile(userId);

  if (isLoading) return <Skeleton />;

  return (
    <div className="space-y-6">
      <UserProfileCard profile={profile} />
      <DetailedStatistics userId={userId} />
      <AttemptHistory userId={userId} />
    </div>
  );
}
```

## Authentication Protection

All dashboard routes are protected by `AuthGuard`:

```tsx
// AuthGuard checks:
// 1. User is authenticated (session exists)
// 2. Session is valid (not expired)
// 3. Redirects to /login if not authenticated

<AuthGuard requireAuth={true}>
  {children}
</AuthGuard>
```

If user is not authenticated:
- Redirected to `/login`
- Return URL preserved for post-login redirect

## Real-time Features

`DashboardNotificationListener` enables:
- Toast notifications for new notifications
- Polling for unread count
- Real-time updates

```tsx
// In layout
<DashboardNotificationListener />

// Shows toasts when new notifications arrive
// Updates notification badge count
```

## Dynamic Route Parameters

### Single Parameter

```tsx
// /quizzes/[id]
const { id } = useParams<{ id: string }>();
```

### Multiple Parameters

```tsx
// /quizzes/[id]/results/[attemptId]
const { id, attemptId } = useParams<{
  id: string;
  attemptId: string;
}>();
```

### Catch-All (if needed)

```tsx
// /docs/[...slug]
const { slug } = useParams<{ slug: string[] }>();
// /docs/a/b/c -> slug = ["a", "b", "c"]
```

## Page Metadata

Each page can define metadata:

```tsx
// Static metadata
export const metadata = {
  title: "Quizzes | QuizNinja",
  description: "Browse and take quizzes",
};

// Dynamic metadata (for dynamic routes)
export async function generateMetadata({ params }: Props): Promise<Metadata> {
  const quiz = await getQuiz(params.id);
  return {
    title: `${quiz.title} | QuizNinja`,
  };
}
```

## Related Documentation

- [Parent: App Router](../README.md)
- [Auth Routes](../\(auth\)/README.md)
- [Layout Components](../../components/layout/README.md)
- [Quiz Components](../../components/quiz/README.md)
- [Auth Components](../../components/auth/README.md)
