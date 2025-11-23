"use client";

import { QuizCard } from "./QuizCard";
import { Skeleton } from "@/components/ui/skeleton";
import { EmptyState } from "@/components/common/EmptyState";
import type { Quiz, QuizAttempt } from "@/types/quiz";
import { BookOpen } from "lucide-react";

interface QuizListProps {
  quizzes: Quiz[];
  isLoading?: boolean;
  error?: Error | null;
  completedQuizMap?: Map<string, QuizAttempt>;
}

export function QuizList({ quizzes, isLoading, error, completedQuizMap }: QuizListProps) {
  // Loading state
  if (isLoading) {
    return (
      <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
        {Array.from({ length: 6 }).map((_, index) => (
          <QuizCardSkeleton key={index} />
        ))}
      </div>
    );
  }

  // Error state
  if (error) {
    return (
      <EmptyState
        icon={BookOpen}
        title="Error Loading Quizzes"
        description={
          error.message || "Failed to load quizzes. Please try again."
        }
      />
    );
  }

  // Empty state
  if (!quizzes || quizzes.length === 0) {
    return (
      <EmptyState
        icon={BookOpen}
        title="No Quizzes Found"
        description="No quizzes match your criteria. Try adjusting your filters."
      />
    );
  }

  // Quiz list
  return (
    <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
      {quizzes.map((quiz) => (
        <QuizCard
          key={quiz.id}
          quiz={quiz}
          completedAttempt={completedQuizMap?.get(quiz.id)}
        />
      ))}
    </div>
  );
}

// Skeleton loader for quiz cards
function QuizCardSkeleton() {
  return (
    <div className="border border-gray-200/60 dark:border-gray-800/60 rounded-xl p-6 space-y-4 bg-white/50 dark:bg-background/50 backdrop-blur-sm shadow-sm">
      {/* Top bar skeleton */}
      <Skeleton className="h-1.5 w-full rounded-full -mt-6 -mx-6 mb-4" style={{ width: 'calc(100% + 3rem)' }} />
      <div className="space-y-2">
        <div className="flex gap-2">
          <Skeleton className="h-5 w-20 rounded-full" />
          <Skeleton className="h-5 w-16 rounded-full" />
        </div>
        <Skeleton className="h-6 w-3/4 rounded-lg" />
        <Skeleton className="h-4 w-full rounded-lg" />
        <Skeleton className="h-4 w-2/3 rounded-lg" />
      </div>
      <div className="grid grid-cols-2 gap-3 py-3 border-t border-b border-gray-100 dark:border-gray-800">
        <Skeleton className="h-4 w-full rounded-lg" />
        <Skeleton className="h-4 w-full rounded-lg" />
        <Skeleton className="h-4 w-full rounded-lg" />
        <Skeleton className="h-4 w-full rounded-lg" />
      </div>
      <div className="flex gap-3 pt-2">
        <Skeleton className="h-10 flex-1 rounded-xl" />
        <Skeleton className="h-10 flex-1 rounded-xl" />
      </div>
    </div>
  );
}